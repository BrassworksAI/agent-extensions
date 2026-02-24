package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/shanepadgett/agent-extensions/internal/sdd"
	"github.com/spf13/cobra"
)

var sddSpecCmd = &cobra.Command{
	Use:   "spec",
	Short: "Spec validation and merge",
}

var sddSpecValidateCmd = &cobra.Command{
	Use:           "validate [name]",
	Short:         "Validate one spec file or all specs in a change set",
	Args:          cobra.MaximumNArgs(1),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE:          runSddSpecValidate,
}

var sddSpecMergeCmd = &cobra.Command{
	Use:           "merge [name]",
	Short:         "Merge change-set specs into canonical specs",
	Args:          cobra.MaximumNArgs(1),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE:          runSddSpecMerge,
}

var sddSpecListCmd = &cobra.Command{
	Use:           "list",
	Short:         "List canonical specs",
	Args:          cobra.NoArgs,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE:          runSddSpecList,
}

var (
	flagSpecValidatePath string
	flagSpecValidateAll  bool
	flagSpecMergeDryRun  bool
	flagSpecMergeDir     string
	flagSpecListDir      string
)

const (
	defaultCanonicalSpecsDir = "docs/specs"
	aeConfigSchemaURL        = "https://raw.githubusercontent.com/shanepadgett/agent-extensions/main/schemas/ae-config.schema.json"
)

func init() {
	sddSpecValidateCmd.Flags().StringVar(&flagSpecValidatePath, "path", "", "Repo-relative path to a single change-set spec file")
	sddSpecValidateCmd.Flags().BoolVar(&flagSpecValidateAll, "all", false, "Validate all specs under changes/<name>/specs")

	sddSpecMergeCmd.Flags().BoolVar(&flagSpecMergeDryRun, "dry-run", false, "Preview merge changes without writing files")
	sddSpecMergeCmd.Flags().StringVar(&flagSpecMergeDir, "specs-dir", "", "Canonical specs root directory (repo-relative). Overrides .ae-config.json specRoot")
	sddSpecListCmd.Flags().StringVar(&flagSpecListDir, "specs-dir", "", "Canonical specs root directory (repo-relative). Overrides .ae-config.json specRoot")

	sddSpecCmd.AddCommand(sddSpecValidateCmd)
	sddSpecCmd.AddCommand(sddSpecMergeCmd)
	sddSpecCmd.AddCommand(sddSpecListCmd)
	sddCmd.AddCommand(sddSpecCmd)
}

func runSddSpecValidate(cmd *cobra.Command, args []string) error {
	if flagSpecValidatePath == "" && !flagSpecValidateAll {
		return fmt.Errorf("must provide exactly one mode: --path <file> or --all [name]")
	}
	if flagSpecValidatePath != "" && flagSpecValidateAll {
		return fmt.Errorf("--path and --all are mutually exclusive")
	}

	repoRootAbs, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("resolving working directory: %w", err)
	}

	if flagSpecValidatePath != "" {
		if len(args) > 0 {
			return fmt.Errorf("positional change-set name is only valid with --all")
		}
		result, err := sdd.ValidateChangeSpecPath(repoRootAbs, flagSpecValidatePath)
		if err != nil {
			return err
		}
		payload, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("encoding validation output: %w", err)
		}
		fmt.Println(string(payload))
		if !result.OK {
			return fmt.Errorf("validation failed")
		}
		return nil
	}

	var changeDir string
	if len(args) > 0 {
		changeDir, err = resolveChangeSet(args[0])
	} else {
		changeDir, err = resolveChangeSet("")
	}
	if err != nil {
		return err
	}

	results, err := sdd.ValidateAllChangeSpecs(repoRootAbs, changeDir)
	if err != nil {
		return err
	}

	ok := true
	for _, result := range results {
		if !result.OK {
			ok = false
			break
		}
	}

	payload := struct {
		OK      bool                   `json:"ok"`
		Change  string                 `json:"change"`
		Results []sdd.ValidationResult `json:"results"`
	}{
		OK:      ok,
		Change:  filepath.Base(changeDir),
		Results: results,
	}

	out, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding validation output: %w", err)
	}
	fmt.Println(string(out))

	if !ok {
		return fmt.Errorf("validation failed")
	}
	return nil
}

func runSddSpecMerge(cmd *cobra.Command, args []string) error {
	var changeDir string
	var err error
	if len(args) > 0 {
		changeDir, err = resolveChangeSet(args[0])
	} else {
		changeDir, err = resolveChangeSet("")
	}
	if err != nil {
		return err
	}

	repoRootAbs, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("resolving working directory: %w", err)
	}

	specsDir, err := resolveCanonicalSpecsDir(repoRootAbs, flagSpecMergeDir)
	if err != nil {
		return err
	}

	summary, err := sdd.MergeChangeSpecs(sdd.MergeOptions{
		RepoRootAbs: repoRootAbs,
		ChangeName:  filepath.Base(changeDir),
		DryRun:      flagSpecMergeDryRun,
		SpecsDir:    specsDir,
	})
	if err != nil {
		return err
	}

	payload, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding merge output: %w", err)
	}
	fmt.Println(string(payload))
	return nil
}

func runSddSpecList(cmd *cobra.Command, args []string) error {
	repoRootAbs, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("resolving working directory: %w", err)
	}

	specsDir, err := resolveCanonicalSpecsDir(repoRootAbs, flagSpecListDir)
	if err != nil {
		return err
	}

	paths, err := sdd.ListCanonicalSpecs(repoRootAbs, specsDir)
	if err != nil {
		return err
	}

	payload := struct {
		SpecsDir string   `json:"specsDir"`
		Count    int      `json:"count"`
		Paths    []string `json:"paths"`
	}{
		SpecsDir: specsDir,
		Count:    len(paths),
		Paths:    paths,
	}

	out, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding list output: %w", err)
	}
	fmt.Println(string(out))
	return nil
}

type aeConfig struct {
	Schema   string  `json:"$schema,omitempty"`
	SpecRoot *string `json:"specRoot,omitempty"`
}

func resolveCanonicalSpecsDir(repoRootAbs, flagValue string) (string, error) {
	flagValue = strings.TrimSpace(flagValue)
	if flagValue != "" {
		return flagValue, nil
	}

	cfg, err := loadAEConfig(repoRootAbs)
	if err != nil {
		return "", err
	}
	if cfg.SpecRoot == nil {
		return defaultCanonicalSpecsDir, nil
	}

	specRoot := strings.TrimSpace(*cfg.SpecRoot)
	if specRoot == "" {
		return "", fmt.Errorf("invalid .ae-config.json: specRoot must be a non-empty string when set")
	}

	return specRoot, nil
}

func loadAEConfig(repoRootAbs string) (aeConfig, error) {
	configPath := filepath.Join(repoRootAbs, ".ae-config.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return aeConfig{}, nil
		}
		return aeConfig{}, fmt.Errorf("reading %s: %w", filepath.Base(configPath), err)
	}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()

	var cfg aeConfig
	if err := decoder.Decode(&cfg); err != nil {
		return aeConfig{}, fmt.Errorf("parsing %s: %w", filepath.Base(configPath), err)
	}

	var extra json.RawMessage
	if err := decoder.Decode(&extra); err != io.EOF {
		return aeConfig{}, fmt.Errorf("parsing %s: must contain exactly one JSON object", filepath.Base(configPath))
	}

	return cfg, nil
}
