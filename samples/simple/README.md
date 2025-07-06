# Simple GPU Power Metrics Sample

This is a basic example that demonstrates how to use the powermetrics package to collect GPU power metrics.

## Features

- Collects 5 GPU power samples with 1-second intervals
- Displays GPU frequency, idle ratio, and energy consumption
- Shows how to access both single and multiple samples
- Uses the same code as the main README example

## Usage

```bash
cd samples/simple
go run main.go
```

## Output

The program will output:
- GPU frequency in Hz
- GPU idle ratio as a percentage
- GPU energy consumption in mJ (if available)

## Example Output

```
GPU Frequency: 338.00 Hz
GPU Idle Ratio: 99.53%
GPU Energy: 9 mJ
Collected 5 samples
Sample 0: GPU Idle Ratio: 99.53%
Sample 1: GPU Idle Ratio: 98.25%
Sample 2: GPU Idle Ratio: 96.97%
Sample 3: GPU Idle Ratio: 91.74%
Sample 4: GPU Idle Ratio: 96.07%
```

## Code Structure

This sample demonstrates:
1. Creating a powermetrics instance
2. Configuring for GPU power metrics
3. Collecting multiple samples (5 samples with 1-second intervals)
4. Accessing GPU data from the result
5. Handling multiple samples (when SampleCount > 1)

This is the simplest way to get started with the powermetrics package. 