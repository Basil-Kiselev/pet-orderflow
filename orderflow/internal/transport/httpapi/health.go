package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (h *Handler) healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Add("Allow", http.MethodGet)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	ok := healthResponse{
		Status: "ok",
	}
	w.Header().Add("Content-type", "application/json")

	encoder := json.NewEncoder(w)
	err := encoder.Encode(ok)
	if err != nil {
		fmt.Printf("fail to encode: %v", err)
		return
	}
}
