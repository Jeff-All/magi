package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"syscall"
	"time"
	"toolbox"

	"github.com/casbin/casbin"
	"github.com/gorilla/sessions"

	"github.com/Jeff-All/magi/auth"
	"github.com/Jeff-All/magi/errors"
	. "github.com/Jeff-All/magi/middleware"
	res "github.com/Jeff-All/magi/resources"
	"github.com/Jeff-All/magi/session"

	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
	"github.com/urfave/cli"

	"github.com/jinzhu/gorm"

	"github.com/Jeff-All/magi/models"

	"strconv"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"golang.org/x/crypto/ssh/terminal"

	data "github.com/Jeff-All/magi/data"

	"github.com/Jeff-All/magi/endpoints/request"

	"io/ioutil"
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
			Name:    "automigrate",
			Aliases: []string{"am"},
			Action:  AutoMigrate,
			Flags:   app.Flags,
		},
		{
			Name:    "init",
			Aliases: []string{"i"},
			Action:  Init,
			Flags:   app.Flags,
		},
		{
			Name:    "auth",
			Aliases: []string{"au"},
			Action:  Auth,
			Flags:   app.Flags,
		},
		{
			Name:    "add",
			Aliases: []string{"ad"},
			Action:  AddUser,
			Flags:   app.Flags,
		},
	}

	app.Action = Run

	err := app.Run(os.Args)
	toolbox.FatalError("unable to launch the app", err)
}

func Common(c *cli.Context) error {
	SetLogLevel(c)

	viper.SetConfigName(c.String("config"))
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	toolbox.FatalError("unable to read config", err)

	return nil
}

func Run(c *cli.Context) error {
	log.Printf("Starting Application")
	Common(c)
	ConnectDatabase()
	defer res.DB.Close()
	Bind()

	LaunchServer(c.Bool("local"))

	return nil
}

func Bind() {
	models.DB = res.DB
	auth.DB = res.DB
}

// SetLogLevel
//
// Sets the current log level based on terminal parameters
// debug -> DebugLevel
// info -> InfoLevel
// default -> ErrorLevel
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

func BuildAddress(
	local bool,
) string {
	toReturn := ""
	if local {
		toReturn = "localhost"
	}
	toReturn += ":" + viper.GetString("server.port")
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

	BuildSessionManager()
	BuildEnforcer()
	ConfigureRoutes(r, local)

	s := BuildServer(r, local)

	// auth := middleware.Authorize(res.Enforcer, res.Session)

	http.Handle("/", r)

	log.Fatal(s.ListenAndServe())
}

