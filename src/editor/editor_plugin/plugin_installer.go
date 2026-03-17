package editor_plugin

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"kaijuengine.com/platform/filesystem"
)

// GitPluginStorage manages Git plugin URLs
type GitPluginStorage struct {
	Plugins []string `json:"plugins"`
}

// getGitPluginsStoragePath returns the path to the Git plugins storage file
func getGitPluginsStoragePath() (string, error) {
	dir, err := filesystem.GameDirectory()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "git_plugins.json"), nil
}

// loadGitPlugins loads the list of Git plugins from storage
func loadGitPlugins() ([]string, error) {
	storageFile, err := getGitPluginsStoragePath()
	if err != nil {
		return nil, err
	}

	if !filesystem.FileExists(storageFile) {
		return []string{}, nil // No plugins stored yet
	}

	data, err := filesystem.ReadFile(storageFile)
	if err != nil {
		return nil, err
	}

	var storage GitPluginStorage
	if err := json.Unmarshal(data, &storage); err != nil {
		return nil, err
	}

	return storage.Plugins, nil
}

// saveGitPlugins saves the list of Git plugins to storage
func saveGitPlugins(plugins []string) error {
	storageFile, err := getGitPluginsStoragePath()
	if err != nil {
		return err
	}

	storage := GitPluginStorage{Plugins: plugins}
	data, err := json.MarshalIndent(storage, "", "  ")
	if err != nil {
		return err
	}

	return filesystem.WriteFile(storageFile, data)
}

// AddGitPluginToStorage adds a Git plugin to persistent storage
func AddGitPluginToStorage(modulePath string) error {
	plugins, err := loadGitPlugins()
	if err != nil {
		return err
	}

	// Check if already exists
	for _, plugin := range plugins {
		if plugin == modulePath {
			return nil // Already exists
		}
	}

	// Add new plugin
	plugins = append(plugins, modulePath)
	return saveGitPlugins(plugins)
}

// RemoveGitPluginFromStorage removes a Git plugin from persistent storage
func RemoveGitPluginFromStorage(modulePath string) error {
	plugins, err := loadGitPlugins()
	if err != nil {
		return err
	}

	// Find and remove the plugin
	for i, plugin := range plugins {
		if plugin == modulePath {
			plugins = append(plugins[:i], plugins[i+1:]...)
			break
		}
	}

	return saveGitPlugins(plugins)
}

// GetStoredGitPlugins returns all Git plugins from storage
func GetStoredGitPlugins() ([]string, error) {
	return loadGitPlugins()
}

// parseGitURL extracts module path and reference (branch/tag) from Git URL
func parseGitURL(gitURL string) (modulePath, ref string, err error) {
	// Clean up URL - remove trailing slashes and fragments
	gitURL = strings.TrimSuffix(gitURL, "/")
	if idx := strings.Index(gitURL, "#"); idx != -1 {
		gitURL = gitURL[:idx]
	}

	// Remove protocol prefix
	cleanURL := gitURL
	if strings.HasPrefix(cleanURL, "https://") {
		cleanURL = cleanURL[8:] // Remove "https://"
	} else if strings.HasPrefix(cleanURL, "http://") {
		return "", "", fmt.Errorf("HTTP URLs are not supported for security reasons")
	} else if strings.HasPrefix(cleanURL, "git@") {
		// SSH format: git@domain.com:owner/repo.git
		cleanURL = strings.TrimPrefix(cleanURL, "git@")
		cleanURL = strings.Replace(cleanURL, ":", "/", 1)
	}

	// Remove .git suffix if present
	cleanURL = strings.TrimSuffix(cleanURL, ".git")

	// Simple pattern matching for common Git providers
	// Expected format: domain.com/owner/repo
	// Examples:
	// github.com/JrSchmidtt/kaiju-fps-counter
	// gitlab.com/owner/repo
	// bitbucket.org/owner/repo

	parts := strings.Split(cleanURL, "/")
	if len(parts) < 3 {
		return "", "", fmt.Errorf("invalid Git URL format: expected domain/owner/repo, got %s", cleanURL)
	}

	// Basic format: domain/owner/repo
	if len(parts) == 3 {
		return cleanURL, "latest", nil
	}

	// Handle special paths like /tree/branch, /releases/tag/, etc.
	if len(parts) >= 5 {
		baseRepo := strings.Join(parts[:3], "/")
		if parts[3] == "tree" || parts[3] == "src" {
			// github.com/owner/repo/tree/branch or bitbucket.org/owner/repo/src/branch
			return baseRepo, parts[4], nil
		} else if parts[3] == "releases" && parts[4] == "tag" && len(parts) >= 6 {
			// github.com/owner/repo/releases/tag/v1.0.0
			return baseRepo, parts[5], nil
		} else if parts[3] == "-" && parts[4] == "tree" {
			// gitlab.com/owner/repo/-/tree/branch
			return baseRepo, parts[5], nil
		}
	}

	// Default to basic repo path with latest
	baseRepo := strings.Join(parts[:3], "/")
	return baseRepo, "latest", nil
}

