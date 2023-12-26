package users

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"encore.dev/beta/auth"
	"encore.dev/beta/errs"
	"github.com/go-playground/validator/v10"

	"encore.app/pkg/middleware"
	"encore.app/pkg/pagination"
	"encore.app/users/store"
)

// Signup is a function that handles the signup process.
//
//	@route POST /signup
//	@param ctx - context.Context
//	@param payload
//	@return response
//	@return error
//
// encore:api public method=POST path=/signup
func Signup(ctx context.Context, payload *store.SignupPayload) (*store.Response, error) {
	// validate user details
	if err := validator.New().Struct(payload); err != nil {
		return &store.Response{}, err
	}

	// create a new user
	user, err := Create(ctx, payload)
	if err != nil {
		return nil, &errs.Error{
			Code:    errs.Internal,
			Message: "authentication failed: could not signup user",
		}
	}

	// generate tokens
	token, err := middleware.GetToken(&middleware.User{
		Id:       user.Id,
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Phone:    user.Phone,
		Roles:    user.Roles,
	})
	if err != nil {
		// return &store.Response{}, errors.New("authentication failed: unable to generate token")
		return &store.Response{}, &errs.Error{
			Code:    errs.Unauthenticated,
			Message: "authentication failed: unable to generate token",
		}
	}

	// return the response
	return &store.Response{
		Message: "Signup successful",
		Token:   token,
		Payload: &store.UserResponse{
			Id:        user.Id,
			Name:      user.Name,
			Username:  user.Username,
			Email:     user.Email,
			Phone:     user.Phone,
			Roles:     user.Roles,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}, nil
}

// Login - Login is a function that handles the login process for a user.
//
//	@route POST /login
//	@param w http.ResponseWriter
//	@param req *http.Request
//
// encore:api public raw method=POST path=/login
func Login(w http.ResponseWriter, req *http.Request) {
	// Set headers for general response
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	// Get the user details from the request
	email, password, ok := req.BasicAuth()
	if !ok {
		writeJSONErrorResponse(w, "authentication failed: invalid credentials", http.StatusUnauthorized)
		return
	}

	// Get the user
	user, err := Get(req.Context(), email)
	if err != nil {
		writeJSONErrorResponse(w, "authentication failed: invalid credentials", http.StatusUnauthorized)
		return
	}

	// Check if the password is correct
	isCorrect, err := middleware.ComparePasswords(user.Password, password)
	if err != nil || !isCorrect {
		writeJSONErrorResponse(w, "authentication failed: invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate tokens
	token, err := middleware.GetToken(&middleware.User{
		Id:       user.Id,
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Phone:    user.Phone,
		Roles:    user.Roles,
	})
	if err != nil {
		writeJSONErrorResponse(w, "authentication failed: unable to generate token", http.StatusInternalServerError)
		return
	}

	// Create the response object
	response := &store.Response{
		Message: "Login successful",
		Token:   token,
		Payload: &store.UserResponse{
			Id:        user.Id,
			Name:      user.Name,
			Username:  user.Username,
			Email:     user.Email,
			Phone:     user.Phone,
			Roles:     user.Roles,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}

	// Convert the response to JSON
	responseJSON, err := json.Marshal(response)
	if err != nil {
		writeJSONErrorResponse(w, "authentication failed: unable to generate token", http.StatusInternalServerError)
		return
	}

	// Write the response
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(responseJSON); err != nil {
		writeJSONErrorResponse(w, "authentication failed: unable to write response", http.StatusInternalServerError)
		return
	}
}

// writeJSONErrorResponse writes the specified error message as a JSON response with the provided status code.
func writeJSONErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	response := map[string]interface{}{
		"message": message,
		"code":    statusCode,
		"token":   "",
		"payload": "",
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		// If unable to marshal the error response, fallback to a generic error message
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error": "An internal server error occurred"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write(responseJSON)
}

// Create - Create a new user
//
//	@param ctx - context.Context
//	@param payload
//	@return user
//	@return error
//
// encore:api private method=POST path=/users/create
func Create(ctx context.Context, payload *store.SignupPayload) (*store.User, error) {
	// create user
	user, err := store.Create(ctx, payload)
	if err != nil {
		return nil, &errs.Error{
			Code:    errs.Internal,
			Message: "authentication failed: could not signup user",
		}
	}

	return user, nil
}

// QueryAll - Get all users
//
//	@param ctx - context.Context
//	@return users
//	@return error
//
// encore:api public method=GET path=/users
func QueryAll(ctx context.Context, options *pagination.Options) (*store.PaginatedUsersResponse, error) {
	// check if user is admin or superadmin
	claims, err := middleware.GetVerifiedClaims(ctx, "")
	if err != nil {
		return &store.PaginatedUsersResponse{}, err
	}

	// check for the roles
	if !claims.HasRole(middleware.RoleAdmin) {
		return &store.PaginatedUsersResponse{}, fmt.Errorf("unauthorized: you are not authorized to perform this action")
	}

	// query users
	users, err := store.GetAll(ctx, options)
	if err != nil {
		return &store.PaginatedUsersResponse{}, fmt.Errorf("querying users: %w", err)
	}

	// return users, nil
	return users, nil
}

// Get - Get a user
//
//	@param ctx - context.Context
//	@param id
//	@return user
//	@return error
//
// encore:api auth method=GET path=/users/:id
func Get(ctx context.Context, id string) (*store.User, error) {
	// check for claims
	claims, err := middleware.GetVerifiedClaims(ctx, "")
	if err != nil {
		return &store.User{}, &errs.Error{
			Code:    errs.Unauthenticated,
			Message: "unauthorized: unable to verify user details",
		}
	}

	// check for the roles
	if !claims.HasRole(middleware.RoleSuperAdmin, middleware.RoleAdmin) || claims.Subject.Id != id {
		// return &store.User{}, fmt.Errorf("unauthorized: you are not authorized to perform this action")
		return &store.User{}, &errs.Error{
			Code:    errs.Unauthenticated,
			Message: "unauthorized: you are not authorized to perform this action",
		}
	}

	// get user
	user, err := store.GetWithID(ctx, id)
	if err != nil {
		return nil, &errs.Error{
			Code:    errs.NotFound,
			Message: "unauthorized: you are not authorized to perform this action",
		}
	}

	// return user
	return user, nil
}

// Delete - Delete a user
//
//	@param ctx - context.Context
//	@param id
//	@return error
//
// encore:api auth method=DELETE path=/users/:id
func Delete(ctx context.Context, id string) error {
	// check for claims
	claims, err := middleware.GetVerifiedClaims(ctx, "")
	if err != nil {
		return err
	}

	// check for the roles
	if !claims.HasRole(middleware.RoleSuperAdmin) || claims.Subject.Id != id {
		return fmt.Errorf("unauthorized: you are not authorized to perform this action")
	}

	// delete user
	if err := store.Delete(ctx, id); err != nil {
		return err
	}

	// return nil on a delete event
	return nil
}

// Update - Updates a user
//
//	@param ctx - context.Context
//	@param payload
//	@return user
//	@return error
//
// encore:api auth method=PATCH path=/users/update/:id
func Update(ctx context.Context, id string, payload store.UpdatePayload) (*store.UserUpdateResponse, error) {
	// check if the user matches the authenticated user
	claims, err := middleware.GetVerifiedClaims(ctx, "")
	if err != nil {
		return &store.UserUpdateResponse{}, err
	}

	// check if the user is authorized to perform this action
	if claims.Subject.Id != id {
		return &store.UserUpdateResponse{}, fmt.Errorf("unauthorized: you are not authorized to perform this action")
	}

	// validate user details
	if err := validator.New().Struct(payload); err != nil {
		return &store.UserUpdateResponse{}, err
	}

	// update user
	if err := store.Update(ctx, id, payload); err != nil {
		return &store.UserUpdateResponse{}, err
	}

	return &store.UserUpdateResponse{
		Message: fmt.Sprintf("user with id %s updated successfully", id),
	}, nil
}

// UpdateRole - Updates a user's role
//
//	@param ctx - context.Context
//	@param id
//	@param role
//	@return user
//	@return error
//
// encore:api auth method=PATCH path=/users/update/:id/toggle-admin
func UpdateRole(ctx context.Context, id string) error {
	// check if user is admin or superadmin
	claims, err := middleware.GetVerifiedClaims(ctx, "")
	if err != nil {
		return err
	}

	// check for the roles
	if !claims.HasRole(middleware.RoleSuperAdmin, middleware.RoleAdmin) {
		return fmt.Errorf("unauthorized: you are not authorized to perform this action")
	}

	// update user
	if err := store.UpdateRole(ctx, id); err != nil {
		return err
	}

	// return user
	return nil
}

// Logout - Logout is a function that handles the logout process for a user.
//
//	@route POST /logout
//	@param ctx - context.Context
//	@return response
//	@return error
//
// encore:api public method=POST path=/logout
func Logout(_ context.Context) (*store.Response, error) {

	return &store.Response{
		Message: "Logout successful",
		Payload: nil,
		Token:   "",
	}, nil
}

// Auth - Auth is a function that handles the authentication process for a user.
//
//	@route POST /auth
//	@param ctx - context.Context
//	@param token - string
//	@return response
//	@return error
//
// encore:authhandler
func Auth(_ context.Context, token string) (auth.UID, *middleware.DataI, error) {
	// check for empty token
	if len(strings.TrimSpace(token)) < 1 {
		return "", &middleware.DataI{}, errors.New("authentication failed: token is empty")
	}

	// validate token
	claims, err := middleware.ValidateToken(token)
	if err != nil {
		return "", &middleware.DataI{}, errors.New("authentication failed: invalid token")
	}

	return auth.UID(claims.User.Id), &middleware.DataI{
		Subject: claims.User,
		Roles:   claims.User.Roles,
	}, nil
}
