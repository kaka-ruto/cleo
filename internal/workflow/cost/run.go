package cost

import "fmt"

func Execute(in Input) (Result, error) {
	switch in.Name {
	case "estimate":
		report, err := Estimate(in.Args)
		if err != nil {
			return Result{}, err
		}
		fmt.Println(report)
		return Result{Name: "estimate"}, nil
	default:
		return Result{}, BuildUnknownError(in.Name)
	}
}
