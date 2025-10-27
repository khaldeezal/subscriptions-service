package subscriptions

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct{ repo *Repository }

func NewHandler(pool *pgxpool.Pool) *Handler { return &Handler{repo: NewRepository(pool)} }

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()

	// коллекционные роуты — только со слэшем
	r.Post("/", h.Create)
	r.Get("/", h.List)

	r.Get("/{id}", h.Get)
	r.Put("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
	return r
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var in CreateInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	id, err := h.repo.Create(r.Context(), in)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{"id": id})
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	res, err := h.repo.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	_ = json.NewEncoder(w).Encode(res)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	var f ListFilter
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
	list, err := h.repo.List(r.Context(), f)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	_ = json.NewEncoder(w).Encode(list)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var in CreateInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if err := h.repo.Update(r.Context(), id, in); err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]any{"status": "ok"})
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.repo.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GET /api/v1/cost?from=YYYY-MM&to=YYYY-MM&user_id=&service_name=
func (h *Handler) GetTotalCost(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	fromS := q.Get("from")
	toS := q.Get("to")
	if fromS == "" || toS == "" {
		http.Error(w, "from/to required", 400)
		return
	}
	from, err := parseYearMonth(fromS)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	to, err := parseYearMonth(toS)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if to.Before(from) {
		http.Error(w, "to must be >= from", http.StatusBadRequest)
		return
	}

	var f ListFilter
	if v := q.Get("user_id"); v != "" {
		f.UserID = &v
	}
	if v := q.Get("service_name"); v != "" {
		f.ServiceName = &v
	}
	f.PageSize = 1000
	list, err := h.repo.List(r.Context(), f)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	res := 0
	for _, s := range list {
		start, _ := parseYearMonth(s.StartDate)
		var e *time.Time
		if s.EndDate != nil && *s.EndDate != "" {
			et, _ := parseYearMonth(*s.EndDate)
			e = &et
		}
		periodEnd := to
		months := monthsOverlapInclusive(start, e, from, &periodEnd)
		if months > 0 {
			res += months * s.Price
		}
	}
	_ = json.NewEncoder(w).Encode(map[string]any{"total": res})
}
