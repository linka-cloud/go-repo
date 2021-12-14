// Copyright 2021 Linka Cloud  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"net/http"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/handlers"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	modules = Modules{}
)

func loadConfig() error {
	conf := &Config{}
	if err := viper.Unmarshal(conf); err != nil {
		return err
	}
	mu.Lock()
	modules = conf.Modules
	modules.Sort()
	mu.Unlock()
	go func() {
		if err := modules.LoadReadme(); err != nil {
			logrus.WithError(err).Error("failed to load all readme")
		}
	}()
	return nil
}

func main() {
	var (
		address     string
		level       string
		noAccessLog bool
	)
	cmd := &cobra.Command{
		Use:     "go-repo [config]",
		Short:   "go-repo serve a minimal go module registry web site",
		Example: "go-repo config.yml",
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if l, err := logrus.ParseLevel(level); err == nil {
				logrus.SetLevel(l)
			}
			path := "config.yaml"
			if len(args) == 1 {
				path = args[0]
			}
			viper.SetConfigType("yaml")
			viper.SetConfigFile(path)
			if err := viper.ReadInConfig(); err != nil {
				logrus.Fatal(err)
			}
			if err := loadConfig(); err != nil {
				logrus.Fatal(err)
			}
			viper.WatchConfig()
			viper.OnConfigChange(func(event fsnotify.Event) {
				if err := loadConfig(); err != nil {
					logrus.Fatal(err)
				}
				logrus.Info("modules reloaded")
			})

			handler := handlers.CustomLoggingHandler(os.Stdout, http.HandlerFunc(modulesHandler), writeLog)
			if noAccessLog {
				handler = http.HandlerFunc(modulesHandler)
			}
			http.Handle("/", handler)
			logrus.Infof("starting server at %s", address)
			if err := http.ListenAndServe(address, nil); err != nil {
				logrus.Fatal(err)
			}
		},
	}
	cmd.Flags().StringVarP(&address, "address", "a", ":8888", "The server address")
	cmd.Flags().StringVar(&level, "logs-level", "info", "")
	cmd.Flags().BoolVar(&noAccessLog, "no-access-log", false, "Disable web server access logs")
	cmd.Execute()
}
