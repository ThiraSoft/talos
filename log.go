package talos

const (
	DEBUG_LEVEL_NONE   = "NONE"
	DEBUG_LEVEL_ERRORS = "ERRORS"
	DEBUG_LEVEL_ALL    = "ALL"
)

var DEBUG_LEVEL DebugLevel = DEBUG_LEVEL_ALL

func logger(message string, level ...DebugLevel) {
	if message == "" {
		return
	}
	if checkDebugLevel(level, DEBUG_LEVEL) {
		println(message)
	}
}

// Verify if the provided debug level is enabled
func checkDebugLevel(s []DebugLevel, level DebugLevel) bool {
	for _, v := range s {
		if v == level {
			return true
		}
	}
	return false
}
