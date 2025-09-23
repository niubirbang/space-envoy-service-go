package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
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
	CoreConfig struct {
		Dir            string `json:"dir"`
		LogLevel       string `json:"logLevel"`
		MixedPort      uint16 `json:"mixedPort"`
		ControllerPort uint16 `json:"controllerPort"`
		DnsPort        uint16 `json:"dnsPort"`
	}
	Rules     []Rule
	Rule      string
	Proxies   []Proxy
	Proxy     string
	BaseParam struct {
		TunEnable   bool    `json:"tunEnable" form:"tunEnable"`
		DirectRules Rules   `json:"directRules" form:"directRules"`
		ProxyRules  Rules   `json:"proxyRules" form:"proxyRules"`
		RejectRules Rules   `json:"rejectRules" form:"rejectRules"`
		Proxies     Proxies `json:"proxies" form:"proxies"`
	}
	GlobalParam struct {
		BaseParam
	}
	Param interface {
		_Param()
	}
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

func (*GlobalParam) _Param()

func (m *Manager) Init(dir, logLevel string, mixedPort, controllerPort, dnsPort uint16) error {
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
	if _, err := m.request(http.MethodPost, "/init", nil, map[string]interface{}{
		"dir":            dir,
		"logLevel":       logLevel,
		"mixedPort":      mixedPort,
		"controllerPort": controllerPort,
		"dnsPort":        dnsPort,
	}); err != nil {
		return err
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

func (m *Manager) Args() (*CoreConfig, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if err := m.checkUninit(); err != nil {
		return nil, err
	}
	body, err := m.request(http.MethodGet, "/args", nil, nil)
	if err != nil {
		return nil, err
	}
	var args CoreConfig
	if err := json.Unmarshal(body, &args); err != nil {
		return nil, err
	}
	return &args, nil
}

func (m *Manager) Config(mode string, param Param) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if err := m.checkUninit(); err != nil {
		return err
	}
	_, err := m.request(http.MethodPost, fmt.Sprintf("/config/%s", mode), nil, param)
	return err
}

func (m *Manager) Up() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if err := m.checkUninit(); err != nil {
		return err
	}
	_, err := m.request(http.MethodPost, "/up", nil, nil)
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

func (m *Manager) Ping(target string, port uint16, timeout uint16) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if err := m.checkUninit(); err != nil {
		return 0, err
	}
	body, err := m.request(http.MethodPost, "/ping", nil, map[string]interface{}{
		"target":  target,
		"port":    port,
		"timeout": timeout,
	})
	delay, err := strconv.ParseInt(string(body), 10, 64)
	if err != nil {
		return 0, err
	}
	return delay, nil
}

func (m *Manager) Log() (string, error) {
	return m.log()
}
