package api

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/orestrepov/metadatahost/app"
	"net/http"
)

func (a *API) SearchDomain(ctx *app.Context, w http.ResponseWriter, r *http.Request) error {

	hostName := chi.URLParam(r, "HostName")
	hostResponse, err := ctx.SearchDomain(hostName)
	if err != nil {
		return err
	}

	data, err := json.Marshal(hostResponse)
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	return err
}

func (a *API) DomainSearchHistory(ctx *app.Context, w http.ResponseWriter, r *http.Request) error {

	historyResponse, err := ctx.DomainSearchHistory()
	if err != nil {
		return err
	}

	data, err := json.Marshal(historyResponse)
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	return err
}