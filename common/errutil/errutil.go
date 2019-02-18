package errutil

import "log"

/*HandleErr displays errors to terminal
@param err - the error in question
*/
func HandleErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
