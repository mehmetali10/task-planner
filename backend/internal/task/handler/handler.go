package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/mehmetali10/task-planner/internal/pkg/payload"
	"github.com/mehmetali10/task-planner/internal/task/service"

	"github.com/go-playground/validator"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Handler interface {
	CreateTask() http.HandlerFunc
	ListTasks() http.HandlerFunc
	ScheduleAssaignments() http.HandlerFunc
	ListDevelopers() http.HandlerFunc
	Metrics() http.HandlerFunc
}

type handler struct {
	service service.Service
}

var (
	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
)

func init() {
	prometheus.MustRegister(requestCount, requestDuration)
}

func NewHandler(service service.Service) Handler {
	return &handler{
		service: service,
	}
}

func metricMiddleware(next http.HandlerFunc, endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rec, r)

		duration := time.Since(start).Seconds()
		requestCount.WithLabelValues(r.Method, endpoint, strconv.Itoa(rec.statusCode)).Inc()
		requestDuration.WithLabelValues(r.Method, endpoint).Observe(duration)
	}
}

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

func validateRequest(req interface{}) error {
	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		return err
	}
	return nil
}

func (h *handler) CreateTask() http.HandlerFunc {
	return metricMiddleware(func(w http.ResponseWriter, r *http.Request) {
		var req payload.CreateTaskRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := validateRequest(req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp, err := h.service.CreateTask(r.Context(), req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}, "/task")
}

func (h *handler) ListTasks() http.HandlerFunc {
	return metricMiddleware(func(w http.ResponseWriter, r *http.Request) {
		limit := r.URL.Query().Get("limit")
		offset := r.URL.Query().Get("offset")

		req := payload.ListTasksRequest{
			Limit:  strToInt(limit),
			Offset: strToInt(offset),
		}

		if err := validateRequest(req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp, err := h.service.ListTasks(r.Context(), req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}, "/tasks")
}

func (h *handler) ScheduleAssaignments() http.HandlerFunc {
	return metricMiddleware(func(w http.ResponseWriter, r *http.Request) {
		req, err := h.service.ScheduleAssaignments(r.Context(), payload.ScheduleAssignmentRequest{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(req)
	}, "/tasks/schedule")
}

func (h *handler) ListDevelopers() http.HandlerFunc {
	return metricMiddleware(func(w http.ResponseWriter, r *http.Request) {
		req, err := h.service.ListDevelopers(r.Context(), payload.ListDevelopersRequest{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(req)
	}, "/developers")
}

func (h *handler) Metrics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		promhttp.Handler().ServeHTTP(w, r)
	}
}

func strToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}
