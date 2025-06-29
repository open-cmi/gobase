package http

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

type UnixSockOption struct {
	UnixSock string // unix sock 地址
}

type TLSOption struct {
	InsecureSkipVerify bool
}

type HTTPClient struct {
	Client  http.Client
	Cookies []*http.Cookie
}

func NewHTTPClient(sockopt *UnixSockOption, tlsopt *TLSOption) *HTTPClient {
	var insecureSkipVerify bool = false
	if tlsopt != nil {
		insecureSkipVerify = tlsopt.InsecureSkipVerify
	}
	dialContext := (&net.Dialer{
		Timeout:   3 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}).DialContext
	if sockopt != nil {
		dialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
			return net.Dial("unix", sockopt.UnixSock)
		}
	}

	var tp = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: insecureSkipVerify,
		},
		Proxy:                 http.ProxyFromEnvironment,
		DisableKeepAlives:     true,
		DialContext:           dialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       30 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return &HTTPClient{
		Client: http.Client{
			Transport: tp,
		},
	}
}

func (t *HTTPClient) LoadCert(certFile string, keyFile string) error {
	cliCrt, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}
	tp, ok := t.Client.Transport.(*http.Transport)
	if !ok {
		return errors.New("trans to http.Transport failed")
	}
	tp.TLSClientConfig.Certificates = []tls.Certificate{cliCrt}
	return nil
}

// ReqAPI func
func (t *HTTPClient) ReqAPI(requrl string, params map[string]string, headers map[string]string, payload interface{}) (respmsg []byte, err error) {
	urlObj, err := url.Parse(requrl)
	if err != nil {
		return respmsg, err
	}
	querys := urlObj.Query()

	for key, value := range params {
		querys.Set(key, value)
	}

	urlObj.RawQuery = querys.Encode()
	connurl := urlObj.String()

	reqbody, err := json.Marshal(payload)
	if err != nil {
		return respmsg, err
	}
	// 通过 http 请求
	req, _ := http.NewRequest("GET", connurl, bytes.NewReader(reqbody))
	for _, cookie := range t.Cookies {
		req.AddCookie(cookie)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := t.Client.Do(req)
	if err != nil {
		return respmsg, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		errmsg := fmt.Sprintf("api not available, status code: %d", resp.StatusCode)
		return respmsg, errors.New(errmsg)
	}

	t.Cookies = resp.Cookies()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return respmsg, err
	}

	return body, nil
}

// PostAPI func
func (t *HTTPClient) PostAPI(requrl string, params map[string]string, headers map[string]string, payload interface{}) (respmsg []byte, err error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return respmsg, err
	}
	// post api
	querys := url.Values{}

	urlObj, err := url.Parse(requrl)
	if err != nil {
		return respmsg, err
	}

	for key, value := range params {
		querys.Set(key, value)
	}

	urlObj.RawQuery = querys.Encode()
	connurl := urlObj.String()

	request, err := http.NewRequest("POST", connurl, bytes.NewReader(data))
	if err != nil {
		return respmsg, err
	}
	for _, cookie := range t.Cookies {
		request.AddCookie(cookie)
	}
	for key, value := range headers {
		request.Header.Set(key, value)
	}

	resp, err := t.Client.Do(request)
	// 这里需要处理出错机制，比如上传失败，需要重新上传等，这里后续完善
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return respmsg, err
	}

	if resp.StatusCode != 200 {
		errmsg := fmt.Sprintf("api not available, status code: %d", resp.StatusCode)
		return respmsg, errors.New(errmsg)
	}
	t.Cookies = resp.Cookies()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return respmsg, err
	}

	return body, nil
}

// DelAPI func
func (t *HTTPClient) DelAPI(requrl string, params map[string]string, headers map[string]string) (respmsg []byte, err error) {
	// post api
	querys := url.Values{}

	urlObj, err := url.Parse(requrl)
	if err != nil {
		return respmsg, err
	}

	for key, value := range params {
		querys.Set(key, value)
	}

	urlObj.RawQuery = querys.Encode()
	connurl := urlObj.String()

	request, err := http.NewRequest("DELETE", connurl, nil)
	if err != nil {
		return respmsg, err
	}

	for _, cookie := range t.Cookies {
		request.AddCookie(cookie)
	}

	resp, err := t.Client.Do(request)
	// 这里需要处理出错机制，比如上传失败，需要重新上传等，这里后续完善
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return respmsg, err
	}

	if resp.StatusCode != 200 {
		errmsg := fmt.Sprintf("api not available, status code: %d", resp.StatusCode)
		return respmsg, errors.New(errmsg)
	}
	t.Cookies = resp.Cookies()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return respmsg, err
	}

	return body, nil
}
