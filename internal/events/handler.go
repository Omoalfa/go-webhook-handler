package events

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	response "github.com/Omoalfa/event-platform/internal/httpx"
	"github.com/Omoalfa/event-platform/internal/webhook"
	"github.com/oapi-codegen/runtime/types"
)

type EventsHandler struct {
	service Service
	secret  string
}

func NewEventsHandler(service *Service, secret string) *EventsHandler {
	return &EventsHandler{
		service: *service,
		secret:  secret,
	}
}

func (h *EventsHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	if limit == "" {
		limit = "20"
	}

	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt <= 0 {
		limitInt = 20
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt <= 1 {
		pageInt = 1
	}

	offSet := (pageInt - 1) * limitInt

	events, err := h.service.List(r.Context(), limitInt, offSet)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Something went wrong", []any{err})
	}
	response.WriteJSON(w, http.StatusOK, "Events retrieved successfully", events)
}

func (h *EventsHandler) PostEvents(w http.ResponseWriter, r *http.Request) {
	var req CreateEventRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.WriteError(
			w,
			http.StatusBadRequest,
			"failed to read request body",
			[]any{err},
		)
		return
	}

	signature := r.Header.Get("X-Webhook-Signature")

	if !webhook.VerifySignature(
		body,
		signature,
		h.secret,
	) {
		response.WriteError(
			w,
			http.StatusUnauthorized,
			"Invalid webhook signature",
			[]any{},
		)
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(
			w,
			http.StatusBadRequest,
			"Invalid request payload",
			[]any{err},
		)
		return
	}

	cEvent, err := toCreateEvent(req)
	if err != nil {
		response.WriteError(
			w,
			http.StatusBadRequest,
			"Invalid request payload",
			[]any{err},
		)
		return
	}

	event, err := h.service.Create(r.Context(), cEvent)

	if err != nil {
		response.WriteError(
			w,
			http.StatusInternalServerError,
			"Failed to create event",
			[]any{err},
		)
		return
	}

	response.WriteJSON(
		w,
		http.StatusCreated,
		"Event created successfully",
		event,
	)
}

func (h *EventsHandler) GetEventsId(
	w http.ResponseWriter,
	r *http.Request,
	id types.UUID,
) {
	event, err := h.service.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrEventNotFound) {
			response.WriteError(
				w,
				http.StatusNotFound,
				"Event not found",
				nil,
			)
			return
		}
		response.WriteError(
			w,
			http.StatusInternalServerError,
			"Failed to retrieve event",
			[]any{err},
		)
		return
	}

	response.WriteJSON(
		w,
		http.StatusOK,
		"Event retrieved successfully",
		event,
	)
}

func (h *EventsHandler) PutEventsId(
	w http.ResponseWriter,
	r *http.Request,
	id types.UUID,
) {

	var req UpdateEventStatusRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(
			w,
			http.StatusBadRequest,
			"Invalid request payload",
			[]any{err},
		)
		return
	}

	event, err := h.service.UpdateStatus(
		r.Context(),
		id,
		string(req.Status),
	)

	if err != nil {
		if errors.Is(err, ErrEventNotFound) {
			response.WriteError(
				w,
				http.StatusNotFound,
				"Event not found",
				nil,
			)
			return
		}

		response.WriteError(
			w,
			http.StatusInternalServerError,
			"Failed to update event status",
			[]any{err},
		)
		return
	}

	response.WriteJSON(
		w,
		http.StatusOK,
		"Event status updated successfully",
		event,
	)
}

func toCreateEvent(req CreateEventRequest) (CreateEvent, error) {
	var payload []byte

	if req.Payload != nil {
		var err error

		payload, err = json.Marshal(*req.Payload)
		if err != nil {
			return CreateEvent{}, err
		}
	}

	return CreateEvent{
		Provider:  req.Provider,
		EventType: req.EventType,
		Payload:   payload,
	}, nil
}
