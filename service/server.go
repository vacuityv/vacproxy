package service

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"strings"
)

type Server struct {
	listener   net.Listener
	addr       string
	credential string
}

/*
*
create proxy server
*/
func NewServer(addr, credential string) *Server {
	return &Server{addr: addr, credential: base64.StdEncoding.EncodeToString([]byte(credential))}
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
func (s *Server) validateCredential(basicCredential string) bool {
	c := strings.Split(basicCredential, " ")
	if len(c) == 2 && strings.EqualFold(c[0], "Basic") && c[1] == s.credential {
		return true
	}
	return false
}

/*
*
stop proxy server
*/
func (s *Server) Stop() {
	s.listener.Close()
}
