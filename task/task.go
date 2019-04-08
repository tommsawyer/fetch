// Package task provides fetch task.
package task

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

// HTTPClient represents a http client.
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

var client HTTPClient = http.DefaultClient

// SetClient sets http client that will be used by all requests.
func SetClient(c HTTPClient) {
	client = c
}

// Status represents status of task.
type Status string

const (
	// StatusScheduled sets when task is scheduled.
	StatusScheduled Status = "scheduled"
	// StatusRunning sets when task is running.
	StatusRunning Status = "running"
	// StatusFinished sets when task is finished.
	StatusFinished Status = "finished"
)

// Task represents fetch task.
type Task struct {
	ID     uuid.UUID
	Status Status

	Method      string
	URL         string
	RequestBody string

	Err                   error
	ResponseStatus        int
	ResponseBody          string
	ResponseContentLength int64

	done   chan struct{}
	ctx    context.Context
	cancel func()
}

// New creates a new task.
func New(method, url string, body string) *Task {
	ctx, cancel := context.WithCancel(context.Background())
	return &Task{
		ID:          uuid.New(),
		Status:      StatusScheduled,
		Method:      method,
		URL:         url,
		RequestBody: body,
		ctx:         ctx,
		cancel:      cancel,
		done:        make(chan struct{}),
	}
}

// Run runs task.
func (t *Task) Run() {
	defer func() {
		t.Status = StatusFinished
		close(t.done)
	}()

	t.Status = StatusRunning
	req, err := http.NewRequest(t.Method, t.URL, strings.NewReader(t.RequestBody))
	if err != nil {
		t.Err = fmt.Errorf("cannot build request: %v", err)
		return
	}

	rsp, err := client.Do(req.WithContext(t.ctx))
	if err != nil {
		t.Err = fmt.Errorf("cannot send request: %v", err)
		return
	}
	defer rsp.Body.Close()

	bs, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		t.Err = fmt.Errorf("cannot read response body: %v", err)
		return
	}

	t.ResponseStatus = rsp.StatusCode
	t.ResponseBody = string(bs)
	t.ResponseContentLength = int64(len(bs))
}

// Wait will wait until task finished.
func (t *Task) Wait() {
	<-t.done
}

// Cancel will cancel task if it was running.
func (t *Task) Cancel() {
	t.cancel()
	<-t.done
}
