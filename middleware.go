package assert

import (
	"log"
	"net/http"
	"os"
)

// This Middleware is required to properly handle the errors thrown using
// this assert package. It must be called before the asserts are used.
func Middleware() func(http.ResponseWriter, *http.Request, http.HandlerFunc) {
	l := log.New(os.Stdout, "[web] ", 0)

	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		defer func() {
			err := recover()
			if err == nil {
				return
			}

			switch assert := err.(type) {
			case assertError:
				if assert.statusCode == http.StatusInternalServerError {
					l.Printf("PANIC: %s\n%s", assert.Error(), assert.stack())
				}

				http.Error(rw, assert.Error(), assert.statusCode)
			default:
				panic(err)
			}
		}()

		next(rw, r)
	}
}
