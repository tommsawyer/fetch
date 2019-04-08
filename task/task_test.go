package task

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type TaskSuite struct {
	suite.Suite

	testTask                 *Task
	testServer               *httptest.Server
	testServerRecievedMethod string
	testServerRecievedBody   string

	slowServerTask     *Task
	slowServer         *httptest.Server
	slowServerResponse *sync.WaitGroup
}

func (s *TaskSuite) SetupSuite() {
	s.testServer = httptest.NewServer(http.HandlerFunc(s.testHandler))
	s.slowServer = httptest.NewServer(http.HandlerFunc(s.slowHandler))
	s.slowServerResponse = new(sync.WaitGroup)
}

func (s *TaskSuite) TearDownSuite() {
	s.testServer.Close()
	s.slowServer.Close()
}

func (s *TaskSuite) SetupTest() {
	s.testTask = New(http.MethodGet, s.testServer.URL, "test")
	s.slowServerTask = New(http.MethodGet, s.slowServer.URL, "test")
}

func (s *TaskSuite) TestTaskConstructAndSendRequestProperly() {
	s.testTask.Run()
	s.NotEqual(uuid.Nil, s.testTask.ID, "should generate unique id")
	s.Nil(s.testTask.Err)
	s.Equal(http.MethodGet, s.testServerRecievedMethod)
	s.Equal("test", s.testServerRecievedBody)
}

func (s *TaskSuite) TestTaskFailsIfCannotBuildRequest() {
	const wrongHTTPMethod = "@"

	task := New(wrongHTTPMethod, "", "")
	task.Run()
	s.EqualError(task.Err, `cannot build request: net/http: invalid method "@"`)
}

func (s *TaskSuite) TestTaskFailsIfCannotSendRequest() {
	task := New(http.MethodGet, "wrong_url", "")
	task.Run()
	s.EqualError(task.Err, `cannot send request: Get wrong_url: unsupported protocol scheme ""`)
}

func (s *TaskSuite) TestTaskSetsRequestResult() {
	s.testTask.Run()
	s.Nil(s.testTask.Err)
	s.Equal(s.testTask.ResponseStatus, http.StatusOK)
	s.Equal(s.testTask.ResponseBody, "test")
	s.Equal(s.testTask.ResponseContentLength, int64(4))
}

func (s *TaskSuite) TestTaskProperlySetStatuses() {
	s.Equal(s.testTask.Status, StatusScheduled)
	s.testTask.Run()
	s.Equal(s.testTask.Status, StatusFinished)
}

func (s *TaskSuite) TestTaskCanCancelRequest() {
	go s.slowServerTask.Run()

	s.slowServerTask.Cancel()
	s.slowServerTask.Wait()

	s.EqualError(s.slowServerTask.Err, fmt.Sprintf("cannot send request: Get %s: context canceled", s.slowServer.URL))
}

func (s *TaskSuite) TestCanSetDifferentHTTPClient() {
	c := &http.Client{}
	SetClient(c)
	s.Equal(c, client)
}

func (s *TaskSuite) testHandler(w http.ResponseWriter, r *http.Request) {
	s.testServerRecievedMethod = r.Method
	body, _ := ioutil.ReadAll(r.Body)
	s.testServerRecievedBody = string(body)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("test"))
}

func (s *TaskSuite) slowHandler(w http.ResponseWriter, r *http.Request) {
	s.slowServerResponse.Wait()
}

func TestTaskSuite(t *testing.T) {
	suite.Run(t, new(TaskSuite))
}
