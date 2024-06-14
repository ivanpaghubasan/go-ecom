package user

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/ivanpaghubasan/go-ecom/service/auth"
	"github.com/ivanpaghubasan/go-ecom/types"
	"github.com/ivanpaghubasan/go-ecom/utils"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRouters(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	//* Get JSON payload
	var payload types.RegisterUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	fmt.Println(payload)
	//* Validate payload
	if err := utils.Validate.Struct(&payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	//* Check if the user exists
	_, err := h.store.GetUserByEmail(payload.Email)
	if err == nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("user with email %s already exist", payload.Email))
		return
	}

	hashPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	//* If doesn't exist, create new user
	errCreateUser := h.store.CreateUser(types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashPassword,
	})

	if errCreateUser != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}
