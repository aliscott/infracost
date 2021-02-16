package output

import (
	"encoding/json"
	"time"

	"github.com/shopspring/decimal"
)

type ReportInput struct {
	Metadata map[string]string
	Root     Root
}

func Load(data []byte) (Root, error) {
	var out Root
	err := json.Unmarshal(data, &out)
	return out, err
}

func Combine(inputs []ReportInput, opts Options) Root {
	var combined Root

	var totalHourlyCost *decimal.Decimal
	var totalMonthlyCost *decimal.Decimal

	projects := make([]Project, 0)
	summaries := make([]*ResourceSummary, 0, len(inputs))

	for _, input := range inputs {

		projects = append(projects, input.Root.Projects...)

		for _, r := range input.Root.Resources {
			for k, v := range input.Metadata {
				if r.Metadata == nil {
					r.Metadata = make(map[string]string)
				}
				r.Metadata[k] = v
			}

			combined.Resources = append(combined.Resources, r)
		}

		if input.Root.TotalHourlyCost != nil {
			if totalHourlyCost == nil {
				totalHourlyCost = decimalPtr(decimal.Zero)
			}

			totalHourlyCost = decimalPtr(totalHourlyCost.Add(*input.Root.TotalHourlyCost))
		}
		if input.Root.TotalMonthlyCost != nil {
			if totalMonthlyCost == nil {
				totalMonthlyCost = decimalPtr(decimal.Zero)
			}

			totalMonthlyCost = decimalPtr(totalMonthlyCost.Add(*input.Root.TotalMonthlyCost))
		}

		summaries = append(summaries, input.Root.ResourceSummary)
	}

	sortResources(combined.Resources, opts.GroupKey)

	combined.Projects = projects
	combined.TotalHourlyCost = totalHourlyCost
	combined.TotalMonthlyCost = totalMonthlyCost
	combined.TimeGenerated = time.Now()
	combined.ResourceSummary = combinedResourceSummaries(summaries)

	return combined
}

func combinedResourceSummaries(summaries []*ResourceSummary) *ResourceSummary {
	combined := &ResourceSummary{}

	for _, s := range summaries {
		if s == nil {
			continue
		}

		combined.SupportedCounts = combineCounts(combined.SupportedCounts, s.SupportedCounts)
		combined.UnsupportedCounts = combineCounts(combined.UnsupportedCounts, s.UnsupportedCounts)
		combined.TotalSupported = addIntPtrs(combined.TotalSupported, s.TotalSupported)
		combined.TotalUnsupported = addIntPtrs(combined.TotalUnsupported, s.TotalUnsupported)
		combined.TotalNoPrice = addIntPtrs(combined.TotalNoPrice, s.TotalNoPrice)
		combined.Total = addIntPtrs(combined.Total, s.Total)
	}

	return combined
}

func combineCounts(c1 *map[string]int, c2 *map[string]int) *map[string]int {
	if c1 == nil && c2 == nil {
		return nil
	}

	res := make(map[string]int)

	if c1 != nil {
		for k, v := range *c1 {
			res[k] = v
		}
	}

	if c2 != nil {
		for k, v := range *c2 {
			res[k] += v
		}
	}

	return &res
}

func addIntPtrs(i1 *int, i2 *int) *int {
	if i1 == nil && i2 == nil {
		return nil
	}

	val1 := 0
	if i1 != nil {
		val1 = *i1
	}

	val2 := 0
	if i2 != nil {
		val2 = *i2
	}

	res := val1 + val2
	return &res
}
