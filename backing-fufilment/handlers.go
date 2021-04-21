package fufilment

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

func getFufilmentStatusHandler(formatter *render.Render) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		// route parameter
		vars := mux.Vars(req)
		sku := vars["sku"]
		// response payload JSON
		_ = formatter.JSON(rw, http.StatusOK, &fufilmentStatus{
			SKU:             sku,
			QuantityInStock: 100,
			ShipsWithin:     14,
		})
	}
}
