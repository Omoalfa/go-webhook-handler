//go:build ignore

package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	fmt.Println("Scanning for OpenAPI files (*.api.yaml)...")

	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(path, ".api.yaml") {
			dir := filepath.Dir(path)
			baseName := strings.TrimSuffix(d.Name(), ".api.yaml")
			outputFile := filepath.Join(dir, fmt.Sprintf("%s.gen.go", baseName))

			// 1. Get the absolute path of the directory to safely extract the last folder name
			absDir, err := filepath.Abs(dir)
			if err != nil {
				return fmt.Errorf("failed to get absolute path for %s: %w", dir, err)
			}
			packageName := filepath.Base(absDir)

			// Clean up package name if the yaml is in the root directory "."
			if packageName == "." || packageName == "/" {
				packageName = "main"
			}

			fmt.Printf("Found: %s -> Package: %s -> Generating: %s\n", path, packageName, outputFile)

			// 2. Execute oapi-codegen passing the configuration file and the dynamic package name flag (-p)
			cmd := exec.Command("go", "tool", "oapi-codegen",
				"-config", "oapi-codegen.yaml",
				"-package", packageName,
				"-o", outputFile,
				path,
			)

			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed generating for %s: %w", path, err)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error during generation: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Code generation complete!")
}
