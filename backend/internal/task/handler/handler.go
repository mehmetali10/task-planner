package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mehmetali10/task-planner/internal/pkg/payload"
	"github.com/mehmetali10/task-planner/internal/task/service"

	"github.com/go-playground/validator"
)

type Handler interface {
	CreateTask() http.HandlerFunc
	ListTasks() http.HandlerFunc
	ScheduleAssaignments() http.HandlerFunc
	ListDevelopers() http.HandlerFunc
}

type handler struct {
	service service.Service
}

func NewHandler(service service.Service) Handler {
	return &handler{
		service: service,
	}
}

func validateRequest(req interface{}) error {
	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		return err
	}
	return nil
}

// CreateTaskHandler godoc
// @Summary Create a task
// @Description Create a new task
// @Tags task
// @Accept json
// @Produce json
// @Param request body payload.CreateTaskRequest true "Create Request"
// @Success 200 {object} payload.CreateTaskResponse "Successfully created task"
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Internal server error"
// @Router /task [post]
func (h *handler) CreateTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req payload.CreateTaskRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Validate request
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
	}
}

// ListTasksHandler godoc
// @Summary List tasks
// @Description Retrieve a list of tasks
// @Tags task
// @Accept json
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} payload.ListTasksResponse "List of tasks"
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Internal server error"
// @Router /tasks [get]
func (h *handler) ListTasks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit := r.URL.Query().Get("limit")
		offset := r.URL.Query().Get("offset")

		req := payload.ListTasksRequest{
			Limit:  strToInt(limit),
			Offset: strToInt(offset),
		}

		// Validate request
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
	}
}

// ScheduleAssignmentsHandler godoc
// @Summary Schedule assignments
// @Description Automatically schedule assignments for tasks
// @Tags task
// @Accept json
// @Produce json
// @Success 200 {object} payload.ScheduleAssignmentResponse "Scheduled assignments"
// @Failure 500 {string} string "Internal server error"
// @Router /tasks/schedule [post]
func (h *handler) ScheduleAssaignments() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := h.service.ScheduleAssaignments(r.Context(), payload.ScheduleAssignmentRequest{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(req)
	}
}

// ListDevelopersHandler godoc
// @Summary List developers
// @Description Retrieve a list of developers
// @Tags developer
// @Accept json
// @Produce json
// @Success 200 {object} payload.ListDevelopersResponse "List of developers"
// @Failure 500 {string} string "Internal server error"
// @Router /developers [get]
func (h *handler) ListDevelopers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := h.service.ListDevelopers(r.Context(), payload.ListDevelopersRequest{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(req)
	}
}

func strToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}
