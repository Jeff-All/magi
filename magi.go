package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
	"toolbox"

	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

func main() {
	log.SetLevel(log.DebugLevel)

	app := cli.NewApp()
	app.Name = "magi"
	app.Usage = "API for handling physical Donations"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: "config",
			Usage: "The config file to use",
		},
		cli.BoolFlag{
			Name: "local, l",
		},
	}

	app.Action = func(c *cli.Context) error {
		viper.SetConfigName(c.String("config"))
		viper.AddConfigPath(".")
		err := viper.ReadInConfig()
		toolbox.FatalError("unable to read config", err)

		fmt.Printf("%s\n", viper.GetString("test"))

		LaunchServer(c.Bool("local"))

		return nil
	}

	err := app.Run(os.Args)
	toolbox.FatalError("unable to launch the app", err)
}

func BuildAddress(
	local bool,
) string {
	toReturn := ""
	if local {
		toReturn = "localhost"
	}
	toReturn += ":" + viper.GetString("server.port")
	log.WithFields(log.Fields{
		"local":   local,
		"address": toReturn,
	}).Debug("BuildAddress")
	return toReturn
}

func BuildServer(
	r *mux.Router,
	local bool,
) *http.Server {
	logrus.WithFields(log.Fields{
		"local":          local,
		"readTimeout":    viper.GetInt("server.readTimeout"),
		"writeTimeout":   viper.GetInt("server.writeTimeout"),
		"maxHeaderBytes": viper.GetInt("server.maxHeaderBytes"),
		"port":           viper.GetInt("server.port"),
	}).Debug("BuildServer")
	return &http.Server{
		Addr:           BuildAddress(local),
		Handler:        r,
		ReadTimeout:    (time.Duration)(viper.GetInt("server.readTimeout")) * time.Second,
		WriteTimeout:   (time.Duration)(viper.GetInt("server.writeTimeout")) * time.Second,
		MaxHeaderBytes: viper.GetInt("server.maxHeaderBytes"),
	}
}

func LaunchServer(
	local bool,
) {
	r := mux.NewRouter()

	ConfigureRoutes(r)

	s := BuildServer(r, local)

	http.Handle("/", r)

	log.Fatal(s.ListenAndServe())
}

func ConfigureRoutes(
	r *mux.Router,
) {

	r.HandleFunc("/gifts", func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		log.Debugf("'/gifts' GET")
	}).Methods("GET")

	r.HandleFunc("/gifts", func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		log.Debugf("'/gifts' PUT")
	}).Methods("PUT")
}
