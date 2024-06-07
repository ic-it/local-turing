package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type TuringClient struct {
	client *http.Client
	url    string
	jar    *cookiejar.Jar
}

func NewTuringClient(url string) (*TuringClient, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	return &TuringClient{
		client: &http.Client{
			Jar: jar,
		},
		url: url,
		jar: jar,
	}, nil
}

func (c *TuringClient) Login(username, password string) error {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	loginRequest := LoginRequest{
		Username: username,
		Password: password,
	}
	data, err := json.Marshal(loginRequest)
	if err != nil {
		logger.Errorw("login failed", "error", err)
		return err
	}
	req, err := http.NewRequest("POST", c.url+"/login", bytes.NewBuffer(data))
	if err != nil {
		logger.Errorw("login failed", "error", err)
		return err
	}
	req.Header.Set("x-cid", uuid.New().String())
	resp, err := c.client.Do(req)
	if err != nil {
		logger.Errorw("login failed", "error", err)
		return err
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorw("login failed", "error", err)
		return err
	}
	if resp.StatusCode != 200 {
		logger.Errorw("login failed", "status", resp.Status, "body", string(bodyText))
		return ErrLoginFailed
	}
	cookies := resp.Cookies()
	if len(cookies) == 0 {
		return ErrLoginFailed
	}
	c.jar.SetCookies(&url.URL{Scheme: "https", Host: c.url}, cookies)
	return nil
}

func (c *TuringClient) SaveAssigment(assignment *LocalTuringAssignment) error {
	type SubmitRunRequest struct {
		Type       string `json:"type"`
		SourceCode string `json:"sourceCode"`
		Stdin      string `json:"stdin"`
	}
	path := filepath.Join(assignment.Dir, assignment.MainFile)
	_, err := os.Stat(path)
	if err != nil {
		return err
	}
	assignmentSrc, err := os.ReadFile(path)
	if err != nil {
		logger.Errorw("read assignment failed", "error", err)
		return err
	}
	submitRunRequest := SubmitRunRequest{
		Type:       "SubmitRun",
		SourceCode: string(assignmentSrc),
		Stdin:      "-",
	}
	data, err := json.Marshal(submitRunRequest)
	if err != nil {
		logger.Errorw("submit run failed", "error", err)
		return err
	}
	req, err := http.NewRequest("POST", c.GetAssignmentLink(assignment), bytes.NewBuffer(data))
	if err != nil {
		logger.Errorw("submit run failed", "error", err)
		return err
	}
	req.Header.Set("x-cid", uuid.New().String())
	resp, err := c.client.Do(req)
	if err != nil {
		logger.Errorw("submit run failed", "error", err)
		return err
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorw("submit run failed", "error", err)
		return err
	}
	if resp.StatusCode != 200 {
		logger.Errorw("submit run failed", "status", resp.Status, "body", string(bodyText))
		return ErrSaveAssignmentFailed
	}
	return nil
}

func (c *TuringClient) GetAssignmentLink(assignment *LocalTuringAssignment) string {
	return fmt.Sprintf("%s/teap?Solve=%s", c.url, assignment.PushName)
}

var (
	ErrLoginFailed          = fmt.Errorf("login failed")
	ErrSaveAssignmentFailed = fmt.Errorf("save assignment failed")
)
