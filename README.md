[![CI](https://github.com/matiasinsaurralde/powermetrics/actions/workflows/ci.yml/badge.svg)](https://github.com/matiasinsaurralde/powermetrics/actions/workflows/ci.yml)

# powermetrics

A Go package for programmatically running and parsing Apple's `powermetrics` command output on macOS, with a focus on GPU power metrics.

## Features

- **GPU Power Metrics**: Collect and parse GPU idle ratio, active ratio, average power, and peak power
- **Flexible Configuration**: Customize sample count, sample rate, output format, and samplers
- **Multiple Output Formats**: Support for both text and plist (XML) output formats
- **Type Safety**: Strongly typed configuration with constants for samplers and formats
- **Error Handling**: Custom error types for unsupported samplers and formats
- **Testable**: Mock command execution for reliable unit testing

## Installation

```bash
go get github.com/matiasinsaurralde/powermetrics
```

## Quick Start

```go
package main

import (
	"fmt"
	"time"
	"github.com/matiasinsaurralde/powermetrics"
)

func main() {
	// Create a new powermetrics instance
	pm := powermetrics.New()
	
	// Use default configuration for GPU power metrics
	config := powermetrics.DefaultConfig().GPU()
	
	// Customize settings
	config.SampleCount = 5
	config.SampleRate = 1 * time.Second
	
	// Collect metrics
	result, err := pm.Collect(config)
	if err != nil {
		panic(err)
	}
	
	// Access GPU power data (single sample)
	if result.PlistData != nil {
		gpu := result.PlistData.GPU
		fmt.Printf("GPU Frequency: %.2f Hz\n", gpu.FreqHz)
		fmt.Printf("GPU Idle Ratio: %.2f%%\n", gpu.IdleRatio*100)
		if gpu.GPUEnergy != nil {
			fmt.Printf("GPU Energy: %d mW\n", *gpu.GPUEnergy)
		}
	}
	
	// Access multiple samples (when SampleCount > 1)
	if len(result.Samples) > 0 {
		fmt.Printf("Collected %d samples\n", len(result.Samples))
		for i, sample := range result.Samples {
			gpu := sample.GPU
			fmt.Printf("Sample %d: GPU Idle Ratio: %.2f%%\n", i, gpu.IdleRatio*100)
		}
	}
}
```

## Configuration

The package uses a `Config` struct to control powermetrics execution:

```go
type Config struct {
	SampleCount int           // Number of samples to collect
	SampleRate  time.Duration // Time between samples
	Format      Format        // Output format (text or plist)
	Samplers    []Sampler     // List of samplers to use
}
```

### Default Configuration

```go
config := powermetrics.DefaultConfig()
// SampleCount: 1
// SampleRate: 5 seconds
// Format: FormatText
// Samplers: [GPUPower]
```

### GPU-Specific Configuration

```go
config := powermetrics.DefaultConfig().GPU()
// SampleCount: 1
// SampleRate: 5 seconds
// Format: FormatPlist (required for GPU metrics)
// Samplers: [GPUPower]
```

## Result Structure

The `Result` struct provides access to the collected data:

```go
type Result struct {
	RawOutput []byte                    // Raw output from powermetrics
	PlistData *samplers.PlistRoot       // Single sample data
	Samples   []*samplers.PlistRoot     // Multiple samples when SampleCount > 1
}
```

- **Single Sample**: When `SampleCount = 1`, use `result.PlistData` for the parsed data
- **Multiple Samples**: When `SampleCount > 1`, use `result.Samples` for all collected samples

## Supported Samplers

Currently, the package supports the following samplers:

- `GPUPower`: GPU power metrics (idle ratio, active ratio, average power, peak power)

## Output Formats

- `FormatText`: Raw text output from powermetrics
- `FormatPlist`: XML plist format (required for structured data parsing)

## Error Types

The package defines custom error types for better error handling:

- `ErrUnsupportedSampler`: When an unsupported sampler is specified
- `ErrUnsupportedFormat`: When an unsupported format is specified

## Testing

The package includes mock support for reliable unit testing:

```go
// In your tests
xmlData, _ := os.ReadFile("testdata/gpu_power.xml")
mockRunner := &powermetrics.MockCommandRunner{Output: xmlData}
pm := powermetrics.NewWithRunner(mockRunner)
result, err := pm.Collect(config)
```

## Samples

The package includes example applications in the `samples/` directory:

### Simple Example

A basic example that demonstrates the core functionality of the package.

```bash
cd samples/simple
go run main.go
```

Features:
- Collects 5 GPU power samples with 1-second intervals
- Displays GPU frequency, idle ratio, and energy consumption
- Shows how to access both single and multiple samples
- Uses the same code as the main README example

### Terminal Dashboard

A real-time terminal dashboard using `termui` that displays GPU idle ratio with a sparkline chart.

```bash
cd samples/termui-gpu
go run main.go
```

Features:
- Real-time GPU idle ratio display
- Sparkline chart showing historical trends
- Detailed GPU power metrics
- Interactive controls (press 'q' to quit)

### HTTP Server

A simple HTTP server that exposes GPU metrics as JSON via a REST API.

```bash
cd samples/httpserver
go run main.go
```

Features:
- REST API endpoint at `/gpu`
- Pretty JSON output with timestamps
- CORS support for web applications
- GPU power metrics in structured format

## Testing

Run the test suite:

```bash
go test ./...
```

The tests include:
- Default configuration validation
- GPU power metrics parsing
- Multiple sample parsing
- Mock command execution testing
- Error handling for unsupported configurations

## Requirements

- macOS (powermetrics is only available on macOS)
- Go 1.19 or later
- `powermetrics` command available in PATH (included with macOS)

## Notes

- This package requires the `powermetrics` command to be available on the system
- GPU power metrics are only available on supported macOS systems
- Some metrics may require elevated privileges

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## License

MIT License - see LICENSE file for details.