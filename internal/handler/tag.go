package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/yabu1121/expense-api/internal/model"
)

type TagStore interface {
	ListTags() ([]model.Tag, error)
	CreateTag(tag model.Tag) (*model.Tag, error)
}

type CreateTagRequest struct {
	Name string `json:"name"`
}
type TagHandler struct {
	store TagStore
}

func NewTagHandler(store TagStore) *TagHandler {
	return &TagHandler{
		store: store,
	}
}

func (h *TagHandler) ListTags(w http.ResponseWriter, r *http.Request) {
	tags, err := h.store.ListTags()
	if err != nil {
		http.Error(
			w,
			"failed to list tags",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(tags); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}
func (h *TagHandler) CreateTag(w http.ResponseWriter, r *http.Request) {
	var tagRequest CreateTagRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&tagRequest); err != nil {
		http.Error(
			w,
			"failed to decode request body",
			http.StatusBadRequest,
		)
		return
	}

	tag := model.Tag{
		Name: tagRequest.Name,
	}
	tag.Normalize()
	if err := tag.Validate(); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	createdTag, err := h.store.CreateTag(tag)
	if err != nil {
		if errors.Is(err, model.ErrTagAlreadyExists) {
			http.Error(
				w,
				err.Error(),
				http.StatusConflict,
			)
			return
		}
		http.Error(
			w,
			"failed to create tag",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(createdTag); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}
