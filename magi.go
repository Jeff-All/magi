package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/casbin/casbin"
	"github.com/google/jsonapi"
	"github.com/rs/cors"

	"github.com/Jeff-All/magi/actions"
	"github.com/Jeff-All/magi/auth"
	"github.com/Jeff-All/magi/data"
	"github.com/Jeff-All/magi/endpoints"
	"github.com/Jeff-All/magi/input"
	"github.com/Jeff-All/magi/models"

	"github.com/Jeff-All/magi/middleware"
	res "github.com/Jeff-All/magi/resources"

	requests "github.com/Jeff-All/magi/endpoints/request"

	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/urfave/cli"

	. "github.com/Jeff-All/magi/errors"
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

	app.Commands = []cli.Command{
		{
			Name:    "init",
			Aliases: []string{"i"},
			Action:  Init,
			Flags:   append(app.Flags, cli.BoolFlag{Name: "overwrite, ow"}),
		},
		{
			Name:    "migrate",
			Aliases: []string{"m"},
			Action:  Migrate,
		},
	}

	app.Action = Run

	if err := app.Run(os.Args); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("unable to launch")
	}
}

func Migrate(c *cli.Context) error {
	Common(c)
	ConnectDatabase()
	defer res.DB.Close()
	actions.DB = res.DB
	actions.AutoMigrate()
	return nil
}

func Init(c *cli.Context) error {
	fmt.Println("Init")
	Common(c)
	log.Debug("Init()")
	if _, err := os.Stat(viper.GetString("database.config.file")); err == nil && !c.Bool("overwrite") {
		log.WithFields(log.Fields{
			"File":  viper.GetString("database.config.file"),
			"error": err,
		}).Error("Cannot Init a database that already exists")
		return fmt.Errorf("Cannot Init a database that already exists")
	}
	ConnectDatabase()
	defer res.DB.Close()

	BindAuth()
	if err := auth.Init(); err != nil {
		return err
	}

	actions.DB = res.DB
	actions.AutoMigrate()
	return nil
}

func Run(context *cli.Context) error {
	log.Printf("Run")
	Common(context)
	// InitEmail(context)
	if err := ConnectDatabase(); err != nil {
		return err
	}
	defer res.DB.Close()
	actions.DB = res.DB
	r := mux.NewRouter()
	ConfigureRoutes(context, r)

	corsConfig := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
		AllowedMethods:   []string{"PUT", "GET", "PATCH"},
		AllowedHeaders:   []string{"Authorization"},
		Debug:            true,
	})

	BuildEnforcer()

	BindAuth()

	server := &http.Server{
		Addr:           viper.GetString("server.domain"),
		Handler:        corsConfig.Handler(r),
		ReadTimeout:    (time.Duration)(viper.GetInt("server.readTimeout")) * time.Second,
		WriteTimeout:   (time.Duration)(viper.GetInt("server.writeTimeout")) * time.Second,
		MaxHeaderBytes: viper.GetInt("server.maxHeaderBytes"),
	}
	http.Handle("/", r)
	log.Fatal(server.ListenAndServeTLS("dev.magi.crt", "dev.magi.key"))
	return nil
}

func Common(c *cli.Context) error {
	SetLogLevel(c)
	viper.SetConfigName(c.String("config"))
	viper.AddConfigPath(".")
	return viper.ReadInConfig()
}

func BindAuth() {
	auth.DB = res.DB
	auth.Enforcer = res.Enforcer
}

