package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func HttpGet(url string) (string, error) {
	resp, err := http.Get(url)
	response := ""
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}
	response = string(body)
	return response, err
}

func HttpPut(url string, requestString string) (string, error) {
	payload := strings.NewReader(requestString)
	req, err := http.NewRequest("PUT", url, payload)
	if err != nil {
		return "", fmt.Errorf("HttpPut err : %s %s ", url, err)
	}
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("HttpPut err : %s %s ", url, err)
	}
	defer response.Body.Close()
	responseText, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("%s", err)
	}
	return string(responseText), nil
}

// HttpPostForm TODO
func HttpPostForm(url string, m map[string]string) (string, error) {
	response := ""
	var r http.Request
	r.ParseForm()
	for k, v := range m {
		r.Form.Add(k, v)
	}
	body := strings.NewReader(r.Form.Encode())
	resp, err := http.Post(url, "application/x-www-form-urlencoded", body)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()
	responseByte, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		return response, err2
	}
	return string(responseByte), nil
}
