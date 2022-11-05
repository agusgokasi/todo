package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"todo/api/models"
	"todo/api/responses"
	"todo/api/utils/formaterror"

	"github.com/gorilla/mux"
)

func (server *Server) CreateTodo(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	todo := models.Todo{}
	err = json.Unmarshal(body, &todo)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	todo.Prepare()
	err = todo.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	todoCreated, err := todo.SaveTodo(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Lacation", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, todoCreated.ID))
	responses.JSON(w, http.StatusCreated, todoCreated)
}

func (server *Server) GetTodos(w http.ResponseWriter, r *http.Request) {

	todo := models.Todo{}

	todos, err := todo.FindAllTodos(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, todos)
}

func (server *Server) GetTodo(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	todo := models.Todo{}

	postReceived, err := todo.FindTodoByID(server.DB, uint32(uid))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, postReceived)
}

func (server *Server) UpdateTodo(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Check if the post id is valid
	uid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Check if the post exist
	todo := models.Todo{}
	err = server.DB.Debug().Model(models.Todo{}).Where("id = ?", uid).Take(&todo).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("todo not found"))
		return
	}

	// Read the data posted
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data
	todoUpdate := models.Todo{}
	err = json.Unmarshal(body, &todoUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	todoUpdate.Prepare()
	err = todoUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	todoUpdate.ID = todo.ID //this is important to tell the model the todo id to update, the other update field are set above

	todoUpdated, err := todoUpdate.UpdateATodo(server.DB, uint32(uid))

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, todoUpdated)
}

func (server *Server) DeleteTodo(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	todo := models.Todo{}

	_, err = todo.DeleteATodo(server.DB, uint32(uid))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", uint32(uid)))
	responses.JSON(w, http.StatusNoContent, "")
}
