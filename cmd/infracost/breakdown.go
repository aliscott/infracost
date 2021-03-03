package main

import (
	"github.com/infracost/infracost/internal/config"
	"github.com/infracost/infracost/internal/ui"
	"github.com/spf13/cobra"
)

func breakdownCmd(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "breakdown",
		Short: "Generates a full breakdown of costs",
		Long: `Generates a full breakdown of costs

Use terraform directory with any required terraform flags:

  infracost breakdown --path /path/to/code --terraform-plan-flags "-var-file=myvars.tfvars"

Use terraform state file:

  infracost breakdown --path /path/to/code --terraform-use-state

Use terraform plan JSON:

  terraform plan -out plan.save .
  terraform show -json plan.save > plan.json
  infracost breakdown --path /path/to/plan.json

Use terraform plan file:

  terraform plan -out plan.save .
  infracost breakdown --path plan.save`,
		PreRun: func(cmd *cobra.Command, args []string) {
			handleDeprecatedEnvVars(deprecatedEnvVarMapping)
			handleDeprecatedFlags(cmd, deprecatedFlagsMapping)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := checkAPIKey(cfg.APIKey, cfg.PricingAPIEndpoint, cfg.DefaultPricingAPIEndpoint); err != nil {
				return err
			}

			err := loadRunFlags(cfg, cmd)
			if err != nil {
				return err
			}

			err = checkRunConfig(cfg)
			if err != nil {
				ui.PrintUsageErrorAndExit(cmd, err.Error())
			}

			// Handle deprecated table output for root command
			if cmd.Name() == "infracost" && (cfg.Format == "" || cfg.Format == "table") {
				cfg.Format = "table_deprecated"
			}

			return runMain(cfg)
		},
	}

	addRunInputFlags(cmd)
	addRunOutputFlags(cmd)

	cmd.Flags().Bool("terraform-use-state", false, "Use Terraform state instead of generating a plan")
	cmd.Flags().String("format", "table", "Output format: json, table, html")

	return cmd
}