package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"student-management/models"

	"github.com/gorilla/mux"
)

// ClassController handles the class-related API endpoints
type ClassController struct {
	DB *sql.DB
}

// NewClassController creates a new ClassController instance
func NewClassController(db *sql.DB) *ClassController {
	return &ClassController{DB: db}
}

// GetClasses handles GET /api/classes to retrieve all classes
func (c *ClassController) GetClasses(w http.ResponseWriter, r *http.Request) {
	// Get classes from the database
	classes, err := models.GetAllClasses(c.DB)
	if err != nil {
		http.Error(w, "Failed to retrieve classes", http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(classes)
}

// GetClassByID handles GET /api/classes/{id} to retrieve a specific class
func (c *ClassController) GetClassByID(w http.ResponseWriter, r *http.Request) {
	// Get class ID from URL
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid class ID", http.StatusBadRequest)
		return
	}

	// Get class from database
	class, err := models.GetClassByID(c.DB, id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Class not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve class", http.StatusInternalServerError)
		}
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(class)
}

// CreateClass handles POST /api/classes to create a new class
func (c *ClassController) CreateClass(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var class models.Class
	if err := json.NewDecoder(r.Body).Decode(&class); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if class.Name == "" {
		http.Error(w, "Name is a required field", http.StatusBadRequest)
		return
	}

	// Create class in database
	id, err := models.CreateClass(c.DB, &class)
	if err != nil {
		http.Error(w, "Failed to create class", http.StatusInternalServerError)
		return
	}

	// Set the class ID and get the full class details
	class.ID = id
	createdClass, err := models.GetClassByID(c.DB, id)
	if err != nil {
		http.Error(w, "Class created but failed to retrieve details", http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdClass)
}

// UpdateClass handles PUT /api/classes/{id} to update a class
func (c *ClassController) UpdateClass(w http.ResponseWriter, r *http.Request) {
	// Get class ID from URL
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid class ID", http.StatusBadRequest)
		return
	}

	// Check if class exists
	_, err = models.GetClassByID(c.DB, id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Class not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve class", http.StatusInternalServerError)
		}
		return
	}

	// Parse request body
	var class models.Class
	if err := json.NewDecoder(r.Body).Decode(&class); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set ID to match the URL parameter
	class.ID = id

	// Update class in database
	err = models.UpdateClass(c.DB, &class)
	if err != nil {
		http.Error(w, "Failed to update class", http.StatusInternalServerError)
		return
	}

	// Get updated class
	updatedClass, err := models.GetClassByID(c.DB, id)
	if err != nil {
		http.Error(w, "Class updated but failed to retrieve details", http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedClass)
}

// DeleteClass handles DELETE /api/classes/{id} to delete a class
func (c *ClassController) DeleteClass(w http.ResponseWriter, r *http.Request) {
	// Get class ID from URL
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid class ID", http.StatusBadRequest)
		return
	}

	// Check if class exists
	_, err = models.GetClassByID(c.DB, id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Class not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve class", http.StatusInternalServerError)
		}
		return
	}

	// Delete class from database
	err = models.DeleteClass(c.DB, id)
	if err != nil {
		if err == sql.ErrNoRows {
			// Custom error that indicates class has students
			http.Error(w, "Cannot delete class with associated students", http.StatusConflict)
		} else {
			http.Error(w, "Failed to delete class", http.StatusInternalServerError)
		}
		return
	}

	// Send response
	w.WriteHeader(http.StatusNoContent)
}

// GetClassStudents handles GET /api/classes/{id}/students to get students in a class
func (c *ClassController) GetClassStudents(w http.ResponseWriter, r *http.Request) {
	// Get class ID from URL
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid class ID", http.StatusBadRequest)
		return
	}

	// Check if class exists
	_, err = models.GetClassByID(c.DB, id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Class not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve class", http.StatusInternalServerError)
		}
		return
	}

	// Get students in the class
	students, err := models.GetStudentsByClassID(c.DB, id)
	if err != nil {
		http.Error(w, "Failed to retrieve students", http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
} 