package internal

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// Represents User who can login to system
type User struct {
	ID       uint     `db:"id" json:"id" example:"123"`                     // User's ID
	Name     string   `db:"name" json:"name" example:"John Doe"`            // User's name
	Email    string   `db:"email" json:"email" example:"johndoe@gmail.com"` // User's Email
	Password string   `db:"password" json:"password" example:"Masked"`      // User's Password
	Roles    []string `db:"roles" json:"roles" example:"['superadmin']"`    // User Roles

	// Timestamps
	CreatedAt time.Time `db:"created_at" json:"createdAt" example:"2024-05-03T15:34:26.460Z"` // User's Creation Time
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt" example:"2024-05-03T15:34:26.460Z"` // User's Updation Time
	DeletedAt time.Time `db:"deleted_at" json:"deletedAt" example:"2024-05-03T15:34:26.460Z"` // Deletion time if User is Deleted
}

// UserService represents a service for managing users.
type UserService interface {
	// Retrieves a user by ID. Returns 404 if user does not exist.
	FindUserByID(id int) (*User, error)

	// Retrieves a user by Email. Returns 404 if user does not exist.
	FindUserByEmail(email string) (*User, error)

	// Retrieves a list of users by filter. Also returns total count of matching
	// users which may differ from returned results if filter.Limit is specified.
	FindUsers(filter UserFilter) ([]*User, error)

	// Creates a new user.
	// CreateUser(u *User) error

	// Updates a user object. Returns 404 if user does not exist.
	// UpdateUser(id int, upd UserUpdate) (*User, error)

	// Soft Deletes User if found. Returns 404 if user does not exist.
	// DeleteUser(id int) error
}

// UserFilter represents a filter passed to FindUsers().
type UserFilter struct {
	// Filtering fields.
	ID    int    `json:"id" example:"123"`                  // User's ID
	Name  string `json:"name" example:"John Doe"`           // User's name
	Email string `json:"email" example:"johndoe@gmail.com"` // User's Email

	// Restrict to subset of results.
	Offset int `json:"offset" example:"0"` // Pagination Offset
	Limit  int `json:"limit" example:"20"` // Number of Records to be fetched
}

// Stringer Interface: Override Default String Method of Struct for Rectangle
func (f UserFilter) String() string {
	return fmt.Sprintf("{ ID: %d, Name: %s, Email: %s, Offset: %d, Limit: %d }", f.ID, f.Name, f.Email, f.Offset, f.Limit)
}

func NewUserFilter(r *http.Request) UserFilter {
	q := r.URL.Query()

	id, err := strconv.Atoi(q.Get("id"))
	if err != nil {
		id = 0
	}
	offset, err := strconv.Atoi(q.Get("offset"))
	if err != nil {
		offset = 0
	}
	limit, err := strconv.Atoi(q.Get("limit"))
	if err != nil {
		limit = 0
	}
	// Parse optional filter object
	return UserFilter{
		ID:     id,
		Name:   q.Get("name"),
		Email:  q.Get("email"),
		Offset: offset,
		Limit:  limit,
	}
}

// UserUpdate represents a set of fields to be updated via UpdateUser().
type UserUpdate struct {
	Name  *string `json:"name" example:"John Doe"`           // User's name
	Email *string `json:"email" example:"johndoe@gmail.com"` // User's Email
}
