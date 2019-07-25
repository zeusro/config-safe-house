package util

import (
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

func HttpPostForm(url string, m map[string]string) (string, error) {
	//TODO:
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
