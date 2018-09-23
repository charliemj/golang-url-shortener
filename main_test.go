package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestMain(t *testing.T) {
	MakeDatabase()
	assert.True(t, true, true) //if MakeDatabase fails, we won't get to this assert
}

func makeRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestGetIndex(t *testing.T) {
	router := MakeRouter()
	w := makeRequest(router, "GET", "/")
	p, _ := ioutil.ReadAll(w.Body)

	title := "<title>Url Shortener</title>"
	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, strings.Index(string(p), title) > 0) //test that title exists in the response

}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router := MakeRouter()
	router.ServeHTTP(rr, req)

	return rr
}

func TestGetForUrlNotInDB(t *testing.T) {
	req, _ := http.NewRequest("GET", "/11", nil)
	response := executeRequest(req)

	assert.Equal(t, http.StatusNotFound, response.Code)
}

func TestMakeAShortURLAndInDB(t *testing.T) {
	data := url.Values{}
	data.Add("longUrl", "http://thetech.com/about/staff")
	req, _ := http.NewRequest("POST", "/create", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response := executeRequest(req)

	node, _ := html.Parse(response.Body)
	document := goquery.NewDocumentFromNode(node)
	var url string
	url = document.Find("#url").Text()

	assert.Equal(t, http.StatusOK, response.Code)
	var s []string
	s = strings.Split(url, "/")
	var slug string
	slug = s[len(s)-1]

	req2, _ := http.NewRequest("GET", "/"+slug, nil)
	response2 := executeRequest(req2)
	assert.Equal(t, 301, response2.Code)
}
