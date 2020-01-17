package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/echoturing/log"
	"github.com/labstack/echo"
)

var (
	defaultTimeout = time.Second * 3
	defaultClient  = http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		Timeout: defaultTimeout,
	}
)

func Get(ctx context.Context, url string, resp interface{}, header http.Header) error {
	return doRequest(ctx, http.MethodGet, url, nil, resp, header)
}

func Post(ctx context.Context, url string, payload interface{}, resp interface{}, header http.Header) error {
	return doRequest(ctx, http.MethodPost, url, payload, resp, header)
}

func doRequest(ctx context.Context, method, url string, payload interface{}, resp interface{}, header http.Header) error {
	var cancel context.CancelFunc
	if _, ok := ctx.Deadline(); !ok {
		ctx, cancel = context.WithTimeout(ctx, defaultTimeout)
		defer cancel()
	}
	var (
		payloadJson []byte
		err         error
	)
	if payload != nil {
		payloadJson, err = json.Marshal(payload)
		if err != nil {
			return err
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(payloadJson))
	if err != nil {
		return err
	}

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	if header != nil {
		for key := range header {
			req.Header[key] = header[key]
		}
	}
	req = req.WithContext(ctx)
	// header首先要把x-request-id传走
	requestIDWithUser := log.FromContext(ctx)
	req.Header.Add(echo.HeaderXRequestID, requestIDWithUser.RequestID)
	response, err := defaultClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	respData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(respData, resp)
}
