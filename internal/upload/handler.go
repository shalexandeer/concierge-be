package upload

import (
	"concierge-be/utils"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

// UploadImage handles image upload
func (h *Handler) UploadImage(c *gin.Context) {
	// Parse multipart form with 10MB max memory
	err := c.Request.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to parse multipart form")
		return
	}

	// Get the file from form data
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "No image file provided")
		return
	}
	defer file.Close()

	// Validate file type
	contentType := header.Header.Get("Content-Type")
	allowedTypes := []string{"image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp"}
	if !contains(allowedTypes, contentType) {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid file type. Only JPEG, PNG, GIF, and WebP are allowed")
		return
	}

	// Validate file size (5MB max)
	if header.Size > 5<<20 { // 5MB
		utils.ErrorResponse(c, http.StatusBadRequest, "File size too large. Maximum size is 5MB")
		return
	}

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".jpg" // Default extension
	}
	
	uniqueID := uuid.New().String()
	filename := fmt.Sprintf("%s%s", uniqueID, ext)

	// Create uploads directory if it doesn't exist
	uploadDir := "uploads/images"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create upload directory")
		return
	}

	// Create the file path
	filePath := filepath.Join(uploadDir, filename)

	// Create the file on disk
	dst, err := os.Create(filePath)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create file")
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the destination
	_, err = io.Copy(dst, file)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to save file")
		return
	}

	// Return the file path
	response := gin.H{
		"filename": filename,
		"path":     filePath,
		"url":      fmt.Sprintf("/api/v1/uploads/images/%s", filename),
		"size":     header.Size,
		"type":     contentType,
	}

	utils.SuccessResponse(c, response)
}

// ServeImage serves uploaded images
func (h *Handler) ServeImage(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Filename is required")
		return
	}

	// Security check: prevent directory traversal
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid filename")
		return
	}

	filePath := filepath.Join("uploads/images", filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		utils.ErrorResponse(c, http.StatusNotFound, "Image not found")
		return
	}

	// Set appropriate headers
	c.Header("Content-Type", "image/jpeg")
	c.Header("Cache-Control", "public, max-age=31536000") // Cache for 1 year

	// Serve the file
	c.File(filePath)
}

// DeleteImage deletes an uploaded image
func (h *Handler) DeleteImage(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Filename is required")
		return
	}

	// Security check: prevent directory traversal
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid filename")
		return
	}

	filePath := filepath.Join("uploads/images", filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		utils.ErrorResponse(c, http.StatusNotFound, "Image not found")
		return
	}

	// Delete the file
	if err := os.Remove(filePath); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete image")
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "Image deleted successfully"})
}

// Helper function to check if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
