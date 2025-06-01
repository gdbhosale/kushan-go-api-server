package http

import (
	"goat/internal"
	"goat/internal/http/middlewares"

	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

// Helper function for registering all auth routes.
func (s *Server) registerAuthRoutes(r *http.ServeMux) {
	sm := http.NewServeMux()

	// Module Middlewares
	stack := middlewares.CreateStack(
		middlewares.Logging,
		middlewares.RateLimiter,
		middlewares.AllowCors,
	)

	sm.HandleFunc("POST /signin", s.Signin)
	sm.HandleFunc("POST /signout", s.Signout)

	r.Handle("/api/v1/auth/", stack(http.StripPrefix("/api/v1/auth", sm)))
}

// Represents User Sign In Request who can login to system
type SigninRequest struct {
	Email    string `json:"email" example:"gdb.sci123@gmail.com"` // Signin Email Address
	Password string `json:"password" example:"12345678"`          // Signin Password
}

// Represents User Sign In Response with Access Token & User Info
type SigninResponse struct {
	ID        uint     `json:"id" example:"1"`                                                                                                                      // User's ID
	Name      string   `json:"name" example:"John Doe"`                                                                                                             // User's Name
	Email     string   `json:"email" example:"gdb.sci123@gmail.com"`                                                                                                // User's Email
	Token     string   `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTQ3NTA0NjYsImlkIjoxfQ.sX29VB_6SQH3Kjpggo89M2QqT5A6PAfMz-r1sR7v99M"` // JWT Auth Token
	ExpiresAt string   `json:"expiresAt" example:"2024-05-03T15:34:26.460Z"`                                                                                        // JWT Token Expiry Time
	Roles     []string `json:"roles" example:"['superadmin']"`                                                                                                      // User Roles
}

// Signin godoc
//
//	@Summary		Sign In API
//	@Description	Sign In API for GOATAdmin
//	@Tags			auth
//	@Accept			json
//	@Param			input	body	SigninRequest	true	"Signin Credentials"
//	@Produce		json
//	@Success		200	{object}	SigninResponse
//	@Failure		400	{object}	internal.ErrorResponse	"Invalid JSON body"
//	@Failure		401	{object}	internal.ErrorResponse	"Invalid Bearer Token"
//	@Failure		404	{object}	internal.ErrorResponse	"User not found"
//	@Failure		500	{object}	internal.ErrorResponse	"Issue with Data Parsing"
//	@Router			/api/v1/auth/signin [post]
func (s *Server) Signin(w http.ResponseWriter, r *http.Request) {

	// Parse signin request into object
	var signinRequest SigninRequest
	if err := json.NewDecoder(r.Body).Decode(&signinRequest); err != nil {
		internal.APIError(w, "Http::Signin", "Invalid JSON body", http.StatusBadRequest, err)
		return
	}
	// Fetch users from database.
	user, err := s.UserService.FindUserByEmail(signinRequest.Email)
	if err != nil {
		internal.APIError(w, "Http::Signin", "User not found", http.StatusNotFound, err)
		return
	}

	// Check the password
	// if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(signinRequest.Password)); err != nil {
	// 	internal.APIError(w, "Http::Signin", "Invalid Credentials", http.StatusUnauthorized, err)
	// 	return
	// }

	expiresAtTime := time.Now().Add(time.Hour * 24)
	accessToken, err := generateAccessToken(user.ID, expiresAtTime)
	if err != nil {
		internal.APIError(w, "Http::Signin", "Failed to generate access token", http.StatusNotFound, err)
		return
	}

	// Mask User Password
	var signinResponse SigninResponse
	signinResponse.ID = user.ID
	signinResponse.Name = user.Name
	signinResponse.Email = user.Email
	signinResponse.Token = accessToken
	signinResponse.ExpiresAt = expiresAtTime.UTC().Format("2006-01-02T15:04:05.000Z")
	signinResponse.Roles = []string{"superadmin"}

	// Response
	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(signinResponse); err != nil {
		internal.APIError(w, "Http::Signin", "Issue with Data Parsing", http.StatusInternalServerError, err)
		return
	}
}

// Represents User Signout Response
type SignoutResponse struct {
	Message string `json:"message" example:"User Signed Out"` // Signout Message
}

// Signout godoc
//
//	@Summary		Signout API
//	@Description	Signout API for GOATAdmin
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	SignoutResponse
//	@Failure		500	{object}	internal.ErrorResponse	"Issue with Data Parsing"
//	@Router			/api/v1/auth/signout [post]
func (s *Server) Signout(w http.ResponseWriter, r *http.Request) {

	var signoutResponse SignoutResponse
	signoutResponse.Message = "User Signed Out"

	authorization := r.Header.Get("Authorization")

	// Check that the header begins with a prefix of Bearer
	if !strings.HasPrefix(authorization, "Bearer ") {
		internal.APIError(w, "Middleware::IsAuthenticated", "Bearer Token not found", http.StatusUnauthorized, nil)
		return
	}

	// Pull out the token
	tokenString := strings.TrimPrefix(authorization, "Bearer ")

	// Add token to blacklist
	s.BlacklistedToken[tokenString] = true

	// Response
	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(signoutResponse); err != nil {
		internal.APIError(w, "Http::Signout", "Issue with Data Parsing", http.StatusInternalServerError, err)
		return
	}
}

// Generate JWT Access Token with userId and token expiry time
func generateAccessToken(userId uint, expiresAtTime time.Time) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	// Set token claims
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = userId
	claims["exp"] = expiresAtTime.Unix()

	// Generate encoded token
	return token.SignedString([]byte("PRIVATE_KEY"))
}

// IsAuthenticated Middleware for authorizing the API Requests based on Bearer JWT Token
func (s *Server) IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")

		// Check that the header begins with a prefix of Bearer
		if !strings.HasPrefix(authorization, "Bearer ") {
			internal.APIError(w, "Middleware::IsAuthenticated", "Bearer Token not found", http.StatusUnauthorized, nil)
			return
		}

		// Pull out the token
		tokenString := strings.TrimPrefix(authorization, "Bearer ")

		// Check if Token is Blacklisted / Invalidated / SignedOut
		if s.BlacklistedToken[tokenString] {
			internal.APIError(w, "Http::IsAuthenticated", "Access Token Expired", http.StatusUnauthorized, nil)
			return
		}

		if tokenString == "" {
			internal.APIError(w, "Http::IsAuthenticated", "Access Token not found", http.StatusUnauthorized, nil)
			return
		}

		// Verify and parse the access token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Verify the signing method and return the secret key
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Http::IsAuthenticated - invalid signing method")
			}

			return []byte("PRIVATE_KEY"), nil
		})
		if err != nil {
			internal.APIError(w, "Http::IsAuthenticated", "Error in parsing access token", http.StatusUnauthorized, err)
			return
		}

		// Check if the token is valid and not expired
		if !token.Valid {
			internal.APIError(w, "Http::IsAuthenticated", "Invalid access token", http.StatusUnauthorized, nil)
			return
		}

		// Extract the user Id from the token claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			internal.APIError(w, "Http::IsAuthenticated", "Cannot extract access token", http.StatusUnauthorized, nil)
			return
		}

		userId := int(claims["id"].(float64))

		user, err := s.UserService.FindUserByID(userId)
		if err != nil {
			internal.APIError(w, "Http::IsAuthenticated", "Token User not found", http.StatusNotFound, err)
			return
		}

		r = r.WithContext(internal.NewContextWithUser(r.Context(), user))

		next.ServeHTTP(w, r)
	})
}
