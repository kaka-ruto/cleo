package cost

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func jsonResponse(code int, body string) *http.Response {
	return &http.Response{
		StatusCode: code,
		Status:     http.StatusText(code),
		Body:       io.NopCloser(bytes.NewBufferString(body)),
		Header:     make(http.Header),
	}
}

func TestResolveLiveRatesByCountryName(t *testing.T) {
	old := httpClient
	httpClient = &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		u := r.URL.String()
		switch {
		case strings.Contains(u, "restcountries.com/v3.1/name/Kenya"):
			return jsonResponse(200, `[{"cca2":"KE","name":{"common":"Kenya"},"currencies":{"KES":{}}}]`), nil
		case strings.Contains(u, "/country/us/indicator/"):
			return jsonResponse(200, `[{},[{"value":85000}]]`), nil
		case strings.Contains(u, "/country/ke/indicator/"):
			return jsonResponse(200, `[{},[{"value":16000}]]`), nil
		default:
			return jsonResponse(404, `{}`), nil
		}
	})}
	defer func() { httpClient = old }()

	rates, country, err := resolveLiveRates("Kenya")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if country.Code2 != "KE" {
		t.Fatalf("unexpected country: %#v", country)
	}
	if rates.Avg <= 0 {
		t.Fatalf("unexpected rates: %#v", rates)
	}
}

func TestResolveCountryISO2NoNetwork(t *testing.T) {
	country, err := resolveCountry("de")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if country.Code2 != "DE" {
		t.Fatalf("unexpected country: %#v", country)
	}
}
