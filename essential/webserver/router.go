package webserver

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

// RouterGroup is a function to api group router
type RouterGroup func(e *gin.Engine)

var mustAuthGroup map[string]string = make(map[string]string)
var authGroup map[string]string = make(map[string]string)
var unauthGroup map[string]string = make(map[string]string)

// UnauthInit no auth router init
func UnauthInit(e *gin.Engine) {
	for mod, groupPath := range unauthGroup {
		g := e.Group(groupPath)
		{
			modPath, found := unauthAPIPath[mod]
			if !found {
				continue
			}

			for _, r := range modPath {
				if r.Method == "POST" {
					g.POST(r.Path, r.Callback)
				} else if r.Method == "GET" {
					g.GET(r.Path, r.Callback)
				} else if r.Method == "DELETE" {
					g.DELETE(r.Path, r.Callback)
				} else if r.Method == "PUT" {
					g.PUT(r.Path, r.Callback)
				}
			}
		}
	}
}

// AuthInit auth router init
func AuthInit(e *gin.Engine) {
	for mod, groupPath := range authGroup {
		g := e.Group(groupPath)
		{
			modPath, found := authAPIPath[mod]
			if !found {
				continue
			}

			for _, r := range modPath {
				if r.Method == "POST" {
					g.POST(r.Path, r.Callback)
				} else if r.Method == "GET" {
					g.GET(r.Path, r.Callback)
				} else if r.Method == "DELETE" {
					g.DELETE(r.Path, r.Callback)
				} else if r.Method == "PUT" {
					g.PUT(r.Path, r.Callback)
				}
			}
		}
	}
}

// AuthInit auth router init
func MustAuthInit(e *gin.Engine) {
	for mod, groupPath := range mustAuthGroup {
		g := e.Group(groupPath)
		{
			modPath, found := mustAuthAPIPath[mod]
			if !found {
				continue
			}

			for _, r := range modPath {
				if r.Method == "POST" {
					g.POST(r.Path, r.Callback)
				} else if r.Method == "GET" {
					g.GET(r.Path, r.Callback)
				} else if r.Method == "DELETE" {
					g.DELETE(r.Path, r.Callback)
				} else if r.Method == "PUT" {
					g.PUT(r.Path, r.Callback)
				}
			}
		}
	}
}

func RegisterAuthRouter(module string, groupPath string) error {
	_, found := authGroup[module]
	if found {
		errMsg := fmt.Sprintf("module %s auth group api has been registered", module)
		return errors.New(errMsg)
	}
	authGroup[module] = groupPath
	return nil
}

func RegisterMustAuthRouter(module string, groupPath string) error {
	_, found := mustAuthGroup[module]
	if found {
		errMsg := fmt.Sprintf("module %s auth group api has been registered", module)
		return errors.New(errMsg)
	}
	mustAuthGroup[module] = groupPath
	return nil
}

func RegisterUnauthRouter(module string, groupPath string) error {
	_, found := unauthGroup[module]
	if found {
		errMsg := fmt.Sprintf("module %s unauth group api has been registered", module)
		return errors.New(errMsg)
	}
	unauthGroup[module] = groupPath
	return nil
}

type API struct {
	Prod     string
	Method   string
	Path     string
	Callback func(c *gin.Context)
}

var mustAuthAPIPath map[string][]API = make(map[string][]API)
var authAPIPath map[string][]API = make(map[string][]API)
var unauthAPIPath map[string][]API = make(map[string][]API)

func RegisterAuthAPI(prod string, method string, path string, proc func(c *gin.Context)) error {
	modPath, found := authAPIPath[prod]
	if !found {
		authAPIPath[prod] = []API{}
		modPath = authAPIPath[prod]
	}

	modPath = append(modPath, API{
		Prod:     prod,
		Method:   method,
		Path:     path,
		Callback: proc,
	})
	authAPIPath[prod] = modPath

	return nil
}

func RegisterUnauthAPI(prod string, method string, path string, proc func(c *gin.Context)) error {
	modPath, found := unauthAPIPath[prod]
	if !found {
		unauthAPIPath[prod] = []API{}
		modPath = unauthAPIPath[prod]
	}

	modPath = append(modPath, API{
		Prod:     prod,
		Method:   method,
		Path:     path,
		Callback: proc,
	})
	unauthAPIPath[prod] = modPath

	return nil
}

func RegisterMustAuthAPI(prod string, method string, path string, proc func(c *gin.Context)) error {
	modPath, found := mustAuthAPIPath[prod]
	if !found {
		mustAuthAPIPath[prod] = []API{}
		modPath = mustAuthAPIPath[prod]
	}

	modPath = append(modPath, API{
		Prod:     prod,
		Method:   method,
		Path:     path,
		Callback: proc,
	})
	mustAuthAPIPath[prod] = modPath

	return nil
}
