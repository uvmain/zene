package handlers

import (
	"encoding/json"
	"net/http"

	"zene/core/logger"
	"zene/core/types"
)

type StandardResponse struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func (r *StandardResponse) SetError(msg string) {
	r.Status = "error"
	r.Error = msg
}

func (r *StandardResponse) SetSuccess() {
	r.Status = "success"
}

type IdResponse struct {
	Id int64 `json:"Id"`
	StandardResponse
}

func (r *IdResponse) SetError(msg string) {
	r.Status = "error"
	r.Error = msg
}

func (r *IdResponse) SetSuccess() {
	r.Status = "success"
}

type UserResponse struct {
	*types.User
	StandardResponse
}

func toUserPointer(user types.User) *types.User {
	return &user
}

func (r *UserResponse) SetError(msg string) {
	r.Status = "error"
	r.Error = msg
}

func (r *UserResponse) SetSuccess() {
	r.Status = "success"
}

type UsersResponse struct {
	Users []*types.User `json:"users"`
	StandardResponse
}

func toUsersPointers(users []types.User) []*types.User {
	pointers := make([]*types.User, len(users))
	for i := range users {
		pointers[i] = &users[i]
	}
	return pointers
}

func (r *UsersResponse) SetError(msg string) {
	r.Status = "error"
	r.Error = msg
}

func (r *UsersResponse) SetSuccess() {
	r.Status = "success"
}

type PlaycountsResponse struct {
	Playcounts []*types.Playcount `json:"playcounts"`
	StandardResponse
}

func (r *PlaycountsResponse) SetError(msg string) {
	r.Status = "error"
	r.Error = msg
}

func (r *PlaycountsResponse) SetSuccess() {
	r.Status = "success"
}

type ErrorResponse interface {
	SetError(string)
}

type SuccessResponse interface {
	SetSuccess()
}

func handleErrorResponse[T ErrorResponse](w http.ResponseWriter, response T, message string, err error, statusCode int) {
	logger.Printf("%s: %v", message, err)
	w.Header().Set("Content-Type", "application/json")
	response.SetError(message)
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func handleSuccessResponse[T SuccessResponse](w http.ResponseWriter, response T) {
	w.Header().Set("Content-Type", "application/json")
	response.SetSuccess()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