func BuildJsonAPI() {
	jsonapi.Instrumentation = func(r *jsonapi.Runtime, eventType jsonapi.Event, callGUID string, dur time.Duration) {
		metricPrefix := r.Value("instrument").(string)

		if eventType == jsonapi.UnmarshalStart {
			log.Debug("%s: id, %s, started at %v\n", metricPrefix+".jsonapi_unmarshal_time", callGUID, time.Now())
		}

		if eventType == jsonapi.UnmarshalStop {
			log.Debug("%s: id, %s, stopped at, %v , and took %v to unmarshal payload\n", metricPrefix+".jsonapi_unmarshal_time", callGUID, time.Now(), dur)
		}

		if eventType == jsonapi.MarshalStart {
			log.Debug("%s: id, %s, started at %v\n", metricPrefix+".jsonapi_marshal_time", callGUID, time.Now())
		}

		if eventType == jsonapi.MarshalStop {
			log.Debug("%s: id, %s, stopped at, %v , and took %v to marshal payload\n", metricPrefix+".jsonapi_marshal_time", callGUID, time.Now(), dur)
		}
	}
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
	handler http.Handler,
) *http.Server {
	log.WithFields(log.Fields{
		"readTimeout":    viper.GetInt("server.readTimeout"),
		"writeTimeout":   viper.GetInt("server.writeTimeout"),
		"maxHeaderBytes": viper.GetInt("server.maxHeaderBytes"),
		"domain":         viper.GetString("server.domain"),
	}).Debug("BuildServer")
	return &http.Server{
		Addr:           viper.GetString("server.domain"),
		Handler:        handler,
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
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {

			claims, err := json.Marshal(token.Claims)

			log.WithFields(log.Fields{
				"claims": string(claims),
				"err":    err,
			}).Debug("Token Claims")

			// Verify 'aud' claim
			checkAud := token.Claims.(jwt.MapClaims).VerifyAudience("cJvGJ3Mpan5HRFFGBBtfg6ch4E2cu10f", true)
			if !checkAud {
				return token, errors.New("Invalid audience.")
			}
			// Verify 'iss' claim
			iss := "https://jefall.auth0.com/"
			checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
			if !checkIss {
				return token, errors.New("Invalid issuer.")
			}

			cert, err := getPemCert(token)
			if err != nil {
				panic(err.Error())
			}

			result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
			log.WithFields(log.Fields{
				"result": result,
			}).Debug("jwtMiddleware")
			return result, nil
		},
		SigningMethod: jwt.SigningMethodRS256,
	})

	authorization := func(next middleware.ErrorHandler) middleware.ErrorHandler {
		return func(w http.ResponseWriter, r *http.Request) error {
			claims := r.Context().Value("user").(*jwt.Token).Claims.(jwt.MapClaims)
			log.WithFields(log.Fields{
				"claims": claims,
			}).Debug("authorization")
			if user, err := auth.GetUser(claims); err != nil {
				return err
			} else if err := user.EnforceRole(r); err != nil {
				return err
			}
			return next(w, r)
		}
	}

	logger := func(next http.Handler) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			log.WithFields(log.Fields{
				"headers": r.Header,
			}).Debug("logging:before")
			var builder strings.Builder
			writer := io.MultiWriter(w, &builder)
			responseWriter := ResponseWriter{
				ResponseWriter: w,
				Writer:         writer,
			}
			next.ServeHTTP(&responseWriter, r)
			log.WithFields(log.Fields{
				"headers": r.Header,
				"body":    builder.String(),
			}).Debug("logging:after")
		}
	}

	authMiddleware := func(next func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			jwtMiddleware.HandlerWithNext(w, r,
				logger(
					middleware.HandleError(
						authorization(next))))
		}
	}

	r.HandleFunc("/requests", authMiddleware(requests.PUT)).Methods("PUT")
	r.HandleFunc("/requests", authMiddleware(
		endpoints.GetPage(models.Request{}, res.DB),
	)).Methods("GET")
	r.HandleFunc("/gifts", authMiddleware(
		endpoints.GetPage(models.Gift{}, res.DB),
	)).Methods("GET")

	r.HandleFunc("/admin/users", authMiddleware(
		endpoints.GetPage(auth.User{}, res.DB),
	)).Methods("GET")
	r.HandleFunc("/admin/users", authMiddleware(
		endpoints.Patch(auth.User{}, input.User{}, res.DB),
	)).Methods("PATCH")

	r.HandleFunc("/roles", authMiddleware(
		endpoints.GetPage(auth.Role{}, res.DB),
	)).Methods("GET")
	r.HandleFunc("/user/role", authMiddleware(
		func(w http.ResponseWriter, r *http.Request) error {
			claims := r.Context().Value("user").(*jwt.Token).Claims.(jwt.MapClaims)
			if user, err := auth.GetUser(claims); err != nil {
				return err
			} else if _, err := w.Write([]byte(user.Role)); err != nil {
				return CodedError{
					Message:  "error writing response",
					HTTPCode: 500,
					Err:      err,
				}
			}
			return nil
		},
	)).Methods("GET")
}

type ResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (rw *ResponseWriter) Write(data []byte) (int, error) {
	return rw.Writer.Write(data)
}

func BuildEnforcer() {
	var err error
	res.Enforcer, err = casbin.NewEnforcerSafe("./auth_model.conf", "./policy.csv")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatal("Unable to build enforcer")
	}
}

func getPemCert(token *jwt.Token) (string, error) {
	log.Debug("getPermCert")
	cert := ""
	resp, err := http.Get("https://jefall.auth0.com/.well-known/jwks.json")

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

// func BuildSessionManager() {
// 	res.Session = &session.Manager{
// 		Store: sessions.NewCookieStore([]byte(viper.GetString("session.private"))),
// 	}
// }

// map[
// 	"updated_at":"2018-09-25T00:01:46.242Z",
// 	"sub":"google-oauth2|111259359070610605241"
// 	"exp":"1.537869706e+09",
// 	"gender":"male"
// 	"locale":"en",
// 	"picture":"https://lh4.googleusercontent.com/-XIEcaF6vX1M/AAAAAAAAAAI/AAAAAAAAAAA/AAN31DUCBV7es-4GmdcZHaF6gJ5iHwpQaw/mo/photo.jpg"
// 	"iss":"https://jefall.auth0.com/"
// 	"iat":"1.537833706e+09"
// 	"nonce":"yf2VnMd_xVdr-rdTnJfmvjuyMtFzcnsx"
// 	"family_name":"Something"
// 	"nickname":"jiffall"
// 	"at_hash":"VK1ZUkv36JUJNXRvmdiBPg"
// 	"name":"Something Something"
// 	"aud":"cJvGJ3Mpan5HRFFGBBtfg6ch4E2cu10f"
// 	"given_name":"Something"
// ]
