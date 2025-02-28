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
	"log"
	"os"
	"path/filepath"
)

type Writer struct {
	TemplateDirectory string
}

func (w *Writer) WriteTemplate(t *Template, component string) {
	var dir = w.TemplateDirectory

	if component != "" {
		dir = filepath.Join(dir, component)
	}

	// create the desired directory if it doesn't exist yet
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			log.Fatalf("Failed to create directory: %v", err)
		}
	}

	// construct the final path for writing
	template := filepath.Join(dir, t.FileName())

	fmt.Printf("\nCreating a new template in: %s\n", template)

	// error on existing file
	if _, err := os.Stat(template); err == nil {
		log.Fatalf("Template already exists: %s", template)
	}

	// render the template
	content := t.Render()

	// write the template to the file
	if err := os.WriteFile(template, []byte(content), 0644); err != nil {
		log.Fatalf("Failed to write path: %v", err)
	}
}
