package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLatestExchangeRate(t *testing.T) {
	// Start a local HTTP serverto  test the client
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		exchangeMap := make(map[string]string)
		exchangeMap["BTC"] = "0.0000159009877296"
		exchangeMap["ETH"] = "0.0003213362446397"

		resp := &Response{Data: struct {
			Currancy string            "json:\"currency\""
			Rates    map[string]string "json:\"rates\""
		}{"USD", exchangeMap}}
		// Marshal the object to JSON.
		jsonBytes, err := json.Marshal(resp)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		// Write the response to the response writer.
		rw.Write(jsonBytes)
	}))
	// Close the server when test finishes
	defer server.Close()

	// Use Client & URL from our local test server
	api := API{server.Client(), server.URL}
	v1, v2, err := api.GetLatestExchangeRate()

	if assert.NotNil(t, v1) {
		assert.Equal(t, v1, 0.0000159009877296)
	}

	if assert.NotNil(t, v2) {
		assert.Equal(t, v2, 0.0003213362446397)
	}
	assert.Nil(t, err)

}

func TestGenerateCryptoPortfolio(t *testing.T) {
	allocation, err := GenerateCryptoPortfolio(1000.00)
	if assert.Nil(t, err) {
		assert.Contains(t, allocation, "BTC")
		assert.Contains(t, allocation, "ETH")
	}

}
