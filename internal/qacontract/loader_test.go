package qacontract

import "testing"

func TestLoadStringAcceptsBDDContract(t *testing.T) {
	raw := `version: 1
name: checkout tax
criteria:
  - id: c1
    title: tax recalculates
    severity: high
    actors: [buyer_web]
    surface: web
    environment: local
    given: buyer has item in cart
    when: buyer updates shipping address
    then:
      - cart total reflects updated tax
    evidence_required:
      - screenshot:cart-summary`
	doc, err := LoadString(raw)
	if err != nil {
		t.Fatalf("LoadString() error = %v", err)
	}
	if got := doc.Criteria[0].Surface; got != "web" {
		t.Fatalf("Surface = %q, want web", got)
	}
	if got := len(doc.Criteria[0].Then); got != 1 {
		t.Fatalf("Then length = %d, want 1", got)
	}
}

func TestLoadStringNormalizesLegacyContract(t *testing.T) {
	raw := `version: 1
name: legacy
criteria:
  - id: c1
    title: legacy criterion
    severity: medium
    actors: [core]
    acceptance:
      goal: cart has taxable line item
      expected_result: tax amount matches rate
    execution:
      surface: api
      environment: local
      steps:
        - action: call_api
          params:
            method: GET
            url: http://localhost:0/placeholder
            output_key: response`
	doc, err := LoadString(raw)
	if err != nil {
		t.Fatalf("LoadString() error = %v", err)
	}
	c := doc.Criteria[0]
	if c.Surface != "api" {
		t.Fatalf("Surface = %q, want api", c.Surface)
	}
	if c.Given != "cart has taxable line item" {
		t.Fatalf("Given = %q", c.Given)
	}
	if len(c.Then) != 1 || c.Then[0] != "tax amount matches rate" {
		t.Fatalf("Then = %#v", c.Then)
	}
}
