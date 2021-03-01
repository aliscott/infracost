package aws

import (
	"fmt"
	"github.com/infracost/infracost/internal/schema"
	log "github.com/sirupsen/logrus"
	"strings"

	"github.com/shopspring/decimal"
)

func GetDXConnectionRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:  "aws_dx_connection",
		RFunc: NewDXConnection,
	}
}

func NewDXConnection(d *schema.ResourceData, u *schema.UsageData) *schema.Resource {
	region := d.Get("region").String()
	fromLocation, ok := regionMapping[region]

	if !ok {
		log.Warnf("Skipping resource %s. Could not find mapping for region %s", d.Address, region)
		return nil
	}

	var gbDataProcessed *decimal.Decimal
	if u != nil && u.Get("monthly_outbound_region_to_dx_location_gb").Exists() {
		gbDataProcessed = decimalPtr(decimal.NewFromFloat(u.Get("monthly_outbound_region_to_dx_location_gb").Float()))
	}

	dxBandwidth := strings.Replace(d.Get("location").String(), "bps", "", 1)
	dxLocation := d.Get("location").String()

	connectionType := "Dedicated"
	if u != nil && u.Get("dx_connection_type").Exists() {
		connectionType = u.Get("dx_connection_type").String()
	}

	virtualInterfaceType := "Private"
	if u != nil && u.Get("dx_virtual_interface_type").Exists() {
		virtualInterfaceType = u.Get("dx_virtual_interface_type").String()
	}

	return &schema.Resource{
		Name: d.Address,
		CostComponents: []*schema.CostComponent{
			{
				Name:           "DX connection",
				Unit:           "hours",
				UnitMultiplier: 1,
				HourlyQuantity: decimalPtr(decimal.NewFromInt(1)),
				ProductFilter: &schema.ProductFilter{
					VendorName:    strPtr("aws"),
					Region:        strPtr(region),
					Service:       strPtr("AWSDirectConnect"),
					ProductFamily: strPtr("Direct Connect"),
					AttributeFilters: []*schema.AttributeFilter{
						{Key: "capacity", ValueRegex: strPtr(dxBandwidth)},
						{Key: "usagetype", ValueRegex: strPtr(fmt.Sprintf("/%s/", dxLocation))},
						{Key: "connectionType", Value: strPtr(connectionType)},
					},
				},
			},
			{
				Name:            fmt.Sprintf("Outbound data transfer to dx location %s", dxLocation),
				Unit:            "GB",
				UnitMultiplier:  1,
				MonthlyQuantity: gbDataProcessed,
				ProductFilter: &schema.ProductFilter{
					VendorName:    strPtr("aws"),
					Service:       strPtr("AWSDirectConnect"),
					ProductFamily: strPtr("Data Transfer"),
					AttributeFilters: []*schema.AttributeFilter{
						{Key: "fromLocation", Value: strPtr(fromLocation)},
						{Key: "usagetype", ValueRegex: strPtr(fmt.Sprintf("/%s-DataXfer-Out/", dxLocation))},
						{Key: "virtualInterfaceType", Value: strPtr(virtualInterfaceType)},
					},
				},
			},
		},
	}
}