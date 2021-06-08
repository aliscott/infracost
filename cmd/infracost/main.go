package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/infracost/infracost/internal/config"
	"github.com/infracost/infracost/internal/events"
	"github.com/infracost/infracost/internal/ui"
	"github.com/infracost/infracost/internal/update"
	"github.com/infracost/infracost/internal/version"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/fatih/color"
)

var spinner *ui.Spinner

func main() {
	var appErr error
	updateMessageChan := make(chan *update.Info)

	cfg := config.DefaultConfig()
	appErr = cfg.LoadFromEnv()

	defer func() {
		if appErr != nil {
			handleAppErr(cfg, appErr)
		}

		unexpectedErr := recover()
		if unexpectedErr != nil {
			handleUnexpectedErr(cfg, unexpectedErr)
		}

		handleUpdateMessage(updateMessageChan)

		if appErr != nil || unexpectedErr != nil {
			os.Exit(1)
		}
	}()

	startUpdateCheck(cfg, updateMessageChan)

	rootCmd := &cobra.Command{
		Use:     "infracost",
		Version: version.Version,
		Short:   "Cloud cost estimates for Terraform",
		Long: fmt.Sprintf(`Infracost - cloud cost estimates for Terraform

%s
  https://infracost.io/docs`, ui.BoldString("DOCS")),
		Example: `  Generate a cost diff from Terraform directory with any required Terraform flags:

      infracost diff --path /path/to/code --terraform-plan-flags "-var-file=my.tfvars"
	
  Generate a full cost breakdown from Terraform directory with any required Terraform flags:

      infracost breakdown --path /path/to/code --terraform-plan-flags "-var-file=my.tfvars"`,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			cfg.Environment.Command = cmd.Name()

			return loadGlobalFlags(cfg, cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Show the help
			return cmd.Help()
		},
	}

	rootCmd.PersistentFlags().Bool("no-color", false, "Turn off colored output")
	rootCmd.PersistentFlags().String("log-level", "", "Log level (trace, debug, info, warn, error, fatal)")

	rootCmd.AddCommand(registerCmd(cfg))
	rootCmd.AddCommand(diffCmd(cfg))
	rootCmd.AddCommand(breakdownCmd(cfg))
	rootCmd.AddCommand(outputCmd(cfg))
	rootCmd.AddCommand(completionCmd())

	rootCmd.SetUsageTemplate(fmt.Sprintf(`%s{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

%s
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

%s
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

%s{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

%s
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

%s
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

%s{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`,
		ui.BoldString("USAGE"),
		ui.BoldString("ALIAS"),
		ui.BoldString("EXAMPLES"),
		ui.BoldString("AVAILABLE COMMANDS"),
		ui.BoldString("FLAGS"),
		ui.BoldString("GLOBAL FLAGS"),
		ui.BoldString("ADDITIONAL HELP TOPICS"),
	))

	rootCmd.SetVersionTemplate("Infracost {{.Version}}\n")

	appErr = rootCmd.Execute()
}

func startUpdateCheck(cfg *config.Config, c chan *update.Info) {
	go func() {
		updateInfo, err := update.CheckForUpdate(cfg)
		if err != nil {
			log.Debugf("error checking for update: %v", err)
		}
		c <- updateInfo
		close(c)
	}()
}

func checkAPIKey(apiKey string, apiEndpoint string, defaultEndpoint string) error {
	if apiEndpoint == defaultEndpoint && apiKey == "" {
		return errors.New(fmt.Sprintf(
			"No INFRACOST_API_KEY environment variable is set.\nWe run a free Cloud Pricing API, to get an API key run %s",
			ui.PrimaryString("infracost register"),
		))
	}

	return nil
}

func handleAppErr(cfg *config.Config, err error) {
	if spinner != nil {
		spinner.Fail()
		fmt.Fprintln(os.Stderr, "")
	}

	if err.Error() != "" {
		ui.PrintError(err.Error())
	}

	msg := ui.StripColor(err.Error())
	var eventsError *events.Error
	if errors.As(err, &eventsError) {
		msg = ui.StripColor(eventsError.Label)
	}
	events.SendReport(cfg, "error", msg)
}

func handleUnexpectedErr(cfg *config.Config, unexpectedErr interface{}) {
	if spinner != nil {
		spinner.Fail()
		fmt.Fprintln(os.Stderr, "")
	}

	stack := string(debug.Stack())

	ui.PrintUnexpectedError(unexpectedErr, stack)

	events.SendReport(cfg, "error", fmt.Sprintf("%s\n%s", unexpectedErr, stack))
}

func handleUpdateMessage(updateMessageChan chan *update.Info) {
	updateInfo := <-updateMessageChan
	if updateInfo != nil {
		msg := fmt.Sprintf("\n%s %s %s → %s\n%s\n",
			ui.WarningString("Update:"),
			"A new version of Infracost is available:",
			ui.PrimaryString(version.Version),
			ui.PrimaryString(updateInfo.LatestVersion),
			ui.Indent(updateInfo.Cmd, "  "),
		)
		fmt.Fprint(os.Stderr, msg)
	}
}

func loadGlobalFlags(cfg *config.Config, cmd *cobra.Command) error {
	if cmd.Flags().Changed("no-color") {
		cfg.NoColor, _ = cmd.Flags().GetBool("no-color")
	}
	color.NoColor = cfg.NoColor

	if cmd.Flags().Changed("log-level") {
		cfg.LogLevel, _ = cmd.Flags().GetString("log-level")
		err := cfg.ConfigureLogger()
		if err != nil {
			return err
		}
	}

	if cmd.Flags().Changed("pricing-api-endpoint") {
		cfg.PricingAPIEndpoint, _ = cmd.Flags().GetString("pricing-api-endpoint")
	}

	cfg.Environment.IsDefaultPricingAPIEndpoint = cfg.PricingAPIEndpoint == cfg.DefaultPricingAPIEndpoint

	flagNames := make([]string, 0)

	cmd.Flags().Visit(func(f *pflag.Flag) {
		flagNames = append(flagNames, f.Name)
	})

	cfg.Environment.Flags = flagNames

	return nil
}
