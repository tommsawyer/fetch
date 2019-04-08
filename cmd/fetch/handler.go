package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tommsawyer/fetch/scheduler"
)

type TaskHandler struct {
	sch *scheduler.Scheduler
}

func (h *TaskHandler) CreateNewTask(ctx *gin.Context) {
	var req CreateNewTaskRequest

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task := h.sch.Schedule(req.Method, req.URL, req.Body)
	if task.Err != nil {
		ctx.JSON(http.StatusOK, gin.H{"error": task.Err.Error()})
		return
	}

	task.Wait()

	ctx.JSON(http.StatusOK, taskResponse(task))
}

func (h *TaskHandler) GetTask(ctx *gin.Context) {
	rawID := ctx.Param("id")
	id, err := uuid.Parse(rawID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task := h.sch.TaskByID(id)
	if task == nil {
		ctx.JSON(http.StatusOK, gin.H{"error": "task not found"})
		return
	}

	ctx.JSON(http.StatusOK, taskResponse(task))
}

func (h *TaskHandler) Tasks(ctx *gin.Context) {
	tasks := h.sch.Tasks()
	tasksResponse := make([]TaskResponse, len(tasks))
	for i, t := range tasks {
		tasksResponse[i] = taskResponse(t)
	}

	ctx.JSON(http.StatusOK, tasksResponse)
}

func (h *TaskHandler) Remove(ctx *gin.Context) {
	rawID := ctx.Param("id")
	id, err := uuid.Parse(rawID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.sch.Delete(id)
	ctx.JSON(http.StatusOK, gin.H{})
}
