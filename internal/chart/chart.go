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

package chart

import (
	"fmt"
	"github.com/tjespers/helm-scaffold/internal/prompt"
	"github.com/tjespers/helm-scaffold/internal/template"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"log"
	"os"
	"path/filepath"
)

type ChartManager struct {
	Root             string
	writer           template.Writer
	chart            *chart.Chart
	parser           template.Parser
	currentComponent string
}

func (c ChartManager) AddResource(template *template.Template) {
	var err error
	// if a component is not set, and the chart has multiple components, ask the user to choose one
	if c.currentComponent == "" && len(c.GetComponents()) > 0 {
		choices := append(c.GetComponents(), "None (use root)")
		if c.currentComponent, err = prompt.ChooseOrNew("In which component would you like to add the resource(s)?", choices); err != nil {
			log.Fatalf("Failed to select component: %v", err)
		}
	}

	// parse the given template and add it to the chart
	c.parser.Parse(template)

	c.writer.WriteTemplate(template, c.currentComponent)
}

func (c ChartManager) Name() string {
	return c.chart.Metadata.Name
}

// NewManagerForChartLocatedInDir creates a new ChartManager instance for the given directory, it returns an error if the directory does not contain a helm chart
func NewManagerForChartLocatedInDir(directory, component string) (*ChartManager, error) {
	chart, err := loader.LoadDir(directory)

	if err != nil {
		return nil, fmt.Errorf("failed to load chart: %v", err)
	}

	cfg := &template.ParserConfig{
		SearchPattern: `%%(\w+)%%`,
		DefaultReplacements: map[string]string{
			"CHARTNAME":      chart.Metadata.Name,
			"COMPONENT_NAME": component,
		},
	}

	writer := template.Writer{
		TemplateDirectory: filepath.Join(directory, "templates", component),
	}

	return &ChartManager{
		Root:   directory,
		writer: writer,
		chart:  chart,
		parser: *template.NewParserWithConfig(cfg),
	}, nil
}

// GetComponents returns a list of all components in the chart based on the directory names in the templates directory
func (c ChartManager) GetComponents() []string {
	components := []string{}
	templatesDir := filepath.Join(c.Root, "templates")

	dirs, err := os.ReadDir(templatesDir)
	if err != nil {
		log.Fatalf("Failed to read templates directory: %v", err)
	}

	for _, dir := range dirs {
		if dir.IsDir() {
			components = append(components, dir.Name())
		}
	}

	return components
}
