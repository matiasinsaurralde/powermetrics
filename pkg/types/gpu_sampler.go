package types

type GPUPowerSample struct {
	BaseSample
	GPU GPUInfo `plist:"gpu"`
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
