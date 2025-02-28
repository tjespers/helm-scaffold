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

package prompt

import (
	"github.com/manifoldco/promptui"
	"strings"
)

// Ask asks the user a question and returns the answer as a string
func Ask(question string) (string, error) {
	prompt := promptui.Prompt{
		Label: question,
	}

	return prompt.Run()
}

// Choose asks the user to choose from a list of items and returns the selected item as a string
func Choose(question string, items []string) (choice string, err error) {
	prompt := promptui.Select{
		Label:             question,
		Items:             items,
		StartInSearchMode: true,
		Searcher: func(input string, index int) bool {
			template := items[index]
			name := strings.Replace(strings.ToLower(template), " ", "", -1)
			input = strings.Replace(strings.ToLower(input), " ", "", -1)
			return strings.Contains(name, input)
		},
	}

	_, choice, err = prompt.Run()

	return
}

func ChooseOrNew(question string, items []string) (choice string, err error) {
	prompt := promptui.SelectWithAdd{
		Label:    question,
		Items:    items,
		AddLabel: "New",
	}

	_, choice, err = prompt.Run()

	return
}
