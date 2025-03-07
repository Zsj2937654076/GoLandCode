# Student Management System

A comprehensive student management system with Go backend, Vue 3 frontend, and MySQL database.

## Features

- **Student Management**
  - List students with filtering and pagination
  - Add, edit, and delete students
  - View student details
  - Associate students with classes

- **Class Management**
  - List all classes
  - Add, edit, and delete classes
  - View class details and students in a class

- **User Authentication**
  - Login/logout functionality
  - User profile management
  - Password change

## Technology Stack

- **Backend**
  - Go (Golang)
  - MySQL database
  - JWT for authentication

- **Frontend**
  - Vue 3
  - Element Plus UI framework
  - Vuex for state management
  - Vue Router for routing

## Project Structure

```
student_management/
├── backend/             # Go backend code
│   ├── api/             # API definitions
│   ├── config/          # Configuration
│   ├── controllers/     # Request handlers
│   ├── middleware/      # Middleware components
│   ├── models/          # Database models
│   ├── routes/          # Route definitions
│   └── utils/           # Utility functions
├── frontend/            # Vue 3 frontend code
│   ├── public/          # Static assets
│   └── src/             # Source code
│       ├── assets/      # Frontend assets
│       ├── components/  # Vue components
│       ├── router/      # Vue Router configuration
│       ├── store/       # Vuex store
│       ├── views/       # Vue views
│       └── services/    # API services
└── db/                  # Database scripts
    ├── migrations/      # Database schema
    └── seeds/           # Sample data
```

## Setup Instructions

### Prerequisites

- Go 1.19 or higher
- Node.js 14 or higher
- MySQL 8.0

### Backend Setup

1. Navigate to the backend directory:
   ```
   cd backend
   ```

2. Install Go dependencies:
   ```
   go mod tidy
   ```

3. Set up the MySQL database:
   ```
   mysql -u root -p < ../db/migrations/schema.sql
   mysql -u root -p < ../db/seeds/sample_data.sql
   ```

4. Start the backend server:
   ```
   go run main.go
   ```

### Frontend Setup

1. Navigate to the frontend directory:
   ```
   cd frontend
   ```

2. Install dependencies:
   ```
   npm install
   ```

3. Start the development server:
   ```
   npm run serve
   ```

## Usage

1. Access the application at `http://localhost:8080` (or the port configured in your environment)
2. Login with the default admin credentials:
   - Username: `admin`
   - Password: `admin123`

## API Endpoints

### Authentication
- `POST /api/auth/login` - User login
- `POST /api/auth/logout` - User logout
- `GET /api/auth/profile` - Get user profile
- `POST /api/auth/change-password` - Change user password

### Students
- `GET /api/students` - List students (with filtering and pagination)
- `GET /api/students/{id}` - Get student details
- `POST /api/students` - Create a new student
- `PUT /api/students/{id}` - Update a student
- `DELETE /api/students/{id}` - Delete a student

### Classes
- `GET /api/classes` - List all classes
- `GET /api/classes/{id}` - Get class details
- `GET /api/classes/{id}/students` - Get students in a class
- `POST /api/classes` - Create a new class
- `PUT /api/classes/{id}` - Update a class
- `DELETE /api/classes/{id}` - Delete a class

## License

This project is licensed under the MIT License. 