package main

import (
	"encoding/json"
	"net/http"

	"github.com/vital-dhaveloose/aldb/examples"
)

func main() {
	http.HandleFunc("/activities", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		bts, _ := json.Marshal(examples.CreateExampleData())
		w.Write(bts)
	})

	http.ListenAndServe(":8080", nil)
}
