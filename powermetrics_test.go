package powermetrics

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/matiasinsaurralde/powermetrics/internal/samplers"
	howett_plist "howett.net/plist"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.SampleCount != 1 {
		t.Errorf("Expected SampleCount to be 1, got %d", config.SampleCount)
	}

	if config.Format != FormatText {
		t.Errorf("Expected Format to be FormatText, got %s", config.Format)
	}

	if len(config.Samplers) != 1 || config.Samplers[0] != "gpu_power" {
		t.Errorf("Expected Samplers to be [gpu_power], got %v", config.Samplers)
	}
}

func TestValidateSamplers(t *testing.T) {
	tests := []struct {
		name      string
		samplers  []string
		expectErr bool
	}{
		{
			name:      "valid sampler",
			samplers:  []string{"gpu_power"},
			expectErr: false,
		},
		{
			name:      "invalid sampler",
			samplers:  []string{"invalid_sampler"},
			expectErr: true,
		},
		{
			name:      "mixed valid and invalid",
			samplers:  []string{"gpu_power", "invalid_sampler"},
			expectErr: true,
		},
		{
			name:      "empty samplers",
			samplers:  []string{},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSamplers(tt.samplers)
			if tt.expectErr && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestGPUPowerXMLUnmarshaling(t *testing.T) {
	// Read the test XML file
	xmlData, err := os.ReadFile("testdata/gpu_power.xml")
	if err != nil {
		t.Fatalf("Failed to read gpu_power.xml: %v", err)
	}

	// Parse the XML data
	var parsed samplers.PlistRoot
	decoder := howett_plist.NewDecoder(bytes.NewReader(xmlData))
	if err := decoder.Decode(&parsed); err != nil {
		t.Fatalf("Failed to decode plist: %v", err)
	}

	// Verify basic structure
	if parsed.HWModel == "" {
		t.Error("Expected HWModel to be non-empty")
	}

	if parsed.KernOSVer == "" {
		t.Error("Expected KernOSVer to be non-empty")
	}

	if parsed.Timestamp.IsZero() {
		t.Error("Expected Timestamp to be non-zero")
	}

	// Verify GPU info
	if parsed.GPU.FreqHz <= 0 {
		t.Errorf("Expected GPU frequency to be positive, got %f", parsed.GPU.FreqHz)
	}

	if parsed.GPU.IdleRatio < 0 || parsed.GPU.IdleRatio > 1 {
		t.Errorf("Expected GPU idle ratio to be between 0 and 1, got %f", parsed.GPU.IdleRatio)
	}

	// Verify DVFM states
	if len(parsed.GPU.DVFMStates) == 0 {
		t.Error("Expected at least one DVFM state")
	}

	for i, state := range parsed.GPU.DVFMStates {
		if state.Freq <= 0 {
			t.Errorf("DVFM state %d: Expected frequency to be positive, got %d", i, state.Freq)
		}
		if state.UsedRatio < 0 || state.UsedRatio > 1 {
			t.Errorf("DVFM state %d: Expected used ratio to be between 0 and 1, got %f", i, state.UsedRatio)
		}
	}

	// Verify SW requested states
	if len(parsed.GPU.SWRequestedState) == 0 {
		t.Error("Expected at least one SW requested state")
	}

	for i, state := range parsed.GPU.SWRequestedState {
		if state.SWReqState == "" {
			t.Errorf("SW requested state %d: Expected state name to be non-empty", i)
		}
		if state.UsedRatio < 0 || state.UsedRatio > 1 {
			t.Errorf("SW requested state %d: Expected used ratio to be between 0 and 1, got %f", i, state.UsedRatio)
		}
	}

	// Verify SW states
	if len(parsed.GPU.SWState) == 0 {
		t.Error("Expected at least one SW state")
	}

	for i, state := range parsed.GPU.SWState {
		if state.SWState == "" {
			t.Errorf("SW state %d: Expected state name to be non-empty", i)
		}
		if state.UsedRatio < 0 || state.UsedRatio > 1 {
			t.Errorf("SW state %d: Expected used ratio to be between 0 and 1, got %f", i, state.UsedRatio)
		}
	}

	// Test JSON marshaling to ensure the structure is complete
	jsonData, err := json.MarshalIndent(parsed, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal to JSON: %v", err)
	}

	if len(jsonData) == 0 {
		t.Error("Expected JSON output to be non-empty")
	}

	// Print first 500 characters of JSON for debugging
	jsonStr := string(jsonData)
	if len(jsonStr) > 500 {
		t.Logf("JSON output (first 500 chars): %s...", jsonStr[:500])
	} else {
		t.Logf("JSON output: %s", jsonStr)
	}

	// Verify that the JSON contains expected fields
	expectedFields := []string{"HWModel", "GPU", "Timestamp", "FreqHz", "IdleRatio"}
	for _, field := range expectedFields {
		if !bytes.Contains(jsonData, []byte(field)) {
			t.Errorf("Expected JSON to contain field: %s", field)
		}
	}

	t.Logf("Successfully parsed GPU power data with %d DVFM states, %d SW requested states, and %d SW states",
		len(parsed.GPU.DVFMStates), len(parsed.GPU.SWRequestedState), len(parsed.GPU.SWState))
}
