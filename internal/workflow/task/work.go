package task

import (
	"fmt"
	"regexp"
	"strings"
)

func branchForTask(id int64, title string) string {
	slug := slugify(title)
	if slug == "" {
		slug = "task"
	}
	return fmt.Sprintf("task/%d-%s", id, slug)
}

func slugify(input string) string {
	value := strings.ToLower(strings.TrimSpace(input))
	value = strings.ReplaceAll(value, "_", "-")
	re := regexp.MustCompile(`[^a-z0-9-]+`)
	value = re.ReplaceAllString(value, "-")
	value = strings.Trim(value, "-")
	for strings.Contains(value, "--") {
		value = strings.ReplaceAll(value, "--", "-")
	}
	if len(value) > 48 {
		value = strings.Trim(value[:48], "-")
	}
	return value
}
