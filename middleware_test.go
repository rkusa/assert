package assert

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMiddleware(t *testing.T) {
	recorder := httptest.NewRecorder()

	mw := Middleware()
	mw(recorder, (*http.Request)(nil), func(res http.ResponseWriter, req *http.Request) {
		Error(fmt.Errorf("Fail"))
	})

	assert.Equal(t, recorder.Code, http.StatusInternalServerError)

	body, err := ioutil.ReadAll(recorder.Body)
	assert.NoError(t, err)

	assert.Equal(t, string(body), "Fail\n")
}
