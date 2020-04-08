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
	ProjectPath string
	ConfigFile  string
	GoPath      string
	GoPrivate   string
	Go111Module string
	Path        string
	UsageHint   string
}

type ProjectConfig struct {
	Go111Module bool   `yaml:"go111module"`
	GoPrivate   string `yaml:"goprivate"`
}

func shellCfgSet(projectName string) (*shellConfig, error) {
	userShell, err := getShell(userShell)
	if err != nil {
		return nil, err
	}
	projectpath := filepath.Join(projectsRoot, projectName)
	gopath := filepath.Join(projectpath, "go")
	//log.Printf("New GOPATH '%s'", gopath)
	//get current
	oldpath := os.Getenv("GOPATH")
	if oldpath == "" {
		oldpath = build.Default.GOPATH
	}
	//log.Printf("Current GOPATH '%s'", oldpath)

	searchPath := os.Getenv("PATH")
	newList := []string{
		filepath.Join(gopath, "bin"),
	}
	list := strings.Split(searchPath, string(os.PathListSeparator))
	for _, p := range list {
		if !strings.Contains(p, oldpath) {
			newList = append(newList, p)
		} else {
			// log.Printf("Dropped '%s' from PATH", p)
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
	default:
		shellCfg.Prefix = "export "
		shellCfg.Suffix = "\"\n"
		shellCfg.Delimiter = "=\""
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

func projectList() ([]string, error) {
	pat := filepath.Join(filepath.Join(projectsRoot, "*"), projectConfigFile)
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

func er(msg string, err error) {
	if err == nil {
		fmt.Println(msg)
	} else {
		fmt.Println(msg, err)
	}
	os.Exit(1)
}

func exitOn(msg string, err error) {
	if err != nil {
		er(msg, err)
	}
}

func readProjectConfig(configFile string) (*ProjectConfig, error) {
	if _, err := os.Stat(configFile); err != nil {
		return nil, err
	}
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	c := &ProjectConfig{}
	err = yaml.Unmarshal(data, c)
	// fmt.Println(string(data), err, c)
	return c, err
}

func (shellCfg *shellConfig) Merge(p *ProjectConfig) {
	fmt.Println(p)
	if p.Go111Module {
		shellCfg.Go111Module = "on"
	}

	shellCfg.GoPrivate = p.GoPrivate
}

func (shellCfg *shellConfig) GetProjectConfig() (*ProjectConfig, error) {
	c := &ProjectConfig{
		Go111Module: shellCfg.Go111Module == "on",
		GoPrivate:   shellCfg.GoPrivate,
	}
	return c, nil
}

func writeProjectConfig(pc *ProjectConfig, filename string) error {
	err := touch(filename)
	if err != nil {
		return err
	}

	out, err := yaml.Marshal(pc)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, out, os.ModeExclusive)

	// file, err := os.Open(filename)
	// if err != nil {
	// 	return err
	// }
	// defer file.Close()
	// e := yaml.NewEncoder(file)
	// defer e.Close()
	// err = e.Encode(c)
	// return err

}
