package aws

import (
	"github.com/infracost/infracost/pkg/schema"

	"github.com/shopspring/decimal"
)

func GetEBSVolumeRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:  "aws_ebs_volume",
		RFunc: NewEBSVolume,
	}
}

func NewEBSVolume(d *schema.ResourceData, u *schema.ResourceData) *schema.Resource {
	region := d.Get("region").String()

	volumeApiName := "gp2"
	if d.Get("type").Exists() {
		volumeApiName = d.Get("type").String()
	}

	gbVal := decimal.NewFromInt(int64(defaultVolumeSize))
	if d.Get("size").Exists() {
		gbVal = decimal.NewFromFloat(d.Get("size").Float())
	}

	iopsVal := decimal.Zero
	if d.Get("iops").Exists() {
		iopsVal = decimal.NewFromFloat(d.Get("iops").Float())
	}

	return &schema.Resource{
		Name:           d.Address,
		CostComponents: ebsVolumeCostComponents(region, volumeApiName, gbVal, iopsVal),
	}
}

func ebsVolumeCostComponents(region string, volumeApiName string, gbVal decimal.Decimal, iopsVal decimal.Decimal) []*schema.CostComponent {
	var name string
	switch volumeApiName {
	case "standard":
		name = "Magnetic storage"
	case "io1":
		name = "Provisioned IOPS SSD storage (io1)"
	case "io2":
		name = "Provisioned IOPS SSD storage (io2)"
	case "st1":
		name = "Throughput Optimized HDD storage (st1)"
	case "sc1":
		name = "Cold HDD storage (sc1)"
	default:
		name = "General Purpose SSD storage (gp2)"
	}

	costComponents := []*schema.CostComponent{
		{
			Name:            name,
			Unit:            "GB-months",
			MonthlyQuantity: &gbVal,
			ProductFilter: &schema.ProductFilter{
				VendorName:    strPtr("aws"),
				Region:        strPtr(region),
				Service:       strPtr("AmazonEC2"),
				ProductFamily: strPtr("Storage"),
				AttributeFilters: []*schema.AttributeFilter{
					{Key: "volumeApiName", Value: strPtr(volumeApiName)},
				},
			},
		},
	}

	if volumeApiName == "io1" || volumeApiName == "io2" {
		costComponents = append(costComponents, &schema.CostComponent{
			Name:            "Provisioned IOPS",
			Unit:            "IOPS-months",
			MonthlyQuantity: &iopsVal,
			ProductFilter: &schema.ProductFilter{
				VendorName:    strPtr("aws"),
				Region:        strPtr(region),
				Service:       strPtr("AmazonEC2"),
				ProductFamily: strPtr("System Operation"),
				AttributeFilters: []*schema.AttributeFilter{
					{Key: "volumeApiName", Value: strPtr(volumeApiName)},
					{Key: "usagetype", ValueRegex: strPtr("/EBS:VolumeP-IOPS/")},
				},
			},
		})
	}

	if volumeApiName == "standard" {
		costComponents = append(costComponents, &schema.CostComponent{
			Name:            "I/O requests",
			Unit:            "Per request",
			MonthlyQuantity: decimalPtr(decimal.Zero),
			ProductFilter: &schema.ProductFilter{
				VendorName:    strPtr("aws"),
				Region:        strPtr(region),
				Service:       strPtr("AmazonEC2"),
				ProductFamily: strPtr("System Operation"),
				AttributeFilters: []*schema.AttributeFilter{
					{Key: "volumeApiName", Value: strPtr(volumeApiName)},
					{Key: "usagetype", ValueRegex: strPtr("/EBS:VolumeIOUsage/")},
				},
			},
		})
	}

	return costComponents
}
