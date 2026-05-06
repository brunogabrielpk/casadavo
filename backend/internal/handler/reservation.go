package handler

import (
	"net/http"

	"casadavo/internal/middleware"
	"casadavo/internal/model"
	"casadavo/internal/repository"
	"casadavo/internal/service"
)

type ReservationHandler struct {
	svc  *service.ReservationService
	repo *repository.ReservationRepo
}

func NewReservationHandler(svc *service.ReservationService, repo *repository.ReservationRepo) *ReservationHandler {
	return &ReservationHandler{svc, repo}
}

func (h *ReservationHandler) List(w http.ResponseWriter, r *http.Request) {
	c := middleware.GetClaims(r)
	var (
		list []model.ReservationDetail
		err  error
	)
	if c.Role == "gerente" {
		list, err = h.repo.ListAll()
	} else {
		list, err = h.repo.ListByUser(c.UserID)
	}
	if err != nil {
		errResp(w, http.StatusInternalServerError, err.Error())
		return
	}
	if list == nil {
		list = []model.ReservationDetail{}
	}
	writeJSON(w, http.StatusOK, list)
}

func (h *ReservationHandler) Create(w http.ResponseWriter, r *http.Request) {
	c := middleware.GetClaims(r)
	var res model.Reservation
	if err := readJSON(r, &res); err != nil {
		errResp(w, http.StatusBadRequest, "invalid body")
		return
	}
	if res.TableID == 0 || res.SlotID == 0 || res.Date == "" || res.PartySize == 0 {
		errResp(w, http.StatusBadRequest, "table_id, slot_id, date and party_size are required")
		return
	}
	res.UserID = c.UserID
	if err := h.svc.Create(&res); err != nil {
		errResp(w, http.StatusConflict, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, res)
}

func (h *ReservationHandler) Update(w http.ResponseWriter, r *http.Request) {
	c := middleware.GetClaims(r)
	id, err := urlID(r, "id")
	if err != nil {
		errResp(w, http.StatusBadRequest, "invalid id")
		return
	}
	var res model.Reservation
	if err := readJSON(r, &res); err != nil {
		errResp(w, http.StatusBadRequest, "invalid body")
		return
	}
	res.ID = id
	if err := h.svc.Update(&res, c.UserID, c.Role == "gerente"); err != nil {
		code := http.StatusInternalServerError
		if err.Error() == "forbidden" {
			code = http.StatusForbidden
		} else if err.Error() == "reservation not found" {
			code = http.StatusNotFound
		}
		errResp(w, code, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *ReservationHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id, err := urlID(r, "id")
	if err != nil {
		errResp(w, http.StatusBadRequest, "invalid id")
		return
	}
	var body struct {
		Status string `json:"status"`
	}
	if err := readJSON(r, &body); err != nil {
		errResp(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.UpdateStatus(id, body.Status); err != nil {
		errResp(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": body.Status})
}

func (h *ReservationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	c := middleware.GetClaims(r)
	id, err := urlID(r, "id")
	if err != nil {
		errResp(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.Delete(id, c.UserID, c.Role == "gerente"); err != nil {
		code := http.StatusInternalServerError
		if err.Error() == "forbidden" {
			code = http.StatusForbidden
		} else if err.Error() == "reservation not found" {
			code = http.StatusNotFound
		}
		errResp(w, code, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
