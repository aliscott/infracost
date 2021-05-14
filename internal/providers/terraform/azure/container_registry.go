package azure

import (
	"fmt"
	"strconv"

	"github.com/infracost/infracost/internal/schema"
	"github.com/shopspring/decimal"
	"github.com/tidwall/gjson"
)

func GetAzureRMContainerRegistryRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:  "azurerm_container_registry",
		RFunc: NewAzureRMContainerRegistry,
	}
}

func NewAzureRMContainerRegistry(d *schema.ResourceData, u *schema.UsageData) *schema.Resource {
	var nambersOflocations int
	var storageGb, includedStorage, monthlyBuildVcpu *decimal.Decimal
	var overStorage decimal.Decimal

	sku := "Classic"
	location := d.Get("location").String()

	if d.Get("sku").Exists() {
		sku = d.Get("sku").String()
	}

	switch sku {
	case "Classic":
		includedStorage = decimalPtr(decimal.NewFromFloat(10))
	case "Basic":
		includedStorage = decimalPtr(decimal.NewFromFloat(10))
	case "Standard":
		includedStorage = decimalPtr(decimal.NewFromFloat(100))
	case "Premium":
		includedStorage = decimalPtr(decimal.NewFromFloat(500))
	}

	if d.Get("georeplication_locations").Exists() {
		nambersOflocations = len(d.Get("georeplication_locations").Array())
	}
	if d.Get("georeplications").Exists() {
		nambersOflocations = len(d.Get("georeplications.0.location").Array())
	}

	costComponents := make([]*schema.CostComponent, 0)

	if d.Get("georeplication_locations").Exists() || d.Get("georeplications").Exists() {
		if nambersOflocations == 1 {
			costComponents = append(costComponents, ContainerRegistryGeolocationCostComponent(fmt.Sprintf("Geo replication (%s location)", strconv.Itoa(nambersOflocations)), location, sku))
		}
		if nambersOflocations > 1 {
			costComponents = append(costComponents, ContainerRegistryGeolocationCostComponent(fmt.Sprintf("Geo replication (%s locations)", strconv.Itoa(nambersOflocations)), location, sku))
		}
	}

	if d.Get("georeplication_locations").Type == gjson.Null && d.Get("georeplications").Type == gjson.Null {
		costComponents = append(costComponents, ContainerRegistryCostComponent(fmt.Sprintf("Registry usage (%s)", sku), location, sku))
	}

	if u != nil && u.Get("storage_gb").Type != gjson.Null {
		storageGb = decimalPtr(decimal.NewFromFloat(u.Get("storage_gb").Float()))
		if storageGb.GreaterThan(*includedStorage) {
			overStorage = storageGb.Sub(*includedStorage)
			storageGb = &overStorage
			costComponents = append(costComponents, ContainerRegistryStorageCostComponent(fmt.Sprintf("Storage (over %s GB)", includedStorage), location, sku, storageGb))
		}
	}

	if u != nil && u.Get("monthly_build_vcpu_hrs").Type != gjson.Null {
		monthlyBuildVcpu = decimalPtr(decimal.NewFromFloat(u.Get("monthly_build_vcpu_hrs").Float() * 3600))
	}

	costComponents = append(costComponents, ContainerRegistryCPUCostComponent("Build vCPU", location, sku, monthlyBuildVcpu))

	return &schema.Resource{
		Name:           d.Address,
		CostComponents: costComponents,
	}
}

func ContainerRegistryCostComponent(name, location, sku string) *schema.CostComponent {

	return &schema.CostComponent{

		Name:            name,
		Unit:            "days",
		UnitMultiplier:  1,
		MonthlyQuantity: decimalPtr(decimal.NewFromInt(30)),
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("azure"),
			Region:        strPtr(location),
			Service:       strPtr("Container Registry"),
			ProductFamily: strPtr("Containers"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "productName", Value: strPtr("Container Registry")},
				{Key: "skuName", Value: strPtr(sku)},
				{Key: "meterName", Value: strPtr(fmt.Sprintf("%s Registry Unit", sku))},
			},
		},
		PriceFilter: &schema.PriceFilter{
			PurchaseOption: strPtr("Consumption"),
		},
	}
}
func ContainerRegistryGeolocationCostComponent(name, location, sku string) *schema.CostComponent {

	return &schema.CostComponent{

		Name:            name,
		Unit:            "days",
		UnitMultiplier:  1,
		MonthlyQuantity: decimalPtr(decimal.NewFromInt(30)),
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("azure"),
			Region:        strPtr(location),
			Service:       strPtr("Container Registry"),
			ProductFamily: strPtr("Containers"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "productName", Value: strPtr("Container Registry")},
				{Key: "skuName", Value: strPtr(sku)},
				{Key: "meterName", Value: strPtr(fmt.Sprintf("%s Registry Unit", sku))},
			},
		},
		PriceFilter: &schema.PriceFilter{
			PurchaseOption: strPtr("Consumption"),
		},
	}
}
func ContainerRegistryStorageCostComponent(name, location, sku string, storage *decimal.Decimal) *schema.CostComponent {

	return &schema.CostComponent{

		Name:            name,
		Unit:            "GB-months",
		UnitMultiplier:  1,
		MonthlyQuantity: storage,
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("azure"),
			Region:        strPtr(location),
			Service:       strPtr("Container Registry"),
			ProductFamily: strPtr("Containers"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "productName", Value: strPtr("Container Registry")},
				{Key: "skuName", Value: strPtr(sku)},
				{Key: "meterName", Value: strPtr("Data Stored")},
			},
		},
		PriceFilter: &schema.PriceFilter{
			PurchaseOption: strPtr("Consumption"),
		},
	}
}
func ContainerRegistryCPUCostComponent(name, location, sku string, monthlyBuildVcpu *decimal.Decimal) *schema.CostComponent {

	return &schema.CostComponent{

		Name:            name,
		Unit:            "seconds",
		UnitMultiplier:  1,
		MonthlyQuantity: monthlyBuildVcpu,
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("azure"),
			Region:        strPtr(location),
			Service:       strPtr("Container Registry"),
			ProductFamily: strPtr("Containers"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "productName", Value: strPtr("Container Registry")},
				{Key: "skuName", Value: strPtr(sku)},
				{Key: "meterName", Value: strPtr("Task vCPU Duration")},
			},
		},
		PriceFilter: &schema.PriceFilter{
			PurchaseOption:   strPtr("Consumption"),
			StartUsageAmount: strPtr("6000"),
		},
	}
}
