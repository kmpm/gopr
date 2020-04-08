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

package project

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// Config contains project specific config
type Config struct {
	Go111Module bool              `yaml:"go111module"`
	GoPrivate   string            `yaml:"goprivate"`
	Env         map[string]string `yaml:"env,flow"`
}

//ReadConfig creates a *Config from a yaml file
func ReadConfig(filename string) (*Config, error) {
	if _, err := os.Stat(filename); err != nil {
		return nil, err
	}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	c := &Config{}
	err = yaml.Unmarshal(data, c)
	// fmt.Printf("Project Config: %+v\n", c)
	return c, err
}

//WriteConfig to save config to yaml file
func WriteConfig(c *Config, filename string) error {
	// err := touch(filename)
	// if err != nil {
	// 	return err
	// }

	out, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, out, 0644)
}
