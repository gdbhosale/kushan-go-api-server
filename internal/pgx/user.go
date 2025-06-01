package pgx

import (
	"goat/internal"

	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

// Ensure service implements interface
var _ internal.UserService = (*UserService)(nil)

// UserService represents a PostgreSQL implementation of models.UserService.
type UserService struct {
	db *sqlx.DB
}

// NewUserService returns a new instance of UserService.
func NewUserService(db *sqlx.DB) *UserService {
	return &UserService{db: db}
}

// Retrieves a user by ID. Returns 404 if user does not exist.
func (s *UserService) FindUserByID(id int) (*internal.User, error) {
	var user internal.User

	row := s.db.QueryRowx(`SELECT id, name, email, created_at, updated_at FROM users WHERE deleted_at IS NULL AND id = $1`, id)

	if err := row.StructScan(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

// Retrieves a user by Email. Returns 404 if user does not exist.
func (s *UserService) FindUserByEmail(email string) (*internal.User, error) {
	var user internal.User

	row := s.db.QueryRowx(`SELECT id, name, email, password, created_at, updated_at FROM users WHERE deleted_at IS NULL AND email = $1`, email)

	if err := row.StructScan(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

// Retrieves a list of users by filter. Also returns total count of matching
// users which may differ from returned results if filter.Limit is specified.
func (s *UserService) FindUsers(filter internal.UserFilter) ([]*internal.User, error) {
	internal.Debug("PGX::FindUsers", "Filter:", filter.String())
	// Build WHERE clause.
	where := []string{}
	if v := filter.ID; v != 0 {
		where = append(where, "id = "+strconv.Itoa(v))
	}
	if v := filter.Name; v != "" {
		where = append(where, "name LIKE '%"+v+"%'")
	}
	if v := filter.Email; v != "" {
		where = append(where, "email LIKE '%"+v+"%'")
	}
	where = append(where, "deleted_at IS NULL")

	// Query
	query := `
		SELECT
		id,
		name,
		email,
		created_at,
		updated_at
		FROM users
		WHERE ` + strings.Join(where, " AND ") + `
		ORDER BY id ASC
		` + FormatLimitOffset(filter.Limit, filter.Offset)

	// Querying Data
	rows, err := s.db.Queryx(query)
	if err != nil {
		return nil, err
	}

	// Prepare Data
	users := make([]*internal.User, 0)
	for rows.Next() {
		var user internal.User
		if err := rows.StructScan(&user); err != nil {
			return users, err
		}
		users = append(users, &user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// Creates a new user.
// func (s *UserService) CreateUser(u internal.User) error

// Updates a user object. Returns 404 if user does not exist.
// func (s *UserService) UpdateUser(id int, upd internal.UserUpdate) (internal.User, error)

// Soft Deletes User if found. Returns 404 if user does not exist.
// func (s *UserService) DeleteUser(id int) error
