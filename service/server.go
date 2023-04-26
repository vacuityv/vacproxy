package service

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"net"
	"os"
	"strings"
)

type Server struct {
	listener   net.Listener
	addr       string
	credential string
	//configFile string
	config   VacConfig
	inMatch  *Node
	outMatch *Node
}

/*
*
create proxy server
*/
func NewServer(config VacConfig, inMatch *Node, outMatch *Node) *Server {

	//var config, inMatch, outMatch = initConfig(configFile)

	credential := ""
	if config.Auth.Enabled && len(config.Auth.User) > 0 && len(config.Auth.Password) > 0 {
		credential = base64.StdEncoding.EncodeToString([]byte(config.Auth.User + ":" + config.Auth.Password))
		log.Printf("init server with auth")
	}

	return &Server{
		addr:       config.Bind,
		credential: credential,
		//configFile: configFile,
		config:   config,
		inMatch:  inMatch,
		outMatch: outMatch,
	}
}

/*
*
start proxy server
*/
func (s *Server) Start() {
	var err error
	s.listener, err = net.Listen("tcp", s.addr)
	if err != nil {
		log.Fatal(err)
	}

	if s.credential != "" {
		log.Printf("use %s for auth\n", s.credential)
	}
	log.Println("vacproxy start successfully\n")
	var showAddr = s.addr
	if strings.HasPrefix(showAddr, ":") {
		showAddr = "http://0.0.0.0" + showAddr
	}
	log.Printf("listen address: %s\n", showAddr)
	// config
	log.Printf("vacproxy config:%v", s.config)

	log.Printf("waiting for connection...\n")

	sf, err := NewSnowflake(1)
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Fatalln(err)
			continue
		}
		logId, err := sf.NextId()
		go s.createConnection(conn, logId).serve()
	}
}

/*
*
create connection
*/
func (s *Server) createConnection(rwc net.Conn, logId int64) *serverConnnection {
	return &serverConnnection{
		server: s,
		rwc:    rwc,
		brc:    bufio.NewReader(rwc),
		logId:  logId,
	}
}

/*
*
check wheather need auth
*/
func (s *Server) isNeedAuth() bool {
	return s.credential != ""
}

/*
*
check proxy auth: Basic credential
*/
func (s *Server) validateCredential(logId int64, basicCredential string) bool {

	c := strings.Split(basicCredential, " ")
	if len(c) == 2 && strings.EqualFold(c[0], "Basic") && c[1] == s.credential {
		return true
	}
	log.Printf("%d-auth failed: %s", logId, basicCredential)
	return false
}

/*
*
stop proxy server
*/
func (s *Server) Stop() {
	log.Printf("server ready to terminating")
	s.listener.Close()
	log.Printf("server terminated")
}

/*
*
reload config file
*/
//func (s *Server) Reload() {
//	log.Println("server ready to reload")
//	var config, inMatch, outMatch = initConfig(s.configFile)
//	s.config = config
//	s.inMatch = inMatch
//	s.outMatch = outMatch
//
//	credential := ""
//	if config.Auth.Enabled && len(config.Auth.User) > 0 && len(config.Auth.Password) > 0 {
//		credential = base64.StdEncoding.EncodeToString([]byte(config.Auth.User + ":" + config.Auth.Password))
//		log.Printf("reload server with auth: %s", credential)
//	}
//	s.credential = credential
//
//	log.Printf("server reload success:%v", config)
//}

func InitConfig(configFile string) (VacConfig, *Node, *Node) {

	var config VacConfig
	file, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Fatalf("Failed to unmarshal file: %v", err)
	}
	var inMatch, outMatch = initIpTrieTree(config)

	return config, inMatch, outMatch
}

func initIpTrieTree(config VacConfig) (*Node, *Node) {

	inMatch := NewNode()

	if len(config.InAllowList) > 0 {
		for _, ip := range config.InAllowList {
			inMatch.Insert(ip)
		}
	}

	outMatch := NewNode()
	if len(config.OutAllowList) > 0 {
		for _, domain := range config.OutAllowList {
			outMatch.Insert(domain)
		}
	}
	return inMatch, outMatch
}
