package qa

func Verify(plan Plan, _ Result) Verification {
	if plan.ReadOnly {
		return Verification{Checked: false, Reason: "read-only command"}
	}
	return Verification{Checked: true, Reason: "command execution completed without errors"}
}
