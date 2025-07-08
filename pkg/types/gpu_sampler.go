package types

type GPUPowerSample struct {
	BaseSample
	GPU GPUInfo `plist:"gpu"`
}

type GPUInfo struct {
}
