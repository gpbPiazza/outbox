package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gpbPiazza/pkg/db"
	"github.com/gpbPiazza/pkg/log"
	"github.com/gpbPiazza/wes"
)

func main() {
	log, ctx := log.New(context.Background())
	dbz, err := db.Init(ctx)
	if err != nil {
		log.Error("listen and server", "err", err)
		return
	}

	fmt.Println(dbz)
	server := http.NewServeMux()

	wesHandler := &wes.WesHandler{}

	server.Handle("/wes", wesHandler)

	log.Info("starting wes server")

	if err := http.ListenAndServe("localhost:8080", server); err != nil {
		log.Error("listen and server", "err", err)
		return
	}
}
