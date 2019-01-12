package owo

import "net/http"

type Client struct {
	token  string
	client *http.Client
}

type Result struct {
	Success bool `json:"success"`
	Files   []struct {
		Hash string `json:"hash"`
		Name string `json:"name"`
		URL  string `json:"url"`
		Size int    `json:"size"`
	} `json:"files"`
}
