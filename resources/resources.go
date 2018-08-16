package resources

import (
	"github.com/Jeff-All/magi/data"
	"github.com/Jeff-All/magi/session"
	"github.com/casbin/casbin"
	"github.com/qor/auth"
)

var DB data.Data
var Auth *auth.Auth
var Session *session.Manager
var Enforcer *casbin.Enforcer
