package controllers

import "todo/api/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	//Todos routes
	s.Router.HandleFunc("/todo", middlewares.SetMiddlewareJSON(s.CreateTodo)).Methods("POST")
	s.Router.HandleFunc("/todo", middlewares.SetMiddlewareJSON(s.GetTodos)).Methods("GET")
	s.Router.HandleFunc("/todo/{id}", middlewares.SetMiddlewareJSON(s.GetTodo)).Methods("GET")
	s.Router.HandleFunc("/todo/{id}", middlewares.SetMiddlewareJSON(s.UpdateTodo)).Methods("PUT")
	s.Router.HandleFunc("/todo/{id}", middlewares.SetMiddlewareJSON(s.DeleteTodo)).Methods("DELETE")

	//TodoList routes
	s.Router.HandleFunc("/todolist", middlewares.SetMiddlewareJSON(s.CreateTodoList)).Methods("POST")
	s.Router.HandleFunc("/todolist", middlewares.SetMiddlewareJSON(s.GetTodoLists)).Methods("GET")
	s.Router.HandleFunc("/todolist/{id}", middlewares.SetMiddlewareJSON(s.GetTodoList)).Methods("GET")
	s.Router.HandleFunc("/todolist/{id}", middlewares.SetMiddlewareJSON(s.UpdateTodoList)).Methods("PUT")
	s.Router.HandleFunc("/todolist/{id}", middlewares.SetMiddlewareJSON(s.DeleteTodoList)).Methods("DELETE")
}
