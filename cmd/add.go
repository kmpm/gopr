/*
Copyright © 2020 Peter Magnusson <code@kmpm.se>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add/Create a go project environment",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatal("You must provide a project name")
			return
		}
		projectName := args[0]

		cfg, err := shellCfgSet(projectName)
		if err != nil {
			fmt.Errorf("Error: %+v\n", err)
			os.Exit(1)
		}
		if _, err := os.Stat(cfg.ProjectPath); !os.IsNotExist(err) {
			fmt.Printf("Project path '%s' exists\n", cfg.ProjectPath)
			os.Exit(1)
		}
		fmt.Println("Creating", cfg.GoPath)
		err = os.MkdirAll(cfg.GoPath, os.ModeDir|os.ModePerm)
		if err != nil {
			fmt.Printf("Error creating %s: %+v\n", cfg.GoPath, err)
			os.Exit(1)
		}

		// if err = touch(cfg.ConfigFile); err != nil {
		// 	fmt.Printf("Could not touch %s: %v\n", cfg.ConfigFile, err)
		// 	os.Exit(1)
		// }

		pc, _ := cfg.GetProjectConfig()
		err = writeProjectConfig(pc, cfg.ConfigFile)
		exitOn("Could not create project configuration", err)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
