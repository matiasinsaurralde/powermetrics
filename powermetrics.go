package powermetrics

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/matiasinsaurralde/powermetrics/internal/samplers"
	howett_plist "howett.net/plist"
)

// Format represents the output format for powermetrics
type Format string

const (
	FormatText  Format = "text"
	FormatPlist Format = "plist"
)

// Custom errors
var (
	ErrUnsupportedSampler = errors.New("unsupported sampler")
	ErrUnsupportedFormat  = errors.New("unsupported format")
)

// Supported samplers
var supportedSamplers = map[string]bool{
	"gpu_power": true,
}

// Config holds the configuration for powermetrics execution
type Config struct {
	SampleCount int
	Format      Format
	Samplers    []string
}

// Result holds the parsed result from powermetrics execution
type Result struct {
	RawOutput []byte
	PlistData *samplers.PlistRoot
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		SampleCount: 1,
		Format:      FormatText,
		Samplers:    []string{"gpu_power"},
	}
}

// GPU returns a new Config configured for GPU power sampling
func (c *Config) GPU() *Config {
	if c == nil {
		c = DefaultConfig()
	}
	return &Config{
		SampleCount: c.SampleCount,
		Format:      FormatPlist,
		Samplers:    []string{"gpu_power"},
	}
}

// GetSupportedSamplers returns a list of supported sampler names
func GetSupportedSamplers() []string {
	samplers := make([]string, 0, len(supportedSamplers))
	for sampler := range supportedSamplers {
		samplers = append(samplers, sampler)
	}
	return samplers
}

// ValidateSamplers checks if the provided samplers are supported
func ValidateSamplers(samplers []string) error {
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
	args := []string{
		fmt.Sprintf("--sample-count=%d", config.SampleCount),
		fmt.Sprintf("--format=%s", config.Format),
		fmt.Sprintf("--samplers=%s", strings.Join(config.Samplers, ",")),
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
		var parsed samplers.PlistRoot
		decoder := howett_plist.NewDecoder(bytes.NewReader(output))
		if err := decoder.Decode(&parsed); err != nil {
			return nil, fmt.Errorf("failed to decode plist output: %w", err)
		}
		result.PlistData = &parsed
	}

	return result, nil
}
