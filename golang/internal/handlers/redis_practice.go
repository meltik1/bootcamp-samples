package handlers

import (
	"context"
	"net/http"
)

type RedisPayload interface {
	Hash(ctx context.Context, smth string) error
}

func DO(db, cache RedisPayload) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		err := cache.Hash(context.Background(), "lol")
		if err != nil {
			JSONError(writer, err, http.StatusBadRequest)
			return
		}
		JSONResponse(writer, MessageResponse{"ok"})
		return
	}
}
