package users

import (
	"net/http"

	"concierge-be/utils"
	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"fullName"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateProfileRequest struct {
	Email    string `json:"email" binding:"omitempty,email"`
	FullName string `json:"fullName"`
	Password string `json:"password" binding:"omitempty,min=6"`
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Check if username already exists
	_, err := h.service.GetUserByUsername(req.Username)
	if err == nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "username already exists")
		return
	}

	// Check if email already exists
	_, err = h.service.GetUserByEmail(req.Email)
	if err == nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "email already exists")
		return
	}

	// Create user
	user := &User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		FullName: req.FullName,
	}

	if err := h.service.CreateUser(user); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Clear password from response
	user.Password = ""
	utils.SuccessResponse(c, gin.H{
		"user": user,
	})
}

// Login handles user login
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Find user
	user, err := h.service.GetUserByUsername(req.Username)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "invalid username or password")
		return
	}

	// Verify password
	if !h.service.VerifyPassword(user, req.Password) {
		utils.ErrorResponse(c, http.StatusUnauthorized, "invalid username or password")
		return
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to generate token")
		return
	}

	// Clear password from response
	user.Password = ""
	utils.SuccessResponse(c, gin.H{
		"token": token,
		"user":  user,
	})
}

// GetCurrentUser gets the current authenticated user
func (h *Handler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	user, err := h.service.GetUserByID(userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	// Clear password from response
	user.Password = ""
	utils.SuccessResponse(c, user)
}

// UpdateCurrentUser updates the current authenticated user
func (h *Handler) UpdateCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.GetUserByID(userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	// Update user information
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Password != "" {
		user.Password = req.Password
	}

	if err := h.service.UpdateUser(user); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Clear password from response
	user.Password = ""
	utils.SuccessResponse(c, user)
}
