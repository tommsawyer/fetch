package main

type CreateNewTaskRequest struct {
	Method string `json:"method"`
	URL    string `json:"url"`
	Body   string `json:"body"`
}
