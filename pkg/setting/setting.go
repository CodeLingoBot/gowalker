// Copyright 2014 Unknwon
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

package setting

import (
	"time"

	"github.com/Unknwon/com"
	log "gopkg.in/clog.v1"
	"gopkg.in/ini.v1"
	"gopkg.in/macaron.v1"
)

var (
	// Application settings
	AppVer   string
	ProdMode bool

	// Server settings
	HTTPPort     int
	FetchTimeout time.Duration
	DocsJSPath   string
	DocsGobPath  string

	// Global settings
	Cfg               *ini.File
	GitHubCredentials string
	RefreshInterval   = 5 * time.Minute
)

func init() {
	log.New(log.CONSOLE, log.ConsoleConfig{})

	sources := []interface{}{"conf/app.ini"}
	if com.IsFile("custom/app.ini") {
		sources = append(sources, "custom/app.ini")
	}

	var err error
	Cfg, err = macaron.SetConfig(sources[0], sources[1:]...)
	if err != nil {
		log.Fatal(2, "Failed to set configuration: %v", err)
	}

	if Cfg.Section("").Key("RUN_MODE").String() == "prod" {
		ProdMode = true
		macaron.Env = macaron.PROD
		macaron.ColorLog = false

		log.New(log.CONSOLE, log.ConsoleConfig{
			Level:      log.INFO,
			BufferSize: 100,
		})
	}

	sec := Cfg.Section("server")
	HTTPPort = sec.Key("HTTP_PORT").MustInt(8080)
	FetchTimeout = time.Duration(sec.Key("FETCH_TIMEOUT").MustInt(60)) * time.Second
	DocsJSPath = sec.Key("DOCS_JS_PATH").MustString("raw/docs/")
	DocsGobPath = sec.Key("DOCS_GOB_PATH").MustString("raw/gob/")

	GitHubCredentials = "client_id=" + Cfg.Section("github").Key("CLIENT_ID").String() +
		"&client_secret=" + Cfg.Section("github").Key("CLIENT_SECRET").String()

}
