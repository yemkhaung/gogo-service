package catalog

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

func getCataglogDetail(formatter *render.Render, client fulfilmentClient) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		sku := vars["sku"]
		fufilmentStatus, err := client.getFulfilmentStatus(sku)
		if err != nil {
			_ = formatter.Text(rw, http.StatusInternalServerError, "Error getting fufilment status of SKU:"+sku)
		} else {
			_ = formatter.JSON(rw, http.StatusOK, catalogDetail{
				SKU:             sku,
				ProductID:       uuid.NewString(),
				Description:     "this is product description",
				QuantityInStock: fufilmentStatus.QuantityInStock,
				ShipsWithin:     fufilmentStatus.ShipsWithin,
			})
		}
	}
}

func getAllCatalogs(formatter *render.Render, client fulfilmentClient) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		_ = formatter.JSON(rw, http.StatusOK, []catalogDetail{createSampleCatalogDetail("MUMBO111"), createSampleCatalogDetail("JUMBO111")})
	}
}
