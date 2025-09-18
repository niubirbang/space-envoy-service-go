package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

func (m *Manager) request(
	method,
	path string,
	query url.Values,
	params interface{},
) ([]byte, error) {
	paramsBody, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(
		method,
		fmt.Sprintf("http://unix%s?%s", path, query.Encode()),
		bytes.NewReader(paramsBody),
	)
	if err != nil {
		return nil, err
	}
	resp, err := m.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return body, nil
	}
	return nil, errors.New(string(body))
}

func (m *Manager) checkUninit() error {
	if !m.inited {
		return errors.New("uninit")
	}
	return nil
}

func (m *Manager) installAfterCheck() error {
	var ok bool
	for i := 0; i < 60; i++ {
		time.Sleep(500 * time.Millisecond)
		if m.isRunning() {
			ok = true
			break
		}
	}
	if !ok {
		return errors.New("service failed to start")
	}
	return nil
}
