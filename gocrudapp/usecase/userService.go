package usecase

import (
	"encoding/json"
	"gocrudapp/dto"
	"gocrudapp/model"
	"gocrudapp/repository"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type UserService struct {
	DBClient repository.UserInterface
}

func (srv UserService) CreateUser(w http.ResponseWriter, r *http.Request){
	res := dto.UserResponse{}

	// extract body from request
	var userRequest dto.UserRequest
	err := json.NewDecoder(r.Body).Decode(&userRequest)
	if err != nil {
		slog.Error(err.Error())
		res.Error = "Invalid request body"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}
	user := model.User{
		Name: 	userRequest.Name,
		Age: 	userRequest.Age,
		Country: userRequest.Country,
	}

	// create user in database
	result, err := srv.DBClient.CreateUser(user)
	if err != nil {
		slog.Error("error creating user", slog.String("error", err.Error()))
		res.Error = "Error creating user"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}


	// success
	slog.Info("user created successfully", slog.String("_id", result))
	res.Data = result
	json.NewEncoder(w).Encode(res)
	
}

func (srv UserService) GetUserByID(w http.ResponseWriter, r *http.Request){
	res := dto.UserResponse{}

	id := chi.URLParam(r, "id")
	if id == "" {
		slog.Error("id is required")
		res.Error = "ID is required"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}

	user, err := srv.DBClient.GetUserByID(id)
	if err != nil {
		slog.Error(err.Error())
		res.Error = "Error getting user"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}

	// success
	slog.Info("user successfully retrieved", slog.String("_id", user.ID.Hex()))
	res.Data = user
	json.NewEncoder(w).Encode(res)
}

func (srv UserService) GetAllUsers(w http.ResponseWriter, r *http.Request){
	res := dto.UserResponse{}

	users, err := srv.DBClient.GetAllUsers()
	if err != nil {
		slog.Error(err.Error())
		res.Error = "Error getting users"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}

	// success
	slog.Info("users successfully fetched")
	res.Data = users
	json.NewEncoder(w).Encode(res)
}

func (srv UserService) UpdateUserAgeByID(w http.ResponseWriter, r *http.Request){
	res := dto.UserResponse{}

	id := chi.URLParam(r, "id")
	if id == "" {
		slog.Error("id is required")
		res.Error = "ID is required"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}

	var userRequest dto.UserRequest
	err := json.NewDecoder(r.Body).Decode(&userRequest)
	if err != nil {
		slog.Error(err.Error())
		res.Error = "Invalid request body"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}

	// update user age in database
	updatedAge, err := srv.DBClient.UpdateUserAgeByID(id, userRequest.Age)
	if err != nil {
		slog.Error(err.Error())
		res.Error = "Error updating user age"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}

	// success
	slog.Info("user age updated successfully")
	res.Data = updatedAge
	json.NewEncoder(w).Encode(res)
}

func (srv UserService) DeleteUserByID(w http.ResponseWriter, r *http.Request){
	res := dto.UserResponse{}

	id := chi.URLParam(r, "id")
	if id == "" {
		slog.Error("id is required")
		res.Error = "ID is required"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}

	deletedCount, err := srv.DBClient.DeleteUserByID(id)
	if err != nil {
		slog.Error(err.Error())
		res.Error = "Error deleting user"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}

	// success
	slog.Info("user deleted successfully")
	res.Data = deletedCount
	json.NewEncoder(w).Encode(res)
}

func (srv UserService) DeleteAllUsers(w http.ResponseWriter, r *http.Request){
	res := dto.UserResponse{}

	deletedCount, err := srv.DBClient.DeleteAllUsers()
	if err != nil {
		slog.Error(err.Error())
		res.Error = "Error deleting users"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}

	// success
	slog.Info("all users deleted successfully")
	res.Data = deletedCount
	json.NewEncoder(w).Encode(res)
}