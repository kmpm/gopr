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
	"errors"
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kmpm/gopr/lib/project"
)

const (
	projectConfigFile string = "project.yaml"
)

type shellConfig struct {
	Prefix      string
	Delimiter   string
	Suffix      string
	Comment     string
	ProjectPath string
	ConfigFile  string
	GoPath      string
	GoPrivate   string
	Go111Module string
	Path        string
	UsageHint   string
	Env         map[string]string
}

var (
	// ErrInvalidProjectName - The given name is not valid for projects
	ErrInvalidProjectName = errors.New("invalid project name")
)

func shellCfgSet(projectName string) (*shellConfig, error) {
	userShell, err := getShell(userShell)
	if err != nil {
		return nil, err
	}
	projectpath := filepath.Join(projectsRoot, projectName)
	gopath := filepath.Join(projectpath, "go")
	//get current
	oldpath := os.Getenv("GOPATH")
	if oldpath == "" {
		oldpath = build.Default.GOPATH
	}

	searchPath := os.Getenv("PATH")
	newList := []string{
		filepath.Join(gopath, "bin"),
	}
	list := strings.Split(searchPath, string(os.PathListSeparator))
	for _, p := range list {
		if !strings.Contains(p, oldpath) {
			newList = append(newList, p)
		}
	}
	searchPath = strings.Join(newList, string(os.PathListSeparator))

	shellCfg := &shellConfig{
		UsageHint:   defaultUsageHinter.GenerateUsageHint(userShell, os.Args),
		Path:        searchPath,
		ProjectPath: projectpath,
		ConfigFile:  filepath.Join(projectpath, projectConfigFile),
		GoPath:      gopath,
		GoPrivate:   defaultGOPRIVATE,
		Go111Module: defaultGO111MODULE,
		Env:         make(map[string]string),
		Prefix:      "export ",
		Suffix:      "\"\n",
		Delimiter:   "=\"",
		Comment:     "#",
	}

	switch userShell {
	case "powershell":
		shellCfg.Prefix = "$Env:"
		shellCfg.Suffix = "\"\n"
		shellCfg.Delimiter = " = \""
	case "cmd":
		shellCfg.Prefix = "SET "
		shellCfg.Suffix = "\n"
		shellCfg.Delimiter = "="
		shellCfg.Comment = "REM "
	}
	return shellCfg, nil
}

func touch(filename string) error {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer file.Close()
	} else {
		currentTime := time.Now().Local()
		err = os.Chtimes(filename, currentTime, currentTime)
		if err != nil {
			return err
		}
	}
	return nil
}

//projectList returns a list of folders which contents match a certain pattern
func projectList() ([]string, error) {
	pat := filepath.Join(filepath.Join(projectsRoot, "*"), "go")
	files, err := filepath.Glob(pat)
	if err != nil {
		return nil, err
	}

	projects := make([]string, 0, len(files))
	for _, f := range files {
		projects = append(projects, filepath.Base(filepath.Dir(f)))
	}
	return projects, nil
}

func projectExists(projectName string) (bool, error) {
	list, err := projectList()
	if err != nil {
		return false, err
	}
	_, found := find(list, projectName)
	return found, nil
}

func find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

//er shows a message and optional error and then os.Exit(1)
func er(msg string, err error) {
	if err == nil {
		fmt.Println(msg)
	} else {
		fmt.Println(msg, err)
	}
	os.Exit(1)
}

//exitOn exits with message and error IF err != nil
func exitOn(msg string, err error) {
	if err != nil {
		er(msg, err)
	}
}

func (shellCfg *shellConfig) Merge(p *project.Config) {
	// fmt.Println(p)
	if p.Go111Module {
		shellCfg.Go111Module = "on"
	}

	shellCfg.GoPrivate = p.GoPrivate
	for k, v := range p.Env {
		shellCfg.Env[k] = v
	}
}

func (shellCfg *shellConfig) GetProjectConfig() (*project.Config, error) {
	c := &project.Config{
		Go111Module: shellCfg.Go111Module == "on",
		GoPrivate:   shellCfg.GoPrivate,
		Env:         make(map[string]string),
	}
	return c, nil
}

func writeProjectConfig(c *project.Config, filename string) error {
	return project.WriteConfig(c, filename)
}

func readProjectConfig(filename string) (*project.Config, error) {
	return project.ReadConfig(filename)
}
