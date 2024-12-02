package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: gin-cli <command> [arguments]")
		os.Exit(1)
	}

	newCmd := flag.NewFlagSet("new", flag.ExitOnError)
	projectName := newCmd.String("name", "", "Project name")

	switch os.Args[1] {
	case "new":
		newCmd.Parse(os.Args[2:])
		if *projectName == "" {
			fmt.Println("Please provide project name")
			os.Exit(1)
		}
		createProject(*projectName)
	case "startapp":
		if len(os.Args) < 3 {
			fmt.Println("Please provide app name")
			os.Exit(1)
		}
		createApp(os.Args[2])
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func createProject(name string) {
	dirs := []string{
		name,
		name + "/cmd",
		name + "/internal",
		name + "/internal/config",
		name + "/internal/middleware",
		name + "/templates",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("Error creating directory %s: %v\n", dir, err)
			os.Exit(1)
		}
	}

	// Create main.go
	mainContent := []byte(`package main

import (
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    r.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "Welcome to Gin",
        })
    })
    r.Run()
}`)

	if err := os.WriteFile(name+"/cmd/main.go", mainContent, 0644); err != nil {
		fmt.Printf("Error creating main.go: %v\n", err)
		os.Exit(1)
	}
}

func createApp(name string) {
	dirs := []string{
		"internal/" + name,
		"internal/" + name + "/handlers",
		"internal/" + name + "/models",
		"internal/" + name + "/routes",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("Error creating directory %s: %v\n", dir, err)
			os.Exit(1)
		}
	}
}
