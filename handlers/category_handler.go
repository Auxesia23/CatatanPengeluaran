package handlers

import (
	"net/http"
	"gorm.io/gorm"
	"github.com/Auxesia23/CatatanPengeluaran/schemas"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/Auxesia23/CatatanPengeluaran/models"
)
type CategoryHandler struct {
	db *gorm.DB
}	

func NewCategoryHandler(db *gorm.DB) *CategoryHandler {
	return &CategoryHandler{
		db: db,
	}
}

func (c *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var input schemas.CategoryInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var existingCategory models.Category
	err = c.db.Where("name = ?", input.Name).First(&existingCategory).Error
	if err == nil {
		http.Error(w, "Category already exists", http.StatusBadRequest)
		return
	}

	NewCategory := models.Category{
		Name: input.Name,
	}	

	err = c.db.Create(&NewCategory).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(NewCategory)
}

func(c *CategoryHandler)GetCategories(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	
	var categories []models.Category
	if name == "" {
		err := c.db.Find(&categories).Error
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		err := c.db.Where("name LIKE ?", "%"+name+"%").Find(&categories).Error
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	responses := []schemas.CategoryResponse{}
	for i := range categories {
		responses = append(responses, schemas.CategoryResponse{
			ID:   categories[i].ID,
			Name: categories[i].Name,
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responses)
}

func (c *CategoryHandler)GetCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var category models.Category
	err := c.db.First(&category, id).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	response := schemas.CategoryResponse{
		ID:   category.ID,
		Name: category.Name,
	}	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (c *CategoryHandler)UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	
	var input schemas.CategoryUpdateInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var category models.Category
	err = c.db.First(&category, id).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = c.db.Model(&category).Update("name", input.Name).Error
	if err != nil {
		http.Error(w, "Failed to update category", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(category)
}

func (c *CategoryHandler)DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var category models.Category
	err := c.db.First(&category, id).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}	

	err = c.db.Delete(&category).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}