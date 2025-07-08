package types

type ResultCollection struct {
	Samples []Sample
}

func (rc *ResultCollection) GetGPUSamples() []*GPUPowerSample {
	var gpuSamples []*GPUPowerSample
	for _, sample := range rc.Samples {
		if gpuSample, ok := sample.(*GPUPowerSample); ok {
			gpuSamples = append(gpuSamples, gpuSample)
		}
	}
	return gpuSamples
}

func (rc *ResultCollection) GetBatterySamples() []*BatterySample {
	var batterySamples []*BatterySample
	for _, sample := range rc.Samples {
		if batterySample, ok := sample.(*BatterySample); ok {
			batterySamples = append(batterySamples, batterySample)
		}
	}
	return batterySamples
}
