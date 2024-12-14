package init

import (
	"context"
	"fmt"
	"net/http"

	"github.com/meltik/study/internal/infra/jsonResponses"
)

const (
	rowsinStorages = 1000
)

type MessageResponse struct {
	Message string `json:"message"`
}

type StoragePayload interface {
	Init(ctx context.Context, rows int) error
}

func InitDBHandler(db, cache StoragePayload) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		err := db.Init(context.Background(), rowsinStorages)
		if err != nil {
			fmt.Printf("Error occured while initing db %s \n", err.Error())
			jsonResponses.JSONError(writer, err, http.StatusBadRequest)
			return
		}
		err = cache.Init(context.Background(), rowsinStorages)
		if err != nil {
			jsonResponses.JSONError(writer, err, http.StatusBadRequest)
			return
		}
		jsonResponses.JSONResponse(writer, MessageResponse{"ok"})
		return
	}
}
