package http

import (
	"go-api/internal"
	"go-api/internal/http/middlewares"

	"encoding/json"
	"net/http"
	"strconv"
)

// Helper function for registering all user routes.
func (s *Server) registerUserRoutes(r *http.ServeMux) {
	sm := http.NewServeMux()

	// Module Middlewares
	stack := middlewares.CreateStack(
		middlewares.Logging,
		middlewares.RateLimiter,
		middlewares.AllowCors,
		s.IsAuthenticated,
	)

	sm.HandleFunc("GET /", s.UserAll)
	sm.HandleFunc("GET /{id}", s.UserFindByID)
	sm.HandleFunc("POST /", s.UserCreate)
	sm.HandleFunc("PUT /{id}", s.UserUpdateByID)
	sm.HandleFunc("DELETE /{id}", s.UserDeleteByID)

	r.Handle("/api/v1/users/", stack(http.StripPrefix("/api/v1/users", sm)))
}

// UserAll godoc
//
//	@Summary		Fetch All Users
//	@Description	Fetch All Users
//	@Tags			users
//	@Accept			json
//	@Param			id		query	integer	false	"Filter by User ID"
//	@Param			name	query	string	false	"Filter by User Name"
//	@Param			email	query	string	false	"Filter by User Email"
//	@Param			offset	query	integer	false	"Pagination Offset"	default(0)
//	@Param			limit	query	integer	false	"Pagination Limit"	default(20)
//	@Produce		json
//	@Success		200	{object}	[]internal.User
//	@Failure		400	{object}	internal.ErrorResponse	"Invalid JSON body"
//	@Failure		401	{object}	internal.ErrorResponse	"Invalid Bearer Token"
//	@Failure		404	{object}	internal.ErrorResponse	"Couldn't find users"
//	@Failure		417	{object}	internal.ErrorResponse	"Issue with Data Parsing"
//	@Failure		500	{object}	internal.ErrorResponse	"Server error"
//	@Router			/api/v1/users [get]
//	@Security		Bearer
func (s *Server) UserAll(w http.ResponseWriter, r *http.Request) {
	// user := internal.UserFromContext(r.Context())

	// Parse optional filter object
	filter := internal.NewUserFilter(r)

	// Fetch users from database.
	users, err := s.UserService.FindUsers(filter)
	if err != nil {
		internal.APIError(w, "Http::UserAll", "Couldn't find users", http.StatusNotFound, err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		internal.APIError(w, "Http::UserAll", "Issue with Data Parsing", http.StatusExpectationFailed, err)
		return
	}
}

// UserFindByID godoc
//
//	@Summary		Fetch User by ID
//	@Description	Fetch User by ID
//	@Tags			users
//	@Accept			json
//	@Param			id	path	integer	true	"User ID" default(1)
//	@Produce		json
//	@Success		200	{object}	[]internal.User
//	@Failure		400	{object}	internal.ErrorResponse	"Invalid JSON body"
//	@Failure		401	{object}	internal.ErrorResponse	"Invalid Bearer Token"
//	@Failure		404	{object}	internal.ErrorResponse	"User not found"
//	@Failure		417	{object}	internal.ErrorResponse	"Issue with Data Parsing"
//	@Failure		500	{object}	internal.ErrorResponse	"Server error"
//	@Router			/api/v1/users/{id} [get]
//	@Security		Bearer
func (s *Server) UserFindByID(w http.ResponseWriter, r *http.Request) {
	// Check if valid `id` present in request
	if idStr, idErr := strconv.Atoi(r.PathValue("id")); idErr == nil {
		if user, err := s.UserService.FindUserByID(idStr); err == nil {

			w.Header().Set("Content-type", "application/json")
			if err := json.NewEncoder(w).Encode(user); err != nil {
				internal.APIError(w, "Http::UserFindByID", "Issue with Data Parsing", http.StatusExpectationFailed, err)
				return
			}
			return
		} else {
			internal.APIError(w, "Http::UserFindByID", "User not found", http.StatusNotFound, err)
			return
		}
	} else {
		internal.APIError(w, "Http::UserFindByID", "Invalid User ID", http.StatusNotFound, idErr)
		return
	}
}

func (s *Server) UserCreate(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created!"))
}

func (s *Server) UserUpdateByID(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("User updated!"))
}

func (s *Server) UserDeleteByID(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("User deleted!"))
}

func (s *Server) UserPatchByID(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("User deleted!"))
}

func (s *Server) UserOptions(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("User options!"))
}
