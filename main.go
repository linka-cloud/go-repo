package main

import (
	"net/http"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	modules = Modules{}
)

func main() {
	var address string
	var level string
	cmd := &cobra.Command{
		Use:        "go-repo [config]",
		Short:      "go-repo serve a minimal go module registry web site",
		Example:    "go-repo config.yml",
		Args:       cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if l, err := logrus.ParseLevel(level); err == nil {
				logrus.SetLevel(l)
			}
			mods, err := NewModules(args[0])
			if err != nil {
			    logrus.Fatal(err)
			}
			mu.Lock()
			modules = mods
			mu.Unlock()

			watcher, err := fsnotify.NewWatcher()
			if err != nil {
				logrus.Fatal(err)
			}
			defer watcher.Close()
			go func() {
				for {
					select {
					case event, ok := <-watcher.Events:
						if !ok {
							return
						}
						logrus.Debugf("file watcher: %s", event.Op.String())
						if event.Op&fsnotify.Write != fsnotify.Write {
							continue
						}
						logrus.Info("reloading modules")
						mu.Lock()
						modules, err = NewModules(args[0])
						if err != nil {
							mu.Unlock()
							logrus.Fatalf("reload modules: %v", err)
						}
						mu.Unlock()
						go func() {
							if err := modules.LoadReadme(); err != nil {
								logrus.WithError(err).Error("failed to load all readme")
							}
						}()
					case err, ok := <-watcher.Errors:
						if !ok {
							return
						}
						logrus.Errorf("file watcher: %v", err)
					}
				}
			}()
			if err := watcher.Add(args[0]); err != nil {
				logrus.Fatal(err)
			}
			go func() {
				if err := modules.LoadReadme(); err != nil {
					logrus.WithError(err).Error("failed to load all readme")
				}
			}()
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
