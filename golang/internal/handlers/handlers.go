package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	nanosecToMillisec = 1000 * 1000
	rowsinStorages    = 1000
)

type MessageResponse struct {
	Message string `json:"message"`
}

type TimingResponse struct {
	WallTimeMSec float64 `json:"wall_time_msec,omitepty"`
	TotalCycles  uint    `json:"total_cycles,omitepty"`
}

type StoragePayload interface {
	Request(ctx context.Context, count int)
	Init(ctx context.Context, rows int) error
}

func JSONError(w http.ResponseWriter, err interface{}, code int) {
	writeJSONResponse(w, err, code)
}

func JSONResponse(w http.ResponseWriter, data interface{}) {
	writeJSONResponse(w, data, http.StatusOK)
}

func writeJSONResponse(w http.ResponseWriter, data interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

func Ok(w http.ResponseWriter, r *http.Request) {
	JSONResponse(w, MessageResponse{"ok"})
}

func Hello(w http.ResponseWriter, r *http.Request) {
	JSONResponse(w, MessageResponse{"Hello, world!"})
}

func InitDBHandler(db, cache StoragePayload) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		err := db.Init(context.Background(), rowsinStorages)
		if err != nil {
			fmt.Printf("Error occured while initing db %s \n", err.Error())
			JSONError(writer, err, http.StatusBadRequest)
			return
		}
		err = cache.Init(context.Background(), rowsinStorages)
		if err != nil {
			JSONError(writer, err, http.StatusBadRequest)
			return
		}
		JSONResponse(writer, MessageResponse{"ok"})
		return
	}
}
