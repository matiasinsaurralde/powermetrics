package types

import "time"

type BaseSample struct {
	IsDelta      bool
	ElapsedNS    int64
	HWModel      string
	KernOSVer    string
	KernBootArgs string
	KernBootTime int64
	Timestamp    time.Time
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