// AddPluginFromGit adds a plugin from Git URL to storage and returns module path
func AddPluginFromGit(gitURL string) (string, error) {
	// Parse the Git URL
	modulePath, ref, err := parseGitURL(gitURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse Git URL: %w", err)
	}

	var fullModuleRef string

	// Check if this is the Kaiju Engine repository - use local path instead
	if strings.Contains(modulePath, "github.com/KaijuEngine/kaiju") {
		// Use local path for Kaiju Engine repository
		fullModuleRef = "kaijuengine.com"
	} else {
		// For Go modules, use latest version to avoid pseudo-version issues
		version := "latest"
		if ref != "latest" && ref != "main" && ref != "master" {
			// For explicit version tags only, use semantic versioning
			tagPattern := regexp.MustCompile(`^v?\d+\.\d+\.\d+$`)
			if tagPattern.MatchString(ref) {
				// It's a proper semantic version tag
				if !strings.HasPrefix(ref, "v") {
					version = "v" + ref
				} else {
					version = ref
				}
			} else {
				// For branches or any other refs, always use latest to avoid issues
				version = "latest"
			}
		}

		// Create full module reference
		fullModuleRef = fmt.Sprintf("%s@%s", modulePath, version)
	}

	// Save to persistent storage
	if err := AddGitPluginToStorage(fullModuleRef); err != nil {
		return "", fmt.Errorf("failed to save Git plugin to storage: %w", err)
	}

	return fullModuleRef, nil
}

// AddPluginFromGitHub is a wrapper for backward compatibility
// Deprecated: Use AddPluginFromGit instead
func AddPluginFromGitHub(githubURL string) (string, error) {
	return AddPluginFromGit(githubURL)
}

// ModuleDependency represents a Go module dependency
type ModuleDependency struct {
	Module  string
	Version string
}

