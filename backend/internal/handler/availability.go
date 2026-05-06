package handler

import (
	"net/http"

	"casadavo/internal/model"
	"casadavo/internal/repository"
)

type AvailabilityHandler struct {
	avail *repository.AvailabilityRepo
	slots *repository.TimeSlotRepo
}

func NewAvailabilityHandler(avail *repository.AvailabilityRepo, slots *repository.TimeSlotRepo) *AvailabilityHandler {
	return &AvailabilityHandler{avail, slots}
}

func (h *AvailabilityHandler) List(w http.ResponseWriter, r *http.Request) {
	avs, err := h.avail.List()
	if err != nil {
		errResp(w, http.StatusInternalServerError, err.Error())
		return
	}
	if avs == nil {
		avs = []model.Availability{}
	}
	writeJSON(w, http.StatusOK, avs)
}

func (h *AvailabilityHandler) Create(w http.ResponseWriter, r *http.Request) {
	var a model.Availability
	if err := readJSON(r, &a); err != nil {
		errResp(w, http.StatusBadRequest, "invalid body")
		return
	}
	if a.Date == "" {
		errResp(w, http.StatusBadRequest, "date is required")
		return
	}
	a.IsOpen = true
	if err := h.avail.Create(&a); err != nil {
		errResp(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, a)
}

func (h *AvailabilityHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := urlID(r, "id")
	if err != nil {
		errResp(w, http.StatusBadRequest, "invalid id")
		return
	}
	var a model.Availability
	if err := readJSON(r, &a); err != nil {
		errResp(w, http.StatusBadRequest, "invalid body")
		return
	}
	a.ID = id
	if err := h.avail.Update(&a); err != nil {
		errResp(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, a)
}

// --- Slots sub-resource ---

func (h *AvailabilityHandler) ListSlots(w http.ResponseWriter, r *http.Request) {
	id, err := urlID(r, "id")
	if err != nil {
		errResp(w, http.StatusBadRequest, "invalid id")
		return
	}
	slots, err := h.slots.ListByAvailability(id)
	if err != nil {
		errResp(w, http.StatusInternalServerError, err.Error())
		return
	}
	if slots == nil {
		slots = []model.TimeSlot{}
	}
	writeJSON(w, http.StatusOK, slots)
}

func (h *AvailabilityHandler) CreateSlot(w http.ResponseWriter, r *http.Request) {
	id, err := urlID(r, "id")
	if err != nil {
		errResp(w, http.StatusBadRequest, "invalid id")
		return
	}
	var s model.TimeSlot
	if err := readJSON(r, &s); err != nil {
		errResp(w, http.StatusBadRequest, "invalid body")
		return
	}
	if s.SlotTime == "" {
		errResp(w, http.StatusBadRequest, "slot_time is required")
		return
	}
	s.AvailabilityID = id
	if err := h.slots.Create(&s); err != nil {
		errResp(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, s)
}

func (h *AvailabilityHandler) DeleteSlot(w http.ResponseWriter, r *http.Request) {
	id, err := urlID(r, "id")
	if err != nil {
		errResp(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.slots.Delete(id); err != nil {
		errResp(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
