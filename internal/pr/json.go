package pr

import "github.com/kaka-ruto/cleo/internal/ghcli"

func ghDecodeRows(out string) ([]int, error) {
	var rows []struct {
		Number int `json:"number"`
	}
	if err := ghcli.DecodeJSON(out, &rows); err != nil {
		return nil, err
	}
	result := make([]int, 0, len(rows))
	for _, row := range rows {
		result = append(result, row.Number)
	}
	return result, nil
}
