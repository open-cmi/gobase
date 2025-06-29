package webserver

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/gobase/essential/logger"
	"github.com/open-cmi/gobase/essential/webserver/middleware"
	"github.com/open-cmi/gobase/pkg/eyas"
)

type Service struct {
	Engine *gin.Engine
}

func New() *Service {
	if !gConf.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	return &Service{
		Engine: gin.New(),
	}
}

func (s *Service) Init() error {
	// init webserver
	middleware.DefaultMiddleware(s.Engine)

	workDir := eyas.GetRootPath()
	dir := fmt.Sprintf("%s/static/", workDir)
	s.Engine.Static("/api-static/", dir)
	middleware.SessionMiddleware(s.Engine)
	middleware.JWTMiddleware(s.Engine)
	UnauthInit(s.Engine)
	if !gConf.StrictAuth {
		AuthInit(s.Engine)
	}
	middleware.AuthMiddleware(s.Engine)
	if gConf.StrictAuth {
		AuthInit(s.Engine)
	}
	MustAuthInit(s.Engine)
	return nil
}

func RunHTTPSServer(eng *gin.Engine, s *Server) {
	rp := eyas.GetRootPath()
	var certFile string
	if strings.HasPrefix(s.CertFile, ".") {
		certFile = filepath.Join(rp, s.CertFile)
	} else {
		certFile = s.CertFile
	}
	var keyFile string
	if strings.HasPrefix(s.KeyFile, ".") {
		keyFile = filepath.Join(rp, s.KeyFile)
	} else {
		keyFile = s.KeyFile
	}
	t := net.JoinHostPort(s.Address, strconv.Itoa(s.Port))
	logger.Debugf("tls server started: %s, cert %s, key %s",
		t, certFile, keyFile)
	err := eng.RunTLS(t, certFile, keyFile)
	logger.Debugf("tls server stopped: %s, cert %s, key %s, err: %s",
		t, certFile, keyFile, err.Error())
}

func RunHTTPServer(eng *gin.Engine, s *Server) {
	t := net.JoinHostPort(s.Address, strconv.Itoa(s.Port))
	logger.Debugf("http server started: %s\n", t)
	err := eng.Run(t)
	logger.Debugf("http server %s stopped: %s\n", t, err.Error())
}

func RunUnixServer(eng *gin.Engine, s *Server) {
	sockAddr := s.Address
	os.Remove(sockAddr)
	unixAddr, err := net.ResolveUnixAddr("unix", sockAddr)
	if err != nil {
		logger.Error(err.Error() + "\n")
		return
	}

	listener, err := net.ListenUnix("unix", unixAddr)
	if err != nil {
		logger.Errorf("listening error: %s\n", err.Error())
		return
	}

	logger.Debugf("unix socket started: %s\n", sockAddr)
	err = http.Serve(listener, eng)
	logger.Debugf("unix socket %s stopped: %s\n", sockAddr, err.Error())
}

func (s *Service) Run() error {
	// unix sock api
	for _, srv := range gConf.Server {
		switch srv.Proto {
		case "unix":
			go RunUnixServer(s.Engine, &srv)
		case "http":
			go RunHTTPServer(s.Engine, &srv)
		case "https":
			go RunHTTPSServer(s.Engine, &srv)
		}
	}

	return nil
}

func (s *Service) Close() {

}
