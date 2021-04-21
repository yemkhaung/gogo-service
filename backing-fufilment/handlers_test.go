package fufilment

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var targetSKU = "MUMBOJUMBO99"

func TestGetFufilmentStatus(t *testing.T) {
	// prepare
	server := NewServer()
	req, err := http.NewRequest("GET", "/skus/"+targetSKU, nil)
	resp := httptest.NewRecorder()
	server.ServeHTTP(resp, req)
	if err != nil {
		t.Errorf("Error creating request: %v", err)
	}

	// assert
	if resp.Code != http.StatusOK {
		t.Errorf("Expected HTTP status code of %d, got %d", http.StatusOK, resp.Code)
	}
	payload, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error in reading response: %v", err)
	}
	var fufilmentStatusResp fufilmentStatus
	err = json.Unmarshal(payload, &fufilmentStatusResp)
	if err != nil {
		t.Errorf("Error in decoding response JSON: %v", err)
	}
	if fufilmentStatusResp.SKU == "" || fufilmentStatusResp.SKU != targetSKU {
		t.Errorf("Expected SKU of '%s', got '%s'", targetSKU, fufilmentStatusResp.SKU)
	}
	if fufilmentStatusResp.ShipsWithin != 14 {
		t.Errorf("Expected ships_within of '14', got %d", fufilmentStatusResp.ShipsWithin)
	}
	if fufilmentStatusResp.QuantityInStock != 100 {
		t.Errorf("Expected quantity_in_stock of '100', got %d", fufilmentStatusResp.QuantityInStock)
	}

}
