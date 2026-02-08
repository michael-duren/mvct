package scaffold

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/michael-duren/mvct/internal/templates"
)

type ProjectConfiguration struct {
	Name   string
	Path   string
	WithDB bool
}

func ScaffoldProject(config ProjectConfiguration) {
	fmt.Printf("Scaffolding project %s at %s (DB: %v)\n", config.Name, config.Path, config.WithDB)

	// Create directories
	dirs := []string{
		config.Path,
		filepath.Join(config.Path, "cmd", "cli"),
		filepath.Join(config.Path, "internal"),
		filepath.Join(config.Path, "components"),
		filepath.Join(config.Path, "controllers"),
	}

	if config.WithDB {
		dirs = append(dirs, filepath.Join(config.Path, "db"))
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("Error creating directory %s: %v\n", dir, err)
			return
		}
	}

	// Generate go.mod
	createFile(filepath.Join(config.Path, "go.mod"), templates.GoModTemplate, config)

	// Generate main.go
	createFile(filepath.Join(config.Path, "cmd", "cli", "main.go"), templates.MainGoTemplate, config)

	// Generate components
	createFile(filepath.Join(config.Path, "components", "layout.go"), templates.LayoutTemplate, config)

	// Generate controllers
	createFile(filepath.Join(config.Path, "controllers", "home.go"), templates.HomeControllerTemplate, config)
	createFile(filepath.Join(config.Path, "controllers", "routing.go"), templates.RoutingTemplate, config)

	// DB specific files
	if config.WithDB {
		createFile(filepath.Join(config.Path, "sqlc.yaml"), templates.SqlcYamlTemplate, config)
		createFile(filepath.Join(config.Path, "db", "schema.sql"), templates.SchemaSqlTemplate, config)
		createFile(filepath.Join(config.Path, "db", "queries.sql"), templates.QueriesSqlTemplate, config)
		createFile(filepath.Join(config.Path, "Makefile"), templates.MakefileTemplate, config)
	}

	// Initialize git
	runCmd(config.Path, "git", "init")

	// Tidy go mod
	// We need to run this inside the project directory, ensuring the go.mod exists first
	runCmd(config.Path, "go", "mod", "tidy")

	fmt.Println("Project scaffolded successfully!")
	if config.WithDB {
		fmt.Println("To generate sqlc code, run: make sqlc")
	}
	fmt.Printf("Run your app: cd %s && go run cmd/cli/main.go\n", config.Path)
}

func createFile(path string, tmpl string, data any) {
	t, err := template.New("file").Parse(tmpl)
	if err != nil {
		fmt.Printf("Error parsing template for %s: %v\n", path, err)
		return
	}

	f, err := os.Create(path)
	if err != nil {
		fmt.Printf("Error creating file %s: %v\n", path, err)
		return
	}
	defer f.Close()

	if err := t.Execute(f, data); err != nil {
		fmt.Printf("Error executing template for %s: %v\n", path, err)
	}
}

func runCmd(dir string, name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running command %s %v: %v\n", name, args, err)
	}
}
