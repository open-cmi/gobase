package webserver

import (
	"github.com/open-cmi/gobase/essential/logger"
	"github.com/open-cmi/gobase/initial"
)

var defaultService *Service

func DefaultService() *Service {
	return defaultService
}

func WaitError() <-chan error {
	if defaultService != nil {
		return defaultService.Error()
	}
	return nil
}

func Close() {
	if defaultService != nil {
		defaultService.Close()
	}
}

func Init() error {
	if len(gConf.Server) == 0 {
		return nil
	}
	// start web service
	s := New()

	// Init
	err := s.Init()
	if err != nil {
		logger.Errorf("web server init %s", err.Error())
		return err
	}

	// Run
	err = s.Run()
	if err != nil {
		logger.Errorf("web server run failed: %s", err.Error())
		return err
	}
	defaultService = s
	return nil
}

func init() {
	initial.Register("webserver", initial.PhaseFinal, Init)
}
