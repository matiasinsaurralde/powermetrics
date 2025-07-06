package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/matiasinsaurralde/powermetrics"
)

type GPUResponse struct {
	Timestamp time.Time `json:"timestamp"`
	GPU       struct {
		FrequencyHz  float64 `json:"frequency_hz"`
		FrequencyMHz float64 `json:"frequency_mhz"`
		IdleRatio    float64 `json:"idle_ratio"`
		IdlePercent  float64 `json:"idle_percent"`
		DVFMStates   []struct {
			FrequencyHz float64 `json:"frequency_hz"`
			UsedRatio   float64 `json:"used_ratio"`
			UsedPercent float64 `json:"used_percent"`
		} `json:"dvfm_states"`
	} `json:"gpu"`
	HardwareModel string `json:"hardware_model"`
	KernelVersion string `json:"kernel_version"`
}

func gpuMetricsHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers for web access
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Collect GPU metrics
	config := &powermetrics.Config{
		SampleCount:    1,
		SampleInterval: 1 * time.Second,
		Format:         powermetrics.FormatPlist,
		Samplers:       []powermetrics.Sampler{powermetrics.GPUPower},
	}

	result, err := powermetrics.Collect(config)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error collecting metrics: %v", err), http.StatusInternalServerError)
		return
	}

	if result.PlistData == nil {
		http.Error(w, "No GPU data available", http.StatusInternalServerError)
		return
	}

	// Build response
	response := GPUResponse{
		Timestamp:     result.PlistData.Timestamp,
		HardwareModel: result.PlistData.HWModel,
		KernelVersion: result.PlistData.KernOSVer,
	}

	// GPU data
	response.GPU.FrequencyHz = result.PlistData.GPU.FreqHz
	response.GPU.FrequencyMHz = result.PlistData.GPU.FreqHz / 1e6
	response.GPU.IdleRatio = result.PlistData.GPU.IdleRatio
	response.GPU.IdlePercent = result.PlistData.GPU.IdleRatio * 100

	// DVFM states
	for _, state := range result.PlistData.GPU.DVFMStates {
		dvfmState := struct {
			FrequencyHz float64 `json:"frequency_hz"`
			UsedRatio   float64 `json:"used_ratio"`
			UsedPercent float64 `json:"used_percent"`
		}{
			FrequencyHz: float64(state.Freq),
			UsedRatio:   state.UsedRatio,
			UsedPercent: state.UsedRatio * 100,
		}
		response.GPU.DVFMStates = append(response.GPU.DVFMStates, dvfmState)
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/gpu", gpuMetricsHandler)

	port := ":8080"
	fmt.Printf("Starting HTTP server on port %s\n", port)
	fmt.Printf("Access GPU metrics at: http://localhost%s/gpu\n", port)
	fmt.Printf("Press Ctrl+C to stop\n")

	log.Fatal(http.ListenAndServe(port, nil))
}
