// Command app chạy một HTTP server nhỏ để minh hoạ việc build ra binary.
// Các endpoint dùng lại logic trong package calculator.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"learnjenkin/calculator"
)

type result struct {
	Operation string `json:"operation"`
	A         int    `json:"a"`
	B         int    `json:"b"`
	Result    int    `json:"result"`
}

func parseParams(r *http.Request) (int, int, error) {
	a, err := strconv.Atoi(r.URL.Query().Get("a"))
	if err != nil {
		return 0, 0, fmt.Errorf("tham số 'a' không hợp lệ: %w", err)
	}
	b, err := strconv.Atoi(r.URL.Query().Get("b"))
	if err != nil {
		return 0, 0, fmt.Errorf("tham số 'b' không hợp lệ: %w", err)
	}
	return a, b, nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func newRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	mux.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		a, b, err := parseParams(r)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, result{"add", a, b, calculator.Add(a, b)})
	})

	mux.HandleFunc("/divide", func(w http.ResponseWriter, r *http.Request) {
		a, b, err := parseParams(r)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		res, err := calculator.Divide(a, b)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, result{"divide", a, b, res})
	})

	return mux
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}

	addr := ":" + port
	log.Printf("server đang chạy tại http://localhost%s", addr)
	if err := http.ListenAndServe(addr, newRouter()); err != nil {
		log.Fatal(err)
	}
}
