package catalog

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var targetSKU = "MUMBOJUMBO99"

// fufilment web client stub
type fakeWebClient struct{}

func (fc fakeWebClient) getFulfilmentStatus(sku string) (status fufilmentStatusResponse, err error) {
	status = fufilmentStatusResponse{
		QuantityInStock: 100,
		ShipsWithin:     14,
		SKU:             targetSKU,
	}
	return status, nil
}

func TestGetCatalogDetail(t *testing.T) {
	// prepare
	fakeClient := fakeWebClient{}
	server := NewServerWithClient(fakeClient)
	req, err := http.NewRequest("GET", "/catalog/"+targetSKU, nil)
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
	var cdetail catalogDetail
	err = json.Unmarshal(payload, &cdetail)
	if err != nil {
		t.Errorf("Error decoding response JSON: %v", err)
	}
	if cdetail.SKU != targetSKU {
		t.Errorf("Expected catalog SKU of %s, got %s", targetSKU, cdetail.SKU)
	}
	if cdetail.QuantityInStock != 100 {
		t.Errorf("Expected catalog SKU of 100, got %d", cdetail.QuantityInStock)
	}
	if cdetail.ShipsWithin != 14 {
		t.Errorf("Expected catalog SKU of 14, got %d", cdetail.ShipsWithin)
	}
}

func TestGetAllCatalogs(t *testing.T) {
	// prepare
	server := NewServer()
	req, err := http.NewRequest("GET", "/catalog", nil)
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
	var catalogs []catalogDetail
	err = json.Unmarshal(payload, &catalogs)
	if err != nil {
		t.Errorf("Error decoding response JSON: %v", err)
	}
	if len(catalogs) != 2 {
		t.Errorf("Expected catalogs size of 1, got %d", len(catalogs))
	}
}
