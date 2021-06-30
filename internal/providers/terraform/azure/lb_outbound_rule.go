package azure

import (
	"strings"

	"github.com/infracost/infracost/internal/schema"
)

func GetAzureRMLoadBalancerOutboundRuleRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:  "azurerm_lb_outbound_rule",
		RFunc: NewAzureRMLoadBalancerRule,
		ReferenceAttributes: []string{
			"loadbalancer_id",
			"resource_group_name",
		},
	}
}

func NewAzureRMLoadBalancerOutboundRule(d *schema.ResourceData, u *schema.UsageData) *schema.Resource {
	region := lookupRegion(d, []string{"loadbalancer_id", "resource_group_name"})
	if strings.Contains(strings.ToLower(region), "usgov") {
		region = "US Gov"
	} else if strings.Contains(strings.ToLower(region), "china") {
		region = "Сhina"
	} else {
		region = "Global"
	}

	var costComponents []*schema.CostComponent
	costComponents = append(costComponents, rulesCostComponent(region))

	return &schema.Resource{
		Name:           d.Address,
		CostComponents: costComponents,
	}
}
