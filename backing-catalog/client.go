package catalog

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type fulfilmentClient interface {
	getFulfilmentStatus(sku string) (fufilmentStatusResponse, error)
}

type fufilmentWebClient struct {
	url string
}

func (c fufilmentWebClient) getFulfilmentStatus(sku string) (status fufilmentStatusResponse, err error) {
	client := &http.Client{}
	req, _ := http.NewRequest(
		"GET",
		"http://"+c.url+"/skus/"+sku,
		nil,
	)

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error getting the fufilment status: %v\n", err)
		return status, err
	}

	defer resp.Body.Close()

	payload, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(payload, &status)
	if err != nil {
		log.Printf("Error decoding response JSON: %v", err)
		return status, err
	}

	return status, nil
}
