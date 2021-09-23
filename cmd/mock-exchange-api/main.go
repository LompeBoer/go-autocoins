package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/LompeBoer/go-autocoins/internal/exchange/binance"
)

func main() {
	http.HandleFunc("/", handleRoot)

	log.Fatal(http.ListenAndServe(":80", nil))
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s%s", "https://fapi.binance.com", r.RequestURI)
	filename := binance.FilenameForURL(url)
	file, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(w, "ERROR: %s: %s\n", url, err.Error())
	}

	fmt.Printf("Serving %s as %s\n", url, filename)

	w.Header().Add("content-type", "application/json")
	w.Header().Add("expires", "0")
	w.Header().Add("pragma", "no-cache")
	fmt.Fprintf(w, "%s", string(file))
}
