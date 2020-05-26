package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"pg-upsert-returning/domain"
	"time"

	"github.com/google/uuid"
)

type customerRepository interface {
	UpsertCustomerCte(ctx context.Context, id uuid.UUID) (res domain.UpsertedRow, err error)
	UpsertCustomerLock(ctx context.Context, id uuid.UUID) (res domain.UpsertedRow, err error)
	UpsertCustomerConflict(ctx context.Context, id uuid.UUID) (res domain.UpsertedRow, err error)
	UpsertCustomerDoNothing(ctx context.Context, id uuid.UUID) (res domain.UpsertedRow, err error)
}

type Server struct {
	HttpAddr     string
	CustomerRepo customerRepository
}

func (s *Server) UpsertCustomer(repoUpsert func(ctx context.Context, id uuid.UUID) (res domain.UpsertedRow, err error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type upsertRequest struct {
			ID uuid.UUID `json:"id"`
		}
		type upsertResponse struct {
			IntID   int    `json:"intID"`
			CTID    string `json:"CTID"`
			XMAX    int    `json:"XMAX"`
			Elapsed int64  `json:"elapsedMilliseconds"`
		}

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var rq upsertRequest
		err = json.Unmarshal(bytes, &rq)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Call the repo
		start := time.Now()
		upserted, err := repoUpsert(r.Context(), rq.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
			return
		}
		t := time.Now()
		elapsed := t.Sub(start)

		respBody, err := json.Marshal(upsertResponse{upserted.ID, upserted.CTID, upserted.XMAX, elapsed.Milliseconds()})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_, _ = w.Write(respBody)
	}
}

func (s *Server) UpsertCustomerRandom(repoUpsert func(ctx context.Context, id uuid.UUID) (res domain.UpsertedRow, err error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type upsertResponse struct {
			IntID   int    `json:"intID"`
			CTID    string `json:"CTID"`
			XMAX    int    `json:"XMAX"`
			Elapsed int64  `json:"elapsedMilliseconds"`
		}

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// Call the repo
		start := time.Now()
		upserted, err := repoUpsert(r.Context(), uuid.New())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
			return
		}

		t := time.Now()
		elapsed := t.Sub(start)

		respBody, err := json.Marshal(upsertResponse{upserted.ID, upserted.CTID, upserted.XMAX, elapsed.Milliseconds()})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_, _ = w.Write(respBody)
	}
}

func (s *Server) ListenAndServe() error {
	routes := http.NewServeMux()
	routes.HandleFunc("/upsert-cte", s.UpsertCustomer(s.CustomerRepo.UpsertCustomerCte))
	routes.HandleFunc("/upsert-cte-random", s.UpsertCustomerRandom(s.CustomerRepo.UpsertCustomerCte))
	routes.HandleFunc("/upsert-lock", s.UpsertCustomer(s.CustomerRepo.UpsertCustomerLock))
	routes.HandleFunc("/upsert-conflict", s.UpsertCustomer(s.CustomerRepo.UpsertCustomerConflict))
	routes.HandleFunc("/upsert-donothing", s.UpsertCustomer(s.CustomerRepo.UpsertCustomerDoNothing))
	fmt.Println("Server UP")
	return http.ListenAndServe(s.HttpAddr, routes)
}
