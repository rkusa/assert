package assert

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOK(t *testing.T) {
	// condition evaluates to true
	func() {
		defer func() {
			err := recover()
			assert.Nil(t, err)
		}()

		OK(true, http.StatusNotFound, "")
	}()

	// condition evaluates to false - empty message provided
	func() {
		defer func() {
			err, ok := recover().(error)
			assert.True(t, ok)

			if assert.Error(t, err) {
				assert.Equal(t, err.Error(), "Not Found")
			}
		}()

		OK(false, http.StatusNotFound, "")
	}()

	// condition evaluates to false - with provided message
	func() {
		defer func() {
			err, ok := recover().(error)
			assert.True(t, ok)

			if assert.Error(t, err) {
				assert.Equal(t, err.Error(), "Invalid input")
			}
		}()

		OK(false, http.StatusBadRequest, "Invalid input")
	}()
}

func TestSuccess(t *testing.T) {
	// no error
	func() {
		defer func() {
			err := recover()
			assert.Nil(t, err)
		}()

		Success(nil, http.StatusNotFound, "")
	}()

	// error - empty message provided
	func() {
		defer func() {
			err, ok := recover().(error)
			assert.True(t, ok)

			if assert.Error(t, err) {
				assert.Equal(t, err.Error(), "Not Found")
			}
		}()

		Success(fmt.Errorf("Fail"), http.StatusNotFound, "")
	}()

	// error - with provided message
	func() {
		defer func() {
			err, ok := recover().(error)
			assert.True(t, ok)

			if assert.Error(t, err) {
				assert.Equal(t, err.Error(), "Invalid input")
			}
		}()

		Success(fmt.Errorf("Fail"), http.StatusBadRequest, "Invalid input")
	}()
}

func TestError(t *testing.T) {
	// no error
	func() {
		defer func() {
			err := recover()
			assert.Nil(t, err)
		}()

		Error(nil)
	}()

	// error
	func() {
		defer func() {
			err, ok := recover().(error)
			assert.True(t, ok)

			if assert.Error(t, err) {
				assert.Equal(t, err.Error(), "Fail")
			}
		}()

		Error(fmt.Errorf("Fail"))
	}()
}

func TestThrow(t *testing.T) {
	func() {
		defer func() {
			err, ok := recover().(error)
			assert.True(t, ok)

			if assert.Error(t, err) {
				assert.Equal(t, err.Error(), "Invalid input")
			}
		}()

		Throw(http.StatusBadRequest, "Invalid input")
	}()
}

func TestOnError(t *testing.T) {
	called := false
	defer func() {
		err, ok := recover().(error)
		assert.True(t, ok)

		if assert.Error(t, err) {
			assert.Equal(t, err.Error(), "Bad Request")
		}

		assert.True(t, called)
	}()

	as := New()
	as.OnError(func() {
		called = true
	})

	as.OK(false, http.StatusBadRequest, "")
}
