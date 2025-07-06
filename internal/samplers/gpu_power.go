package samplers

import "time"

// PlistRoot maps the root structure of the powermetrics plist output
// generated with --samplers=gpu_power
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

type DVFMState struct {
	Freq      int64   `plist:"freq"`
	UsedNS    int64   `plist:"used_ns"`
	UsedRatio float64 `plist:"used_ratio"`
}

type SWReqState struct {
	SWReqState string  `plist:"sw_req_state"`
	UsedNS     int64   `plist:"used_ns"`
	UsedRatio  float64 `plist:"used_ratio"`
}

type SWState struct {
	SWState   string  `plist:"sw_state"`
	UsedNS    int64   `plist:"used_ns"`
	UsedRatio float64 `plist:"used_ratio"`
}
