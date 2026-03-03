package qatools

import (
	"fmt"
	"os/exec"
	"strings"
)

type Check struct {
	Name   string
	Status string
	Detail string
}

func Doctor(required []string) []Check {
	set := map[string]struct{}{}
	for _, tool := range required {
		set[strings.TrimSpace(tool)] = struct{}{}
	}
	var checks []Check
	if _, ok := set["api"]; ok {
		checks = append(checks, binCheck("curl", "api"))
	}
	if _, ok := set["browser"]; ok {
		checks = append(checks, binCheck("npx", "browser (playwright via npx)"))
	}
	return checks
}

func binCheck(bin string, name string) Check {
	if _, err := exec.LookPath(bin); err != nil {
		return Check{Name: name, Status: "missing", Detail: fmt.Sprintf("%s not found in PATH", bin)}
	}
	return Check{Name: name, Status: "ready", Detail: fmt.Sprintf("%s available", bin)}
}
