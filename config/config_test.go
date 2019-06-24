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
	"fmt"
	"reflect"
	"testing"
)

func TestConfig(t *testing.T) {
	var testConfigData = []byte(
		`
addr : 0.0.0.0:9696
user : root
password : root

node :
  name : node1 
  max_conns_limit : 16
  user: root
  password: root
  address : 127.0.0.1:3306
`)

	cfg, err := ParseConfigData(testConfigData)
	if err != nil {
		t.Fatal(err)
	}

	testNode := NodeConfig{
		Name:       "node1",
		MaxConnNum: 16,

		User:     "root",
		Password: "root",

		Address: "127.0.0.1:3306",
	}

	if !reflect.DeepEqual(cfg.Node, testNode) {
		fmt.Printf("%v\n", cfg.Node)
		t.Fatal("node1 must equal")
	}
}
