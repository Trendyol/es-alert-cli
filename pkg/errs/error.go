package errs

import (
	"github.com/labstack/gommon/log"
)

func LogError(err error, message string) bool {
	if err != nil {
		log.Errorf("%s: %v\n", message, err)
		return true
	}
	return false
}
