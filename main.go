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
  stop — shutdown, same as -q
  reload — reloading the configuration file`)
)

var (
	stop   = make(chan struct{})
	done   = make(chan struct{})
	reload = make(chan struct{})
)

func main() {

	addr := flag.String("bind", "0.0.0.0:7777", "proxy bind address")
	logf := flag.String("log", "./vacproxy.log", "the log file path")
	pidf := flag.String("pid", "./vacproxy.pid", "the pid file path")
	quit := flag.Bool("q", false, "quit proxy")
	configFile := flag.String("config", "./config.yml", "config file")
	flag.Parse()

	if *quit || *signal == "stop" {
		*signal = "stop"
		*quit = true
	}

	daemon.AddCommand(daemon.StringFlag(signal, "stop"), syscall.SIGTERM, termHandler)
	daemon.AddCommand(daemon.StringFlag(signal, "reload"), syscall.SIGHUP, reloadHandler)

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
	go startWorker(*addr, *configFile)

	err = daemon.ServeSignals()
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}
	log.Println("daemon terminated")

}

func startWorker(addr string, configFile string) {
	// main worker
	server := service.NewServer(addr, configFile)
	// watch the signal
	go watchSig(server)
	// start server
	server.Start()

}

func watchSig(s *service.Server) {
	for {
		select {
		case <-stop:
			s.Stop()
		case <-reload:
			s.Reload()
		default:
		}
	}
}

func reloadHandler(sig os.Signal) error {
	reload <- struct{}{}
	return nil
}

func termHandler(sig os.Signal) error {
	stop <- struct{}{}
	if sig == syscall.SIGQUIT {
		<-done
	}
	return daemon.ErrStop
}
