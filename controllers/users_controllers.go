package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"rest/db"
	"rest/models"
)

func CreateUser(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-type", "application/json")
	var user models.User
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		res.WriteHeader(http.StatusBadGateway)
		return
	}
	user.Password = string(hashedPass)
	_, err = db.GetInstance().Exec("INSERT INTO rest.users (username,password) VALUES (?,?);",
		user.Username, user.Password)
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	res.WriteHeader(http.StatusCreated)
}

func GetUsers(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-type", "application/json")
	result, err := db.GetInstance().Query("SELECT * FROM rest.users;")
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
	}
	var users []models.User
	for result.Next() {
		var user models.User
		err := result.Scan(&user.Id, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			res.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		users = append(users, user)
	}
	_ = json.NewEncoder(res).Encode(users)
}

func GetUser(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-type", "application/json")
	params := mux.Vars(req)
	var user models.User
	err := db.GetInstance().
		QueryRow("SELECT * FROM rest.users WHERE username=?;", params["username"]).
		Scan(&user.Id, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		res.WriteHeader(http.StatusNoContent)
	}
	_ = json.NewEncoder(res).Encode(user)
}

func UpdateUser(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-type", "application/json")
	var reqBody models.User
	_ = json.NewDecoder(req.Body).Decode(&reqBody)
	params := mux.Vars(req)
	if reqBody.Password != "" {
		hashedPass, err := bcrypt.GenerateFromPassword([]byte(params["password"]), bcrypt.MinCost)
		if err != nil {
			res.WriteHeader(http.StatusBadGateway)
			return
		}
		reqBody.Password = string(hashedPass)
	}
	_, err := db.GetInstance().Exec("UPDATE rest.users SET username=?, password=? WHERE username=?;",
		reqBody.Username, reqBody.Password, params["username"])
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	res.WriteHeader(http.StatusOK)
}

func DeleteUser(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-type", "application/json")
	params := mux.Vars(req)
	_, err := db.GetInstance().Exec("DELETE FROM rest.users WHERE username=?", params["username"])
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	res.WriteHeader(http.StatusOK)
}
