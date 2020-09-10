package terraform

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/infracost/infracost/pkg/schema"
	"github.com/pkg/errors"

	"github.com/urfave/cli/v2"
)

type terraformProvider struct {
	jsonFile string
	planFile string
	dir      string
}

// New returns new Terraform Provider
func New() schema.Provider {
	return &terraformProvider{}
}

func (p *terraformProvider) ProcessArgs(c *cli.Context) error {
	p.jsonFile = c.String("tfjson")
	p.planFile = c.String("tfplan")
	p.dir = c.String("tfdir")

	if p.jsonFile != "" && p.planFile != "" {
		return errors.New("Please provide either a Terraform Plan JSON file (tfjson) or a Terraform Plan file (tfplan)")
	}

	if p.planFile != "" && p.dir == "" {
		return errors.New("Please provide a path to the Terraform project (tfdir) if providing a Terraform Plan file (tfplan)")
	}

	return nil
}

func (p *terraformProvider) LoadResources() ([]*schema.Resource, error) {
	var planJSON []byte
	var err error

	if p.jsonFile != "" {
		planJSON, err = loadPlanJSON(p.jsonFile)
	} else {
		planJSON, err = generatePlanJSON(p.dir, p.planFile)
	}

	if err != nil {
		return []*schema.Resource{}, errors.Wrap(err, "error loading resources")
	}

	schemaResources := parsePlanJSON(planJSON)
	return schemaResources, nil
}

func loadPlanJSON(path string) ([]byte, error) {
	planFile, err := os.Open(path)
	if err != nil {
		return []byte{}, errors.Wrapf(err, "error opening file '%v'", path)
	}
	defer planFile.Close()

	out, err := ioutil.ReadAll(planFile)
	if err != nil {
		return []byte{}, errors.Wrapf(err, "error reading file '%v'", path)
	}

	return out, nil
}

func generatePlanJSON(tfdir string, planPath string) ([]byte, error) {
	var err error

	opts := &cmdOptions{
		TerraformDir: tfdir,
	}

	if planPath == "" {
		_, err = terraformCmd(opts, "init")
		if err != nil {
			return []byte{}, errors.Wrap(err, "error initializing terraform working directory")
		}

		planfile, err := ioutil.TempFile(os.TempDir(), "tfplan")
		if err != nil {
			return []byte{}, errors.Wrap(err, "error creating temporary file 'tfplan'")
		}
		defer os.Remove(planfile.Name())

		_, err = terraformCmd(opts, "plan", "-input=false", "-lock=false", fmt.Sprintf("-out=%s", planfile.Name()))
		if err != nil {
			return []byte{}, errors.Wrap(err, "error generating terraform execution plan")
		}

		planPath = planfile.Name()
	}

	out, err := terraformCmd(opts, "show", "-json", planPath)
	if err != nil {
		return []byte{}, errors.Wrap(err, "error inspecting terraform plan")
	}

	return out, nil
}
