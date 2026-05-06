package handler

import (
	"net/http"

	"casadavo/internal/model"
	"casadavo/internal/repository"
)

type TableHandler struct{ repo *repository.TableRepo }

func NewTableHandler(repo *repository.TableRepo) *TableHandler { return &TableHandler{repo} }

func (h *TableHandler) List(w http.ResponseWriter, r *http.Request) {
	tables, err := h.repo.List()
	if err != nil {
		errResp(w, http.StatusInternalServerError, err.Error())
		return
	}
	if tables == nil {
		tables = []model.Table{}
	}
	writeJSON(w, http.StatusOK, tables)
}

func (h *TableHandler) Create(w http.ResponseWriter, r *http.Request) {
	var t model.Table
	if err := readJSON(r, &t); err != nil {
		errResp(w, http.StatusBadRequest, "invalid body")
		return
	}
	if t.Number == 0 || t.Location == "" || t.Capacity == 0 {
		errResp(w, http.StatusBadRequest, "number, location and capacity are required")
		return
	}
	t.IsActive = true
	if err := h.repo.Create(&t); err != nil {
		errResp(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, t)
}

func (h *TableHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := urlID(r, "id")
	if err != nil {
		errResp(w, http.StatusBadRequest, "invalid id")
		return
	}
	var t model.Table
	if err := readJSON(r, &t); err != nil {
		errResp(w, http.StatusBadRequest, "invalid body")
		return
	}
	t.ID = id
	if err := h.repo.Update(&t); err != nil {
		errResp(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, t)
}

func (h *TableHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := urlID(r, "id")
	if err != nil {
		errResp(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.repo.Delete(id); err != nil {
		errResp(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
