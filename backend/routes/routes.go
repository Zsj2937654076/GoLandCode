package routes

import (
	"database/sql"
	"net/http"
	"student-management/controllers"
	"student-management/middleware"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// SetupRouter configures all the routes for the API
func SetupRouter(db *sql.DB) http.Handler {
	// Create new router
	router := mux.NewRouter()
	
	// Create API subrouter
	api := router.PathPrefix("/api").Subrouter()

	// Initialize controllers
	studentController := controllers.NewStudentController(db)
	classController := controllers.NewClassController(db)
	authController := controllers.NewAuthController(db)

	// Auth routes (public)
	authRoutes := api.PathPrefix("/auth").Subrouter()
	authRoutes.HandleFunc("/login", authController.Login).Methods("POST")
	authRoutes.HandleFunc("/logout", authController.Logout).Methods("POST")
	
	// Protected auth routes
	protectedAuthRoutes := authRoutes.NewRoute().Subrouter()
	protectedAuthRoutes.Use(middleware.AuthMiddleware)
	protectedAuthRoutes.HandleFunc("/profile", authController.Profile).Methods("GET")
	protectedAuthRoutes.HandleFunc("/change-password", authController.ChangePassword).Methods("POST")

	// Protected API routes
	protectedAPI := api.NewRoute().Subrouter()
	protectedAPI.Use(middleware.AuthMiddleware)

	// Student routes
	students := protectedAPI.PathPrefix("/students").Subrouter()
	students.HandleFunc("", studentController.GetStudents).Methods("GET")
	students.HandleFunc("/{id:[0-9]+}", studentController.GetStudentByID).Methods("GET")
	students.HandleFunc("", studentController.CreateStudent).Methods("POST")
	students.HandleFunc("/{id:[0-9]+}", studentController.UpdateStudent).Methods("PUT")
	students.HandleFunc("/{id:[0-9]+}", studentController.DeleteStudent).Methods("DELETE")

	// Class routes
	classes := protectedAPI.PathPrefix("/classes").Subrouter()
	classes.HandleFunc("", classController.GetClasses).Methods("GET")
	classes.HandleFunc("/{id:[0-9]+}", classController.GetClassByID).Methods("GET")
	classes.HandleFunc("", classController.CreateClass).Methods("POST")
	classes.HandleFunc("/{id:[0-9]+}", classController.UpdateClass).Methods("PUT")
	classes.HandleFunc("/{id:[0-9]+}", classController.DeleteClass).Methods("DELETE")
	classes.HandleFunc("/{id:[0-9]+}/students", classController.GetClassStudents).Methods("GET")
	
	// Set up CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://www.zsjurl.top"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		AllowCredentials: true,
	})

	// Wrap router with CORS middleware
	handler := c.Handler(router)
	
	return handler
} 