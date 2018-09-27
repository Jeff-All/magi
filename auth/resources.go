package auth

import (
	"github.com/Jeff-All/magi/data"
	"github.com/casbin/casbin"
)

var DB data.Data
var Salt string
var PrivateKey string
var Enforcer *casbin.Enforcer
