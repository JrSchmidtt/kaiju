package editor_plugin

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"kaijuengine.com/platform/filesystem"
)

func AddGitPluginToStorage(modulePath string) error {
	plugFolder, err := PluginsFolder()
	if err != nil {
		return err
	}

	plugins, err := GetStoredGitPlugins()
	if err != nil {
		return err
	}
	for _, plugin := range plugins {
		if plugin == modulePath {
			return nil
		}
	}

	module := strings.Split(modulePath, "@")[0]
	parts := strings.Split(module, "/")
	packageName := parts[len(parts)-1]
	folderName := "git_" + strings.NewReplacer("/", "_", "@", "_", ":", "_").Replace(modulePath)
	folderPath := filepath.Join(plugFolder, folderName)

	if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
		return err
	}

	cfg := PluginConfig{
		Name:        fmt.Sprintf("Git Plugin: %s", packageName),
		PackageName: packageName,
		Description: fmt.Sprintf("Git plugin from %s", module),
		Version:     0.0,
		Author:      "Git Repository",
		Website:     "https://" + module,
		Enabled:     true,
		GitModule:   modulePath,
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return filesystem.WriteFile(filepath.Join(folderPath, pluginConfigFile), data)
}

func RemoveGitPluginFromStorage(modulePath string) error {
	plugFolder, err := PluginsFolder()
	if err != nil {
		return err
	}

	dirs, err := os.ReadDir(plugFolder)
	if err != nil {
		return err
	}

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
		cfgPath := filepath.Join(plugFolder, dir.Name(), pluginConfigFile)
		if s, err := os.Stat(cfgPath); err != nil || s.IsDir() {
			continue
		}
		data, err := filesystem.ReadFile(cfgPath)
		if err != nil {
			continue
		}
		var cfg PluginConfig
		if err := json.Unmarshal(data, &cfg); err != nil {
			continue
		}
		if cfg.GitModule == modulePath {
			return os.RemoveAll(filepath.Join(plugFolder, dir.Name()))
		}
	}

	return nil
}

func GetStoredGitPlugins() ([]string, error) {
	plugFolder, err := PluginsFolder()
	if err != nil {
		return nil, err
	}

	dirs, err := os.ReadDir(plugFolder)
	if err != nil {
		return nil, err
	}

	plugins := make([]string, 0)
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
		cfgPath := filepath.Join(plugFolder, dir.Name(), pluginConfigFile)
		if s, err := os.Stat(cfgPath); err != nil || s.IsDir() {
			continue
		}
		data, err := filesystem.ReadFile(cfgPath)
		if err != nil {
			continue
		}
		var cfg PluginConfig
		if err := json.Unmarshal(data, &cfg); err != nil {
			continue
		}
		if cfg.GitModule != "" {
			plugins = append(plugins, cfg.GitModule)
		}
	}

	return plugins, nil
}

func parseGitURL(gitURL string) (modulePath, ref string) {
	cleanURL := strings.TrimSpace(gitURL)
	if idx := strings.IndexAny(cleanURL, "?#"); idx != -1 {
		cleanURL = cleanURL[:idx]
	}

	if strings.HasPrefix(cleanURL, "git@") {
		cleanURL = strings.TrimPrefix(cleanURL, "git@")
		cleanURL = strings.Replace(cleanURL, ":", "/", 1)
	}
	cleanURL = strings.TrimPrefix(cleanURL, "https://")
	cleanURL = strings.TrimPrefix(cleanURL, "http://")
	cleanURL = strings.TrimPrefix(cleanURL, "git://")

	cleanURL = strings.TrimSuffix(cleanURL, ".git")
	cleanURL = strings.TrimSuffix(cleanURL, "/")

	ref = "latest"
	if idx := strings.LastIndex(cleanURL, "@"); idx != -1 {
		candidate := cleanURL[idx+1:]
		cleanURL = cleanURL[:idx]
		if candidate != "" {
			ref = candidate
		}
	}

	modulePath = cleanURL
	return modulePath, ref
}

func AddPluginFromGit(gitURL string) (string, error) {
	modulePath, ref := parseGitURL(gitURL)

	if strings.Contains(modulePath, "github.com/KaijuEngine/kaiju") {
		modulePath = "kaijuengine.com"
		ref = ""
	}

	fullModuleRef := modulePath
	if ref != "" {
		fullModuleRef = fmt.Sprintf("%s@%s", modulePath, ref)
	}

	if err := AddGitPluginToStorage(fullModuleRef); err != nil {
		return "", fmt.Errorf("failed to save Git plugin to storage: %w", err)
	}

	return fullModuleRef, nil
}

func AddPluginFromGitHub(githubURL string) (string, error) {
	return AddPluginFromGit(githubURL)
}
