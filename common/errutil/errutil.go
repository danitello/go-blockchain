package errutil

import "log"

// Handle displays errors to terminal
func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
