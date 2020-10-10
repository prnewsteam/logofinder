package main

import (
	"encoding/json"
	"io"
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
		width = 150
	}

	height, err := strconv.ParseUint(query.Get("height"), 10, 64)
	if err != nil {
		height = 150
	}

	log.Printf("start search: %s", domain)

	logo, err := finder.FindLogo(domain)
	if err != nil {
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	rLogo, err := logo.Resize(uint(width), uint(height))
	if err != nil {
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	defer rLogo.Clear()

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=logo.png")
	w.WriteHeader(http.StatusOK)

	p := make([]byte, 1024)
	for {
		n, err := rLogo.File.Read(p)
		if err == io.EOF {
			break
		}
		w.Write(p[:n])
	}
}
