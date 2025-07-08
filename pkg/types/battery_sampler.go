package types

type BatterySample struct {
	BaseSample
	Battery BatteryInfo `plist:"battery"`
}

type BatteryInfo struct {
	PercentCharge int `plist:"percent_charge"`
}
