// Copyright 2016 The kingshard Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

//用于通过api保存配置
var configFileName string

//Config 整个config文件对应的结构
type Config struct {
	Addr string `yaml:"addr"`

	Charset string     `yaml:"proxy_charset"`
	Node    NodeConfig `yaml:"node"`
}

//NodeConfig node节点对应的配置
type NodeConfig struct {
	Name       string `yaml:"name"`
	MaxConnNum int    `yaml:"max_conns_limit"`

	User     string `yaml:"user"`
	Password string `yaml:"password"`

	Address string `yaml:"address"`
}

// ParseConfigData parse config data
func ParseConfigData(data []byte) (*Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// ParseConfigFile parse config file
func ParseConfigFile(fileName string) (*Config, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	configFileName = fileName

	return ParseConfigData(data)
}

// WriteConfigFile write config file
func WriteConfigFile(cfg *Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(configFileName, data, 0755)
	if err != nil {
		return err
	}

	return nil
}
