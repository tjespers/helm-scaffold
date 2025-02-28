/*
 * Copyright 2025 Tim Jespers
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at:
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package template

import (
	"fmt"
	prompt "github.com/tjespers/helm-scaffold/internal/prompt"
	"log"
	"os"
	"path/filepath"
	"slices"
)

type Catalog struct {
	directory string
	templates []string
}

func FromDirectory(templateDirectory string) *Catalog {
	templates, err := findTemplateFiles(templateDirectory)

	if err != nil {
		log.Fatalf("Failed to parse templates directory: %v", err)
	}

	return &Catalog{
		directory: templateDirectory,
		templates: templates,
	}
}

func findTemplateFiles(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip directories
		if info.IsDir() {
			return nil
		}
		// Check file extension
		ext := filepath.Ext(path)
		if slices.Contains([]string{".yaml", ".yml", ".tpl"}, ext) {
			// Store relative path from root
			relPath, err := filepath.Rel(root, path)
			if err != nil {
				return fmt.Errorf("failed to get relative path: %w", err)
			}
			files = append(files, relPath)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}
	return files, nil
}

func (t *Catalog) selectTemplate(templates []string, allowCancel bool) *Template {
	cancel := "⬅️ Cancel"

	if allowCancel {
		templates = append([]string{cancel}, templates...)
	}

	choice, err := prompt.Choose("Choose Template to add, press Enter to finish", templates)

	if err != nil {
		log.Fatalf("Template selection failed: %v", err)
	}

	if choice == cancel {
		return nil
	}

	return NewTemplateFromFile(filepath.Join(t.directory, choice))
}

// SelectMultiple prompts the user to select multiple templates from the current catalog
func (t *Catalog) SelectMultiple() (selectedTemplates []*Template) {
	availableTemplates := t.templates

	for {
		selectedTemplate := t.selectTemplate(availableTemplates, true)

		if selectedTemplate == nil {
			break
		}

		selectedTemplates = append(selectedTemplates, selectedTemplate)

		// Remove the selected template from the list of available templates
		for i, template := range availableTemplates {
			if template == selectedTemplate.FileName() {
				availableTemplates = append(availableTemplates[:i], availableTemplates[i+1:]...)
				break
			}
		}

		// If no more templates are available, stop prompting
		if len(availableTemplates) == 0 {
			break
		}
	}

	return selectedTemplates
}

func (t *Catalog) SelectOne() *Template {
	return t.selectTemplate(t.templates, false)
}

func (t *Catalog) Read(templateName string) []byte {
	// parse the selected template path
	templatePath := filepath.Join(t.directory, templateName)

	content, err := os.ReadFile(templatePath)
	if err != nil {
		log.Fatalf("Failed to parse template path: %v", err)
	}

	return content
}
