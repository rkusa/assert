// Simplified error handling for http routes using assert with status code.
//
//  func Login(ctx web.Context) {
//    username := ctx.FormValue("email")
//    password := ctx.FormValue("password")
//
//    assert.OK(username != "", 400, "No username given")
//    assert.OK(password != "", 400, "No password given")
//
//    ds := datastore.FromContext(ctx)
//
//    user, err := ds.Users.GetByUsername(username)
//    assert.OK(err == nil, 400, "Invalid username")
//
//    err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
//    assert.OK(err == nil, 400, "Invalid password")
//
//    session := web.GetSession(ctx)
//    session["user"] = user.Id
//
//    ctx.Redirect("/app")
//  }
package assert

import (
	"net/http"
	"runtime"
	"strings"
)

type assertError struct {
	statusCode int
	message    string
}

func (err assertError) Error() string {
	return err.message
}

func (err assertError) stack() string {
	buf := make([]byte, 32)
	for {
		n := runtime.Stack(buf, false)
		if n < len(buf) {
			break
		}
		buf = make([]byte, len(buf)*2)
	}

	stack := string(buf)
	lines := strings.Split(stack, "\n")
	return strings.Join(lines, "\n")
}

func ok(condition bool, statusCode int, message string) error {
	if !condition {
		if len(message) == 0 {
			message = http.StatusText(statusCode)
		}

		return assertError{statusCode, message}
	}

	return nil
}

// Success throws with the given statusCode and message if the provided
// condition evaluates to false. If message is an empty string, the default
// status description is used.
func OK(condition bool, statusCode int, message string) {
	if err := ok(condition, statusCode, message); err != nil {
		panic(err)
	}
}

// Success throws with the given statusCode and message if the provided error
// exists. If message is an empty string, the default status description is used.
func Success(err error, statusCode int, message string) {
	if e := ok(err == nil, statusCode, message); e != nil {
		panic(e)
	}
}

// Error throws and responds with an 500 Internal Server Error if the provided
// error exists.
func Error(err error) {
	if err != nil {
		panic(ok(false, http.StatusInternalServerError, err.Error()))
	}
}

// Assert represents an encapsulation for the assertions to provide an OnError
// hook.
type Assert interface {
	OnError(func())
	OK(bool, int, string)
	Success(error, int, string)
	Error(error)
}

type assertEncapsulation struct {
	onError func()
}

func (a *assertEncapsulation) throw(err error) {
	if a.onError != nil {
		a.onError()
	}

	panic(err)
}

// Register a callback that will be called once a assertion throws.
func (a *assertEncapsulation) OnError(fn func()) {
	a.onError = fn
}

// Success throws with the given statusCode and message if the provided
// condition evaluates to false. If message is an empty string, the default
// status description is used.
func (a *assertEncapsulation) OK(condition bool, statusCode int, message string) {
	if err := ok(condition, statusCode, message); err != nil {
		a.throw(err)
	}
}

// Success throws with the given statusCode and message if the provided error
// exists. If message is an empty string, the default status description is used.
func (a *assertEncapsulation) Success(err error, statusCode int, message string) {
	if e := ok(err == nil, statusCode, message); e != nil {
		a.throw(e)
	}
}

// Error throws and responds with an 500 Internal Server Error if the provided
// error exists.
func (a *assertEncapsulation) Error(err error) {
	if err != nil {
		a.throw(ok(false, http.StatusInternalServerError, err.Error()))
	}
}

// Create a new assertion encapsulation.
func New() Assert {
	return &assertEncapsulation{nil}
}
