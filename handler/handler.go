package handler

import (
	// "fmt"
	"context"
	"strings"
	"errors"
	"net/http"
	"io/ioutil"
	"encoding/json"

	"axiata_test/usecase"
	"axiata_test/model"
	"axiata_test/utils"
	
	"github.com/google/uuid"
)

type Handler struct {
	Usecase usecase.IUsecase
}

type IHandler interface {
	PostsHandler(w http.ResponseWriter, r *http.Request)
	PostByTagHandler(w http.ResponseWriter, r *http.Request)
	PostsByIDHandler(w http.ResponseWriter, r *http.Request)

	RegisterHandler(w http.ResponseWriter, r *http.Request)
	LoginHandler(w http.ResponseWriter, r *http.Request)
}

func NewHandler() IHandler {
	return &Handler{Usecase: usecase.New()}
}

func validateToken(r *http.Request) (string, string, error) {
	// get headers
	authHeader := r.Header.Get("Authorization")
	token := strings.Split(authHeader, " ")
	if len(token) < 2 || token[1] == "" {
		return "", "", errors.New("Unauthorized")
	}
	// validate token
	userData, err := utils.DecodeTokenJWT(token[1])
	if err != nil {
		return "", "", errors.New("Unauthorized")
	}
	return userData.Username, userData.Role, nil
}

// handler for posts
func (h *Handler) PostsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	username, role, err := validateToken(r)
	if role == "" || err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// if token is valid, continue
	// Add role to context
    ctx = context.WithValue(ctx, utils.RoleKey, role)
	ctx = context.WithValue(ctx, utils.UsernameKey, username)
	r = r.WithContext(ctx)
	if r.Method == http.MethodGet {
		// get all posts
		data := h.Usecase.GetAllPost(ctx)
		jsonData, err := json.Marshal(data)
		if err != nil {
			http.Error(w, "Invalid response data", http.StatusBadRequest)
		}
		w.Write(jsonData)
		return
	}
	if r.Method == http.MethodPost {
		// save the post
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Unable to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var payload model.PostRequest
		if err := json.Unmarshal(body, &payload); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		resp := h.Usecase.StorePosts(ctx, payload)
		jsonData, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, "Invalid response data", http.StatusBadRequest)
		}
		w.Write(jsonData)
		return
	}
	
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	return
}

func (h *Handler) PostByTagHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	_, role, err := validateToken(r)
	if role == "" || err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// if token is valid, continue
	// Add role to context
    ctx = context.WithValue(ctx, utils.RoleKey, role)
	r = r.WithContext(ctx)
	
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	tagLabel := r.URL.Query().Get("tag")
    if tagLabel == "" {
        http.Error(w, "query tag is required...", http.StatusBadRequest)
        return
    }

	data := h.Usecase.GetPostByTag(ctx, tagLabel)
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Invalid response data", http.StatusBadRequest)
	}
	w.Write(jsonData)
	return
}

func (h *Handler) PostsByIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	username, role, err := validateToken(r)
	if role == "" || err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// if token is valid, continue
	// Add role to context
    ctx = context.WithValue(ctx, utils.RoleKey, role)
	ctx = context.WithValue(ctx, utils.UsernameKey, username)
	r = r.WithContext(ctx)

	idParam := strings.TrimPrefix(r.URL.Path, "/posts/")
	postID, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	if r.Method == http.MethodGet {
		// handle get by id
		data := h.Usecase.GetPostByID(ctx, postID)
		jsonData, err := json.Marshal(data)
		if err != nil {
			http.Error(w, "Invalid response data", http.StatusBadRequest)
		}
		w.Write(jsonData)
		return
	}
	if r.Method == http.MethodPut {
		// handle edit by id
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Unable to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var payload model.PostRequest
		if err := json.Unmarshal(body, &payload); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		resp := h.Usecase.UpdatePostByID(ctx, postID, payload)
		jsonData, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, "Invalid response data", http.StatusBadRequest)
		}
		w.Write(jsonData)
		return
	}
	if r.Method == http.MethodDelete {
		// handle delete by id
		resp := h.Usecase.DeletePostByID(ctx, postID)
		jsonData, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, "Invalid response data", http.StatusBadRequest)
		}
		w.Write(jsonData)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	return
}


// handler for accounts
func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {	
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// register the account
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var payload model.PayloadRegister
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	resp := h.Usecase.RegisterAccount(ctx, payload)
	jsonData, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Invalid response data", http.StatusBadRequest)
	}
	w.Write(jsonData)
	return
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {	
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// register the account
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var payload model.PayloadLogin
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	resp := h.Usecase.Login(ctx, payload)
	jsonData, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Invalid response data", http.StatusBadRequest)
	}
	w.Write(jsonData)
	return
}