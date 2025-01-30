package handlers

import (
	"encoding/json"
	"net/http"
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
	userID, ok := r.Context().Value(middlewares.UserIdContextKey).(uint)
	if !ok {
		http.Error(w, "Invalid user claims", http.StatusUnauthorized)
		return
	}

	var transactions []models.Transaction
	err := t.db.Where("user_id = ?", userID).Find(&transactions).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transactions)
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
