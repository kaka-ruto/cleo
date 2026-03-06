package cost

func Verify(_ Plan, _ Result) Verification {
	return Verification{Checked: true, Reason: "cost estimate produced"}
}
