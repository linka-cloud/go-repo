package main

import (
	"net/http"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	modules = Modules{}
)

func main() {
	var address string
	var level string
	cmd := &cobra.Command{
		Use:     "go-repo [config]",
		Short:   "go-repo serve a minimal go module registry web site",
		Example: "go-repo config.yml",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if l, err := logrus.ParseLevel(level); err == nil {
				logrus.SetLevel(l)
			}
			viper.SetConfigType("yaml")
			viper.SetConfigFile(args[0])
			if err := viper.ReadInConfig(); err != nil {
				logrus.Fatal(err)
			}
			conf := &Config{}
			if err := viper.Unmarshal(conf); err != nil {
				mu.Unlock()
				logrus.Fatal(err)
			}
			mu.Lock()
			modules = conf.Modules
			mu.Unlock()
			go func() {
				if err := modules.LoadReadme(); err != nil {
					logrus.WithError(err).Error("failed to load all readme")
				}
			}()
			viper.WatchConfig()
			viper.OnConfigChange(func(event fsnotify.Event) {
				mu.Lock()
				var err error
				modules, err = NewModules(event.Name)
				if err != nil {
					logrus.Fatal(err)
				}
				mu.Unlock()
				logrus.Info("modules reloaded")
				go func() {
					if err := modules.LoadReadme(); err != nil {
						logrus.WithError(err).Error("failed to load all readme")
					}
				}()
			})

			http.HandleFunc("/", modulesHandler)
			logrus.Infof("starting server at %s", address)
			if err := http.ListenAndServe(address, nil); err != nil {
				logrus.Fatal(err)
			}
		},
	}
	cmd.Flags().StringVarP(&address, "address", "a", ":8888", "The server address")
	cmd.Flags().StringVar(&level, "logs-level", "info", "")
	cmd.Execute()
}
