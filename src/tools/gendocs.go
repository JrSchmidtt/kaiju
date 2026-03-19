package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type DocGenerator struct {
	sourceDir string
	outputDir string
	packages  []string
}

func NewDocGenerator() *DocGenerator {
	return &DocGenerator{
		sourceDir: "src",
		outputDir: "./docs/api",
		packages: []string{
			".", "bootstrap", "engine", "matrix", "klib", "rendering",
			"platform/windowing", "platform/hid", "platform/audio",
			"platform/filesystem", "engine/ui", "registry/shader_data_registry", "debug",
		},
	}
}

func main() {
	generator := NewDocGenerator()
	generator.Generate()
}

func (g *DocGenerator) Generate() {
	fmt.Printf("Generating docs from: %s\n", g.sourceDir)

	if err := os.MkdirAll(g.outputDir, 0755); err != nil {
		fmt.Printf("Error creating output dir: %v\n", err)
		os.Exit(1)
	}

	g.generateIndex()

	successCount := 0
	for _, pkg := range g.packages {
		fmt.Printf("Processing: %s... ", pkg)
		if g.generatePackageDoc(pkg) {
			fmt.Println("✓")
			successCount++
		} else {
			fmt.Println("✗")
		}
	}

	fmt.Printf("Generated! %d/%d packages processed.\n", successCount, len(g.packages))
}

func (g *DocGenerator) generateIndex() {
	indexPath := filepath.Join(g.outputDir, "index.md")
	file, err := os.Create(indexPath)
	if err != nil {
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	writer.WriteString("# Kaiju Engine API Documentation\n\n")
	writer.WriteString("Auto-generated using go doc.\n\n")

	categories := g.categorizePackages()
	for title, pkgs := range categories {
		if len(pkgs) == 0 {
			continue
		}
		writer.WriteString(fmt.Sprintf("## %s\n\n", title))
		for _, pkg := range pkgs {
			link := g.getPackageLink(pkg)
			displayName := g.getDisplayName(pkg)
			writer.WriteString(fmt.Sprintf("- [%s](%s)\n", displayName, link))
		}
		writer.WriteString("\n")
	}
}

func (g *DocGenerator) categorizePackages() map[string][]string {
	categories := map[string][]string{
		"Core":      {},
		"Engine":    {},
		"Platform":  {},
		"Rendering": {},
		"Other":     {},
	}

	for _, pkg := range g.packages {
		switch {
		case g.isCorePackage(pkg):
			categories["Core"] = append(categories["Core"], pkg)
		case strings.HasPrefix(pkg, "engine"):
			categories["Engine"] = append(categories["Engine"], pkg)
		case strings.HasPrefix(pkg, "platform"):
			categories["Platform"] = append(categories["Platform"], pkg)
		case g.isRenderingPackage(pkg):
			categories["Rendering"] = append(categories["Rendering"], pkg)
		default:
			categories["Other"] = append(categories["Other"], pkg)
		}
	}

	return categories
}

func (g *DocGenerator) isCorePackage(pkg string) bool {
	return pkg == "." || pkg == "bootstrap" || pkg == "matrix" || pkg == "klib"
}

func (g *DocGenerator) isRenderingPackage(pkg string) bool {
	return strings.HasPrefix(pkg, "rendering") || strings.HasPrefix(pkg, "registry")
}

func (g *DocGenerator) getPackageLink(pkg string) string {
	if pkg == "." {
		return "root.md"
	}
	return strings.ReplaceAll(pkg, "/", "_") + ".md"
}

func (g *DocGenerator) getDisplayName(pkg string) string {
	if pkg == "." {
		return "root"
	}
	return pkg
}

func (g *DocGenerator) generatePackageDoc(pkg string) bool {
	output, err := g.executeGoDoc(pkg)
	if err != nil {
		return false
	}

	filename := g.getPackageLink(pkg)
	filePath := filepath.Join(g.outputDir, filename)

	file, err := os.Create(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	g.writePackageHeader(writer, pkg)
	g.convertDocToMarkdown(writer, string(output))

	return true
}

func (g *DocGenerator) executeGoDoc(pkg string) ([]byte, error) {
	relativePkg := "./" + pkg
	if pkg == "." {
		relativePkg = "."
	}

	cmd := exec.Command("go", "doc", "-all", relativePkg)
	cmd.Dir = g.sourceDir
	return cmd.Output()
}

func (g *DocGenerator) writePackageHeader(writer *bufio.Writer, pkg string) {
	displayName := g.getDisplayName(pkg)
	writer.WriteString(fmt.Sprintf("# Package %s\n\n", displayName))

	importPath := "kaijuengine.com"
	if pkg != "." {
		importPath += "/" + pkg
	}
	writer.WriteString(fmt.Sprintf("**Import path:** `%s`\n\n", importPath))
}

func (g *DocGenerator) convertDocToMarkdown(writer *bufio.Writer, output string) {
	lines := strings.Split(output, "\n")
	inCodeBlock := false

	for _, line := range lines {
		if strings.HasPrefix(line, "package ") {
			continue
		}

		if g.isCodeLine(line) {
			if !inCodeBlock {
				writer.WriteString("```go\n")
				inCodeBlock = true
			}
			writer.WriteString(strings.TrimPrefix(line, "\t") + "\n")
		} else {
			if inCodeBlock {
				writer.WriteString("```\n\n")
				inCodeBlock = false
			}
			g.writeFormattedLine(writer, line)
		}
	}

	if inCodeBlock {
		writer.WriteString("```\n")
	}
}

func (g *DocGenerator) isCodeLine(line string) bool {
	return strings.HasPrefix(line, "\t") || (len(line) > 0 && strings.HasPrefix(line, "    "))
}

func (g *DocGenerator) writeFormattedLine(writer *bufio.Writer, line string) {
	if g.isHeaderLine(line) {
		writer.WriteString(fmt.Sprintf("## %s\n\n", strings.Title(strings.ToLower(line))))
	} else if g.isDeclarationLine(line) {
		writer.WriteString("### " + line + "\n\n")
	} else {
		writer.WriteString(line + "\n")
	}
}

func (g *DocGenerator) isHeaderLine(line string) bool {
	return line != "" && strings.ToUpper(line) == line && !strings.Contains(line, " ") && len(line) > 1
}

func (g *DocGenerator) isDeclarationLine(line string) bool {
	return strings.HasPrefix(line, "func ") || strings.HasPrefix(line, "type ") ||
		strings.HasPrefix(line, "var ") || strings.HasPrefix(line, "const ")
}
