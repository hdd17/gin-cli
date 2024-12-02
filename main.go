package main

import (
	"fmt"
	"os"
	"os/exec"
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
	cmd := exec.Command("go", "run", "main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func createProject(name string) {
	// Create project directory
	os.MkdirAll(name, 0755)

	// Create main.go
	mainContent := []byte(`package main

import "github.com/gin-gonic/gin"

func main() {
    r := gin.Default()
    r.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "Welcome to Gin",
        })
    })
    r.Run()
}`)

	if err := os.WriteFile(name+"/main.go", mainContent, 0644); err != nil {
		fmt.Printf("Error creating main.go: %v\n", err)
		os.Exit(1)
	}

	// Initialize go module
	cmd := exec.Command("go", "mod", "init", name)
	cmd.Dir = name
	cmd.Run()

	// Get Gin dependency
	cmd = exec.Command("go", "get", "github.com/gin-gonic/gin")
	cmd.Dir = name
	cmd.Run()

	fmt.Printf("Project %s created successfully!\n", name)
	fmt.Println("To start your project:")
	fmt.Printf("1. cd %s\n", name)
	fmt.Println("2. gin-cli start")
}
