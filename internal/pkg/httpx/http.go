package httpx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type ErrorRequests struct {
	Result  bool   `json:"result"`
	Message string `json:"message"`
}

func HttpGet(url string, object interface{}) (errorCode int, err error) {
	var Client = http.Client{
		Timeout: time.Second * 20, // 20秒超时
	}

	res, err := Client.Get(url)
	if err != nil {
		return http.StatusBadGateway, err
	}
	defer func(Body io.ReadCloser) {
		BErr := Body.Close()
		if BErr != nil {
			err = BErr
		}
	}(res.Body)
	if res.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		var ErrMessage ErrorRequests
		err = json.Unmarshal(body, &ErrMessage)
		if err != nil {
			return http.StatusBadGateway, err
		}
		return res.StatusCode, fmt.Errorf(ErrMessage.Message)
	} else {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		err = json.Unmarshal(body, object)
		if err != nil {
			return http.StatusBadGateway, err
		}
		return http.StatusOK, nil
	}
}

func HttpHead(url string) (errorCode int, err error) {
	var Client http.Client = http.Client{
		Timeout: time.Second * 5, // 5秒超时
	}

	res, err := Client.Head(url)
	if err != nil {
		return http.StatusBadGateway, err
	}
	if res.StatusCode != http.StatusOK {
		return http.StatusNotFound, fmt.Errorf("未找到资源")
	}
	return http.StatusOK, nil
}

func HttpPost(url string, payload interface{}, object interface{}) (errorCode int, err error) {
	var Client http.Client = http.Client{
		Timeout: time.Second * 5, // 5秒超时
	}

	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	res, err := Client.Post(url, "application/json", bytes.NewReader(payloadJson))
	if err != nil {
		return http.StatusBadGateway, err
	}
	if res.StatusCode != http.StatusOK {
		return http.StatusNotFound, fmt.Errorf("未找到资源")
	}
	defer func(Body io.ReadCloser) {
		BErr := Body.Close()
		if BErr != nil {
			err = BErr
		}
	}(res.Body)
	if res.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		err = json.Unmarshal(body, object)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}
	return http.StatusOK, nil
}
