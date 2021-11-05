package openhab

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func InitClient(url string, timeout int) *Client {
	client := http.Client{
		Timeout: time.Second * time.Duration(timeout),
	}
	return &Client{http: client, url: url}
}

func (cl *Client) GetItems() ([]Item, error) {
	resp, err := cl.http.Get(fmt.Sprintf("%s/rest/items", cl.url))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	result := []Item{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("can not unmarshal JSON: %s", err.Error())
	}
	return result, nil
}
