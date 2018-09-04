package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"syscall"
	"time"
	"toolbox"

	"github.com/Jeff-All/magi/endpoints/login"

	"github.com/Jeff-All/magi/endpoints/admin/user/application"

	"github.com/casbin/casbin"
	"github.com/gorilla/sessions"

	"github.com/Jeff-All/magi/endpoints"

	"github.com/Jeff-All/magi/auth"
	"github.com/Jeff-All/magi/errors"
	"github.com/Jeff-All/magi/mail"
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
			Flags:   append(app.Flags, cli.BoolFlag{Name: "overwrite, ow"}),
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

func InitEmail(c *cli.Context) {
	password, err := ReadPassword()
	if err != nil {
		log.Fatal("Unable to read password")
	}
	mail.Init(
		viper.GetString("mail.host"),
		viper.GetInt("mail.port"),
		viper.GetString("mail.user"),
		password,
	)
}

func Run(c *cli.Context) error {
	log.Printf("Starting Application")
	Common(c)
	InitEmail(c)
	ConnectDatabase()
	defer res.DB.Close()
	Bind()

	LaunchServer(c)

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

func BuildServer(
	r *mux.Router,
) *http.Server {
	logrus.WithFields(log.Fields{
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

func LaunchServer(
	c *cli.Context,
) {
	r := mux.NewRouter()

	BuildSessionManager()
	BuildEnforcer()
	ConfigureRoutes(c, r)

	s := BuildServer(r)

	http.Handle("/", r)

	log.Fatal(s.ListenAndServeTLS("dev.magi.crt", "dev.magi.key"))
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
	c *cli.Context,
	r *mux.Router,
) {
	log.Debugf("ConfigureRoutes")

	middleware := func(final ErrorHandler) func(http.ResponseWriter, *http.Request) {
		return Log(HandleError(Authorize(
			res.Enforcer,
			res.Session,
			"/login",
		)(final))).ServeHTTP
	}

	noauth := func(final ErrorHandler) func(http.ResponseWriter, *http.Request) {
		return Log(HandleError(final)).ServeHTTP
	}

	r.HandleFunc("/api/requests", middleware(request.Request.PUT)).Methods("PUT")
	r.HandleFunc("/api/requests", middleware(request.Request.GETPAGE)).Methods("GET")
	r.HandleFunc("/api/requests/upload", middleware(GetHTML("frontend/requests_upload"))).Methods("GET")
	r.HandleFunc("/api/requests/{id}", middleware(request.Request.GET)).Methods("GET")
	r.HandleFunc("/api/requests/{id}", middleware(request.Request.DELETE)).Methods("DELETE")

	r.HandleFunc("/api/requests/{id}/gifts", middleware(request.Request.PUTGift)).Methods("PUT")
	// r.HandleFunc("api/requests/{id}/gifts", middleware(request.Request.GETPAGEGift)).Methods("GET")
	// r.HandleFunc("api/requests/{id}/gifts/{gift_id}", middleware(request.Request.GET)).Methods("GET")
	// r.HandleFunc("api/requests/{id}/gifts/{gift_id}", middleware(request.Request.DELETE)).Methods("DELETE")

	r.HandleFunc("/admin/user/application.html", middleware(GetHTML("frontend/add_user"))).Methods("GET")
	r.HandleFunc("/admin/user/application", middleware(application.PUT())).Methods("PUT")

	type User struct {
		ID     uint64
		Active bool
		Email  string
		Roles  []auth.Role `gorm:"many2many:user_roles"`
	}
	r.HandleFunc("/admin/users", middleware(endpoints.GETPage(User{}, res.DB, "Roles"))).Methods("GET")
	// r.HandleFunc("/admin/users", middleware(endpoints.UPDATE()))
	if GetHTML, err := endpoints.GetHTML("frontend/users", c.Bool("debug")); err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatal("Error Getting HTML for 'frontend/users'")
	} else {
		r.HandleFunc("/admin/users.html", middleware(GetHTML)).Methods("GET")
	}

	r.HandleFunc("/login", noauth(res.Session.Login)).Methods("PUT")
	r.HandleFunc("/login", noauth(GetHTML("frontend/login"))).Methods("GET")

	r.HandleFunc("/logout", noauth(res.Session.LogOut)).Methods("GET")

	r.HandleFunc("/js/validation.js", noauth(endpoints.GetResource)).Methods("GET")

	GETApplication, err := application.GET("frontend/register", true)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Error building /login/{hash} handler")
	}

	r.HandleFunc(`/login/{hash:[A-Za-z0-9\-_\=]+}.html`, noauth(GETApplication)).Methods("GET")
	r.HandleFunc(`/login/{hash:[A-Za-z0-9\-_\=]+}`, noauth(login.PUTHash())).Methods("PUT")

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
	if _, err := os.Stat(viper.GetString("database.config.file")); err == nil && !c.Bool("overwrite") {
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
	auth.AddRole("manager")
	auth.AddRole("recorder")
	auth.AddRole("shopper")
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
