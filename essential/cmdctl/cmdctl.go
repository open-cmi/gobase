package cmdctl

import (
	"errors"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/open-cmi/gobase/essential/logger"
)

// Status the status of process

const (
	// Waiting wait to start
	Waiting int = iota

	// Stopped the stopped status
	Stopped

	// Starting the starting status
	Starting

	// Running the running status
	Running

	// Backoff the backoff status
	Backoff

	// Stopping the stopping status
	Stopping

	// Exited the Exited status
	Exited

	// Fatal the Fatal status
	Fatal

	// Unknown the unknown status
	Unknown
)

// Process struct
type Process struct {
	Mutex     sync.Mutex
	Config    Config
	cmd       *exec.Cmd
	Status    int
	ForceStop bool
}

type Manager struct {
	Mutex sync.Mutex
	Procs map[string]*Process
}

// Start start process
func (p *Process) Start() error {
	cmdstring := p.Config.ExecStart

	// 这里有bug，当参数值中含有空格时，会导致split出问题
	args := strings.Split(cmdstring, " ")

	p.ForceStop = false

	go func() {
		for !p.ForceStop {
			var cmd *exec.Cmd
			if len(args) > 1 {
				cmd = exec.Command(args[0], args[1:]...)
			} else {
				cmd = exec.Command(args[0])
			}
			cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
			p.Status = Starting
			err := cmd.Start()
			if err == nil {
				p.Mutex.Lock()
				p.cmd = cmd
				p.Mutex.Unlock()
				// 等待退出
				err = cmd.Wait()
				if err != nil {
					logger.Errorf(err.Error())
				}
				p.Mutex.Lock()
				p.cmd = nil
				p.Mutex.Unlock()
				p.Status = Exited
			} else {
				p.Status = Exited
				logger.Errorf(err.Error())
			}
			if p.Config.RestartSec != 0 {
				time.Sleep(time.Second * time.Duration(p.Config.RestartSec))
			} else {
				time.Sleep(time.Second * 1)
			}
		}
		p.Status = Stopped
	}()

	return nil
}

// Stop stop process
func (p *Process) Stop() (err error) {
	p.ForceStop = true
	if p.cmd != nil && p.cmd.Process != nil {
		err = p.cmd.Process.Signal(syscall.SIGINT)
	}

	return err
}

// GetStatus get status
func (p *Process) GetStatus() int {
	return p.Status
}

// IsRunning is running
func (p *Process) IsRunning() bool {
	if p.cmd != nil && p.cmd.Process != nil {
		if runtime.GOOS == "windows" {
			proc, err := os.FindProcess(p.cmd.Process.Pid)
			return proc != nil && err == nil
		}
		return p.cmd.Process.Signal(syscall.Signal(0)) == nil
	}
	return false
}

func NewManager() *Manager {
	return &Manager{
		Procs: make(map[string]*Process),
	}
}

func (m *Manager) AddProcess(conf *Config) error {
	if m.Procs[conf.Name] != nil {
		return errors.New("process exist")
	}
	p := new(Process)
	p.Config = *conf
	p.Status = Waiting
	m.Procs[conf.Name] = p
	return nil
}

func (m *Manager) DelProcess(name string) error {
	p := m.Procs[name]
	if p == nil {
		return errors.New("process is not existing")
	}
	p.Stop()
	m.Procs[name] = nil
	return nil
}

func (m *Manager) StartProcess(name string) error {
	p := m.Procs[name]
	if p == nil {
		return errors.New("process is not existing")
	}
	return p.Start()
}

func (m *Manager) StopProcess(name string) error {
	p := m.Procs[name]
	if p == nil {
		return errors.New("process is not existing")
	}
	return p.Stop()
}

func (m *Manager) IsRunning(name string) bool {
	p := m.Procs[name]
	if p == nil {
		return false
	}
	return p.IsRunning()
}

// Exist process exist
func (m *Manager) Exist(name string) bool {
	return m.Procs[name] != nil
}
