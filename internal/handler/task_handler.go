package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/NITHISH-2006/taskflow-go/internal/model"
	"github.com/NITHISH-2006/taskflow-go/internal/repository"
	"github.com/NITHISH-2006/taskflow-go/internal/service"
	"github.com/NITHISH-2006/taskflow-go/internal/utils"
)

type TaskHandler struct {
	service *service.TaskService
}

func NewTaskHandler(service *service.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

func (h *TaskHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", h.Health)
	mux.HandleFunc("/tasks", h.Tasks)
	mux.HandleFunc("/tasks/", h.TaskByID)
}

func (h *TaskHandler) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.JSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	utils.JSONResponse(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *TaskHandler) Tasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetTasks(w, r)
	case http.MethodPost:
		h.CreateTask(w, r)
	default:
		utils.JSONError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *TaskHandler) TaskByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/tasks/")
	if id == "" {
		utils.JSONError(w, http.StatusNotFound, "task not found")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetTask(w, r, id)
	case http.MethodPut:
		h.UpdateTask(w, r, id)
	case http.MethodDelete:
		h.DeleteTask(w, r, id)
	default:
		utils.JSONError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	filter, err := parseTaskFilter(r)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	tasks, err := h.service.ListTasks(filter)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed to load tasks")
		return
	}

	utils.JSONResponse(w, http.StatusOK, tasks)
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request, id string) {
	task, err := h.service.GetTask(id)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, task)
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req model.CreateTaskRequest
	if err := decodeJSONBody(w, r, &req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	task, err := h.service.CreateTask(req)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	utils.JSONResponse(w, http.StatusCreated, task)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request, id string) {
	var req model.UpdateTaskRequest
	if err := decodeJSONBody(w, r, &req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	task, err := h.service.UpdateTask(id, req)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, task)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.service.DeleteTask(id); err != nil {
		h.handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TaskHandler) handleServiceError(w http.ResponseWriter, err error) {
	if errors.Is(err, repository.ErrTaskNotFound) {
		utils.JSONError(w, http.StatusNotFound, "task not found")
		return
	}

	if errors.Is(err, utils.ErrTitleRequired) || errors.Is(err, utils.ErrPriorityOutOfRange) || errors.Is(err, utils.ErrStatusInvalid) || errors.Is(err, utils.ErrDueDateInvalid) {
		utils.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSONError(w, http.StatusInternalServerError, "internal server error")
}

func parseTaskFilter(r *http.Request) (model.TaskFilter, error) {
	query := r.URL.Query()
	var filter model.TaskFilter

	if status := strings.TrimSpace(query.Get("status")); status != "" {
		if err := utils.ValidateStatus(status); err != nil {
			return filter, err
		}
		statusValue := model.Status(strings.ToLower(status))
		filter.Status = &statusValue
	}

	if priorityRaw := strings.TrimSpace(query.Get("priority")); priorityRaw != "" {
		priority, err := strconv.Atoi(priorityRaw)
		if err != nil {
			return filter, utils.ErrPriorityOutOfRange
		}
		if err := utils.ValidatePriority(priority); err != nil {
			return filter, err
		}
		filter.Priority = &priority
	}

	if dueBefore := strings.TrimSpace(query.Get("due_before")); dueBefore != "" {
		if err := utils.ValidateTimestamp(dueBefore); err != nil {
			return filter, err
		}
		parsed, err := utils.ParseTime(dueBefore)
		if err != nil {
			return filter, err
		}
		filter.DueBefore = &parsed
	}

	return filter, nil
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	contentType := strings.ToLower(r.Header.Get("Content-Type"))
	if !strings.HasPrefix(contentType, "application/json") {
		return errors.New("content type must be application/json")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return errors.New("unable to read request body")
	}
	defer r.Body.Close()

	if len(body) == 0 {
		return errors.New("request body cannot be empty")
	}

	if err := json.NewDecoder(strings.NewReader(string(body))).Decode(dst); err != nil {
		return errors.New("request body contains invalid JSON")
	}

	return nil
}
