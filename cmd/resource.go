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

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tjespers/helm-scaffold/internal/chart"
	"github.com/tjespers/helm-scaffold/internal/template"
)

// resourceCmd represents the resource command
var resourceCmd = &cobra.Command{
	Use:   "resource",
	Short: "Create a new resource template in the current chart",
	Long: `Adds one or more resources to the templates directory of the current chart

Examples:

# Create a new resource template in the current chart
helm-scaffold resource

# Create multiple resources interactively in the current chart
helm-scaffold resource --multiple
`,
	Run: func(cmd *cobra.Command, args []string) {
		// get the template directory from the config file
		dir := viper.GetString("templatesDir")

		// open the helmchart
		mgr, err := chart.NewManagerForChartLocatedInDir(dir, cmd.Flag("component").Value.String())

		if err != nil {
			panic(err)
		}

		// create the catalog
		templateCatalog := template.FromDirectory(viper.GetString("templatesDir"))

		var resources []*template.Template

		if cmd.Flag("multi").Value.String() == "true" {
			resources = templateCatalog.SelectMultiple()
		} else {
			template := templateCatalog.SelectOne()
			resources = append(resources, template)
		}

		// add the resources to the chart
		for _, resource := range resources {
			mgr.AddResource(resource)
		}
	},
}

func init() {
	rootCmd.AddCommand(resourceCmd)

	resourceCmd.Flags().StringP("component", "c", "", "Component to generate the new resource(s) under")
	resourceCmd.Flags().BoolP("multi", "m", false, "Create multiple resources interactively")
}