// parseGoModFile parses a go.mod file and extracts external dependencies
func parseGoModFile(goModPath string) ([]ModuleDependency, error) {
	file, err := os.Open(goModPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open go.mod file: %w", err)
	}
	defer file.Close()

	var dependencies []ModuleDependency
	scanner := bufio.NewScanner(file)
	inRequireBlock := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip comments and empty lines
		if strings.HasPrefix(line, "//") || line == "" {
			continue
		}

		// Handle require block
		if strings.HasPrefix(line, "require (") {
			inRequireBlock = true
			continue
		}
		if inRequireBlock && strings.Contains(line, ")") {
			inRequireBlock = false
			continue
		}

		// Parse dependencies
		if inRequireBlock || strings.HasPrefix(line, "require ") {
			// Remove "require " prefix if present
			if after, ok :=strings.CutPrefix(line, "require "); ok  {
				line = after
			}

			// Parse module and version
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				module := parts[0]
				version := parts[1]

				// Skip kaijuengine.com modules (internal)
				if !strings.HasPrefix(module, "kaijuengine.com") {
					dependencies = append(dependencies, ModuleDependency{
						Module:  module,
						Version: version,
					})
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading go.mod file: %w", err)
	}

	return dependencies, nil
}

// GetPluginDependencies analyzes a plugin directory and returns its external dependencies
func GetPluginDependencies(pluginPath string) ([]ModuleDependency, error) {
	goModPath := filepath.Join(pluginPath, "go.mod")

	// Check if go.mod exists
	if !filesystem.FileExists(goModPath) {
		return []ModuleDependency{}, nil // No dependencies
	}

	return parseGoModFile(goModPath)
}

// GetAllPluginsDependencies gets all external dependencies from a list of plugins
func GetAllPluginsDependencies(plugins []PluginInfo) ([]ModuleDependency, error) {
	allDeps := make(map[string]string) // module -> version

	for _, plugin := range plugins {
		if !plugin.Config.Enabled {
			continue
		}

		deps, err := GetPluginDependencies(plugin.Path)
		if err != nil {
			return nil, fmt.Errorf("failed to get dependencies for plugin %s: %w", plugin.Config.Name, err)
		}

		// Merge dependencies, handling version conflicts
		for _, dep := range deps {
			if existingVersion, exists := allDeps[dep.Module]; exists {
				// For now, use the later version lexicographically
				// In the future, we could implement proper semantic versioning
				if dep.Version > existingVersion {
					allDeps[dep.Module] = dep.Version
				}
			} else {
				allDeps[dep.Module] = dep.Version
			}
		}
	}

	// Convert map back to slice
	var result []ModuleDependency
	for module, version := range allDeps {
		result = append(result, ModuleDependency{
			Module:  module,
			Version: version,
		})
	}

	return result, nil
}

// AddDependenciesToGoMod adds dependencies to a go.mod file
func AddDependenciesToGoMod(goModPath string, dependencies []ModuleDependency) error {
	if len(dependencies) == 0 {
		return nil // Nothing to add
	}

	// Read existing go.mod
	data, err := filesystem.ReadFile(goModPath)
	if err != nil {
		return fmt.Errorf("failed to read go.mod: %w", err)
	}

	content := string(data)
	lines := strings.Split(content, "\n")

	// Find where to insert the require block or add to existing one
	requireBlockExists := false
	requireBlockEnd := -1

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "require (") {
			requireBlockExists = true
			// Find the end of the require block
			for j := i + 1; j < len(lines); j++ {
				if strings.Contains(strings.TrimSpace(lines[j]), ")") {
					requireBlockEnd = j
					break
				}
			}
			break
		}
	}

	// Build dependency lines
	var depLines []string
	for _, dep := range dependencies {
		depLines = append(depLines, fmt.Sprintf("\t%s %s", dep.Module, dep.Version))
	}

	var newLines []string
	if requireBlockExists {
		// Insert into existing require block
		newLines = append(newLines, lines[:requireBlockEnd]...)
		newLines = append(newLines, depLines...)
		newLines = append(newLines, lines[requireBlockEnd:]...)
	} else {
		// Add new require block at the end
		newLines = append(newLines, lines...)
		newLines = append(newLines, "")
		newLines = append(newLines, "require (")
		newLines = append(newLines, depLines...)
		newLines = append(newLines, ")")
	}

	// Write back to file
	newContent := strings.Join(newLines, "\n")
	return filesystem.WriteFile(goModPath, []byte(newContent))
}

// CreateCleanGoMod creates a clean go.mod file from scratch with all dependencies and replaces
func CreateCleanGoMod(goModPath string, pluginPackages []string, dependencies []ModuleDependency) error {
	var content strings.Builder

	// Write module declaration
	content.WriteString("module kaijuengine.com\n\n")
	content.WriteString("go 1.25.0\n\n")

	// Write require block
	if len(dependencies) > 0 || len(pluginPackages) > 0 {
		content.WriteString("require (\n")

		// Add external dependencies
		for _, dep := range dependencies {
			content.WriteString(fmt.Sprintf("\t%s %s\n", dep.Module, dep.Version))
		}

		// Add plugin dependencies
		for _, pluginPkg := range pluginPackages {
			pluginImport := fmt.Sprintf("kaijuengine.com/editor/editor_plugin/developer_plugins/%s", pluginPkg)
			content.WriteString(fmt.Sprintf("\t%s v0.0.0-00010101000000-000000000000\n", pluginImport))
		}

		content.WriteString(")\n\n")
	}

	// Write replace directives
	content.WriteString("// Replace with local development version\n")
	content.WriteString("replace kaijuengine.com => .\n")

	// Add plugin replace directives
	for _, pluginPkg := range pluginPackages {
		pluginImport := fmt.Sprintf("kaijuengine.com/editor/editor_plugin/developer_plugins/%s", pluginPkg)
		content.WriteString(fmt.Sprintf("replace %s => ./editor/editor_plugin/developer_plugins/%s\n", pluginImport, pluginPkg))
	}

	// Write to file
	return filesystem.WriteFile(goModPath, []byte(content.String()))
}
