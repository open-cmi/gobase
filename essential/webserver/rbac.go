package webserver

import (
	"github.com/gin-gonic/gin"
)

type RBACInterface interface {
	Init(*gin.Engine)
}

var globaRBAC RBACInterface

func SetRBAC(rbac RBACInterface) {
	globaRBAC = rbac
}

func GetRBAC() RBACInterface {
	return globaRBAC
}
