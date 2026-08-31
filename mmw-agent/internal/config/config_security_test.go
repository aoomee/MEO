package config

import "testing"

func TestValidateRequiresManagementToken(t *testing.T) {
	if err := (&Config{}).Validate(); err == nil {
		t.Fatal("empty token configuration was accepted")
	}
	if err := (&Config{Token: "secret"}).Validate(); err != nil {
		t.Fatalf("valid token configuration rejected: %v", err)
	}
}
