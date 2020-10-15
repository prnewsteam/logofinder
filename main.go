package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/prnewsteam/logofinder/finder"
)

func main() {
	http.HandleFunc("/logo", findLogo)
	http.ListenAndServe(":8099", nil)
}

func findLogo(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	domain := query.Get("domain")

	width, err := strconv.ParseUint(query.Get("width"), 10, 64)
	if err != nil {
		width = 0
	}

	height, err := strconv.ParseUint(query.Get("height"), 10, 64)
	if err != nil {
		height = 0
	}

	log.Printf("start search: %s", domain)

	logo, err := finder.FindLogo(domain)
	defer logo.Clear()
	if err != nil {
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	if (width == 0 || height == 0) {
		logo.WriteResponse(w)
		return
	}

	rLogo, err := logo.Resize(uint(width), uint(height))
	defer rLogo.Clear()
	if err != nil {
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	rLogo.WriteResponse(w)
}
