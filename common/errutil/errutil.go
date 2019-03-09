package errutil

import "log"

/*Handle displays errors to terminal
@param err - the error in question
*/
func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
