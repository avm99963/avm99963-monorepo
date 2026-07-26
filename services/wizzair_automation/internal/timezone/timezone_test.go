package timezone_test

import (
	"testing"

	"avm99963-monorepo/services/wizzair_automation/internal/timezone"
)

func TestGetLocation(t *testing.T) {
	tests := []struct {
		name     string
		iata     string
		fallback string
		expected string
	}{
		{name: "WAW_resolves_to_Europe_Warsaw", iata: "WAW", fallback: "Europe/Warsaw", expected: "Europe/Warsaw"},
		{name: "BCN_resolves_to_Europe_Madrid", iata: "BCN", fallback: "Europe/Warsaw", expected: "Europe/Madrid"},
		{name: "LTN_resolves_to_Europe_London", iata: "LTN", fallback: "Europe/Warsaw", expected: "Europe/London"},
		{name: "AUH_resolves_to_Asia_Dubai", iata: "AUH", fallback: "Europe/Warsaw", expected: "Asia/Dubai"},
		{name: "TFS_resolves_to_Atlantic_Canary", iata: "TFS", fallback: "Europe/Warsaw", expected: "Atlantic/Canary"},
		{name: "unknown_iata_with_empty_fallback_defaults_to_Europe_Warsaw", iata: "UNKNOWN", fallback: "", expected: "Europe/Warsaw"},
		{name: "unknown_iata_with_custom_fallback_uses_provided_fallback", iata: "UNKNOWN", fallback: "America/New_York", expected: "America/New_York"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loc := timezone.GetLocation(tt.iata, tt.fallback)
			if loc == nil {
				t.Fatalf("expected non-nil location for IATA %s", tt.iata)
			}
			if loc.String() != tt.expected {
				t.Errorf("GetLocation(%q, %q) = %q; want %q", tt.iata, tt.fallback, loc.String(), tt.expected)
			}
		})
	}
}
