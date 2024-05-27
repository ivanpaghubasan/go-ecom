package user

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
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
	if err := utils.ParseJSON(r, payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	//* Check if the user exists
	_, err := h.store.GetUserByEmail(payload.Email)
	if err == nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("user with email %s already exist", payload.Email))
		return
	}

	//* If doesn't exist, create new user
	errCreateUser := h.store.CreateUser(types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  payload.Password,
	})

	if errCreateUser != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}
