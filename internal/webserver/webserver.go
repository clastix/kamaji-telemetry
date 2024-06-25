// Copyright 2024 Clastix Labs
// SPDX-License-Identifier: Apache-2.0

//nolint:dupl
package webserver

import (
	"context"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/clastix/kamaji-telemetry/api"
	kamajictx "github.com/clastix/kamaji-telemetry/internal/ctx"
	"github.com/clastix/kamaji-telemetry/internal/db"
)

type WebServer struct {
	json   jsoniter.API
	conn   *pgx.Conn
	logger logr.Logger
}

func New(ctx context.Context, logger logr.Logger, conn *pgx.Conn, port int) error {
	ws := WebServer{
		json:   jsoniter.ConfigCompatibleWithStandardLibrary,
		logger: logger,
		conn:   conn,
	}

	r := mux.NewRouter()
	r.Use(handlers.RecoveryHandler())
	r.HandleFunc("/create", ws.recordCreate).Methods("POST")
	r.HandleFunc("/update", ws.recordUpdate).Methods("POST")
	r.HandleFunc("/delete", ws.recordDelete).Methods("POST")
	r.PathPrefix("/stats").HandlerFunc(ws.recordStats).Methods("POST")
	r.PathPrefix("/").HandlerFunc(ws.NoOp)

	srv := &http.Server{
		Handler:           r,
		Addr:              ":" + strconv.Itoa(port),
		ReadTimeout:       1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	<-ctx.Done()

	logger.Info("engaging server shutdown")

	gracefulCtx, cancel := kamajictx.GracefulTimeout(5 * time.Second)
	defer cancel()

	if err := srv.Shutdown(gracefulCtx); err != nil { //nolint:contextcheck
		return errors.Wrap(err, "failed to shutdown http server")
	}

	return nil
}

func (w *WebServer) NoOp(http.ResponseWriter, *http.Request) {
	// no-op
}

func (w *WebServer) recordCreate(writer http.ResponseWriter, request *http.Request) {
	logger := w.logger.WithValues("path", request.URL.Path, "ua", request.Header.Get("User-Agent"), "rid", uuid.New().String())

	data, err := io.ReadAll(request.Body)
	if err != nil {
		logger.Error(err, "error reading request body")
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	var payload api.Create

	if err = w.json.Unmarshal(data, &payload); err != nil {
		logger.Error(err, "error unmarshalling request body")
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	if payload.KamajiVersion == "" || payload.KubernetesVersion == "" || payload.TenantVersion == "" {
		logger.Info("request body does not contain required fields")
		writer.WriteHeader(http.StatusBadRequest)

		return
	}

	go db.Create(w.conn, data, logger) //nolint:contextcheck

	writer.WriteHeader(http.StatusOK)
}

func (w *WebServer) recordUpdate(writer http.ResponseWriter, request *http.Request) {
	logger := w.logger.WithValues("path", request.URL.Path, "ua", request.Header.Get("User-Agent"), "rid", uuid.New().String())

	data, err := io.ReadAll(request.Body)
	if err != nil {
		logger.Error(err, "error reading request body")
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	var payload api.Update

	if err = w.json.Unmarshal(data, &payload); err != nil {
		logger.Error(err, "error unmarshalling request body")
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	if payload.KamajiVersion == "" || payload.KubernetesVersion == "" || payload.OldTenantVersion == "" || payload.NewTenantVersion == "" {
		logger.Info("request body does not contain required fields")
		writer.WriteHeader(http.StatusBadRequest)

		return
	}

	go db.Update(w.conn, data, logger) //nolint:contextcheck

	writer.WriteHeader(http.StatusOK)
}

func (w *WebServer) recordDelete(writer http.ResponseWriter, request *http.Request) {
	logger := w.logger.WithValues("path", request.URL.Path, "ua", request.Header.Get("User-Agent"), "rid", uuid.New().String())

	data, err := io.ReadAll(request.Body)
	if err != nil {
		logger.Error(err, "error reading request body")
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	var payload api.Delete

	if err = w.json.Unmarshal(data, &payload); err != nil {
		logger.Error(err, "error unmarshalling request body")
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	if payload.KamajiVersion == "" || payload.KubernetesVersion == "" || payload.TenantVersion == "" || payload.Status == "" {
		logger.Info("request body does not contain required fields")
		writer.WriteHeader(http.StatusBadRequest)

		return
	}

	go db.Delete(w.conn, data, logger) //nolint:contextcheck

	writer.WriteHeader(http.StatusOK)
}

func (w *WebServer) recordStats(writer http.ResponseWriter, request *http.Request) {
	logger := w.logger.WithValues("path", request.URL.Path, "ua", request.Header.Get("User-Agent"), "rid", uuid.New().String())

	data, err := io.ReadAll(request.Body)
	if err != nil {
		logger.Error(err, "error reading request body")
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	var payload api.Stats

	if err = w.json.Unmarshal(data, &payload); err != nil {
		logger.Error(err, "error unmarshalling request body")
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	if payload.UUID == "" || payload.KamajiVersion == "" || payload.KubernetesVersion == "" {
		logger.Info("request body does not contain required fields")
		writer.WriteHeader(http.StatusBadRequest)

		return
	}

	go db.Stats(w.conn, data, logger, payload.UUID) //nolint:contextcheck

	writer.WriteHeader(http.StatusOK)
}
