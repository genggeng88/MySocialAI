package handler

import (
	"encoding/json"
	"fmt"
	"mysocialai/model"
	"net/http"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one upload request")
	decoder := json.NewDecoder(r.Body)
	var p model.Post
	if err := decoder.Decode(&p); err != nil {
		panic(err)
	}

	fmt.Fprintf(w, "Post received: %s\n", p.Message)
}