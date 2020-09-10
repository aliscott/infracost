package aws

import (
	"github.com/infracost/infracost/pkg/schema"
	"github.com/shopspring/decimal"
)

func NewElasticsearchDomain(d *schema.ResourceData, u *schema.ResourceData) *schema.Resource {

	region := d.Get("region").String()
	clusterConfig := d.Get("cluster_config").Array()[0]
	instanceType := clusterConfig.Get("instance_type").String()
	instanceCount := clusterConfig.Get("instance_count").Int()
	ebsOptions := d.Get("ebs_options").Array()[0]

	ebsTypeMap := map[string]string{
		"gp2":      "GP2",
		"io1":      "PIOPS-Storage",
		"standard": "Magnetic",
	}

	gbVal := decimal.NewFromInt(int64(defaultVolumeSize))
	if ebsOptions.Get("volume_size").Exists() {
		gbVal = decimal.NewFromFloat(ebsOptions.Get("volume_size").Float())
	}

	ebsType := "gp2"
	if ebsOptions.Get("volume_type").Exists() {
		ebsType = ebsOptions.Get("volume_type").String()
	}

	ebsFilter := "gp2"
	if val, ok := ebsTypeMap[ebsType]; ok {
		ebsFilter = val
	}

	iopsVal := decimal.NewFromInt(1)
	if ebsOptions.Get("iops").Exists() {
		iopsVal = decimal.NewFromFloat(ebsOptions.Get("iops").Float())

		if iopsVal.LessThan(decimal.NewFromInt(1)) {
			iopsVal = decimal.NewFromInt(1)
		}
	}

	costComponents := []*schema.CostComponent{
		{
			Name:           "Per instance hour",
			Unit:           "hours",
			HourlyQuantity: decimalPtr(decimal.NewFromInt(instanceCount)),
			ProductFilter: &schema.ProductFilter{
				VendorName:    strPtr("aws"),
				Region:        strPtr(region),
				Service:       strPtr("AmazonES"),
				ProductFamily: strPtr("Elastic Search Instance"),
				AttributeFilters: []*schema.AttributeFilter{
					{Key: "usagetype", ValueRegex: strPtr("/ESInstance/")},
					{Key: "instanceType", Value: &instanceType},
				},
			},
			PriceFilter: &schema.PriceFilter{
				PurchaseOption: strPtr("on_demand"),
			},
		},
		{
			Name:            "Storage",
			Unit:            "GB-months",
			MonthlyQuantity: &gbVal,
			ProductFilter: &schema.ProductFilter{
				VendorName:    strPtr("aws"),
				Region:        strPtr(region),
				Service:       strPtr("AmazonES"),
				ProductFamily: strPtr("Elastic Search Volume"),
				AttributeFilters: []*schema.AttributeFilter{
					{Key: "usagetype", ValueRegex: strPtr("/ES.+-Storage/")},
					{Key: "storageMedia", Value: strPtr(ebsFilter)},
				},
			},
			PriceFilter: &schema.PriceFilter{
				PurchaseOption: strPtr("on_demand"),
			},
		},
	}

	if ebsType == "io1" || ebsType == "io2" {
		costComponents = append(costComponents, &schema.CostComponent{
			Name:            "Storage IOPS",
			Unit:            "IOPS-months",
			MonthlyQuantity: &iopsVal,
			ProductFilter: &schema.ProductFilter{
				VendorName:    strPtr("aws"),
				Region:        strPtr(region),
				Service:       strPtr("AmazonES"),
				ProductFamily: strPtr("Elastic Search Volume"),
				AttributeFilters: []*schema.AttributeFilter{
					{Key: "usagetype", ValueRegex: strPtr("/ES:PIOPS/")},
					{Key: "storageMedia", Value: strPtr("PIOPS")},
				},
			},
			PriceFilter: &schema.PriceFilter{
				PurchaseOption: strPtr("on_demand"),
			},
		})
	}

	return &schema.Resource{
		Name:           d.Address,
		CostComponents: costComponents,
	}
}
