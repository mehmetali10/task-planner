package handler

import (
	"net/http"

	"github.com/mehmetali10/task-planner/internal/task/service"

	"github.com/go-playground/validator"
)

type Handler interface {
	CreateTask() http.HandlerFunc
	GetTask() http.HandlerFunc
	ListTasks() http.HandlerFunc
	ScheduleAssaignments() http.HandlerFunc
	ListAssignments() http.HandlerFunc
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

func (h *handler) CreateTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		panic("unimplemented")
	}
}

func (h *handler) GetTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		panic("unimplemented")
	}
}

func (h *handler) ListTasks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		panic("unimplemented")
	}
}

func (h *handler) ScheduleAssaignments() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		panic("unimplemented")
	}
}

func (h *handler) ListAssignments() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		panic("unimplemented")
	}
}

func (h *handler) ListDevelopers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		panic("unimplemented")
	}
}
