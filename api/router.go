package api

import (
	"github.com/bdzhalalov/resilient-scatter-gather/internal/handler"
	"github.com/bdzhalalov/resilient-scatter-gather/internal/service"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

func Router(logger *logrus.Logger, service *service.MockServices) *mux.Router {
	router := mux.NewRouter()

	h := handler.New(logger, service)

	group := router.PathPrefix("/api/v1").Subrouter()
	group.HandleFunc("/chat/summary", h.GetSummary).Methods(http.MethodGet)

	return router
}
