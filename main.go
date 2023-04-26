package main

import (
	"flag"
	"fmt"
	"github.com/sevlyar/go-daemon"
	"io"
	"log"
	"os"
	"runtime"
	"syscall"
	"vacproxy/service"
)

var (
	signal = flag.String("s", "", `Send signal to the daemon[Windows platform not support]:
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
	pidf := flag.String("pid", "./vacproxy.pid", "the pid file path[Windows platform not support]")
	quit := flag.Bool("q", false, "quit proxy[Windows platform not support]")
	configFile := flag.String("config", "./config.yml", "config file")
	flag.Parse()

	osName := runtime.GOOS
	if osName == "windows" {
		startWindowsWorker(*addr, *configFile, *logf)
	} else {
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
		go startUnixWorker(*addr, *configFile)

		err = daemon.ServeSignals()
		if err != nil {
			log.Fatalf("Error: %s", err.Error())
		}
		log.Println("daemon terminated")
	}

}

func startUnixWorker(addr string, configFile string) {
	// main worker
	server := service.NewServer(addr, configFile)

	// watch the signal
	go watchSig(server)

	// start server
	server.Start()

}

func startWindowsWorker(addr string, configFile string, logFile string) {

	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %s", err.Error())
		fmt.Printf("Failed to open log file: %s", err.Error())
	}
	defer file.Close()

	// 组合一下即可，os.Stdout代表标准输出流
	multiWriter := io.MultiWriter(os.Stdout, file)
	log.SetOutput(multiWriter)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Printf("windows start server at %s", addr)
	// main worker
	server := service.NewServer(addr, configFile)
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
