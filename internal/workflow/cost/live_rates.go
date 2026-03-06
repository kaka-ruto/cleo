package cost

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

type countryInfo struct {
	Code2    string
	Name     string
	Currency string
}

func resolveLiveRates(countryInput string) (rateTable, countryInfo, error) {
	country, err := resolveCountry(countryInput)
	if err != nil {
		return rateTable{}, countryInfo{}, err
	}
	us, err := latestPPP("US")
	if err != nil {
		return rateTable{}, countryInfo{}, err
	}
	target, err := latestPPP(country.Code2)
	if err != nil {
		return rateTable{}, countryInfo{}, err
	}
	if us <= 0 || target <= 0 {
		return rateTable{}, countryInfo{}, fmt.Errorf("unable to resolve live PPP data")
	}
	ratio := target / us
	if ratio < 0.2 {
		ratio = 0.2
	}
	if ratio > 1.8 {
		ratio = 1.8
	}
	base := rateTable{Low: 95, Avg: 145, High: 220}
	return rateTable{
		Low:  round2(base.Low * ratio),
		Avg:  round2(base.Avg * ratio),
		High: round2(base.High * ratio),
	}, country, nil
}

func resolveCountry(input string) (countryInfo, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		trimmed = "US"
	}
	if len(trimmed) == 2 && isAlpha(trimmed) {
		code := strings.ToUpper(trimmed)
		return countryInfo{Code2: code, Name: code, Currency: ""}, nil
	}
	return countryByName(trimmed)
}

func countryByName(name string) (countryInfo, error) {
	endpoint := "https://restcountries.com/v3.1/name/" + url.PathEscape(name) + "?fullText=false"
	var rows []map[string]any
	if err := getJSON(endpoint, &rows); err != nil {
		return countryInfo{}, err
	}
	if len(rows) == 0 {
		return countryInfo{}, fmt.Errorf("country not found: %s", name)
	}
	return parseCountry(rows[0])
}

func parseCountry(row map[string]any) (countryInfo, error) {
	var out countryInfo
	if cca2, ok := row["cca2"].(string); ok {
		out.Code2 = strings.ToUpper(cca2)
	}
	if name, ok := row["name"].(map[string]any); ok {
		if common, ok := name["common"].(string); ok {
			out.Name = common
		}
	}
	if currencies, ok := row["currencies"].(map[string]any); ok {
		for k := range currencies {
			out.Currency = strings.ToUpper(k)
			break
		}
	}
	if out.Code2 == "" {
		return countryInfo{}, fmt.Errorf("invalid country payload")
	}
	if out.Name == "" {
		out.Name = out.Code2
	}
	return out, nil
}

func latestPPP(code2 string) (float64, error) {
	endpoint := "https://api.worldbank.org/v2/country/" + url.PathEscape(strings.ToLower(code2)) + "/indicator/NY.GDP.PCAP.PP.CD?format=json&per_page=70"
	var payload []any
	if err := getJSON(endpoint, &payload); err != nil {
		return 0, err
	}
	if len(payload) < 2 {
		return 0, fmt.Errorf("unexpected world bank response")
	}
	rows, ok := payload[1].([]any)
	if !ok {
		return 0, fmt.Errorf("unexpected world bank data rows")
	}
	for _, item := range rows {
		row, ok := item.(map[string]any)
		if !ok {
			continue
		}
		v := row["value"]
		switch t := v.(type) {
		case float64:
			if t > 0 {
				return t, nil
			}
		case string:
			f, err := strconv.ParseFloat(strings.TrimSpace(t), 64)
			if err == nil && f > 0 {
				return f, nil
			}
		}
	}
	return 0, fmt.Errorf("no non-null PPP value for %s", code2)
}

func getJSON(endpoint string, out any) error {
	resp, err := httpClient.Get(endpoint)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("request failed: %s", resp.Status)
	}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(out); err != nil {
		return err
	}
	return nil
}

func isAlpha(s string) bool {
	for _, r := range s {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') {
			return false
		}
	}
	return true
}

func round2(v float64) float64 {
	return float64(int(v*100+0.5)) / 100
}
