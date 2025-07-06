[![CI](https://github.com/matiasinsaurralde/powermetrics/actions/workflows/ci.yml/badge.svg)](https://github.com/matiasinsaurralde/powermetrics/actions/workflows/ci.yml)

# Go Powermetrics

A Go package for programmatically executing and parsing Apple's `powermetrics` command output. This package provides a clean interface to collect GPU power metrics and other system information on macOS.

## Features

- **Programmatic Execution**: Run `powermetrics` commands from Go code
- **Flexible Configuration**: Configure sample count, output format, and samplers
- **Structured Output**: Parse plist output into Go structs
- **Validation**: Built-in validation for supported samplers
- **GPU Power Metrics**: Full support for GPU power sampling

## Installation

```bash
go get github.com/matiasinsaurralde/powermetrics
```

## Requirements

- macOS (required for `powermetrics` command)
- Go 1.24 or later

## Usage

### One-liner Example

```go
import "github.com/matiasinsaurralde/powermetrics/powermetrics"

result, err := powermetrics.Collect(nil)
```

### GPU Power Sampling

```go
import "github.com/matiasinsaurralde/powermetrics/powermetrics"

config := &powermetrics.Config{SampleCount: 3}
result, err := powermetrics.Collect(config.GPU())
```

### Custom Configuration

```go
import "github.com/matiasinsaurralde/powermetrics/powermetrics"

config := &powermetrics.Config{
    SampleCount: 5,
    Format:      powermetrics.FormatPlist,
    Samplers:    []string{"gpu_power"},
}
result, err := powermetrics.Collect(config)
```

## Configuration Options

### Config Struct

```go
type Config struct {
    SampleCount int           // Number of samples to collect
    Format      Format        // Output format (text or plist)
    Samplers    []string      // List of samplers to use
}
```

### Config Methods

```go
// Configure for GPU power sampling
config := &powermetrics.Config{SampleCount: 1}
gpuConfig := config.GPU()
result, err := powermetrics.Collect(gpuConfig)
```

### Supported Formats

- `powermetrics.FormatText` - Plain text output
- `powermetrics.FormatPlist` - Property list (plist) output

### Supported Samplers

Currently supported samplers:
- `gpu_power` - GPU power and frequency metrics

## Output Structure

When using `FormatPlist`, the output is parsed into structured data:

```go
type PlistRoot struct {
    IsDelta      bool      `plist:"is_delta"`
    ElapsedNS    int64     `plist:"elapsed_ns"`
    HWModel      string    `plist:"hw_model"`
    KernOSVer    string    `plist:"kern_osversion"`
    KernBootArgs string    `plist:"kern_bootargs"`
    KernBootTime int64     `plist:"kern_boottime"`
    Timestamp    time.Time `plist:"timestamp"`
    GPU          GPUInfo   `plist:"gpu"`
}

type GPUInfo struct {
    FreqHz           float64      `plist:"freq_hz"`
    IdleNS           int64        `plist:"idle_ns"`
    IdleRatio        float64      `plist:"idle_ratio"`
    DVFMStates       []DVFMState  `plist:"dvfm_states"`
    SWRequestedState []SWReqState `plist:"sw_requested_state"`
    SWState          []SWState    `plist:"sw_state"`
    GPUEnergy        *int64       `plist:"gpu_energy,omitempty"`
}
```

## Examples

### Collecting GPU Metrics

```go
import (
    "fmt"
    "github.com/matiasinsaurralde/powermetrics/powermetrics"
)

func main() {
    config := &powermetrics.Config{SampleCount: 1}
    result, err := powermetrics.Collect(config.GPU())
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    if result.PlistData != nil {
        gpu := result.PlistData.GPU
        fmt.Printf("GPU Frequency: %.2f MHz\n", gpu.FreqHz/1e6)
        fmt.Printf("GPU Idle Ratio: %.1f%%\n", gpu.IdleRatio*100)
        fmt.Printf("Hardware Model: %s\n", result.PlistData.HWModel)
        
        // Print all DVFM states
        for i, state := range gpu.DVFMStates {
            fmt.Printf("DVFM State %d: %d Hz (%.1f%% used)\n", 
                i, state.Freq, state.UsedRatio*100)
        }
    }
}
```

### Error Handling

```go
config := &powermetrics.Config{
    Samplers: []string{"gpu_power", "invalid_sampler"},
}

result, err := powermetrics.Collect(config)
if err != nil {
    // This will fail validation
    fmt.Printf("Configuration error: %v\n", err)
    return
}
```

## Testing

Run the tests:

```bash
go test ./...
```

The test suite includes:
- Configuration validation
- Sampler validation
- GPU power XML parsing with real data
- JSON marshaling verification

## Building

```bash
go build -o powermetrics main.go
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## Notes

- This package requires the `powermetrics` command to be available on the system
- GPU power metrics are only available on supported macOS systems
- Some metrics may require elevated privileges 