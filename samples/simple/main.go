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
