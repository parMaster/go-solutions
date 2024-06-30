package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func NewTestServer() *Server {
	s := Server{}

	return &s
}

func TestSetGet(t *testing.T) {

	s := NewTestServer()

	// SET
	req := httptest.NewRequest("GET", "/?id=1", nil)
	w := httptest.NewRecorder()
	s.HandleSetId(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// GET
	getReq := httptest.NewRequest("GET", "/", nil)
	w = httptest.NewRecorder()
	s.HandleGetId(w, getReq)
	resp = w.Result()

	// Validate
	var body map[string]any
	json.NewDecoder(resp.Body).Decode(&body)
	// fmt.Println(body)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	id, ok := body["Id"].(float64)
	assert.True(t, ok)
	assert.InDelta(t, 1, id, 0.00001)

	result, ok := body["Result"].(string)
	assert.True(t, ok)
	assert.Equal(t, "OK", result)
}

func NewTestClient() *Client {
	return &Client{}
}

func Test_ClientRequest(t *testing.T) {

	// how to test a client that makes a json request? and returns a json encoded response?

	expected := "{'data': 'dummy'}"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprint(w, expected)

	}))
	defer srv.Close()

	c := NewTestClient()
	res, err := c.echoRequest(srv.URL)

	assert.NoError(t, err)

	log.Println(res)
	assert.Equal(t, res, expected)

}
