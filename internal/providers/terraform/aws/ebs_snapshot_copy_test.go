package aws_test

import (
	"testing"

	"github.com/infracost/infracost/internal/schema"
	"github.com/infracost/infracost/internal/testutil"

	"github.com/infracost/infracost/internal/providers/terraform/tftest"

	"github.com/shopspring/decimal"
)

func TestEBSSnapshotCopy(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	tf := `
		resource "aws_ebs_volume" "gp2" {
			availability_zone = "us-east-1a"
			size              = 10
		}

		resource "aws_ebs_snapshot" "gp2" {
			volume_id = aws_ebs_volume.gp2.id
		}

		resource "aws_ebs_snapshot_copy" "gp2" {
			source_snapshot_id = aws_ebs_snapshot.gp2.id
			source_region      = "us-east-1"
		}`

	resourceChecks := []testutil.ResourceCheck{
		{
			Name:      "aws_ebs_volume.gp2",
			SkipCheck: true,
		},
		{
			Name:      "aws_ebs_snapshot.gp2",
			SkipCheck: true,
		},
		{
			Name: "aws_ebs_snapshot_copy.gp2",
			CostComponentChecks: []testutil.CostComponentCheck{
				{
					Name:            "EBS snapshot storage",
					PriceHash:       "63a6765e67e0ebcd29f15f1570b5e692-ee3dd7e4624338037ca6fea0933a662f",
					HourlyCostCheck: testutil.MonthlyPriceMultiplierCheck(decimal.NewFromInt(10)),
				},
			},
		},
	}

	tftest.ResourceTests(t, tf, schema.NewEmptyUsageMap(), resourceChecks)
}
