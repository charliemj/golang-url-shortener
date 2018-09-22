package main

import (
	//"fmt"
	//"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func makeRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

//Testing strategy:
// - Test routes
// - Test that slugs are unique

func TestGetIndex(t *testing.T) {
	router := MakeRouter()
	w := makeRequest(router, "GET", "/")
	p, _ := ioutil.ReadAll(w.Body)

	title := "<title>Url Shortener</title>"
	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, strings.Index(string(p), title) > 0) //test that title exists in the response

}

//test a URL that is in the database and one that isn't
//func TestExpandUrl() {}

//test that a long URL can be made into short URL 1, and then a
//the same long URL can be put in again and be mapped to a different
//short URL 2.
//func TestCreateShortUrl() {}

//test that we cannot make a slug that already exists
//idk how to actually definitively test this given the way the
//code is structured
//func TestMakeUniqueSlug() {}
