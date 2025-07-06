package main

import (
	"context"
	"fmt"
	"math"
	"os"
	"os/signal"
	"time"

	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/matiasinsaurralde/powermetrics"
)

func main() {
	if err := termui.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize termui: %v\n", err)
		os.Exit(1)
	}
	defer termui.Close()

	// Get terminal size
	termWidth, termHeight := termui.TerminalDimensions()

	// Create a grid to use the full screen
	grid := termui.NewGrid()
	grid.SetRect(0, 0, termWidth, termHeight)

	// Create title
	title := widgets.NewParagraph()
	title.Text = "GPU Power Metrics Dashboard"
	title.TextStyle = termui.NewStyle(termui.ColorYellow, termui.ColorClear, termui.ModifierBold)
	title.BorderStyle = termui.NewStyle(termui.ColorBlue)
	title.Border = true

	// Create current metrics display
	metrics := widgets.NewParagraph()
	metrics.Title = "Current GPU Metrics"
	metrics.Text = "Initializing..."
	metrics.Border = true

	// Create sparkline for GPU idle ratio
	spark := widgets.NewSparkline()
	spark.LineColor = termui.ColorGreen
	spark.Data = []float64{}

	sparkGroup := widgets.NewSparklineGroup(spark)
	sparkGroup.Title = "GPU Idle Ratio History"
	sparkGroup.Border = true

	// Set up the grid layout to use full screen
	grid.Set(
		termui.NewRow(0.1, title),
		termui.NewRow(0.3, metrics),
		termui.NewRow(0.6, sparkGroup),
	)

	// Initial render
	termui.Render(grid)

	uiEvents := termui.PollEvents()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Collect a new sample
			pm := powermetrics.New()
			config := powermetrics.DefaultConfig().GPU()
			config.SampleCount = 1
			config.SampleRate = 1 * time.Second
			result, err := pm.Collect(config)
			if err != nil {
				metrics.Text = "Error: " + err.Error()
				sparkGroup.Title = "GPU Idle Ratio History (Error)"
			} else if result.PlistData != nil {
				gpu := result.PlistData.GPU
				// Show frequency in Hz (the raw value from powermetrics)
				freqStr := fmt.Sprintf("%.0f Hz", gpu.FreqHz)

				energyStr := "N/A"
				if gpu.GPUEnergy != nil {
					energyStr = fmt.Sprintf("%d mW", *gpu.GPUEnergy)
				}

				var idleStr string
				// More robust bounds checking to prevent NaN
				if !math.IsNaN(gpu.IdleRatio) && !math.IsInf(gpu.IdleRatio, 0) && gpu.IdleRatio >= 0 && gpu.IdleRatio <= 1 {
					idlePercent := gpu.IdleRatio * 100
					// Additional check to ensure the percentage is valid
					if !math.IsNaN(idlePercent) && !math.IsInf(idlePercent, 0) && idlePercent >= 0 && idlePercent <= 100 {
						idleStr = fmt.Sprintf("%.2f%%", idlePercent)
						// Update sparkline data
						data := spark.Data
						if len(data) >= 30 {
							data = data[1:]
						}
						data = append(data, idlePercent)
						spark.Data = data
						sparkGroup.Title = fmt.Sprintf("GPU Idle Ratio History (%d samples) - Latest: %.2f%%", len(spark.Data), idlePercent)
					} else {
						idleStr = "N/A"
						sparkGroup.Title = "GPU Idle Ratio History (Invalid data)"
					}
				} else {
					idleStr = "N/A"
					sparkGroup.Title = "GPU Idle Ratio History (No valid data)"
				}

				metrics.Text = fmt.Sprintf(
					"GPU Frequency: %s\nGPU Idle Ratio: %s\nGPU Energy: %s",
					freqStr, idleStr, energyStr,
				)
			}
			termui.Render(grid)
		case e := <-uiEvents:
			if e.Type == termui.KeyboardEvent {
				if e.ID == "q" || e.ID == "<C-c>" {
					return
				}
			}
			// Handle resize
			if e.Type == termui.ResizeEvent {
				termWidth, termHeight = termui.TerminalDimensions()
				grid.SetRect(0, 0, termWidth, termHeight)
				termui.Clear()
				termui.Render(grid)
			}
		}
	}
}
