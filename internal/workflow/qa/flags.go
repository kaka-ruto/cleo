package qa

import (
	"fmt"
	"strconv"
	"strings"
)

func hasFlag(args []string, key string) bool {
	for _, a := range args {
		if a == key {
			return true
		}
	}
	return false
}

func flagValue(args []string, key string) string {
	for i := 0; i < len(args); i++ {
		if args[i] == key && i+1 < len(args) {
			return args[i+1]
		}
	}
	return ""
}

func int64Flag(args []string, key string) (int64, error) {
	raw := flagValue(args, key)
	if raw == "" {
		return 0, fmt.Errorf("%s is required", key)
	}
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s must be an integer", key)
	}
	return id, nil
}

func profileList(raw string) []string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil
	}
	parts := strings.Split(trimmed, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		name := strings.TrimSpace(part)
		if name != "" {
			out = append(out, name)
		}
	}
	return out
}
