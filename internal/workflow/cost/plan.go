package cost

import "fmt"

func BuildPlan(in Input) (Plan, error) {
	switch in.Name {
	case "estimate":
		rateSource := flagValue(in.Args, "--rates-source")
		if rateSource != "" && rateSource != "cached" && rateSource != "manual" && rateSource != "live" {
			return Plan{}, fmt.Errorf("--rates-source must be cached|manual|live")
		}
		if rateSource == "manual" && flagValue(in.Args, "--hourly-rate") == "" {
			return Plan{}, fmt.Errorf("--hourly-rate is required when --rates-source manual")
		}
		return Plan{Name: "estimate", Description: "Estimate project development cost from codebase metrics", ReadOnly: true}, nil
	default:
		return Plan{}, BuildUnknownError(in.Name)
	}
}
