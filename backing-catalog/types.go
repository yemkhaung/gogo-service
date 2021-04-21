package catalog

import "github.com/google/uuid"

type catalogDetail struct {
	SKU             string `json:"sku"`
	ProductID       string `json:"product_id"`
	Price           int    `json:"price"`
	Description     string `json:"description"`
	ShipsWithin     int    `json:"ships_within"`
	QuantityInStock int    `json:"quantity_in_stock"`
}

type fufilmentStatusResponse struct {
	SKU             string `json:"sku"`
	ShipsWithin     int    `json:"ships_within"`
	QuantityInStock int    `json:"quantity_in_stock"`
}

func createSampleCatalogDetail(sku string) catalogDetail {
	return catalogDetail{
		SKU:             sku,
		ProductID:       uuid.NewString(),
		Description:     "this is a sample",
		Price:           1000,
		ShipsWithin:     14,
		QuantityInStock: 100,
	}
}
