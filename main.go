package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  new <project-name> - Create new project")
		fmt.Println("  start - Run the server")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "new":
		if len(os.Args) < 3 {
			fmt.Println("Please provide project name")
			os.Exit(1)
		}
		createProject(os.Args[2])
	case "start":
		startProject()
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func startProject() {
	cmd := exec.Command("go", "run", "./cmd/main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func createDirectoryStructure(projectPath string) error {
	dirs := []string{
		"cmd",
		"internal/handlers",
		"internal/middleware",
		"internal/models",
		"internal/repository",
		"internal/service",
		"pkg/config",
		"pkg/utils",
		"api",
		"config",
		"migrations",
	}

	for _, dir := range dirs {
		err := os.MkdirAll(filepath.Join(projectPath, dir), 0755)
		if err != nil {
			return fmt.Errorf("error creating directory %s: %v", dir, err)
		}
	}
	return nil
}

func writeFile(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

func createProject(name string) {
	projectPath := name

	// Create project directory and structure
	if err := createDirectoryStructure(projectPath); err != nil {
		fmt.Printf("Error creating directory structure: %v\n", err)
		os.Exit(1)
	}

	// Initialize go module
	cmd := exec.Command("go", "mod", "init", name)
	cmd.Dir = projectPath
	cmd.Run()

	// Get dependencies
	dependencies := []string{
		"github.com/gin-gonic/gin",
		"github.com/joho/godotenv",
		"gorm.io/gorm",
		"gorm.io/driver/postgres",
	}

	for _, dep := range dependencies {
		cmd = exec.Command("go", "get", dep)
		cmd.Dir = projectPath
		cmd.Run()
	}

	// Create main.go
	mainContent := `package main

import (
    "log"
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
)

func main() {
    // Load .env file
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }

    // Initialize router
    r := gin.Default()

    // Add middleware
    r.Use(gin.Logger())
    r.Use(gin.Recovery())

    // Routes
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "status": "OK",
        })
    })

    // Start server
    r.Run(":8080")
}`

	if err := writeFile(filepath.Join(projectPath, "cmd", "main.go"), mainContent); err != nil {
		fmt.Printf("Error creating main.go: %v\n", err)
		os.Exit(1)
	}

	// Create .env file
	envContent := `DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=dbname
SERVER_PORT=8080`

	if err := writeFile(filepath.Join(projectPath, ".env"), envContent); err != nil {
		fmt.Printf("Error creating .env: %v\n", err)
	}

	// Create basic handler
	handlerContent := `package handlers

import "github.com/gin-gonic/gin"

type Handler struct {
    // Add service interfaces here
}

func NewHandler() *Handler {
    return &Handler{}
}

func (h *Handler) InitRoutes() *gin.Engine {
    router := gin.Default()

    api := router.Group("/api")
    {
        v1 := api.Group("/v1")
        {
            v1.GET("/health", h.healthCheck)
        }
    }

    return router
}

func (h *Handler) healthCheck(c *gin.Context) {
    c.JSON(200, gin.H{
        "status": "OK",
    })
}`

	if err := writeFile(filepath.Join(projectPath, "internal/handlers", "handler.go"), handlerContent); err != nil {
		fmt.Printf("Error creating handler.go: %v\n", err)
	}

	// Create basic model
	modelContent := `package models

import "time"

type Base struct {
    ID        uint      ` + "`json:\"id\" gorm:\"primaryKey\"`" + `
    CreatedAt time.Time ` + "`json:\"created_at\"`" + `
    UpdatedAt time.Time ` + "`json:\"updated_at\"`" + `
}`

	if err := writeFile(filepath.Join(projectPath, "internal/models", "base.go"), modelContent); err != nil {
		fmt.Printf("Error creating base.go: %v\n", err)
	}

	// Create README.md
	readmeContent := fmt.Sprintf(`# %s

## Description
A Go web application using the Gin framework.

## Setup
1. Copy .env.example to .env and configure your environment variables
2. Run the application:
   ~~~bash
   gin-cli start
   ~~~

## Project Structure
- cmd/: Application entry points
- internal/: Private application code
  - handlers/: HTTP request handlers
  - middleware/: Custom middleware
  - models/: Data models
  - repository/: Data access layer
  - service/: Business logic
- pkg/: Public libraries
- api/: API documentation
- config/: Configuration files
- migrations/: Database migrations
`, name)

	if err := writeFile(filepath.Join(projectPath, "README.md"), readmeContent); err != nil {
		fmt.Printf("Error creating README.md: %v\n", err)
	}

	fmt.Printf("Project %s created successfully!\n", name)
	fmt.Println("\nTo start your project:")
	fmt.Printf("1. cd %s\n", name)
	fmt.Println("2. Configure your .env file")
	fmt.Println("3. gin-cli start")
}
