// Package wes is the mother fucking wes
package wes

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gpbPiazza/pkg/log"
)

type WesHandler struct {
	repo Repository
}

func NewHandler(repo Repository) *WesHandler {
	return &WesHandler{
		repo: repo,
	}
}

func RegisterHandlers(ctx context.Context, dbx *sql.DB, server *http.ServeMux) {
	repo := NewRepository(dbx)
	wesHandler := NewHandler(repo)

	server.HandleFunc("POST /wes", ctxMiddleware(wesHandler.post(), ctx))
	server.HandleFunc("GET /wes/{id}", ctxMiddleware(wesHandler.get(), ctx))
}

func ctxMiddleware(next http.Handler, ctx context.Context) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

type PostBody struct {
	Moves []Move `json:"moves"`
}

func (wes *WesHandler) post() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := log.FromContext(ctx)

		bodyReader := r.Body

		bodyByte, err := io.ReadAll(bodyReader)
		if err != nil {
			logger.Error("read body fail", "err", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		defer func() { _ = bodyReader.Close() }()

		body := new(PostBody)
		if err := json.Unmarshal(bodyByte, body); err != nil {
			logger.Error("unmarshal body fail", "err", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		theMan := makeWes(body.Moves)

		var wesDB Wes
		err = wes.repo.RunInTx(ctx, func(query Querier) error {
			wesDB, err = query.CreateWes(ctx, theMan.Height)
			if err != nil {
				logger.Error("fail to create wes row", "err", err)
				return err
			}

			var mvsDB []Move
			var mvDB Move
			for _, mv := range theMan.Moves {
				mvDB, err = query.CreateMove(ctx, wesDB.ID, mv.Status, mv.Description)
				if err != nil {
					logger.Error("fail to create wes row", "err", err)
					return err
				}
				mvsDB = append(mvsDB, mvDB)
			}
			wesDB.Moves = mvsDB

			return nil
		})

		if err != nil {
			logger.Error("fail to start tx", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		wesInBytes, err := json.Marshal(wesDB)
		if err != nil {
			logger.Error("marshalling wes fail", "err", err)
			w.WriteHeader(http.StatusTeapot)
			_, _ = w.Write([]byte("bro wes não consegue nem virar JSON, logo ele é um teapot"))
		}

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write(wesInBytes)
	}
}

func (wes *WesHandler) get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wesID := r.PathValue("id")

		id, err := uuid.Parse(wesID)

		log.FromContext(r.Context()).Info("get wes start", "wes_id", id)

		if err != nil {
			log.FromContext(r.Context()).Error("get wes fail", "err", err)

			w.WriteHeader(http.StatusTeapot)
			_, _ = w.Write([]byte("bro manda o wes ID certo po, logo você é um teapot"))
			return
		}

		wesDB, err := wes.repo.GetWesByID(r.Context(), id)
		if err != nil {
			log.FromContext(r.Context()).Error("get wes fail", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		wesInBytes, err := json.Marshal(wesDB)
		if err != nil {
			log.FromContext(r.Context()).Error("marshalling wes fail", "err", err)
			w.WriteHeader(http.StatusTeapot)
			_, _ = w.Write([]byte("bro wes não consegue nem virar JSON, logo ele é um teapot"))
		}

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write(wesInBytes)
	}
}

type Wes struct {
	ID        uuid.UUID `json:"id"`
	Height    string    `json:"hight"`
	Moves     []Move    `json:"moves"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MoveStatus string

const (
	BagreMS            MoveStatus = "Wes sendo Bagre"
	MakingTheGood      MoveStatus = "Wes fez a boa"
	BookUnderTheArms   MoveStatus = "Wes com o livro de baixo do braço"
	TheNoneStoryTeller MoveStatus = "Wes esquecendo a historia enquanto conta"
)

type Move struct {
	ID          uuid.UUID  `json:"id"`
	WesID       uuid.UUID  `json:"wes_id"`
	Status      MoveStatus `json:"status"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
}

func makeWes(movesToAdd []Move) *Wes {
	moves := func() []Move {
		var moreMoves []Move
		for _, m := range movesToAdd {
			moreMoves = append(moreMoves, Move{
				ID:          uuid.New(),
				Status:      m.Status,
				Description: m.Description,
			})
		}

		return append([]Move{{ID: uuid.New(), Status: BagreMS, Description: "wes nasceu, logo é um bagre move"}},
			moreMoves...,
		)
	}

	w := &Wes{
		ID:     uuid.New(),
		Height: "baixo, wes é pequeno, minúsculo! por isso usamos uma string aqui, não tem medida na terra que pega essa pequenesa toda dele",
		Moves:  moves(),
	}

	return w
}
