package orm

// log is the global logger that will be used if not struct-specific logger was defined
var log Logger

// Logf is used by orm implementations to log if no implementation specific logger was set
func Logf(format string, args ...interface{}) {
	if log == nil {
		return
	}
	log(format, args...)
}
