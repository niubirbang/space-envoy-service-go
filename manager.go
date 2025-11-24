package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path"
	"runtime"
	"sync"
	"time"
)

const (
	StatusEnable  Status = "enable"
	StatusDisable Status = "disable"

	UseFixed Use = "fixed"
	UseAuto  Use = "auto"

	ModeGlobal    Mode = "global"
	ModeAbroad    Mode = "abroad"
	ModeReturning Mode = "returning"
)

var (
	serverFileName string
	StatusAll      = map[Status]bool{
		StatusEnable:  true,
		StatusDisable: true,
	}
	UseAll = map[Use]bool{
		UseFixed: true,
		UseAuto:  true,
	}
	ModeAll = map[Mode]bool{
		ModeGlobal:    true,
		ModeAbroad:    true,
		ModeReturning: true,
	}
)

func init() {
	switch runtime.GOOS {
	case "windows":
		serverFileName = "service.exe"
	case "darwin":
		serverFileName = "service"
	case "linux":
		serverFileName = "service"
	}
}

type (
	State struct {
		ServerInstallerExists bool `json:"serverInstallerExists"`
		ServerFileExists      bool `json:"serverFileExists"`
		ServerIsRunning       bool `json:"serverIsRunning"`
	}
	Option struct {
		Dir            string `json:"dir"`            // set invalid
		LogLevel       string `json:"logLevel"`       // set valid
		MixedPort      uint16 `json:"mixedPort"`      // set invalid
		ControllerPort uint16 `json:"controllerPort"` // set invalid
		DnsPort        uint16 `json:"dnsPort"`        // set invalid
	}
	URIInfo struct {
		Name   string `json:"name"`
		Type   string `json:"type"`
		Server string `json:"server"`
		Port   uint16 `json:"port"`
	}
	Status  string
	Use     string
	Mode    string
	Rules   []Rule
	Rule    string
	Proxies []Proxy
	Proxy   string
	Param   struct {
		Mode        Mode    `json:"mode"`
		Use         Use     `json:"use"`
		DirectRules Rules   `json:"directRules"`
		ProxyRules  Rules   `json:"proxyRules"`
		RejectRules Rules   `json:"rejectRules"`
		Proxies     Proxies `json:"proxies"`
	}

	Manager struct {
		mu                    sync.RWMutex
		serverName            string
		serverInstaller       string
		serverInstallerExists bool
		serverFile            string
		serverFileExists      bool
		serverIsRunning       bool
		stateListeners        []func(State)
		client                *http.Client
	}
)

func NewManager(serverName, serverDir string) *Manager {
	if serverName == "" {
		serverName = "space_envoy"
	}
	m := &Manager{
		serverName:      serverName,
		serverInstaller: path.Join(serverDir, "service.zip"),
		serverFile:      path.Join(serverDir, serverFileName),
	}
	if _, err := os.Stat(m.serverInstaller); err == nil {
		m.serverInstallerExists = true
	}
	if _, err := os.Stat(m.serverFile); err == nil {
		m.serverFileExists = true
	}
	m.initClient()
	m.intervalRefreshState()
	m.listenServerIsRunningByClient()
	return m
}

func (m *Manager) GetServerIsRunningByServer() bool {
	return m.getServerIsRunningByServer()
}

func (m *Manager) Download(downloadFunc func(serverInstaller string) error) error {
	return downloadFunc(m.serverInstaller)
}

func (m *Manager) Install() error {
	return m.install()
}

func (m *Manager) Uninstall() error {
	return m.uninstall()
}

func (m *Manager) ListenState(callback func(State)) {
	if callback != nil {
		m.addStateListeners(callback)
	}
}

func (m *Manager) Version() (string, error) {
	if err := m.check(); err != nil {
		return "", err
	}
	body, err := m.request(http.MethodGet, "/version", nil, nil)
	if err != nil {
		return "", err
	}
	var data string
	err = json.Unmarshal(body, &data)
	return data, err
}

func (m *Manager) Option() (*Option, error) {
	if err := m.check(); err != nil {
		return nil, err
	}
	body, err := m.request(http.MethodGet, "/option", nil, nil)
	if err != nil {
		return nil, err
	}
	var data Option
	err = json.Unmarshal(body, &data)
	return &data, err
}

func (m *Manager) SetOption(opt Option) error {
	if err := m.check(); err != nil {
		return err
	}
	_, err := m.request(http.MethodPost, "/option", nil, opt)
	return err
}

