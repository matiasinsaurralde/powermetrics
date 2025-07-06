package powermetrics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/matiasinsaurralde/powermetrics/internal/samplers"
	howett_plist "howett.net/plist"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.SampleCount != 1 {
		t.Errorf("Expected SampleCount to be 1, got %d", config.SampleCount)
	}

	if config.SampleRate != 5*time.Second {
		t.Errorf("Expected SampleRate to be 5 seconds, got %v", config.SampleRate)
	}

	if config.Format != FormatText {
		t.Errorf("Expected Format to be FormatText, got %s", config.Format)
	}

	if len(config.Samplers) != 1 || config.Samplers[0] != GPUPower {
		t.Errorf("Expected Samplers to be [GPUPower], got %v", config.Samplers)
	}
}

func TestValidateSamplers(t *testing.T) {
	tests := []struct {
		name      string
		samplers  []Sampler
		expectErr bool
	}{
		{
			name:      "valid sampler",
			samplers:  []Sampler{GPUPower},
			expectErr: false,
		},
		{
			name:      "invalid sampler",
			samplers:  []Sampler{"invalid_sampler"},
			expectErr: true,
		},
		{
			name:      "mixed valid and invalid",
			samplers:  []Sampler{GPUPower, "invalid_sampler"},
			expectErr: true,
		},
		{
			name:      "empty samplers",
			samplers:  []Sampler{},
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

func TestSamplesXMLUnmarshaling(t *testing.T) {
	// Read the test XML file with multiple samples
	xmlData, err := os.ReadFile("testdata/gpu_power_multiple_samples.xml")
	if err != nil {
		t.Fatalf("Failed to read gpu_power_multiple_samples.xml: %v", err)
	}

	// Parse the XML data using the parseMultipleSamples function directly
	samples, err := parseMultipleSamples(xmlData)
	if err != nil {
		t.Fatalf("Failed to parse multiple samples: %v", err)
	}

	// Verify we got the expected number of samples
	expectedSampleCount := 5
	if len(samples) != expectedSampleCount {
		t.Errorf("Expected %d samples, got %d", expectedSampleCount, len(samples))
	}

	// Verify each sample has valid data
	for i, sample := range samples {
		t.Run(fmt.Sprintf("sample_%d", i), func(t *testing.T) {
			// Verify basic structure
			if sample.HWModel == "" {
				t.Errorf("Sample %d: Expected HWModel to be non-empty", i)
			}

			if sample.KernOSVer == "" {
				t.Errorf("Sample %d: Expected KernOSVer to be non-empty", i)
			}

			if sample.Timestamp.IsZero() {
				t.Errorf("Sample %d: Expected Timestamp to be non-zero", i)
			}

			// Verify GPU info
			if sample.GPU.FreqHz <= 0 {
				t.Errorf("Sample %d: Expected GPU frequency to be positive, got %f", i, sample.GPU.FreqHz)
			}

			if sample.GPU.IdleRatio < 0 || sample.GPU.IdleRatio > 1 {
				t.Errorf("Sample %d: Expected GPU idle ratio to be between 0 and 1, got %f", i, sample.GPU.IdleRatio)
			}

			// Verify DVFM states
			if len(sample.GPU.DVFMStates) == 0 {
				t.Errorf("Sample %d: Expected at least one DVFM state", i)
			}

			// Verify SW requested states
			if len(sample.GPU.SWRequestedState) == 0 {
				t.Errorf("Sample %d: Expected at least one SW requested state", i)
			}

			// Verify SW states
			if len(sample.GPU.SWState) == 0 {
				t.Errorf("Sample %d: Expected at least one SW state", i)
			}
		})
	}

	// Test that samples have different timestamps (indicating they're different samples)
	firstTimestamp := samples[0].Timestamp
	lastTimestamp := samples[len(samples)-1].Timestamp
	if firstTimestamp.Equal(lastTimestamp) {
		t.Error("Expected samples to have different timestamps")
	}

	t.Logf("Successfully parsed %d samples with timestamps from %s to %s",
		len(samples), firstTimestamp.Format(time.RFC3339), lastTimestamp.Format(time.RFC3339))
}

func TestCollectWithMock(t *testing.T) {
	// Read test XML data
	xmlData, err := os.ReadFile("testdata/gpu_power_multiple_samples.xml")
	if err != nil {
		t.Fatalf("Failed to read test XML: %v", err)
	}

	// Create mock runner
	mockRunner := &MockCommandRunner{Output: xmlData}
	pm := NewWithRunner(mockRunner)

	// Test configuration
	config := &Config{
		SampleCount: 3,
		SampleRate:  1 * time.Second,
		Format:      FormatPlist,
		Samplers:    []Sampler{GPUPower},
	}

	result, err := pm.Collect(config)
	if err != nil {
		t.Fatalf("Collect failed: %v", err)
	}

	// Check that we have multiple samples
	if len(result.Samples) == 0 {
		t.Error("Expected multiple samples, got none")
	}

	// Check that PlistData is set to the first sample for backward compatibility
	if result.PlistData == nil {
		t.Error("Expected PlistData to be set for backward compatibility")
	}

	// Verify that PlistData matches the first sample
	if len(result.Samples) > 0 && result.PlistData != result.Samples[0] {
		t.Error("Expected PlistData to be the first sample")
	}
}

func TestCollectWithMockSingleSample(t *testing.T) {
	// Read test XML data for single sample
	xmlData, err := os.ReadFile("testdata/gpu_power.xml")
	if err != nil {
		t.Fatalf("Failed to read test XML: %v", err)
	}

	// Create mock runner
	mockRunner := &MockCommandRunner{Output: xmlData}
	pm := NewWithRunner(mockRunner)

	// Test configuration
	config := &Config{
		SampleCount: 1,
		SampleRate:  1 * time.Second,
		Format:      FormatPlist,
		Samplers:    []Sampler{GPUPower},
	}

	result, err := pm.Collect(config)
	if err != nil {
		t.Fatalf("Collect failed: %v", err)
	}

	// Check that we have a single sample
	if result.PlistData == nil {
		t.Error("Expected PlistData to be set")
	}

	// Check that GPU data is available
	if result.PlistData.GPU.FreqHz <= 0 {
		t.Error("Expected GPU frequency to be positive")
	}
}

func TestNewAndNewWithRunner(t *testing.T) {
	// Test New() creates a real runner
	pm1 := New()
	if pm1.runner == nil {
		t.Error("Expected New() to create a runner")
	}

	// Test NewWithRunner() uses the provided runner
	mockRunner := &MockCommandRunner{}
	pm2 := NewWithRunner(mockRunner)
	if pm2.runner != mockRunner {
		t.Error("Expected NewWithRunner() to use the provided runner")
	}
}
