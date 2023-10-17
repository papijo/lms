package main

import (
	"os"
	"os/signal"

	"github.com/papijo/lms/application"
)

func main() {
	//Start the Application
	e, db := application.Start()

	//Stop the Application
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	application.Stop(e, db)

}
