package webserver

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/gobase/essential/logger"
	"github.com/open-cmi/gobase/pkg/eyas"
)

type serverInstance struct {
	httpServer *http.Server
	listener   net.Listener
	proto      string
}

type Service struct {
	Engine    *gin.Engine
	servers   []*serverInstance
	errChan   chan error
	closeOnce sync.Once
}

func New() *Service {
	if !gConf.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	return &Service{
		Engine:  gin.New(),
		errChan: make(chan error, 16),
	}
}

func (s *Service) Error() <-chan error {
	return s.errChan
}

func (s *Service) Init() error {

	rbacInstance := GetRBAC()
	if rbacInstance != nil {
		rbacInstance.Init(s.Engine)
	}

	s.initNoRoute()

	return nil
}

func RunHTTPSServer(eng *gin.Engine, s *Server) error {
	rp := eyas.GetWorkingDir()
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

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		logger.Errorf("tls load keypair failed: %s\n", err.Error())
		return err
	}

	t := net.JoinHostPort(s.Address, strconv.Itoa(s.Port))
	ln, err := net.Listen("tcp", t)
	if err != nil {
		logger.Errorf("tls server listen on %s failed: %s\n", t, err.Error())
		return err
	}
	defer ln.Close()

	tlsListener := tls.NewListener(ln, &tls.Config{
		Certificates: []tls.Certificate{cert},
	})
	defer tlsListener.Close()

	srv := &http.Server{
		Addr:    t,
		Handler: eng,
	}

	logger.Debugf("tls server started: %s, cert %s, key %s",
		t, certFile, keyFile)
	err = srv.Serve(tlsListener)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Debugf("tls server stopped: %s, cert %s, key %s, err: %s",
			t, certFile, keyFile, err.Error())
		return err
	}
	logger.Debugf("tls server %s stopped\n", t)
	return nil
}

func RunHTTPServer(eng *gin.Engine, s *Server) error {
	t := net.JoinHostPort(s.Address, strconv.Itoa(s.Port))
	ln, err := net.Listen("tcp", t)
	if err != nil {
		logger.Errorf("http server listen on %s failed: %s\n", t, err.Error())
		return err
	}
	defer ln.Close()

	srv := &http.Server{
		Addr:    t,
		Handler: eng,
	}

	logger.Debugf("http server started: %s\n", t)
	err = srv.Serve(ln)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Debugf("http server %s stopped: %s\n", t, err.Error())
		return err
	}
	logger.Debugf("http server %s stopped\n", t)
	return nil
}

func RunUnixServer(eng *gin.Engine, s *Server) error {
	sockAddr := s.Address
	_ = os.Remove(sockAddr)
	unixAddr, err := net.ResolveUnixAddr("unix", sockAddr)
	if err != nil {
		logger.Errorf("resolve unix addr error: %s\n", err.Error())
		return err
	}

	listener, err := net.ListenUnix("unix", unixAddr)
	if err != nil {
		logger.Errorf("listening unix error: %s\n", err.Error())
		return err
	}
	defer listener.Close()

	srv := &http.Server{
		Handler: eng,
	}

	logger.Debugf("unix socket started: %s\n", sockAddr)
	err = srv.Serve(listener)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Debugf("unix socket %s stopped: %s\n", sockAddr, err.Error())
		return err
	}
	logger.Debugf("unix socket %s stopped\n", sockAddr)
	return nil
}

func (s *Service) createServerInstance(srvConf *Server) (*serverInstance, error) {
	t := net.JoinHostPort(srvConf.Address, strconv.Itoa(srvConf.Port))
	switch srvConf.Proto {
	case "http":
		ln, err := net.Listen("tcp", t)
		if err != nil {
			return nil, fmt.Errorf("http listen on %s failed: %w", t, err)
		}
		httpSrv := &http.Server{
			Addr:    t,
			Handler: s.Handler(),
		}
		return &serverInstance{httpServer: httpSrv, listener: ln, proto: "http"}, nil

	case "https":
		rp := eyas.GetWorkingDir()
		var certFile string
		if strings.HasPrefix(srvConf.CertFile, ".") {
			certFile = filepath.Join(rp, srvConf.CertFile)
		} else {
			certFile = srvConf.CertFile
		}
		var keyFile string
		if strings.HasPrefix(srvConf.KeyFile, ".") {
			keyFile = filepath.Join(rp, srvConf.KeyFile)
		} else {
			keyFile = srvConf.KeyFile
		}

		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return nil, fmt.Errorf("load tls cert (%s, %s) failed: %w", certFile, keyFile, err)
		}

		ln, err := net.Listen("tcp", t)
		if err != nil {
			return nil, fmt.Errorf("https listen on %s failed: %w", t, err)
		}

		tlsLn := tls.NewListener(ln, &tls.Config{
			Certificates: []tls.Certificate{cert},
		})
		httpSrv := &http.Server{
			Addr:    t,
			Handler: s.Handler(),
		}
		return &serverInstance{httpServer: httpSrv, listener: tlsLn, proto: "https"}, nil

	case "unix":
		sockAddr := srvConf.Address
		_ = os.Remove(sockAddr)
		unixAddr, err := net.ResolveUnixAddr("unix", sockAddr)
		if err != nil {
			return nil, fmt.Errorf("resolve unix addr %s failed: %w", sockAddr, err)
		}
		ln, err := net.ListenUnix("unix", unixAddr)
		if err != nil {
			return nil, fmt.Errorf("listen unix %s failed: %w", sockAddr, err)
		}
		httpSrv := &http.Server{
			Handler: s.Handler(),
		}
		return &serverInstance{httpServer: httpSrv, listener: ln, proto: "unix"}, nil

	default:
		return nil, fmt.Errorf("unsupported protocol: %s", srvConf.Proto)
	}
}

func (s *Service) Run() error {
	for _, srv := range gConf.Server {
		srvConf := srv
		inst, err := s.createServerInstance(&srvConf)
		if err != nil {
			logger.Errorf("webserver create server instance failed: %s", err.Error())
			s.Close()
			return err
		}
		s.servers = append(s.servers, inst)
	}

	for _, inst := range s.servers {
		instance := inst
		go func() {
			logger.Debugf("%s server started: %s\n", instance.proto, instance.listener.Addr().String())
			err := instance.httpServer.Serve(instance.listener)
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				logger.Errorf("%s server %s stopped with error: %s\n", instance.proto, instance.listener.Addr().String(), err.Error())
				select {
				case s.errChan <- err:
				default:
				}
			} else {
				logger.Debugf("%s server %s stopped\n", instance.proto, instance.listener.Addr().String())
			}
		}()
	}

	return nil
}

func (s *Service) Close() {
	s.closeOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		for _, inst := range s.servers {
			if inst.httpServer != nil {
				_ = inst.httpServer.Shutdown(ctx)
			}
			if inst.listener != nil {
				_ = inst.listener.Close()
			}
		}
		s.servers = nil
	})
}
