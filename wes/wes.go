// Package wes is the mother fucking wes
package wes

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gpbPiazza/pkg/db"
	"github.com/gpbPiazza/pkg/log"
	q "github.com/gpbPiazza/wes/db"
)

type WesHandler struct {
	Querier q.Querier
}

func NewHandler(querier q.Querier) *WesHandler {
	return &WesHandler{
		Querier: querier,
	}
}

func (wes *WesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		wes.post(w, r)
	case http.MethodGet:
		wes.get(w, r)
	default:
		log.FromContext(r.Context()).Error("http request received not implemented", "method_got", r.Method)
		w.WriteHeader(http.StatusNotImplemented)
	}
}

type PostBody struct {
	Moves []Move `json:"moves"`
}

func (wes *WesHandler) post(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := log.FromContext(ctx)
	if logger == nil {
		panic("ai middlweare nao foi")
	}

	bodyReader := r.Body

	bodyByte, err := io.ReadAll(bodyReader)
	if err != nil {
		log.FromContext(r.Context()).Error("read body fail", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer func() { _ = bodyReader.Close() }()

	body := new(PostBody)
	if err := json.Unmarshal(bodyByte, body); err != nil {
		log.FromContext(r.Context()).Error("unmarshal body fail", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	theMan := makeWes(body.Moves)

	dbx := db.FromContext(ctx)

	tx, err := dbx.Begin()
	if err != nil {
		log.FromContext(r.Context()).Error("fail to start tx", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer func() {
		if err := tx.Rollback(); err != nil {
			log.FromContext(r.Context()).Error("fail to rollback tx", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

	wesInBytes, err := json.Marshal(theMan)
	if err != nil {
		log.FromContext(r.Context()).Error("marshalling wes fail", "err", err)
		_, _ = w.Write([]byte("bro wes não consegue nem virar JSON, logo ele é um teapot"))
		w.WriteHeader(http.StatusTeapot)
	}

	_, _ = w.Write(wesInBytes)
	w.WriteHeader(http.StatusCreated)
}

func (wes *WesHandler) get(w http.ResponseWriter, r *http.Request) {
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
