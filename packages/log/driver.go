package log

// Driver identifies a logging channel driver.
type Driver string

const (
	DriverSingle   Driver = "single"
	DriverDaily    Driver = "daily"
	DriverSyslog   Driver = "syslog"
	DriverErrorlog Driver = "errorlog"
	DriverStack    Driver = "stack"
	DriverNull     Driver = "null"
	DriverCustom   Driver = "custom"
)
