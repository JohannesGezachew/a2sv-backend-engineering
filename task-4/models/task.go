package models

import "time"

// Task represents a task in the task management system
type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	Status      string    `json:"status" binding:"required"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TaskRequest represents the request payload for creating/updating tasks
type TaskRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
	Status      string `json:"status" binding:"required"`
}

// TaskResponse represents the response format for task operations
type TaskResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse represents error response format
type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}