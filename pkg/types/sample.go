package types

import "time"

type BaseSample struct {
	IsDelta      bool      `plist:"is_delta"`
	ElapsedNS    int64     `plist:"elapsed_ns"`
	HWModel      string    `plist:"hw_model"`
	KernOSVer    string    `plist:"kern_osversion"`
	KernBootArgs string    `plist:"kern_bootargs"`
	KernBootTime int64     `plist:"kern_boottime"`
	Timestamp    time.Time `plist:"timestamp"`
}

type Sample interface {
	GetTimestamp() time.Time
	GetElapsedNS() int64
	GetHWModel() string
	GetKernOSVer() string
	GetKernBootArgs() string
	GetKernBootTime() int64
	GetIsDelta() bool
}

func (s *BaseSample) GetTimestamp() time.Time {
	return s.Timestamp
}

func (s *BaseSample) GetElapsedNS() int64 {
	return s.ElapsedNS
}

func (s *BaseSample) GetHWModel() string {
	return s.HWModel
}

func (s *BaseSample) GetKernOSVer() string {
	return s.KernOSVer
}

func (s *BaseSample) GetKernBootArgs() string {
	return s.KernBootArgs
}

func (s *BaseSample) GetKernBootTime() int64 {
	return s.KernBootTime
}

func (s *BaseSample) GetIsDelta() bool {
	return s.IsDelta
}
