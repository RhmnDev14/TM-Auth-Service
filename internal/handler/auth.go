package handler

import (
	"auth-service/internal/dto"
	"auth-service/internal/helper"
	"auth-service/internal/usecase"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type AuthHandler struct {
	authUc usecase.AuthUc
	rg     *http.ServeMux
}

func NewAuthHandler(authUc usecase.AuthUc, rg *http.ServeMux) *AuthHandler {
	return &AuthHandler{authUc: authUc, rg: rg}
}

func (h *AuthHandler) SetupRoutes() {
	h.rg.HandleFunc(helper.Register, h.Register)
	h.rg.HandleFunc(helper.Login, h.Login)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔘 API: %s %s", r.Method, r.URL.Path)

	if r.Method != http.MethodPost {
		helper.WriteJSON(w, http.StatusMethodNotAllowed, helper.MethodNotAllowed)
		return
	}

	var reqBody dto.Register
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, fmt.Sprintf(helper.InvalidJson, err))
		return
	}

	if err := h.authUc.Register(r.Context(), reqBody); err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, dto.Response{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf(helper.GagalRegister, err),
		})
		return
	}

	helper.WriteJSON(w, http.StatusCreated, dto.Response{
		Status:  http.StatusCreated,
		Message: helper.Succes,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔘 API: %s %s", r.Method, r.URL.Path)

	if r.Method != http.MethodPost {
		helper.WriteJSON(w, http.StatusMethodNotAllowed, helper.MethodNotAllowed)
		return
	}

	var reqBody dto.Login
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, fmt.Sprintf(helper.InvalidJson, err))
		return
	}

	result, err := h.authUc.Login(r.Context(), reqBody)
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, dto.Response{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf(helper.GagalLogin, err),
		})
		return
	}

	helper.WriteJSON(w, http.StatusOK, dto.Response{
		Status:  http.StatusOK,
		Message: helper.Succes,
		Data:    result,
	})
}
