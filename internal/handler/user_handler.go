package handler

import (
	"strings"
	"time"
	"wallet/internal/usecase"

	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(uu usecase.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: uu}
}

// UserResponse defines the user data returned by the API.
type UserResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	Message   string    `json:"message,omitempty"`
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	DNI      string `json:"dni"`
}

// @Summary Create a new user
// @Description Creates a new user and an associated empty wallet.
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "User to create"
// @Success 201 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users [post]
func (h *UserHandler) CreateUser(c fiber.Ctx) error {
	// 1. Parse request body
	var req CreateUserRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse request"})
	}

	// 2. Call the use case
	user, err := h.userUsecase.Create(c.Context(), req.Username, req.Name, req.DNI)
	if err != nil {
		// 3. Map domain errors to HTTP errors
		if strings.Contains(err.Error(), "username already exists") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		}
		// For any other unexpected error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	response := UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		Message:   "User created successfully with an empty wallet",
	}

	// 4. Return success response
	return c.Status(fiber.StatusCreated).JSON(response)
}

// ErrorResponse defines the structure for API error responses.
type ErrorResponse struct {
	Error string `json:"error"`
}
