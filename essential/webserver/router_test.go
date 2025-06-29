package webserver

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func BlacklistList(c *gin.Context) {
}

func TestRegister(t *testing.T) {

	RegisterUnauthAPI("ngfw", "GET", "/blacklist", BlacklistList)

	prod, found := unauthAPIPath["ngfw"]
	if !found {
		t.Errorf("api not found")
	}
	var apiFound bool = false
	for _, r := range prod {
		if r.Method == "GET" && r.Path == "/blacklist" {
			apiFound = true
		}
	}
	if !apiFound {
		t.Errorf("api not found")
	}
}
