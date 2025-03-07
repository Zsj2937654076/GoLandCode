package models

import (
	"database/sql"
	"time"
)

// Class represents a class in the school
type Class struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	StudentCount int      `json:"student_count,omitempty"` // Not stored in DB, calculated when needed
}

// GetAllClasses retrieves all classes from the database
func GetAllClasses(db *sql.DB) ([]Class, error) {
	query := `
		SELECT c.id, c.name, c.description, c.created_at, c.updated_at,
		(SELECT COUNT(*) FROM students s WHERE s.class_id = c.id) as student_count
		FROM classes c
		ORDER BY c.name
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classes []Class
	for rows.Next() {
		var c Class
		err := rows.Scan(
			&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt, &c.StudentCount,
		)
		if err != nil {
			return nil, err
		}
		classes = append(classes, c)
	}
	return classes, nil
}

// GetClassByID retrieves a class by ID
func GetClassByID(db *sql.DB, id int64) (Class, error) {
	var class Class
	query := `
		SELECT c.id, c.name, c.description, c.created_at, c.updated_at,
		(SELECT COUNT(*) FROM students s WHERE s.class_id = c.id) as student_count
		FROM classes c
		WHERE c.id = ?
	`
	err := db.QueryRow(query, id).Scan(
		&class.ID, &class.Name, &class.Description, &class.CreatedAt, &class.UpdatedAt, &class.StudentCount,
	)
	return class, err
}

// CreateClass inserts a new class into the database
func CreateClass(db *sql.DB, class *Class) (int64, error) {
	query := `
		INSERT INTO classes (name, description, created_at, updated_at)
		VALUES (?, ?, NOW(), NOW())
	`
	result, err := db.Exec(query, class.Name, class.Description)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// UpdateClass updates an existing class
func UpdateClass(db *sql.DB, class *Class) error {
	query := `
		UPDATE classes
		SET name = ?, description = ?, updated_at = NOW()
		WHERE id = ?
	`
	_, err := db.Exec(query, class.Name, class.Description, class.ID)
	return err
}

// DeleteClass removes a class from the database
func DeleteClass(db *sql.DB, id int64) error {
	// First check if there are students in this class
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM students WHERE class_id = ?", id).Scan(&count)
	if err != nil {
		return err
	}
	
	if count > 0 {
		return sql.ErrNoRows // Using standard error to indicate class has students
	}
	
	query := "DELETE FROM classes WHERE id = ?"
	_, err = db.Exec(query, id)
	return err
} 