package fufilment

type fufilmentStatus struct {
	SKU             string `json:"sku"`
	ShipsWithin     int    `json:"ships_within"`
	QuantityInStock int    `json:"quantity_in_stock"`
}
