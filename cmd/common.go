package cmd

import (
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
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

type projectConfig struct {
	Go111Module string `yaml:"go111module"`
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

func readProjectConfig(configFile string) (*projectConfig, error) {
	if _, err := os.Stat(configFile); err != nil {
		return nil, err
	}
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	var c projectConfig
	err = yaml.Unmarshal(data, &c)
	fmt.Println(data, err, c)
	return &c, err
}

func (shellCfg *shellConfig) Merge(p *projectConfig) {
	fmt.Println(p)
	if p.Go111Module != "NULL" {
		shellCfg.Go111Module = p.Go111Module
	}
	if p.Go111Module != "NULL" {
		shellCfg.GoPrivate = p.GoPrivate
	}
}
