package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Auxesia23/CatatanPengeluaran/middlewares"
	"github.com/Auxesia23/CatatanPengeluaran/models"
	"github.com/Auxesia23/CatatanPengeluaran/schemas"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type TransactionHandler struct {
	db *gorm.DB
}

func NewTransactionHandler(db *gorm.DB) *TransactionHandler {
	return &TransactionHandler{
		db: db,
	}
}

func (t *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middlewares.UserIdContextKey).(uint)
	if !ok {
		http.Error(w, "Invalid user claims", http.StatusUnauthorized)
		return
	}

	var input schemas.InputTransaction
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validasi apakah CategoryID ada
	var category models.Category
	if err := t.db.First(&category, input.CategoryID).Error; err != nil {
		http.Error(w, "invalid category_id", http.StatusBadRequest)
		return
	}

	// Validasi apakah MethodID ada
	var method models.Method
	if err := t.db.First(&method, input.MethodID).Error; err != nil {
		http.Error(w, "invalid method_id", http.StatusBadRequest)
		return
	}

	newTransaction := models.Transaction{
		UserID:      userID,
		CategoryID:  input.CategoryID,
		MethodID:    input.MethodID,
		Amount:      input.Amount,
		Description: input.Description,
		Date:        input.Date,
	}

	err = t.db.Create(&newTransaction).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newTransaction)
}

func (t *TransactionHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	fromDate := r.URL.Query().Get("from")
	toDate := r.URL.Query().Get("to")

	userID, ok := r.Context().Value(middlewares.UserIdContextKey).(uint)
	if !ok {
		http.Error(w, "Invalid user claims", http.StatusUnauthorized)
		return
	}

	var transactions []models.Transaction
	if fromDate == "" || toDate == "" {
		now := time.Now()
		oneMonthAgo := now.AddDate(0, -1, 0)
		err := t.db.Preload("Category").Preload("Method").Where("user_id = ? AND date BETWEEN ? AND ?", userID, oneMonthAgo, now).Or("user_id = ?", userID).Find(&transactions).Error
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		err := t.db.Preload("Category").Preload("Method").Where("user_id = ? AND date BETWEEN ? AND ?", userID, fromDate, toDate).Find(&transactions).Error
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	var response []schemas.TransactionResponse
	for _, transaction := range transactions {
		response = append(response, schemas.TransactionResponse{
			ID:          transaction.ID,
			UserID:      transaction.UserID,
			Category:    transaction.Category.Name,
			Method:      transaction.Method.Name,
			Amount:      transaction.Amount,
			Description: transaction.Description,
			Date:        transaction.Date,
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (t *TransactionHandler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID, ok := r.Context().Value(middlewares.UserIdContextKey).(uint)
	if !ok {
		http.Error(w, "Invalid user claims", http.StatusUnauthorized)
		return
	}

	var transaction models.Transaction
	err := t.db.Where("user_id = ? AND id = ?", userID, id).First(&transaction).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transaction)
}

func (t *TransactionHandler) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var input schemas.InputTransaction
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var transaction models.Transaction
	err = t.db.First(&transaction, id).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = t.db.Model(&transaction).Updates(input).Error
	if err != nil {
		http.Error(w, "Failed to update transaction", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transaction)
}

func (t *TransactionHandler) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var transaction models.Transaction
	err := t.db.First(&transaction, id).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = t.db.Delete(&transaction).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