func (m *Manager) ParseURI(uri string) (*URIInfo, error) {
	if err := m.check(); err != nil {
		return nil, err
	}
	body, err := m.request(http.MethodPost, "/parseuri", nil, map[string]interface{}{
		"uri": uri,
	})
	if err != nil {
		return nil, err
	}
	var data URIInfo
	err = json.Unmarshal(body, &data)
	return &data, err
}

func (m *Manager) Ping(target string, port uint16, timeout uint16) (time.Duration, error) {
	if timeout == 0 {
		timeout = 2000
	}
	if err := m.check(); err != nil {
		return 0, err
	}
	body, err := m.request(http.MethodPost, "/ping", nil, map[string]interface{}{
		"target":  target,
		"port":    port,
		"timeout": timeout,
	})
	if err != nil {
		return 0, err
	}
	var data int64
	err = json.Unmarshal(body, &data)
	return time.Duration(data) * time.Millisecond, err
}

func (m *Manager) Status() (Status, error) {
	if err := m.check(); err != nil {
		return "", err
	}
	body, err := m.request(http.MethodGet, "/status", nil, nil)
	if err != nil {
		return "", err
	}
	var data Status
	err = json.Unmarshal(body, &data)
	return data, err
}

func (m *Manager) Enable(param Param) error {
	if err := m.check(); err != nil {
		return err
	}
	_, err := m.request(http.MethodPost, "/enable", nil, param)
	return err
}

func (m *Manager) Disable() error {
	if err := m.check(); err != nil {
		return err
	}
	_, err := m.request(http.MethodPost, "/disable", nil, nil)
	return err
}

func (m *Manager) Log() (string, error) {
	return m.log()
}

func (m *Manager) setServerInstallerExists(exists bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.serverInstallerExists = exists
}
func (m *Manager) GetServerInstallerExists() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.serverInstallerExists
}
func (m *Manager) setServerFileExists(exists bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.serverFileExists = exists
}
func (m *Manager) GetServerFileExists() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.serverFileExists
}
func (m *Manager) setServerIsRunning(isRunning bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.serverIsRunning = isRunning
}
func (m *Manager) GetServerIsRunning() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.serverIsRunning
}
func (m *Manager) addStateListeners(callback func(State)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.stateListeners = append(m.stateListeners, callback)
}
func (m *Manager) getStateListeners() []func(State) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.stateListeners
}

func (m *Manager) check() error {
	if !m.GetServerIsRunning() {
		return errors.New("server_not_run")
	}
	return nil
}

func (m *Manager) notifyState() {
	state := State{
		ServerInstallerExists: m.GetServerInstallerExists(),
		ServerFileExists:      m.GetServerFileExists(),
		ServerIsRunning:       m.GetServerIsRunning(),
	}
	for _, cb := range m.getStateListeners() {
		go cb(state)
	}
}

func (m *Manager) refreshState() {
	var (
		serverInstallerExists, serverFileExists, serverIsRunning, notify bool
	)
	if _, err := os.Stat(m.serverInstaller); err == nil {
		serverInstallerExists = true
	}
	if _, err := os.Stat(m.serverFile); err == nil {
		serverFileExists = true
	}
	serverIsRunning = m.getServerIsRunningByServer()

	if m.GetServerInstallerExists() != serverInstallerExists {
		m.setServerInstallerExists(serverInstallerExists)
		notify = true
	}
	if m.GetServerFileExists() != serverFileExists {
		m.setServerFileExists(serverFileExists)
		notify = true
	}
	if m.GetServerIsRunning() != serverIsRunning {
		m.setServerIsRunning(serverIsRunning)
		notify = true
	}
	if notify {
		m.notifyState()
	}
}

func (m *Manager) intervalRefreshState() {
	go func() {
		for {
			m.refreshState()
			time.Sleep(200 * time.Millisecond)
		}
	}()
}

func (m *Manager) getServerIsRunningByClient() bool {
	_, err := m.request(http.MethodGet, "", nil, nil)
	return err == nil
}

func (m *Manager) listenServerIsRunningByClient() {
	go func() {
		for {
			m.setServerIsRunning(m.getServerIsRunningByClient())
			time.Sleep(time.Second)
		}
	}()
}

func (m *Manager) installServerAfterCheck() error {
	var ok bool
	for i := 0; i < 60; i++ {
		time.Sleep(500 * time.Millisecond)
		isRunningByServer := m.getServerIsRunningByServer()
		isRunningByClient := m.getServerIsRunningByClient()
		if isRunningByServer && isRunningByClient {
			ok = true
			break
		}
	}
	if !ok {
		return errors.New("server_not_run")
	}
	m.refreshState()
	return nil
}
