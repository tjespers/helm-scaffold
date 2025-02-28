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
	"github.com/manifoldco/promptui"
	"github.com/tjespers/helm-scaffold/internal/variables"
	"log"
	"regexp"
	"slices"
)

type Parser struct {
	config *ParserConfig
}

type ParserConfig struct {
	SearchPattern       string
	DefaultReplacements map[string]string
}

func NewParserWithConfig(config *ParserConfig) *Parser {
	return &Parser{
		config: config,
	}
}

// parseVariables extracts all variables from the given content
func (p *Parser) parseVariables(content string) (vars []*variables.Variable) {
	pattern := regexp.MustCompile(p.config.SearchPattern)
	matches := pattern.FindAllString(content, -1)
	var patterns []string

	for _, match := range matches {
		// if we already have a variable with the same name, skip it
		if slices.Contains(patterns, pattern.FindStringSubmatch(match)[0]) {
			continue
		}
		patterns = append(patterns, pattern.FindStringSubmatch(match)[0])
		vars = append(vars, &variables.Variable{
			Pattern: pattern.FindStringSubmatch(match)[0],
			Name:    pattern.FindStringSubmatch(match)[1],
			Value:   "",
		})
	}

	return vars
}

func (p *Parser) getDefaultValueForVariable(variable string) (value string, found bool) {
	if value, found = p.config.DefaultReplacements[variable]; found {
		return value, true
	}

	return "", false
}

// Parse takes a template and replaces all variables with values either from defaults or user input
func (p *Parser) Parse(template *Template) {
	vars := p.parseVariables(template.String())

	header := fmt.Sprintf("Supply missing values for variables found in template: %s", template.FileName())
	headerPrinted := false

	// for each variable in the template we need to come up with a value
	for _, variable := range vars {
		// for any variables not replaced yet prompt the user for a value
		if variable.Value == "" {

			// if a default value is present for the variable, use it
			if defaultValue, ok := p.getDefaultValueForVariable(variable.Name); ok {
				variable.Value = defaultValue
				continue
			}

			// if not prompt the user for it
			if !headerPrinted {
				fmt.Println(header)
				headerPrinted = true
			}
			value := p.promptUserForValue(template, variable.Name)

			variable.Value = value
		}
	}
	template.SetVariables(vars)
}

// promptUserForValue prompts the user for a value for the given variable
func (p *Parser) promptUserForValue(template *Template, name string) string {
	// if not prompt the user for it
	prompt := promptui.Prompt{
		Label: name,
	}

	answer, err := prompt.Run()

	if err != nil {
		log.Fatalf("Prompt failed: %v", err)
	}

	return answer
}
