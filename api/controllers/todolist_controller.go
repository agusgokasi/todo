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

func (server *Server) CreateTodoList(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	todolist := models.TodoList{}
	err = json.Unmarshal(body, &todolist)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	todolist.Prepare()
	err = todolist.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	todoCreated, err := todolist.SaveTodoList(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Lacation", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, todoCreated.ID))
	responses.JSON(w, http.StatusCreated, todoCreated)
}

func (server *Server) GetTodoLists(w http.ResponseWriter, r *http.Request) {

	todolist := models.TodoList{}

	todos, err := todolist.FindAllTodoLists(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, todos)
}

func (server *Server) GetTodoList(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	todolist := models.TodoList{}

	todolistReceived, err := todolist.FindTodoListByID(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, todolistReceived)
}

func (server *Server) UpdateTodoList(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Check if the post id is valid
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Check if the post exist
	todolist := models.TodoList{}
	err = server.DB.Debug().Model(models.TodoList{}).Where("id = ?", pid).Take(&todolist).Error
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
	todoListUpdate := models.TodoList{}
	err = json.Unmarshal(body, &todoListUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	todoListUpdate.Prepare()
	err = todoListUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	todoListUpdate.ID = todolist.ID //this is important to tell the model the todo id to update, the other update field are set above

	todoListUpdated, err := todoListUpdate.UpdateATodoList(server.DB, pid)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, todoListUpdated)
}

func (server *Server) DeleteTodoList(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	todolist := models.TodoList{}

	_, err = todolist.DeleteATodoList(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	responses.JSON(w, http.StatusNoContent, "")
}
