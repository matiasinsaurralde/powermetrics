package powermetrics

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/matiasinsaurralde/powermetrics/internal/samplers"
	howett_plist "howett.net/plist"
)

// Format represents the output format for powermetrics
type Format string

const (
	FormatText  Format = "text"
	FormatPlist Format = "plist"
)

// Sampler represents a powermetrics sampler
type Sampler string

const (
	GPUPower Sampler = "gpu_power"
)

// Custom errors
var (
	ErrUnsupportedSampler = errors.New("unsupported sampler")
	ErrUnsupportedFormat  = errors.New("unsupported format")
)

// Supported samplers
var supportedSamplers = map[Sampler]bool{
	GPUPower: true,
}

// Config holds the configuration for powermetrics execution
type Config struct {
	SampleCount    int
	SampleInterval time.Duration
	Format         Format
	Samplers       []Sampler
}

// Result holds the parsed result from powermetrics execution
type Result struct {
	RawOutput []byte
	PlistData *samplers.PlistRoot
	// For multiple samples
	MultipleSamples []*samplers.PlistRoot
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		SampleCount:    1,
		SampleInterval: 5 * time.Second, // Default to 5 seconds like powermetrics CLI
		Format:         FormatText,
		Samplers:       []Sampler{GPUPower},
	}
}

// GPU returns a new Config configured for GPU power sampling
func (c *Config) GPU() *Config {
	if c == nil {
		c = DefaultConfig()
	}
	return &Config{
		SampleCount:    c.SampleCount,
		SampleInterval: c.SampleInterval,
		Format:         FormatPlist,
		Samplers:       []Sampler{GPUPower},
	}
}

// GetSupportedSamplers returns a list of supported sampler names
func GetSupportedSamplers() []string {
	samplers := make([]string, 0, len(supportedSamplers))
	for sampler := range supportedSamplers {
		samplers = append(samplers, string(sampler))
	}
	return samplers
}

// ValidateSamplers checks if the provided samplers are supported
func ValidateSamplers(samplers []Sampler) error {
	for _, sampler := range samplers {
		if !supportedSamplers[sampler] {
			return fmt.Errorf("%w: %s", ErrUnsupportedSampler, sampler)
		}
	}
	return nil
}

// ValidateFormat checks if the provided format is supported
func ValidateFormat(format Format) error {
	switch format {
	case FormatText, FormatPlist:
		return nil
	default:
		return fmt.Errorf("%w: %s", ErrUnsupportedFormat, format)
	}
}

// Collect executes powermetrics with the given configuration
func Collect(config *Config) (*Result, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// Validate format
	if err := ValidateFormat(config.Format); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	// Validate samplers
	if err := ValidateSamplers(config.Samplers); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	// Build command arguments
	samplerStrings := make([]string, len(config.Samplers))
	for i, sampler := range config.Samplers {
		samplerStrings[i] = string(sampler)
	}

	args := []string{
		fmt.Sprintf("--sample-count=%d", config.SampleCount),
		fmt.Sprintf("--format=%s", config.Format),
		fmt.Sprintf("--samplers=%s", strings.Join(samplerStrings, ",")),
	}

	// Add sample interval if specified
	if config.SampleInterval > 0 {
		args = append(args, fmt.Sprintf("--sample-rate=%d", int(config.SampleInterval.Milliseconds())))
	}

	// Execute powermetrics command
	cmd := exec.Command("powermetrics", args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute powermetrics: %w", err)
	}

	result := &Result{
		RawOutput: output,
	}

	// Parse plist output if format is plist
	if config.Format == FormatPlist {
		// Try to parse as multiple samples first
		samples, err := parseMultipleSamples(output)
		if err == nil && len(samples) > 0 {
			result.MultipleSamples = samples
			// Set the first sample as the main PlistData for backward compatibility
			if len(samples) > 0 {
				result.PlistData = samples[0]
			}
		} else {
			// Fall back to single sample parsing
			var parsed samplers.PlistRoot
			decoder := howett_plist.NewDecoder(bytes.NewReader(output))
			if err := decoder.Decode(&parsed); err != nil {
				return nil, fmt.Errorf("failed to decode plist output: %w", err)
			}
			result.PlistData = &parsed
		}
	}

	return result, nil
}

// parseMultipleSamples attempts to parse multiple plist documents from the output
func parseMultipleSamples(output []byte) ([]*samplers.PlistRoot, error) {
	var samples []*samplers.PlistRoot

	// Split by XML declaration to find multiple plist documents
	xmlDeclarations := []byte("<?xml version")
	parts := bytes.Split(output, xmlDeclarations)

	for i, part := range parts {
		if i == 0 && len(part) == 0 {
			// Skip empty part before first XML declaration
			continue
		}

		if len(part) == 0 {
			continue
		}

		// Reconstruct the XML document
		xmlDoc := append(xmlDeclarations, part...)

		var parsed samplers.PlistRoot
		decoder := howett_plist.NewDecoder(bytes.NewReader(xmlDoc))
		if err := decoder.Decode(&parsed); err != nil {
			// Skip invalid plist documents
			continue
		}

		samples = append(samples, &parsed)
	}

	if len(samples) == 0 {
		return nil, fmt.Errorf("no valid plist documents found")
	}

	return samples, nil
}
