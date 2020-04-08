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
	"os"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

const (
	//envTmpl = `{{ .Prefix }}DOCKER_TLS_VERIFY{{ .Delimiter }}{{ .DockerTLSVerify }}{{ .Suffix }}{{ .Prefix }}DOCKER_HOST{{ .Delimiter }}{{ .DockerHost }}{{ .Suffix }}{{ .Prefix }}DOCKER_CERT_PATH{{ .Delimiter }}{{ .DockerCertPath }}{{ .Suffix }}{{ .Prefix }}DOCKER_MACHINE_NAME{{ .Delimiter }}{{ .MachineName }}{{ .Suffix }}{{ if .ComposePathsVar }}{{ .Prefix }}COMPOSE_CONVERT_WINDOWS_PATHS{{ .Delimiter }}true{{ .Suffix }}{{end}}{{ if .NoProxyVar }}{{ .Prefix }}{{ .NoProxyVar }}{{ .Delimiter }}{{ .NoProxyValue }}{{ .Suffix }}{{end}}{{ .UsageHint }}`
	//envTmpl contains the template to show
	envTmpl = `{{ .Prefix }}GOPATH{{ .Delimiter }}{{ .GoPath }}{{ .Suffix }}{{ .Prefix }}GO111MODULE{{ .Delimiter }}{{ .Go111Module }}{{ .Suffix }}{{ .Prefix }}GOPRIVATE{{ .Delimiter }}{{ .GoPrivate }}{{ .Suffix }}{{.Prefix}}PATH{{.Delimiter}}{{.Path}}{{.Suffix}}{{.Comment}}
{{ range $key, $value := .Env }}{{$.Prefix}}{{$key}}{{$.Delimiter}}{{$value}}{{$.Suffix}}{{end}}{{.Comment}}
{{ .UsageHint }}`
)

var (
	userShell          string
	defaultUsageHinter UsageHintGenerator
)

// envCmd represents the env command
var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Display commands to set up environment for the go project",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			exitOn("parameter error", ErrInvalidProjectName)
		}
		projectName := args[0]

		found, err := projectExists(projectName)
		exitOn("Can not list projects", err)

		if !found {
			exitOn("Invalid project", fmt.Errorf("project not in list"))
		}

		cfg, err := shellCfgSet(projectName)
		exitOn("Unexpected error", err)

		pc, err := readProjectConfig(cfg.ConfigFile)
		if err == nil {
			cfg.Merge(pc)
		}

		err = executeTemplateStdout(cfg)
		exitOn("Unexpected error", err)

	},
}

func init() {
	rootCmd.AddCommand(envCmd)
	defaultUsageHinter = &EnvUsageHintGenerator{}
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// envCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// envCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	envCmd.Flags().StringVar(&userShell, "shell", "", "set custom shell")
}

func executeTemplateStdout(shellCfg *shellConfig) error {
	t := template.New("envConfig")
	tmpl, err := t.Parse(envTmpl)
	if err != nil {
		return err
	}
	return tmpl.Execute(os.Stdout, shellCfg)
}

type UsageHintGenerator interface {
	GenerateUsageHint(string, []string) string
}

type EnvUsageHintGenerator struct{}

func (g *EnvUsageHintGenerator) GenerateUsageHint(userShell string, args []string) string {
	cmd := ""
	comment := "#"

	projectPath := args[0]
	if strings.Contains(projectPath, " ") || strings.Contains(projectPath, `\`) {
		args[0] = fmt.Sprintf("\"%s\"", projectPath)
	}

	commandLine := strings.Join(args, " ")

	switch userShell {
	case "fish":
		cmd = fmt.Sprintf("eval (%s)", commandLine)
	case "powershell":
		cmd = fmt.Sprintf("& %s | Invoke-Expression", commandLine)
	case "cmd":
		cmd = fmt.Sprintf("\t@FOR /f \"tokens=*\" %%i IN ('%s') DO @%%i", commandLine)
		comment = "REM"
	case "emacs":
		cmd = fmt.Sprintf("(with-temp-buffer (shell-command \"%s\" (current-buffer)) (eval-buffer))", commandLine)
		comment = ";;"
	case "tcsh":
		cmd = fmt.Sprintf("eval `%s`", commandLine)
		comment = ":"
	default:
		cmd = fmt.Sprintf("eval $(%s)", commandLine)
	}

	return fmt.Sprintf("%s Run this command to configure your shell: \n%s %s\n", comment, comment, cmd)
}
