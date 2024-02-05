package errs

import "fmt"

func HandleError(err error, message string) bool {
	if err != nil {
		fmt.Printf("%s: %v\n", message, err)
		return true
	}
	return false
}
