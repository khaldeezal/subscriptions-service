package controller

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/khaldeezal/subscriptions-service/internal/domain"
	"github.com/khaldeezal/subscriptions-service/internal/utils"
	"net/http"
	"strconv"
	"time"
)

type Controller struct {
	subscriptionService serviceSubscription
}

func NewController(subscriptionService serviceSubscription) *Controller {
	return &Controller{
		subscriptionService: subscriptionService,
	}
}

func (h *Controller) InitRoutes(r chi.Router) http.Handler {
	r.Route("/api/v1/subscriptions", func(r chi.Router) {
		r.Post("/", h.CreateSubscription)
		r.Get("/list", h.ListSubscription)

		r.Get("/{id}", h.Subscription)
		r.Put("/{id}", h.UpdateSubscription)
		r.Delete("/{id}", h.DeleteSubscription)

		r.Get("/cost", h.GetTotalSubscriptionCost)
	})
	return r
}

func (h *Controller) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	in := &domain.CreateInput{}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.subscriptionService.Create(r.Context(), in)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{"id": id})
}

func (h *Controller) Subscription(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	res, err := h.subscriptionService.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(res)
}

func (h *Controller) ListSubscription(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	f := &domain.ListFilter{}
	if v := q.Get("user_id"); v != "" {
		f.UserID = &v
	}
	if v := q.Get("service_name"); v != "" {
		f.ServiceName = &v
	}
	if v := q.Get("page"); v != "" {
		if n, _ := strconv.Atoi(v); n > 0 {
			f.Page = n
		}
	}
	if v := q.Get("page_size"); v != "" {
		if n, _ := strconv.Atoi(v); n > 0 {
			f.PageSize = n
		}
	}

	list, err := h.subscriptionService.List(r.Context(), f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(list)
}

func (h *Controller) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	in := &domain.CreateInput{}

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.subscriptionService.Update(r.Context(), id, in); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]any{"status": "ok"})
}

func (h *Controller) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.subscriptionService.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Controller) GetTotalSubscriptionCost(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	fromS := q.Get("from")
	toS := q.Get("to")
	if fromS == "" || toS == "" {
		http.Error(w, "from/to required", http.StatusBadRequest)
		return
	}

	from, err := utils.ParseYearMonth(fromS)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	to, err := utils.ParseYearMonth(toS)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if to.Before(from) {
		http.Error(w, "to must be >= from", http.StatusBadRequest)
		return
	}

	f := &domain.ListFilter{}
	if v := q.Get("user_id"); v != "" {
		f.UserID = &v
	}

	if v := q.Get("service_name"); v != "" {
		f.ServiceName = &v
	}

	f.PageSize = 1000
	list, err := h.subscriptionService.List(r.Context(), f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := 0
	for _, s := range list {
		start, _ := utils.ParseYearMonth(s.StartDate)

		var e *time.Time
		if s.EndDate != nil && *s.EndDate != "" {
			et, _ := utils.ParseYearMonth(*s.EndDate)
			e = &et
		}

		periodEnd := to
		months := utils.MonthsOverlapInclusive(start, e, from, &periodEnd)
		if months > 0 {
			res += months * s.Price
		}
	}

	_ = json.NewEncoder(w).Encode(map[string]any{"total": res})
}
