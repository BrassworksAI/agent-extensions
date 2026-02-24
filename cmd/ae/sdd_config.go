package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func runSddConfigInit(cmd *cobra.Command, args []string) error {
	repoRootAbs, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("resolving working directory: %w", err)
	}

	specRoot := strings.TrimSpace(flagConfigInitSpecDir)
	if specRoot == "" {
		return fmt.Errorf("invalid --spec-root: must be a non-empty repo-relative path")
	}

	configPath := filepath.Join(repoRootAbs, ".ae-config.json")
	if !flagConfigInitForce {
		if _, err := os.Stat(configPath); err == nil {
			return fmt.Errorf("%s already exists (use --force to overwrite)", filepath.Base(configPath))
		} else if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("checking %s: %w", filepath.Base(configPath), err)
		}
	}

	cfg := aeConfig{
		Schema:   aeConfigSchemaURL,
		SpecRoot: &specRoot,
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding config: %w", err)
	}
	data = append(data, '\n')

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("writing %s: %w", filepath.Base(configPath), err)
	}

	fmt.Printf("✓ Initialized %s\n", filepath.Base(configPath))
	fmt.Printf("  schema: %s\n", aeConfigSchemaURL)
	fmt.Printf("  specRoot: %s\n", specRoot)
	return nil
}
