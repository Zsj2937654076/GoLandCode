package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"student-management/models"

	"github.com/gorilla/mux"
)

// StudentController handles the student-related API endpoints
type StudentController struct {
	DB *sql.DB
}

// NewStudentController creates a new StudentController instance
func NewStudentController(db *sql.DB) *StudentController {
	return &StudentController{DB: db}
}

// GetStudents handles GET /api/students to retrieve the student list with filtering and pagination
func (c *StudentController) GetStudents(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	classIDStr := r.URL.Query().Get("class_id")
	studentID := r.URL.Query().Get("student_id")
	name := r.URL.Query().Get("name")
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")

	// Default values for pagination
	page := 1
	pageSize := 10

	// Parse class_id
	var classID int64
	if classIDStr != "" {
		var err error
		classID, err = strconv.ParseInt(classIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid class_id parameter", http.StatusBadRequest)
			return
		}
	}

	// Parse page and page_size
	if pageStr != "" {
		var err error
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			http.Error(w, "Invalid page parameter", http.StatusBadRequest)
			return
		}
	}
	if pageSizeStr != "" {
		var err error
		pageSize, err = strconv.Atoi(pageSizeStr)
		if err != nil || pageSize < 1 || pageSize > 100 {
			http.Error(w, "Invalid page_size parameter", http.StatusBadRequest)
			return
		}
	}

	// Get students from the database
	students, total, err := models.GetAllStudents(c.DB, classID, studentID, name, page, pageSize)
	if err != nil {
		http.Error(w, "Failed to retrieve students", http.StatusInternalServerError)
		return
	}

	// Prepare response with pagination metadata
	response := map[string]interface{}{
		"data": students,
		"pagination": map[string]interface{}{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetStudentByID handles GET /api/students/{id} to retrieve a specific student
func (c *StudentController) GetStudentByID(w http.ResponseWriter, r *http.Request) {
	// Get student ID from URL
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	// Get student from database
	student, err := models.GetStudentByID(c.DB, id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Student not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve student", http.StatusInternalServerError)
		}
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}

// CreateStudent handles POST /api/students to create a new student
func (c *StudentController) CreateStudent(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var student models.Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if student.Name == "" || student.StudentID == "" {
		http.Error(w, "Name and StudentID are required fields", http.StatusBadRequest)
		return
	}

	// Create student in database
	id, err := models.CreateStudent(c.DB, &student)
	if err != nil {
		http.Error(w, "Failed to create student", http.StatusInternalServerError)
		return
	}

	// Set the student ID and get the full student details
	student.ID = id
	createdStudent, err := models.GetStudentByID(c.DB, id)
	if err != nil {
		http.Error(w, "Student created but failed to retrieve details", http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdStudent)
}

// UpdateStudent handles PUT /api/students/{id} to update a student
func (c *StudentController) UpdateStudent(w http.ResponseWriter, r *http.Request) {
	// Get student ID from URL
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	// Check if student exists
	_, err = models.GetStudentByID(c.DB, id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Student not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve student", http.StatusInternalServerError)
		}
		return
	}

	// Parse request body
	var student models.Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set ID to match the URL parameter
	student.ID = id

	// Update student in database
	err = models.UpdateStudent(c.DB, &student)
	if err != nil {
		http.Error(w, "Failed to update student", http.StatusInternalServerError)
		return
	}

	// Get updated student
	updatedStudent, err := models.GetStudentByID(c.DB, id)
	if err != nil {
		http.Error(w, "Student updated but failed to retrieve details", http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedStudent)
}

// DeleteStudent handles DELETE /api/students/{id} to delete a student
func (c *StudentController) DeleteStudent(w http.ResponseWriter, r *http.Request) {
	// Get student ID from URL
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	// Check if student exists
	_, err = models.GetStudentByID(c.DB, id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Student not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve student", http.StatusInternalServerError)
		}
		return
	}

	// Delete student from database
	err = models.DeleteStudent(c.DB, id)
	if err != nil {
		http.Error(w, "Failed to delete student", http.StatusInternalServerError)
		return
	}

	// Send response
	w.WriteHeader(http.StatusNoContent)
} 