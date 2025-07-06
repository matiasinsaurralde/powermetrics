package powermetrics

import (
	"bytes"
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

// ValidateSamplers checks if the provided samplers are supported
func ValidateSamplers(samplers []string) error {
	supportedSamplers := map[string]bool{
		"gpu_power": true,
	}

	for _, sampler := range samplers {
		if !supportedSamplers[sampler] {
			return fmt.Errorf("unsupported sampler: %s", sampler)
		}
	}
	return nil
}

// Collect executes powermetrics with the given configuration
func Collect(config *Config) (*Result, error) {
	if config == nil {
		config = DefaultConfig()
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
