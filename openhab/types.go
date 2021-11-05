package openhab

import "net/http"

type Client struct {
	http http.Client
	url  string
}
