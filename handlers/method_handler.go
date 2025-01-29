package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Auxesia23/CatatanPengeluaran/models"
	"github.com/Auxesia23/CatatanPengeluaran/schemas"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type MethodHandler struct {
	db *gorm.DB
}

func NewMethodHandler(db *gorm.DB) *MethodHandler {
	return &MethodHandler{
		db: db,
	}
}

func (m *MethodHandler) CreateMethod(w http.ResponseWriter, r *http.Request) {
	var input schemas.MethodInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}


	var existingMethod models.Method
	err = m.db.Where("name = ?", input.Name).First(&existingMethod).Error
	if err == nil {
		http.Error(w, "Method already exists", http.StatusBadRequest)
		return
	}

	NewMethod := models.Method{
		Name: input.Name,
	}

	err = m.db.Create(&NewMethod).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(NewMethod)
}

func (m *MethodHandler) GetMethods(w http.ResponseWriter, r *http.Request) {
	var methods []models.Method
	err := m.db.Find(&methods).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	responses := []schemas.MethodResponse{}
	for i := range methods {
		responses = append(responses, schemas.MethodResponse{
			ID:   methods[i].ID,
			Name: methods[i].Name,
		})
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responses)
}

func (m *MethodHandler) GetMethod(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var method models.Method
	err := m.db.First(&method, id).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	var response schemas.MethodResponse
	response.ID = method.ID
	response.Name = method.Name
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (m *MethodHandler) UpdateMethod(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	
	var input schemas.MethodUpdateInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var method models.Method
	err = m.db.First(&method, id).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = m.db.Model(&method).Update("name", input.Name).Error
	if err != nil {
		http.Error(w, "Failed to update method", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(method)
}

func (m *MethodHandler) DeleteMethod(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var method models.Method
	err := m.db.First(&method, id).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = m.db.Delete(&method).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
