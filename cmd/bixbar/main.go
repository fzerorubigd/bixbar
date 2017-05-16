package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fzerorubigd/bixbar"
)

func main() {
	bar := bixbar.NewBar(time.Second, os.Stdout, os.Stdin)
	bar.AddBlock(bixbar.NewTextBlock("Example", "text", "ins"))
	bar.Start()
	quit := make(chan os.Signal, 6)
	signal.Notify(quit, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT)
	<-quit

	bar.Stop()
}
