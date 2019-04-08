package main

import (
	"flag"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/tommsawyer/fetch/scheduler"
)

func main() {
	listen := flag.String("listen", ":9090", "a port to listen to")
	flag.Parse()

	app := gin.Default()

	sch := scheduler.New()
	go sch.Run()
	handler := TaskHandler{
		sch: sch,
	}

	app.POST("/task", handler.CreateNewTask)
	app.GET("/task/", handler.Tasks)
	app.GET("/task/:id", handler.GetTask)
	app.DELETE("/task/:id", handler.Remove)

	if err := app.Run(*listen); err != nil {
		log.Fatalf("cannot run app: %v", err)
	}
}
