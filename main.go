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

			handler := handlers.LoggingHandler(os.Stdout, http.HandlerFunc(modulesHandler))
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