func BuildSessionManager() {
	res.Session = &session.Manager{
		Store: sessions.NewCookieStore([]byte(viper.GetString("session.private"))),
	}
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

func ConfigureRoutes(
	r *mux.Router,
	local bool,
) {
	log.Debugf("ConfigureRoutes")

	middleware := func(final ErrorHandler) func(http.ResponseWriter, *http.Request) {
		return HandleError(Authorize(
			res.Enforcer,
			res.Session,
			"/login",
		)(final)).ServeHTTP
	}

	r.HandleFunc("/requests", middleware(request.Request.PUT)).Methods("PUT")
	r.HandleFunc("/requests", middleware(request.Request.GETPAGE)).Methods("GET")
	r.HandleFunc("/requests/{id}", middleware(request.Request.GET)).Methods("GET")
	r.HandleFunc("/requests/{id}", middleware(request.Request.DELETE)).Methods("DELETE")

	r.HandleFunc("/requests/{id}/gifts", middleware(request.Request.PUTGift)).Methods("PUT")
	// r.HandleFunc("/requests/{id}/gifts", middleware(request.Request.GETPAGEGift)).Methods("GET")
	// r.HandleFunc("/requests/{id}/gifts/{gift_id}", middleware(request.Request.GET)).Methods("GET")
	// r.HandleFunc("/requests/{id}/gifts/{gift_id}", middleware(request.Request.DELETE)).Methods("DELETE")

	r.HandleFunc("/login", HandleError(res.Session.Login).ServeHTTP).Methods("PUT")
	r.HandleFunc("/login", HandleError(GetHTML("frontend/login")).ServeHTTP).Methods("GET")

	r.HandleFunc("/processor.html", middleware(GetHTML("frontend/processor"))).Methods("GET")
	r.HandleFunc("/", middleware(GetHTML("frontend/index"))).Methods("GET")
}

func GetHTML(filename string) func(w http.ResponseWriter, r *http.Request) error {
	filename = "./" + filename + ".html"
	return func(w http.ResponseWriter, r *http.Request) error {
		file, err := ioutil.ReadFile(filename)
		if err != nil {
			return errors.CodedError{
				Message:  "Internal Server Error",
				HTTPCode: http.StatusInternalServerError,
				Err:      err,
			}
		}
		_, err = w.Write(file)
		if err != nil {
			return errors.CodedError{
				Message:  "Internal Server Error",
				HTTPCode: http.StatusInternalServerError,
				Err:      err,
			}
		}
		return nil
	}

}

func createUsers() models.Users {
	users := models.Users{}
	users = append(users, models.User{ID: 1, Name: "Admin", Role: "admin"})
	users = append(users, models.User{ID: 2, Name: "Sabine", Role: "member"})
	users = append(users, models.User{ID: 3, Name: "Sepp", Role: "member"})
	return users
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

func AutoMigrate(c *cli.Context) error {
	fmt.Println("AutoMigrate")
	Common(c)
	log.Debug("AutoMigrate()")
	ConnectDatabase()
	defer res.DB.Close()

	models.AutoMigrate()
	res.DB.Close()

	return nil
}

func Init(c *cli.Context) error {
	fmt.Println("Init")
	Common(c)
	log.Debug("Init()")
	if _, err := os.Stat(viper.GetString("database.config.file")); err == nil {
		log.WithFields(log.Fields{
			"File": viper.GetString("database.config.file"),
		}).Error("Cannot Init a database that already exists")
		return fmt.Errorf("Cannot Init a database that already exists")
	}
	ConnectDatabase()
	Bind()

	password, _ := ReadPassword()
	if confirmed, err := ConfirmPassword(password); !confirmed || err != nil {
		log.Info("Passwords did not match")
		return err
	}

	models.AutoMigrate()
	auth.Init()
	auth.AddRootUser(password)
	return nil
}

func Auth(c *cli.Context) error {
	fmt.Println("Auth")
	Common(c)
	log.Debug("Auth()")

	username, _ := ReadLine("Username")
	password, err := ReadPassword()
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error reading password")
	}
	ConnectDatabase()
	Bind()
	user, err := auth.BasicAuthentication(
		username,
		password,
	)
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Authenticating User")
		return err
	}
	if user == nil {
		log.Info("User could not be authenticated")
		return nil
	}
	log.Infof("User ID = %d", user.ID)
	return nil
}

func ReadPassword() (string, error) {
	fmt.Print("Password: ")
	password, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Print("\n")
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error reading password")
		return "", err
	}
	return string(password), nil
}

func ReadInt(mes string) (int, error) {
	str, err := ReadLine(mes)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(str)
}

func ReadLine(mes string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s: ", mes)
	line, _, err := reader.ReadLine()
	return string(line), err
}

func ConfirmPassword(password string) (bool, error) {
	confirmation, err := ReadPassword()
	return confirmation == password, err
}

func AddUser(c *cli.Context) error {
	fmt.Println("AddUser")
	Common(c)
	log.Debug("AddUser()")
	ConnectDatabase()
	// Get User
	username, _ := ReadLine("Username")
	password, err := ReadPassword()
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error reading password")
	}
	user, err := auth.BasicAuthentication(
		username,
		password,
	)

	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error authenticating user")
		return err
	}
	if user == nil {
		log.Info("Unable to authenticate user")
		return nil
	}

	// Add User
	// newUserName, _ := ReadLine("New Username")
	// newPassword, _ := ReadPassword()
	// if confirmed, err := ConfirmPassword(password); !confirmed || err != nil {
	// 	log.Info("Passwords did not match")
	// 	return err
	// }
	// level, _ := ReadInt("Level")

	// newUser, err := user.AddUser(
	// 	newUserName,
	// 	newPassword,
	// 	auth.Level(level),
	// )
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error adding user")
		return err
	}
	// if newUser == nil {
	// 	log.Info("User wasn't created")
	// 	return nil
	// }
	log.Info("User successfuly created")
	return nil
}
