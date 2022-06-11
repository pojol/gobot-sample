package main

import (
	"context"
	"flag"
	"fmt"
	"runtime"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pojol/gobot/factory"
	"github.com/pojol/gobot/server"
)

var (
	help bool

	dbmode      bool
	reportLimit int
	batchSize   int
	scriptPath  string
)

func initFlag() {
	flag.BoolVar(&help, "h", false, "this help")

	flag.BoolVar(&dbmode, "no_database", false, "Run in local mode")
	flag.IntVar(&reportLimit, "report_limit", 100, "Report retention limit")
	flag.IntVar(&batchSize, "batch_size", 1024, "The maximum number of robots in parallel")
	flag.StringVar(&scriptPath, "script_path", "script/", "Path to bot script")
}

func main() {

	defer func() {
		if err := recover(); err != nil {
			var buf [4096]byte
			n := runtime.Stack(buf[:], false)
			fmt.Println("panic:", string(buf[:n]))
		}
	}()

	initFlag()
	flag.Parse()
	if help {
		flag.Usage()
		return
	}

	_, err := factory.Create(
		factory.WithNoDatabase(dbmode),
		factory.WithReportLimit(reportLimit),
		factory.WithScriptPath(scriptPath),
		factory.WithBatchSize(batchSize),
	)
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.Use(middleware.CORS())
	server.Route(e)
	e.Start(":8888")

	// Stop the service gracefully.
	if err := e.Shutdown(context.TODO()); err != nil {
		panic(err)
	}
}
