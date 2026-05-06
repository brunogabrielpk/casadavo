package handler

import (
	"net/http"

	"casadavo/internal/model"
	"casadavo/internal/repository"
)

type LayoutHandler struct {
	repo *repository.LayoutRepo
}

func NewLayoutHandler(repo *repository.LayoutRepo) *LayoutHandler {
	return &LayoutHandler{repo}
}

func (h *LayoutHandler) ListExclusions(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	if date == "" {
		errResp(w, http.StatusBadRequest, "date query param required")
		return
	}
	exclusions, err := h.repo.ListByDate(date)
	if err != nil {
		errResp(w, http.StatusInternalServerError, err.Error())
		return
	}
	if exclusions == nil {
		exclusions = []model.TableExclusion{}
	}
	writeJSON(w, http.StatusOK, exclusions)
}

func (h *LayoutHandler) AddExclusion(w http.ResponseWriter, r *http.Request) {
	var e model.TableExclusion
	if err := readJSON(r, &e); err != nil {
		errResp(w, http.StatusBadRequest, "invalid body")
		return
	}
	if e.TableID == 0 || e.Date == "" {
		errResp(w, http.StatusBadRequest, "table_id and date are required")
		return
	}
	if err := h.repo.Add(&e); err != nil {
		errResp(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, e)
}

func (h *LayoutHandler) RemoveExclusion(w http.ResponseWriter, r *http.Request) {
	id, err := urlID(r, "id")
	if err != nil {
		errResp(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.repo.Remove(id); err != nil {
		errResp(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
