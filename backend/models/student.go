package models

import (
	"database/sql"
	"time"
)

// Student represents a student in the system
type Student struct {
	ID        int64     `json:"id"`
	StudentID string    `json:"student_id"` // University/School ID
	Name      string    `json:"name"`
	ClassID   int64     `json:"class_id"`
	ClassName string    `json:"class_name,omitempty"` // Not stored in DB, populated when joining with class
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetAllStudents retrieves all students with optional filters and pagination
func GetAllStudents(db *sql.DB, classID int64, studentID, name string, page, pageSize int) ([]Student, int, error) {
	query := `
		SELECT s.id, s.student_id, s.name, s.class_id, c.name as class_name, 
		s.email, s.phone, s.address, s.created_at, s.updated_at
		FROM students s
		LEFT JOIN classes c ON s.class_id = c.id
		WHERE 1=1
	`
	countQuery := `SELECT COUNT(*) FROM students WHERE 1=1`
	params := []interface{}{}

	// Apply filters
	if classID > 0 {
		query += " AND s.class_id = ?"
		countQuery += " AND class_id = ?"
		params = append(params, classID)
	}
	if studentID != "" {
		query += " AND s.student_id LIKE ?"
		countQuery += " AND student_id LIKE ?"
		params = append(params, "%"+studentID+"%")
	}
	if name != "" {
		query += " AND s.name LIKE ?"
		countQuery += " AND name LIKE ?"
		params = append(params, "%"+name+"%")
	}

	// Apply pagination
	query += " ORDER BY s.id DESC LIMIT ? OFFSET ?"
	offset := (page - 1) * pageSize
	params = append(params, pageSize, offset)

	// Execute the count query
	var total int
	err := db.QueryRow(countQuery, params[:len(params)-2]...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Execute the main query
	rows, err := db.Query(query, params...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var students []Student
	for rows.Next() {
		var s Student
		err := rows.Scan(
			&s.ID, &s.StudentID, &s.Name, &s.ClassID, &s.ClassName,
			&s.Email, &s.Phone, &s.Address, &s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		students = append(students, s)
	}

	return students, total, nil
}

// GetStudentByID retrieves a student by ID
func GetStudentByID(db *sql.DB, id int64) (Student, error) {
	var student Student
	query := `
		SELECT s.id, s.student_id, s.name, s.class_id, c.name as class_name, 
		s.email, s.phone, s.address, s.created_at, s.updated_at
		FROM students s
		LEFT JOIN classes c ON s.class_id = c.id
		WHERE s.id = ?
	`
	err := db.QueryRow(query, id).Scan(
		&student.ID, &student.StudentID, &student.Name, &student.ClassID, &student.ClassName,
		&student.Email, &student.Phone, &student.Address, &student.CreatedAt, &student.UpdatedAt,
	)
	return student, err
}

// CreateStudent inserts a new student into the database
func CreateStudent(db *sql.DB, student *Student) (int64, error) {
	query := `
		INSERT INTO students (student_id, name, class_id, email, phone, address, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())
	`
	result, err := db.Exec(query,
		student.StudentID, student.Name, student.ClassID,
		student.Email, student.Phone, student.Address,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// UpdateStudent updates an existing student
func UpdateStudent(db *sql.DB, student *Student) error {
	query := `
		UPDATE students
		SET student_id = ?, name = ?, class_id = ?, 
		    email = ?, phone = ?, address = ?, updated_at = NOW()
		WHERE id = ?
	`
	_, err := db.Exec(query,
		student.StudentID, student.Name, student.ClassID,
		student.Email, student.Phone, student.Address, student.ID,
	)
	return err
}

// DeleteStudent removes a student from the database
func DeleteStudent(db *sql.DB, id int64) error {
	query := "DELETE FROM students WHERE id = ?"
	_, err := db.Exec(query, id)
	return err
}

// GetStudentsByClassID retrieves all students in a specific class
func GetStudentsByClassID(db *sql.DB, classID int64) ([]Student, error) {
	query := `
		SELECT s.id, s.student_id, s.name, s.class_id, c.name as class_name, 
		s.email, s.phone, s.address, s.created_at, s.updated_at
		FROM students s
		LEFT JOIN classes c ON s.class_id = c.id
		WHERE s.class_id = ?
		ORDER BY s.name
	`
	rows, err := db.Query(query, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []Student
	for rows.Next() {
		var s Student
		err := rows.Scan(
			&s.ID, &s.StudentID, &s.Name, &s.ClassID, &s.ClassName,
			&s.Email, &s.Phone, &s.Address, &s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		students = append(students, s)
	}
	return students, nil
} 