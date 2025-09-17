package ses

import "net/http"

func Check() error {
	if _, err := Version(); err == nil {
		return nil
	}
	return Install()
}

func Version() (string, error) {
	body, err := request(http.MethodGet, "/version", nil, nil)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func Start(homeDir, configFile string) error {
	_, err := request(http.MethodPost, "/start", nil, map[string]interface{}{
		"homeDir":    homeDir,
		"configFile": configFile,
	})
	return err
}

func Stop() error {
	_, err := request(http.MethodPost, "/stop", nil, nil)
	return err
}
