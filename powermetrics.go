package powermetrics

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/matiasinsaurralde/powermetrics/internal/samplers"
	howett_plist "howett.net/plist"
)

// CommandRunner interface for executing external commands
type CommandRunner interface {
	Run(name string, args ...string) ([]byte, error)
}

// RealCommandRunner implements CommandRunner using exec.Command
type RealCommandRunner struct{}

func (r *RealCommandRunner) Run(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).Output()
}

// MockCommandRunner implements CommandRunner for testing
type MockCommandRunner struct {
	Output []byte
	Err    error
}

func (m *MockCommandRunner) Run(name string, args ...string) ([]byte, error) {
	return m.Output, m.Err
}

// Powermetrics struct holds the command runner
type Powermetrics struct {
	runner CommandRunner
}

// New creates a new Powermetrics instance with the real command runner
func New() *Powermetrics {
	return &Powermetrics{
		runner: &RealCommandRunner{},
	}
}

// NewWithRunner creates a new Powermetrics instance with a custom command runner
func NewWithRunner(runner CommandRunner) *Powermetrics {
	return &Powermetrics{
		runner: runner,
	}
}

// Sampler represents a powermetrics sampler
type Sampler string

const (
	GPUPower Sampler = "gpu_power"
)

// Format represents the output format
type Format string

const (
	FormatText  Format = "text"
	FormatPlist Format = "plist"
)

// Error types
var (
	ErrUnsupportedSampler = fmt.Errorf("unsupported sampler")
	ErrUnsupportedFormat  = fmt.Errorf("unsupported format")
)

// Supported samplers
var supportedSamplers = map[Sampler]bool{
	GPUPower: true,
}

// Config holds the configuration for powermetrics execution
type Config struct {
	SampleCount int
	SampleRate  time.Duration
	Format      Format
	Samplers    []Sampler
}

// Result holds the parsed result from powermetrics execution
type Result struct {
	RawOutput []byte
	PlistData *samplers.PlistRoot
	// For multiple samples
	Samples []*samplers.PlistRoot
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		SampleCount: 1,
		SampleRate:  5 * time.Second, // Default to 5 seconds like powermetrics CLI
		Format:      FormatText,
		Samplers:    []Sampler{GPUPower},
	}
}

// GPU returns a new Config configured for GPU power sampling
func (c *Config) GPU() *Config {
	if c == nil {
		c = DefaultConfig()
	}
	return &Config{
		SampleCount: c.SampleCount,
		SampleRate:  c.SampleRate,
		Format:      FormatPlist,
		Samplers:    []Sampler{GPUPower},
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
func (p *Powermetrics) Collect(config *Config) (*Result, error) {
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

	// Add sample rate if specified
	if config.SampleRate > 0 {
		args = append(args, fmt.Sprintf("--sample-rate=%d", int(config.SampleRate.Milliseconds())))
	}

	// Execute powermetrics command
	output, err := p.runner.Run("powermetrics", args...)
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
			result.Samples = samples
			// Set the first sample as the main PlistData
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
	// Split by plist boundaries
	parts := bytes.Split(output, []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>"))

	var samples []*samplers.PlistRoot

	for i, part := range parts {
		if i == 0 {
			// Skip the first empty part
			continue
		}

		// Reconstruct the XML with the header
		xmlData := append([]byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>"), part...)

		var parsed samplers.PlistRoot
		decoder := howett_plist.NewDecoder(bytes.NewReader(xmlData))
		if err := decoder.Decode(&parsed); err != nil {
			continue // Skip invalid plists
		}

		samples = append(samples, &parsed)
	}

	if len(samples) == 0 {
		return nil, fmt.Errorf("no valid plist documents found")
	}

	return samples, nil
}
