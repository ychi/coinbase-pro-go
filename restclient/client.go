package restclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ychi/coinbase-pro-go/signature"
	"net/http"
	"strconv"
	"time"
)

type RestClient interface {
	Request(
		method string,
		path string,
		query map[string]string,
		params, result interface{},
	) (*http.Response, error)
}

type restClient struct {
	client        *http.Client
	baseURL       string
	apiKey        string
	apiSecret     string
	apiPassphrase string
	retryCount    int
	timeout       int
}

func NewRestClient() *restClient {
	return &restClient{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		baseURL:    "https://api-public.sandbox.pro.coinbase.com",
		retryCount: 3,
	}
}

func (rc *restClient) SetBaseURL(baseURL string) *restClient {
	rc.baseURL = baseURL
	return rc
}

func (rc *restClient) SetApiKey(apiKey string) *restClient {
	rc.apiKey = apiKey
	return rc
}

func (rc *restClient) SetApiSecret(apiSecret string) *restClient {
	rc.apiSecret = apiSecret
	return rc
}

func (rc *restClient) SetApiPassphrase(apiPassphrase string) *restClient {
	rc.apiPassphrase = apiPassphrase
	return rc
}

func (rc *restClient) SetRetryCount(retryCount int) *restClient {
	rc.retryCount = retryCount
	return rc
}

func (rc *restClient) SetTimeout(timeout int) *restClient {
	rc.client.Timeout = time.Duration(timeout) * time.Second
	return rc
}

func (rc *restClient) Request(
	method string,
	path string,
	query map[string]string,
	params, result interface{},
) (
	res *http.Response,
	err error,
) {
	for i := 0; i <= rc.retryCount; i++ {
		res, err = rc.request(method, path, query, params, result)

		if res != nil && res.StatusCode == 429 {
			time.Sleep(1000)
			continue
		} else {
			break
		}
	}
	return res, err
}

func (rc *restClient) request(
	method string,
	path string,
	query map[string]string,
	params, result interface{},
) (
	res *http.Response,
	err error,
) {
	fullURL := fmt.Sprintf("%s%s", rc.baseURL, path)
	var pBytes []byte
	if params != nil {
		pBytes, err = json.Marshal(params)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, fullURL, bytes.NewReader(pBytes))
	if err != nil {
		return nil, err
	}

	if query != nil {
		q := req.URL.Query()
		for k, v := range query {
			q.Add(k, v)
		}
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	p := fmt.Sprintf("%s%s%s%s", timestamp, method, path, string(pBytes))
	signed, err := signature.SignPayload(p, rc.apiSecret)
	if err != nil {
		return nil, err
	}

	req.Header.Add("CB-ACCESS-KEY", rc.apiKey)
	req.Header.Add("CB-ACCESS-TIMESTAMP", timestamp)
	req.Header.Add("CB-ACCESS-PASSPHRASE", rc.apiPassphrase)
	req.Header.Add("CB-ACCESS-SIGN", signed)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "Coinbase Pro go")

	res, err = rc.client.Do(req)

	if err != nil {
		return res, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		serverErr := Error{}
		decoder := json.NewDecoder(res.Body)
		err := decoder.Decode(&serverErr)
		if err != nil {
			return res, err
		}
		return res, error(serverErr)
	}

	if result != nil {
		decoder := json.NewDecoder(res.Body)
		err := decoder.Decode(result)
		if err != nil {
			return res, err
		}
	}

	return res, err
}
