package scheduler

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type SchedulerSuite struct {
	suite.Suite

	scheduler  *Scheduler
	testServer *httptest.Server
	slowServer *httptest.Server
}

func (s *SchedulerSuite) SetupSuite() {
	s.testServer = httptest.NewServer(http.HandlerFunc(s.testHandler))
	s.slowServer = httptest.NewServer(http.HandlerFunc(s.slowHandler))
}

func (s *SchedulerSuite) TearDownSuite() {
	s.testServer.Close()
	s.slowServer.Close()
}

func (s *SchedulerSuite) SetupTest() {
	s.scheduler = New()
	go s.scheduler.Run()
}

func (s *SchedulerSuite) TestSchedulerSchedulesNewTask() {
	task := s.scheduler.Schedule(http.MethodGet, s.testServer.URL, "")
	task.Wait()
	s.Equal(task.ResponseStatus, http.StatusOK)
}

func (s *SchedulerSuite) TestTaskShouldBeStored() {
	task := s.scheduler.Schedule(http.MethodGet, s.testServer.URL, "")
	task.Wait()
	s.Equal(task, s.scheduler.TaskByID(task.ID))
}

func (s *SchedulerSuite) TestSchedulerReturnsAllTask() {
	task := s.scheduler.Schedule(http.MethodGet, s.testServer.URL, "")
	task.Wait()
	tasks := s.scheduler.Tasks()
	s.Equal(1, len(tasks))
	s.Equal(task, tasks[0])
}

func (s *SchedulerSuite) TestSchedulerDeleteAndCancelTask() {
	task := s.scheduler.Schedule(http.MethodGet, s.slowServer.URL, "")
	s.scheduler.Delete(task.ID)
	task.Wait()
	s.NotNil(task.Err)
}

func (s *SchedulerSuite) testHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("test"))
}

func (s *SchedulerSuite) slowHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(5 * time.Second)
}

func TestSchedulerSuiteSuite(t *testing.T) {
	suite.Run(t, new(SchedulerSuite))
}
