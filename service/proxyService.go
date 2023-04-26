package service

import (
	"fmt"
	"github.com/kardianos/service"
	"io"
	"log"
	"os"
	"path/filepath"
)

type program struct {
	//configFile  string
	proxyServer *Server
	config      VacConfig
	inMatch     *Node
	outMatch    *Node
}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) run() {

	server := NewServer(p.config, p.inMatch, p.outMatch)
	p.proxyServer = server
	StartProxy(server, false)
}

func (p *program) Stop(s service.Service) error {
	p.proxyServer.Stop()
	return nil
}

func ProxyService(config VacConfig, inMatch *Node, outMatch *Node) {
	svcConfig := &service.Config{
		Name:        "vacproxy",
		DisplayName: "vacproxy",
		Description: "http proxy",
	}

	prg := &program{
		config:   config,
		inMatch:  inMatch,
		outMatch: outMatch,
	}

	svc, err := service.New(prg, svcConfig)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	if len(os.Args) > 1 {
		if os.Args[1] == "install" {
			err := svc.Install()
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
			fmt.Println("Service installed")
			return
		}
		if os.Args[1] == "uninstall" {
			err := svc.Uninstall()
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
			fmt.Println("Service uninstalled")
			return
		}
		if os.Args[1] == "stop" {
			err := svc.Stop()
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
			fmt.Println("Service stop")
			return
		}
	}

	err = svc.Run()
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
}

func GetExecPath() (string, error) {
	fullExecPath, err := os.Executable()
	if err != nil {
		return "", err
	}
	dir, _ := filepath.Split(fullExecPath)
	return dir, nil
}

func GetConfigPath() (string, error) {
	dir, err := GetExecPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.yaml"), nil
}

func StartProxy(server *Server, consoleFlag bool) {

	logPath := server.config.Log
	if len(logPath) == 0 {
		dir, _ := GetExecPath()
		logPath = filepath.Join(dir, "vacproxy.log")
	}

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %s", err.Error())
		fmt.Printf("Failed to open log file: %s", err.Error())
	}
	defer file.Close()

	// 组合一下即可，os.Stdout代表标准输出流
	if consoleFlag {
		multiWriter := io.MultiWriter(os.Stdout, file)
		log.SetOutput(multiWriter)
	} else {
		log.SetOutput(file)
	}
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	server.Start()
}
