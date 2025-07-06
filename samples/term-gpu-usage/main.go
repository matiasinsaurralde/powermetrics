package main

import (
	"context"
	"fmt"
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

	spark := widgets.NewSparkline()
	spark.LineColor = termui.ColorGreen
	spark.Data = []float64{}

	sparkGroup := widgets.NewSparklineGroup(spark)
	sparkGroup.Title = "GPU Idle Ratio (waiting for samples...)"
	sparkGroup.SetRect(0, 0, termWidth, termHeight)

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
			config := &powermetrics.Config{SampleCount: 1}
			result, err := powermetrics.Collect(config.GPU())
			if err != nil {
				sparkGroup.Title = "Error: " + err.Error()
			} else if result.PlistData != nil {
				idle := result.PlistData.GPU.IdleRatio * 100 // percent
				data := spark.Data
				if len(data) >= 30 {
					data = data[1:]
				}
				data = append(data, idle)
				spark.Data = data
				sparkGroup.Title = fmt.Sprintf("GPU Idle Ratio (%d samples) - Latest: %.2f%%", len(spark.Data), idle)
			}
			termui.Render(sparkGroup)
		case e := <-uiEvents:
			if e.Type == termui.KeyboardEvent {
				if e.ID == "q" || e.ID == "<C-c>" {
					return
				}
			}
			// Handle resize
			if e.Type == termui.ResizeEvent {
				termWidth, termHeight = termui.TerminalDimensions()
				sparkGroup.SetRect(0, 0, termWidth, termHeight)
				termui.Clear()
				termui.Render(sparkGroup)
			}
		}
	}
}
