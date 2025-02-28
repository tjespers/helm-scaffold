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
	"github.com/tjespers/helm-scaffold/internal/variables"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Template struct {
	path      string
	content   []byte
	pattern   *regexp.Regexp
	variables []*variables.Variable
}

func (t *Template) String() string {
	return string(t.content)
}

func NewTemplateFromFile(file string) *Template {
	// parse the selected template path
	content, err := os.ReadFile(file)

	if err != nil {
		log.Fatalf("Failed to parse template path: %v", err)
	}

	pattern := regexp.MustCompile(`%%(\w+)%%`)

	return &Template{
		path:    file,
		content: content,
		pattern: pattern,
	}
}

func (t *Template) FileName() string {
	return filepath.Base(t.path)
}

func (t *Template) Render() string {
	renderedContent := string(t.content)
	// Replace each pattern with its corresponding value
	for _, variable := range t.variables {
		renderedContent = strings.ReplaceAll(renderedContent, variable.Pattern, variable.Value)
	}
	return renderedContent
}

func (t *Template) SetVariables(vars []*variables.Variable) {
	t.variables = vars
}
