package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run generate_service.go <service-name>")
		os.Exit(1)
	}

	serviceName := os.Args[1]
	if err := GenerateService(serviceName); err != nil {
		fmt.Printf("Error generating service: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully generated service: %s\n", serviceName)
}

// Validate service name (only lowercase letters and hyphens)
func validateServiceName(name string) error {
	// Check if name is empty
	if name == "" {
		return fmt.Errorf("service name cannot be empty")
	}

	// Check if name contains only lowercase letters and hyphens
	matched, err := regexp.MatchString("^[a-z]+(-[a-z]+)*$", name)
	if err != nil {
		return fmt.Errorf("error validating service name: %v", err)
	}
	if !matched {
		return fmt.Errorf("invalid service name '%s': must contain only lowercase letters and hyphens, and cannot start or end with a hyphen", name)
	}

	return nil
}

// Check if service already exists
func serviceExists(name string) bool {
	serviceDir := filepath.Join("services", name)
	if _, err := os.Stat(serviceDir); err == nil {
		return true
	}
	return false
}

// Convert kebab-case to CamelCase
func toCamelCase(kebab string) string {
	var result strings.Builder
	capitalize := true

	for _, char := range kebab {
		if char == '-' || char == '_' {
			capitalize = true
			continue
		}

		if capitalize {
			result.WriteString(strings.ToUpper(string(char)))
			capitalize = false
		} else {
			result.WriteString(strings.ToLower(string(char)))
		}
	}

	return result.String()
}

// Convert hyphens to underscores
func toUnderscoreCase(s string) string {
	return strings.ReplaceAll(s, "-", "_")
}

type ServiceData struct {
	ServiceName           string // Original name with hyphens (e.g., "test-service")
	ServiceNameCamel      string // CamelCase version (e.g., "TestService")
	ServiceNameUnderscore string // Underscore version for proto package (e.g., "test_service")
}

func GenerateService(serviceName string) error {
	// Validate service name
	if err := validateServiceName(serviceName); err != nil {
		return err
	}

	// Check if service already exists
	if serviceExists(serviceName) {
		return fmt.Errorf("service '%s' already exists", serviceName)
	}

	data := ServiceData{
		ServiceName:           serviceName,
		ServiceNameCamel:      toCamelCase(serviceName),
		ServiceNameUnderscore: toUnderscoreCase(serviceName),
	}

	// Get the template directory
	templateDir := filepath.Join("core", "boilerplate", "service")
	targetDir := filepath.Join("services", serviceName)

	// Create the target directory
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create service directory: %v", err)
	}

	// Walk through the template directory
	return filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Get the relative path from template directory
		relPath, err := filepath.Rel(templateDir, path)
		if err != nil {
			return err
		}

		// Create target file path with special handling for CI file and swagger file
		targetPath := filepath.Join(targetDir, strings.TrimSuffix(relPath, ".template"))
		if strings.HasSuffix(targetPath, "service.ci.yml") {
			// Rename service.ci.yml to <service_name>.ci.yml
			dir := filepath.Dir(targetPath)
			targetPath = filepath.Join(dir, toUnderscoreCase(serviceName)+".ci.yml")
		} else if strings.HasSuffix(targetPath, "handler.go") {
			// Rename handler.go to <service_name>_handler.go
			dir := filepath.Dir(targetPath)
			targetPath = filepath.Join(dir, toUnderscoreCase(serviceName)+"_handler.go")
		} else if strings.HasSuffix(targetPath, "service.proto") {
			// Rename service.proto to <service_name>.proto
			dir := filepath.Dir(targetPath)
			targetPath = filepath.Join(dir, serviceName+".proto")
		} else if strings.HasSuffix(targetPath, "service.go") {
			// Rename service.proto to <service_name>.proto
			dir := filepath.Dir(targetPath)
			targetPath = filepath.Join(dir, serviceName+".go")
		}

		// Create parent directories if they don't exist
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return err
		}

		// Parse and execute the template
		tmpl, err := template.ParseFiles(path)
		if err != nil {
			return err
		}

		// Create the target file
		f, err := os.Create(targetPath)
		if err != nil {
			return err
		}
		defer f.Close()

		// Execute the template
		if err := tmpl.Execute(f, data); err != nil {
			return err
		}

		return nil
	})
}
