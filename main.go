package main

import (
	"flag"
	"github.com/sevlyar/go-daemon"
	"log"
	"os"
	"syscall"
	"vacproxy/service"
)

var (
	signal = flag.String("s", "", `Send signal to the daemon:
  stop — shutdown
  reload — reloading the configuration file`)
)

var (
	stop = make(chan struct{})
	done = make(chan struct{})
)

func main() {

	addr := flag.String("bind", "0.0.0.0:7777", "proxy bind address")
	auth := flag.String("auth", "", "basic credentials(username:password)")
	logf := flag.String("log", "./vacproxy.log", "the log file path")
	pidf := flag.String("pid", "./vacproxy.pid", "the pid file path")
	quit := flag.Bool("q", false, "quit proxy")
	nd := flag.Bool("nd", false, "not run as daemon")
	flag.Parse()

	if *quit || *signal == "stop" {
		*signal = "stop"
		*nd = false
		*quit = true
	}

	if *nd {
		startWorker(*addr, *auth)
		return
	}

	daemon.AddCommand(daemon.StringFlag(signal, "stop"), syscall.SIGTERM, termHandler)
	//daemon.AddCommand(daemon.StringFlag(signal, "reload"), syscall.SIGHUP, reloadHandler)

	cntxt := &daemon.Context{
		PidFileName: *pidf,
		PidFilePerm: 0644,
		LogFileName: *logf,
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
	}

	if len(daemon.ActiveFlags()) > 0 {
		d, err := cntxt.Search()
		if err != nil {
			log.Fatalf("Unable send signal to the daemon: %s", err.Error())
		}
		daemon.SendCommands(d)
		return
	}

	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatalln(err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()

	log.Println("- - - - - - - - - - - - - - -")
	log.Println("daemon started")

	log.Println(*addr)
	startWorker(*addr, *auth)

	err = daemon.ServeSignals()
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}

	log.Println("daemon terminated")

}

func startWorker(addr string, auth string) {
	// main worker
	server := service.NewServer(addr, auth)
	server.Start()
	select {
	case <-stop:
		server.Stop()
	default:
	}
}

func termHandler(sig os.Signal) error {
	log.Println("terminating...")
	stop <- struct{}{}
	if sig == syscall.SIGQUIT {
		<-done
	}
	return daemon.ErrStop
}
