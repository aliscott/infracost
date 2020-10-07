package terraform

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/infracost/infracost/internal/spin"
	"github.com/infracost/infracost/pkg/schema"
	"github.com/kballard/go-shellquote"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/mod/semver"

	"github.com/urfave/cli/v2"
)

var minTerraformVer = "v0.12"

type terraformProvider struct {
	jsonFile  string
	planFile  string
	dir       string
	planFlags string
}

// New returns new Terraform Provider
func New() schema.Provider {
	return &terraformProvider{}
}

func (p *terraformProvider) ProcessArgs(c *cli.Context) error {
	p.jsonFile = c.String("tfjson")
	p.planFile = c.String("tfplan")
	p.dir = c.String("tfdir")
	p.planFlags = c.String("tfflags")

	if p.jsonFile != "" && p.planFile != "" {
		return errors.New("Please provide either a Terraform Plan JSON file (tfjson) or a Terraform Plan file (tfplan)")
	}

	if p.planFile != "" && p.dir == "" {
		return errors.New("Please provide a path to the Terraform project (tfdir) if providing a Terraform Plan file (tfplan)")
	}

	return nil
}

func (p *terraformProvider) LoadResources() ([]*schema.Resource, error) {
	plan, err := p.loadPlanJSON()
	if err != nil {
		return []*schema.Resource{}, err
	}

	resources, err := parsePlanJSON(plan)
	if err != nil {
		return resources, errors.Wrap(err, "Error parsing plan JSON")
	}

	return resources, nil
}

func (p *terraformProvider) loadPlanJSON() ([]byte, error) {
	if p.jsonFile == "" {
		return p.generatePlanJSON()
	}

	f, err := os.Open(p.jsonFile)
	if err != nil {
		return []byte{}, errors.Wrapf(err, "Error reading plan file")
	}
	defer f.Close()

	out, err := ioutil.ReadAll(f)
	if err != nil {
		return []byte{}, errors.Wrapf(err, "Error reading plan file")
	}

	return out, nil
}

func (p *terraformProvider) generatePlanJSON() ([]byte, error) {
	err := p.terraformPreChecks()
	if err != nil {
		return []byte{}, err
	}

	opts := &CmdOptions{
		TerraformDir: p.dir,
	}

	var spinner *spin.Spinner

	if p.planFile == "" {
		if !p.isTerraformInitRun() {
			spinner = spin.NewSpinner("Running terraform init")
			_, err := TerraformCmd(opts, "init", "-no-color")
			if err != nil {
				spinner.Fail()
				terraformError(err)
				return []byte{}, errors.Wrap(err, "Error running terraform init")
			}
			spinner.Success()
		}

		spinner = spin.NewSpinner("Running terraform plan")
		f, err := ioutil.TempFile(os.TempDir(), "tfplan")
		if err != nil {
			spinner.Fail()
			return []byte{}, errors.Wrap(err, "Error creating temporary file 'tfplan'")
		}
		defer os.Remove(f.Name())

		flags, err := shellquote.Split(p.planFlags)
		if err != nil {
			return []byte{}, errors.Wrap(err, "Error parsing terraform plan flags")
		}
		args := []string{"plan", "-input=false", "-lock=false", "-no-color"}
		args = append(args, flags...)
		args = append(args, fmt.Sprintf("-out=%s", f.Name()))
		_, err = TerraformCmd(opts, args...)
		if err != nil {
			spinner.Fail()
			terraformError(err)
			return []byte{}, errors.Wrap(err, "Error running terraform plan")
		}
		spinner.Success()

		p.planFile = f.Name()
	}

	spinner = spin.NewSpinner("Running terraform show")
	out, err := TerraformCmd(opts, "show", "-no-color", "-json", p.planFile)
	if err != nil {
		spinner.Fail()
		terraformError(err)
		return []byte{}, errors.Wrap(err, "Error running terraform show")
	}
	spinner.Success()

	return out, nil
}

func (p *terraformProvider) terraformPreChecks() error {
	if p.jsonFile == "" {
		_, err := exec.LookPath(terraformBinary())
		if err != nil {
			return errors.Errorf("Could not find Terraform binary \"%s\" in path.\nYou can set a custom Terraform binary using the environment variable TERRAFORM_BINARY.", terraformBinary())
		}

		if v, ok := checkTerraformVersion(); !ok {
			return errors.Errorf("Terraform %s is not supported. Please use Terraform version >= %s.", v, minTerraformVer)
		}

		if !p.inTerraformDir() {
			return errors.Errorf("Directory \"%s\" does not have any .tf files.\nYou can pass a path to a Terraform directory using the --tfdir option.", p.dir)
		}
	}
	return nil
}

func checkTerraformVersion() (string, bool) {
	out, err := TerraformVersion()
	if err != nil {
		// If we encounter any errors here we just return true
		// since it might be caused by a custom Terraform binary
		return "", true
	}
	p := strings.Split(out, " ")
	v := p[len(p)-1]

	// Allow any non-terraform binaries, e.g. terragrunt
	if !strings.HasPrefix(out, "Terraform ") {
		return v, true
	}

	return v, semver.Compare(v, minTerraformVer) >= 0
}

func (p *terraformProvider) inTerraformDir() bool {
	matches, err := filepath.Glob(filepath.Join(p.dir, "*.tf"))
	return matches != nil && err == nil
}

func (p *terraformProvider) isTerraformInitRun() bool {
	_, err := os.Stat(filepath.Join(p.dir, ".terraform"))
	if err != nil {
		if !os.IsNotExist(err) {
			log.Errorf("error checking if .terraform directory exists: %v", err)
		}
		return false
	}
	return true
}

func terraformError(err error) {
	if e, ok := err.(*TerraformCmdError); ok {
		fmt.Fprintln(os.Stderr, indent(color.HiRedString("Terraform command failed with:"), "  "))
		stderr := stripBlankLines(string(e.Stderr))
		fmt.Fprintln(os.Stderr, indent(color.HiRedString(stderr), "    "))
		if strings.HasPrefix(stderr, "Error: No value for required variable") {
			fmt.Fprintln(os.Stderr, color.HiRedString("You can pass any Terraform args using the --tfflags option."))
			fmt.Fprintln(os.Stderr, color.HiRedString("For example: infracost --tfdir=path/to/terraform --tfflags=\"-var-file=myvars.tf\"\n"))
		}
		if strings.HasPrefix(stderr, "Error: Failed to read variables file") {
			fmt.Fprintln(os.Stderr, color.HiRedString("You should specify the -var-file flag as a path relative to your Terraform directory."))
			fmt.Fprintln(os.Stderr, color.HiRedString("For example: infracost --tfdir=path/to/terraform --tfflags=\"-var-file=myvars.tf\"\n"))
		}
	}
}

func indent(s, indent string) string {
	result := ""
	for _, j := range strings.Split(strings.TrimRight(s, "\n"), "\n") {
		result += indent + j + "\n"
	}
	return result
}

func stripBlankLines(s string) string {
	return regexp.MustCompile(`[\t\r\n]+`).ReplaceAllString(strings.TrimSpace(s), "\n")
}
