package util

import (
	"io/ioutil"
	"net/http"
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

func HttpPostForm(url string,m[string]string)(string, error) {
	//TODO:
	response:=""

	var r http.Request

	r.ParseForm()
	for k,v := range m{
		r.Form.Add(k, v)
	}
body := strings.NewReader(r.Form.Encode())
resp, err :=http.Post("xxxx", "application/x-www-form-urlencoded", body)


    if err != nil {
		return response, err
    }
 
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
		return response, err
    }
	return string(body), nil
}
