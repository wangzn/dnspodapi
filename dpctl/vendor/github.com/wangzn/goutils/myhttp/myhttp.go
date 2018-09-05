// @Author: wangzn04@gmail.com
// @Date: 2017-08-30 13:46:04

package myhttp

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
)

// HRunBasicAuth runs a http request, with basic auth info
func HRunBasicAuth(method, u string, params url.Values, reqBody []byte,
	user, pwd string) (
	int, []byte, error,
) {
	var req *http.Request
	var err error
	client := &http.Client{}
	if len(reqBody) > 0 {
		req, err = http.NewRequest(method, u, bytes.NewReader(reqBody))
		req.ContentLength = int64(len(reqBody))
	} else {
		req, err = http.NewRequest(method, u, nil)
	}
	if params != nil {
		req.URL.RawQuery = params.Encode()
	}
	if user != "" || pwd != "" {
		req.SetBasicAuth(user, pwd)
	}
	resp, err := client.Do(req)
	if err != nil {
		if resp != nil {
			return resp.StatusCode, nil, err
		}
		return 0, nil, err
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, err
	}
	return resp.StatusCode, contents, err
}

// HRun runs a http request, with method, url and request body
func HRun(method, u string, params url.Values, reqBody []byte) (
	int, []byte, error,
) {
	return HRunBasicAuth(method, u, params, reqBody, "", "")
}

// HGet makes a HTTP Get request
func HGet(u string, params url.Values) (int, []byte, error) {
	return HRun("GET", u, params, nil)
}

// HPost makes a HTTP Post request
func HPost(u string, params url.Values, reqBody []byte, bodyType string) (
	int, []byte, error,
) {
	if bodyType == "" {
		bodyType = "application/x-www-form-urlencoded"
	}
	resp, err := http.Post(u, bodyType, bytes.NewReader(reqBody))
	if err != nil {
		if resp != nil {
			return resp.StatusCode, nil, err
		}
		return 0, nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, err
	}
	return resp.StatusCode, body, nil
}

// HPostForm makes a HTTP Post request in form
func HPostForm(u string, params url.Values, form url.Values) (
	int, []byte, error,
) {
	var err error
	client := &http.Client{}
	if params != nil {
		u = u + "?" + params.Encode()
	}
	resp, err := client.PostForm(u, form)
	if err != nil {
		if resp != nil {
			return resp.StatusCode, nil, err
		}
		return 0, nil, err
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, err
	}
	return resp.StatusCode, contents, err
}

// HPut makes a HTTP Put request
func HPut(url string, params url.Values, reqBody []byte) (int, []byte, error) {
	return HRun("PUT", url, params, reqBody)
}

// HDelete makes a HTTP Delete request
func HDelete(url string, params url.Values, reqBody []byte) (int, []byte, error) {
	return HRun("DELETE", url, params, reqBody)
}
