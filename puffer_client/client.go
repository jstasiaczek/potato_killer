package puffer_client

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type PufferClient struct {
	PufferUrl  string
	HttpClient *http.Client
	Token      string
	ServerId   string
}

func NewPufferClient(url string, serverId string) *PufferClient {
	client := PufferClient{}
	client.PufferUrl = url
	client.ServerId = serverId
	client.HttpClient = &http.Client{
		Timeout: time.Second * 3,
	}
	return &client
}

func (client *PufferClient) newHTTPRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+client.Token)
	return req, nil
}

type loginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Session string `json:"session"`
}

func newLoginDataJson(email, password string) io.Reader {
	data := loginData{
		Email:    email,
		Password: password,
	}
	raw, _ := json.Marshal(data)
	return bytes.NewReader(raw)
}

func (client *PufferClient) StartServer() error {
	return client.doServerAction("start")
}

func (client *PufferClient) StopServer() error {
	return client.doServerAction("stop")
}

func (client *PufferClient) KillServer() error {
	return client.doServerAction("kill")
}

func (client *PufferClient) doServerAction(action string) error {
	url := client.PufferUrl + "/proxy/daemon/server/" + client.ServerId + "/" + action
	req, err := client.newHTTPRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}
	_, err = client.HttpClient.Do(req)
	if err != nil {
		return err
	}
	return nil
}

func (client *PufferClient) Login(email, password string) error {
	data := newLoginDataJson(email, password)
	resp, err := client.HttpClient.Post(client.PufferUrl+"/auth/login/", "application/json", data)
	if err != nil {
		return err
	}
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	respData := loginResponse{}
	err = json.Unmarshal(raw, &respData)
	if err != nil {
		return err
	}
	client.Token = respData.Session

	return nil
}

func (client *PufferClient) WillTokenExpire() bool {
	if client.Token == "" {
		return true
	}
	token, _ := jwt.Parse(client.Token, nil)
	expDate := int32(token.Claims.(jwt.MapClaims)["exp"].(float64))
	now := int32(time.Now().Unix())
	if expDate-now < 300 {
		return true
	}
	return false
}
