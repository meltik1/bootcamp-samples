package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
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

type SleepPayload interface {
	Sleep(msec float64) (uint, error)
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

func SleepHandler(cpuPayload, ioPayload SleepPayload, db, cache StoragePayload) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		var totalCycles uint

		// cpu sleep
		if cpuMsecStr := r.URL.Query().Get("cpu_msec"); cpuMsecStr != "" {
			msec, err := strconv.ParseFloat(cpuMsecStr, 32)
			if err != nil {
				JSONError(w, err, http.StatusBadRequest)
				return
			}

			cycles, err := cpuPayload.Sleep(msec)
			if err != nil {
				//zap.L().Error("cpu sleep error", zap.Error(err))
				JSONError(w, err, http.StatusInternalServerError)
				return
			}
			totalCycles += cycles
		}

		// io sleep
		if ioMsecStr := r.URL.Query().Get("io_msec"); ioMsecStr != "" {
			msec, err := strconv.ParseFloat(ioMsecStr, 32)
			if err != nil {
				JSONError(w, err, http.StatusBadRequest)
				return
			}

			cycles, err := ioPayload.Sleep(msec)
			if err != nil {
				//zap.L().Error("io sleep error", zap.Error(err))
				JSONError(w, err, http.StatusInternalServerError)
				return
			}
			totalCycles += cycles
		}

		// pk DB
		if bdReqCount := r.URL.Query().Get("bd_pk"); bdReqCount != "" {
			reqCount, err := strconv.ParseInt(bdReqCount, 10, 32)
			if err != nil {
				JSONError(w, err, http.StatusBadRequest)
				return
			}
			db.Request(context.Background(), int(reqCount))
		}

		if cacheReqCount := r.URL.Query().Get("cache_pk"); cacheReqCount != "" {
			reqCount, err := strconv.ParseInt(cacheReqCount, 10, 32)
			if err != nil {
				JSONError(w, err, http.StatusBadRequest)
				return
			}
			cache.Request(context.Background(), int(reqCount))
		}

		end := time.Now()

		JSONResponse(w, TimingResponse{
			WallTimeMSec: float64(end.Sub(start).Nanoseconds()) / nanosecToMillisec,
			TotalCycles:  totalCycles,
		})
	}
}
