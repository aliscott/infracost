package azure

import (
	"fmt"
	"strings"

	"github.com/infracost/infracost/internal/schema"

	"github.com/shopspring/decimal"
)

func GetAzureMariaDBServerRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:  "azurerm_mariadb_server",
		RFunc: NewAzureMariaDBServer,
	}
}

func NewAzureMariaDBServer(d *schema.ResourceData, u *schema.UsageData) *schema.Resource {
	var costComponents []*schema.CostComponent

	region := d.Get("location").String()
	sku := d.Get("sku_name").String()
	tier := strings.Split(sku, "_")[0]
	family := strings.Split(sku, "_")[1]
	cores := strings.Split(sku, "_")[2]

	tierName := map[string]string{
		"B":  "Basic",
		"GP": "General Purpose",
		"MO": "Memory Optimized",
	}[tier]

	productNameRegex := fmt.Sprintf("/%s - Compute %s/", tierName, family)
	skuName := fmt.Sprintf("%s vCore", cores)

	costComponents = append(costComponents, &schema.CostComponent{
		Name:           fmt.Sprintf("Compute (%s)", sku),
		Unit:           "hours",
		UnitMultiplier: 1,
		HourlyQuantity: decimalPtr(decimal.NewFromInt(1)),
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("azure"),
			Region:        strPtr(region),
			Service:       strPtr("Azure Database for MariaDB"),
			ProductFamily: strPtr("Databases"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "productName", ValueRegex: strPtr(productNameRegex)},
				{Key: "skuName", Value: strPtr(skuName)},
			},
		},
	})

	storageGB := d.Get("storage_mb").Int()

	// MO and GP storage cost are the same, and we don't have cost component for MO Storage now
	if tier == "MO" {
		tierName = "General Purpose"
	}
	productNameRegex = fmt.Sprintf("/%s - Storage/", tierName)

	costComponents = append(costComponents, &schema.CostComponent{
		Name:            "Storage",
		Unit:            "GB-months",
		UnitMultiplier:  1,
		MonthlyQuantity: decimalPtr(decimal.NewFromInt(storageGB)),
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("azure"),
			Region:        strPtr(region),
			Service:       strPtr("Azure Database for MariaDB"),
			ProductFamily: strPtr("Databases"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "productName", ValueRegex: strPtr(productNameRegex)},
			},
		},
	})

	var backupStorageGB *decimal.Decimal

	if u != nil && u.Get("additional_backup_storage_gb").Exists() {
		backupStorageGB = decimalPtr(decimal.NewFromInt(u.Get("additional_backup_storage_gb").Int()))
	}

	skuName = "Backup LRS"
	if d.Get("geo_redundant_backup_enabled").Exists() {
		if d.Get("geo_redundant_backup_enabled").String() == "true" {
			skuName = "Backup GRS"
		}
	}

	costComponents = append(costComponents, &schema.CostComponent{
		Name:            "Additional backup storage",
		Unit:            "GB-months",
		UnitMultiplier:  1,
		MonthlyQuantity: backupStorageGB,
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("azure"),
			Region:        strPtr(region),
			Service:       strPtr("Azure Database for MariaDB"),
			ProductFamily: strPtr("Databases"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "productName", ValueRegex: strPtr("/Backup Storage/")},
				{Key: "skuName", Value: strPtr(skuName)},
			},
		},
	})

	return &schema.Resource{
		Name:           d.Address,
		CostComponents: costComponents,
	}
}