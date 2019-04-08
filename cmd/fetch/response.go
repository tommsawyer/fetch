package main

import (
	"github.com/google/uuid"
	"github.com/tommsawyer/fetch/task"
)

type TaskResponse struct {
	ID     uuid.UUID `json:"id"`
	Status string    `json:"status"`

	Error                 string `json:"error,omitempty"`
	ResponseStatus        int    `json:"response_status,omitempty"`
	ResponseBody          string `json:"response_body,omitempty"`
	ResponseContentLength int64  `json:"response_content_length"`
}

func taskResponse(t *task.Task) TaskResponse {
	errMsg := ""
	if t.Err != nil {
		errMsg = t.Err.Error()
	}

	return TaskResponse{
		ID:                    t.ID,
		Status:                string(t.Status),
		Error:                 errMsg,
		ResponseStatus:        t.ResponseStatus,
		ResponseContentLength: t.ResponseContentLength,
		ResponseBody:          t.ResponseBody,
	}
}
