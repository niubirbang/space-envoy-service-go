package service

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sync"
)

var (
	currentDir string
)

func init() {
	var err error
	currentDir, err = os.Getwd()
	if err != nil {
		panic(err)
	}
}

type (
	Manager struct {
		mu          sync.Mutex
		serviceName string
		serviceFile string
		client      *http.Client
		inited      bool
	}
)

func NewManager(serviceName, serviceFile string) *Manager {
	if serviceName == "" {
		serviceName = "space_envoy"
	}
	if serviceFile == "" {
		serviceFile = path.Join(currentDir, "space-envoy")
	}
	if !filepath.IsAbs(serviceFile) {
		serviceFile = filepath.Join(currentDir, serviceFile)
	}
	m := &Manager{
		serviceName: serviceName,
		serviceFile: serviceFile,
	}
	m.initClient()
	return m
}

func (m *Manager) Init() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.inited {
		return nil
	}
	if !m.isRunning() {
		if err := m.install(); err != nil {
			return err
		}
	}
	m.inited = true
	return nil
}

func (m *Manager) Uninstall() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.uninstall()
}

func (m *Manager) Version() (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if err := m.checkUninit(); err != nil {
		return "", err
	}
	body, err := m.request(http.MethodGet, "/version", nil, nil)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (m *Manager) Status() (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if err := m.checkUninit(); err != nil {
		return "", err
	}
	body, err := m.request(http.MethodGet, "/status", nil, nil)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (m *Manager) Up(homeDir, configFile string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if err := m.checkUninit(); err != nil {
		return err
	}
	if !filepath.IsAbs(homeDir) {
		homeDir = filepath.Join(currentDir, homeDir)
	}
	if !filepath.IsAbs(configFile) {
		configFile = filepath.Join(currentDir, configFile)
	}
	_, err := m.request(http.MethodPost, "/up", nil, map[string]interface{}{
		"homeDir":    homeDir,
		"configFile": configFile,
	})
	return err
}

func (m *Manager) Down() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if err := m.checkUninit(); err != nil {
		return err
	}
	_, err := m.request(http.MethodPost, "/down", nil, nil)
	return err
}

func (m *Manager) Log() (string, error) {
	return m.log()
}
