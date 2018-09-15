package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"
	"toolbox"

	"github.com/Jeff-All/magi/data"

	"github.com/Jeff-All/magi/middleware"
	"github.com/Jeff-All/magi/models"
	res "github.com/Jeff-All/magi/resources"

	requests "github.com/Jeff-All/magi/endpoints/request"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

func main() {
	log.Printf("Starting...")
	app := cli.NewApp()
	app.Name = "magi"
	app.Usage = "API for handling physical donations"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: "config",
			Usage: "The config file to use",
		},
		cli.BoolFlag{
			Name: "local, l",
		},
		cli.BoolFlag{
			Name: "debug, d",
		},
		cli.BoolFlag{
			Name: "info, i",
		},
	}

	app.Commands = []cli.Command{}

	app.Action = Run

	if err := app.Run(os.Args); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("unable to launch")
	}
}

func Run(c *cli.Context) error {
	log.Printf("Run")
	Common(c)
	// InitEmail(c)
	if err := ConnectDatabase(); err != nil { return err }
	models.DB = res.DB
	r := mux.NewRouter()
	ConfigureRoutes(c, r)
	s := BuildServer(r)
	http.Handle("/", r)
	log.Fatal(s.ListenAndServeTLS("dev.magi.crt", "dev.magi.key"))
	return nil
}

func Common(c *cli.Context) error {
	SetLogLevel(c)
	viper.SetConfigName(c.String("config"))
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	toolbox.FatalError("unable to read config", err)
	return nil
}

// func InitEmail(c *cli.Context) {
// 	password, err := ReadPassword()
// 	if err != nil {
// 		log.Fatal("Unable to read password")
// 	}
// 	mail.Init(
// 		viper.GetString("mail.host"),
// 		viper.GetInt("mail.port"),
// 		viper.GetString("mail.user"),
// 		password,
// 	)
// }

func SetLogLevel(
	c *cli.Context,
) {
	if c.Bool("debug") {
		log.SetLevel(log.DebugLevel)
		log.Info("Log level set to debug")
		return
	}
	if c.Bool("info") {
		log.SetLevel(log.InfoLevel)
		log.Info("Log level set to info")
		return
	}
	log.Info("Log level set to error")
	log.SetLevel(log.ErrorLevel)
}

func BuildServer(
	r *mux.Router,
) *http.Server {
	log.WithFields(log.Fields{
		"readTimeout":    viper.GetInt("server.readTimeout"),
		"writeTimeout":   viper.GetInt("server.writeTimeout"),
		"maxHeaderBytes": viper.GetInt("server.maxHeaderBytes"),
		"domain":         viper.GetString("server.domain"),
	}).Debug("BuildServer")
	return &http.Server{
		Addr:           viper.GetString("server.domain"),
		Handler:        r,
		ReadTimeout:    (time.Duration)(viper.GetInt("server.readTimeout")) * time.Second,
		WriteTimeout:   (time.Duration)(viper.GetInt("server.writeTimeout")) * time.Second,
		MaxHeaderBytes: viper.GetInt("server.maxHeaderBytes"),
	}
}

type Response struct {
	Message string `json:"message"`
}

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

func ConfigureRoutes(
	c *cli.Context,
	r *mux.Router,
) {
	r.Handle("/", http.FileServer(http.Dir(viper.GetString("frontend.views"))))

	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir(viper.GetString("frontend.static")))),
	)

	r.PathPrefix("/node_modules/").Handler(
		http.StripPrefix("/node_modules/",
			http.FileServer(http.Dir(viper.GetString("frontend.node_modules"))),
		),
	)

	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			log.Info("ValidationKeyGetter")
			// Verify 'aud' claim
			aud := "https://jeffall.com/magi"
			checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
			if !checkAud {
				return token, errors.New("Invalid audience.")
			}
			// Verify 'iss' claim
			iss := "https://magiadmen.auth0.com/"
			checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
			if !checkIss {
				return token, errors.New("Invalid issuer.")
			}

			cert, err := getPemCert(token)
			if err != nil {
				panic(err.Error())
			}

			result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
			return result, nil
		},
		SigningMethod: jwt.SigningMethodRS256,
	})

	authMiddleware := func(next func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			jwtMiddleware.HandlerWithNext(w, r, middleware.HandleError(next).ServeHTTP)
		}
	}

	r.HandleFunc("/requests", authMiddleware(requests.PUT)).Methods("PUT")
}

func getPemCert(token *jwt.Token) (string, error) {
	log.Info("getPermCert")
	cert := ""
	resp, err := http.Get("https://magiadmen.auth0.com/.well-known/jwks.json")

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k, _ := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("Unable to find appropriate key.")
		return cert, err
	}

	return cert, nil
}

// func ReadPassword() (string, error) {
// 	fmt.Print("Password: ")
// 	password, err := terminal.ReadPassword(int(syscall.Stdin))
// 	fmt.Print("\n")
// 	if err != nil {
// 		log.WithFields(log.Fields{
// 			"Error": err,
// 		}).Error("Error reading password")
// 		return "", err
// 	}
// 	return string(password), nil
// }

func ConnectDatabase() error {
	db, err := gorm.Open(
		viper.GetString("database.config.driver"),
		viper.GetString("database.config.file"),
	)

	res.DB = &data.Gorm{DB: db}

	if err != nil {
		log.WithFields(
			log.Fields{
				"package":  "magi",
				"function": "ConnectDatabase()",
				"error":    err.Error(),
				"driver":   viper.GetString("database.config.driver"),
				"file":     viper.GetString("database.config.file"),
			},
		).Fatal("Failed to open file")
		return err
	}

	return nil
}
