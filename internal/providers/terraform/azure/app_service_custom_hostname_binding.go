package azure

import (
	"fmt"
	"strings"

	"github.com/infracost/infracost/internal/schema"
	"github.com/shopspring/decimal"
	"github.com/tidwall/gjson"
)

func GetAzureRMAppServiceCustomHostnameBindingRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:  "azurerm_app_service_custom_hostname_binding",
		RFunc: NewAzureRMAppServiceCustomHostnameBinding,
		ReferenceAttributes: []string{
			"resource_group_name",
		},
	}
}

func NewAzureRMAppServiceCustomHostnameBinding(d *schema.ResourceData, u *schema.UsageData) *schema.Resource {
	var sslType, sslState, location string

	group := d.References("resource_group_name")
	location = group[0].Get("location").String()

	if d.Get("ssl_state").Type != gjson.Null {
		sslState = d.Get("ssl_state").String()
	} else {
		// returning directly since SNI is currently defined as free in the Azure cost page
		return &schema.Resource{
			NoPrice:   true,
			IsSkipped: true,
		}
	}

	// The two approved values are IpBasedEnabled or SniEnabled
	sslState = strings.ToUpper(sslState)[0:2]

	if strings.ToLower(sslState) == "ip" {
		sslType = "IP"
	} else {
		// returning directly since SNI is currently defined as free in the Azure cost page
		return &schema.Resource{
			NoPrice:   true,
			IsSkipped: true,
		}
	}

	var instanceCount int64 = 1

	costComponents := []*schema.CostComponent{
		{
			Name:            "IP SSL certificate",
			Unit:            "months",
			UnitMultiplier:  1,
			MonthlyQuantity: decimalPtr(decimal.NewFromInt(instanceCount)),
			ProductFilter: &schema.ProductFilter{
				VendorName:    strPtr("azure"),
				Region:        strPtr(location),
				Service:       strPtr("Azure App Service"),
				ProductFamily: strPtr("Compute"),
				AttributeFilters: []*schema.AttributeFilter{
					{Key: "skuName", Value: strPtr(fmt.Sprintf("%s SSL", sslType))},
				},
			},
			PriceFilter: &schema.PriceFilter{
				PurchaseOption: strPtr("Consumption"),
			},
		},
	}

	return &schema.Resource{
		Name:           d.Address,
		CostComponents: costComponents,
	}
}
