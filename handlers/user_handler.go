package handlers

import (
	"encoding/json"
	"net/http"

	"gorm.io/gorm"
	"github.com/Auxesia23/CatatanPengeluaran/schemas"
	"github.com/Auxesia23/CatatanPengeluaran/models"
	"github.com/Auxesia23/CatatanPengeluaran/utils"
	"github.com/Auxesia23/CatatanPengeluaran/middlewares"
)

type Userhandler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *Userhandler {
	return &Userhandler{
		db: db,
	}
}

func (u *Userhandler) Register(w http.ResponseWriter, r *http.Request) {
	var input schemas.RegisterInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var existingUser models.User
	err = u.db.Where("username = ?", input.Username).First(&existingUser).Error
	if err == nil {
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}

	if len(input.Password) < 6 {
		http.Error(w, "Password must be at least 8 characters long", http.StatusBadRequest)
		return
	}

	NewUser := models.User{
		Username: input.Username,
		Password: utils.HashPassword(input.Password),
	}

	err = u.db.Create(&NewUser).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := schemas.UserResponse{
		ID:   NewUser.ID,
		Username: NewUser.Username,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (u *Userhandler) Login(w http.ResponseWriter, r *http.Request) {
	var input schemas.LoginInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user models.User
	err = u.db.Where("username = ?", input.Username).First(&user).Error
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}	

	if !utils.CheckPasswordHash(input.Password, user.Password) {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateToken(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

func (u *Userhandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middlewares.UserIdContextKey).(uint)
	if !ok {
		http.Error(w, "Invalid user claims", http.StatusUnauthorized)
		return
	}

	var user models.User
	err := u.db.First(&user, userID).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	var response schemas.UserResponse
	response.ID = user.ID
	response.Username = user.Username
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}