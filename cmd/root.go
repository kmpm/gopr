/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"path/filepath"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var projectsRoot string
var userHome string
var defaultGOPRIVATE string
var defaultGO111MODULE string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goflip",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		viperGetStringP(&projectsRoot, "root")
		viperGetStringP(&defaultGOPRIVATE, "goprivate")
		viperGetStringP(&defaultGO111MODULE, "go111module")
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	home, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}
	userHome = home

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.goflip.yaml)")
	rootCmd.PersistentFlags().StringVar(&projectsRoot, "root", filepath.Join(userHome, ".goflip"), "$HOME/.goflip")
	rootCmd.PersistentFlags().StringVar(&defaultGOPRIVATE, "goprivate", "", "private go repostiories GOPRIVATE")
	rootCmd.PersistentFlags().StringVar(&defaultGO111MODULE, "go111module", "on", "GO111MODULE")

	viper.BindPFlag("goprivate", rootCmd.PersistentFlags().Lookup("goprivate"))
	viper.BindPFlag("root", rootCmd.PersistentFlags().Lookup("root"))
	viper.BindPFlag("go111module", rootCmd.PersistentFlags().Lookup("go111module"))
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".goflip" (without extension).
		viper.AddConfigPath(userHome)
		viper.SetConfigName(".goflip")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

}

func viperGetStringP(p *string, name string) string {
	v := viper.GetString(name)
	if v != "" {
		*p = v
	}
	return *p
}
