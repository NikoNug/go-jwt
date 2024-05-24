package productcontroller

import (
	"golang-jwt-practice/helpers"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	data := []map[string]interface{}{
		{
			"id":           1,
			"nama_product": "Kemeja",
			"Stock":        1000,
		},
		{
			"id":           2,
			"nama_product": "Celana",
			"Stock":        10000,
		},
		{
			"id":           3,
			"nama_product": "Sepatu",
			"Stock":        500,
		},
	}

	helpers.ResponseJSON(w, http.StatusOK, data)
}
