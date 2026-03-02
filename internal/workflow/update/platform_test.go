package update

import "testing"

func TestAssetName(t *testing.T) {
	name, err := assetName("v1.2.3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if name == "" {
		t.Fatal("expected non-empty asset name")
	}
}
