package task

import (
	"fmt"
	"strconv"
)

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

func hasFlag(args []string, key string) bool {
	for _, a := range args {
		if a == key {
			return true
		}
	}
	return false
}
