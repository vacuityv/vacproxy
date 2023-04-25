package service

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/textproto"
	"net/url"
	"strings"
)

type serverConnnection struct {
	rwc    net.Conn
	brc    *bufio.Reader
	server *Server
	logId  int64
}

/*
*
process server connection
*/
func (c *serverConnnection) serve() {
	defer c.rwc.Close()
	rawHttpRequestHeader, remote, credential, isHttps, err := c.preProcessRequest()
	if err != nil {
		log.Println(err)
		return
	}

	if c.checkAuth(credential) == false {
		log.Printf("%d-proxy auth fail", c.logId)
		c.rwc.Write([]byte("Error: proxy auth fail"))
		return
	}

	log.Printf("%d-connecting to %s", c.logId, remote)
	remoteConn, err := net.Dial("tcp", remote)
	if err != nil {
		log.Printf("%d-%s", c.logId, err)
		c.rwc.Write([]byte("Error: " + err.Error()))
		return
	}

	localForClientIp, localForClientPort, clientIp, clientPort := getIpAddr(c.rwc)
	localForServerIp, localForServerPort, serverIp, serverPort := getIpAddr(remoteConn)

	// check client ip
	if len(c.server.config.InAllowList) > 0 && !c.server.inMatch.Match(clientIp) {
		log.Printf("%d-clientIp not in allow list: %s", c.logId, clientIp)
		_, err = c.rwc.Write([]byte("HTTP/1.1 401 clientIp not in allow list " + clientIp + "\r\n\r\n"))
		if err != nil {
			log.Printf("%s", err.Error())
			return
		}
	}
	// check server ip
	domain := strings.Split(remote, ":")[0]
	if len(c.server.config.OutAllowList) > 0 && !c.server.outMatch.Match(domain) {
		log.Printf("%d-server host not in allow list: %s", c.logId, domain)
		_, err = c.rwc.Write([]byte("HTTP/1.1 401 server host not in allow list " + domain + "\r\n\r\n"))
		if err != nil {
			log.Printf("%s", err.Error())
			return
		}
	}

	log.Printf("%d-client connect %s:%d to %s:%d", c.logId, clientIp, clientPort, localForClientIp, localForClientPort)
	log.Printf("%d-server connect %s:%d to %s:%d", c.logId, localForServerIp, localForServerPort, serverIp, serverPort)

	if isHttps {
		// if https, should sent 200 to client
		_, err = c.rwc.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
		if err != nil {
			log.Printf("%s", err.Error())
			return
		}
	} else {
		// if not https, should sent the request header to remote
		_, err = rawHttpRequestHeader.WriteTo(remoteConn)
		if err != nil {
			log.Printf("%d-%s", c.logId, err)
			return
		}
	}

	// build bidirectional-streams
	log.Printf("%d-begin tunnel", c.logId)
	c.tunnel(remoteConn)
}

/*
* propess the request
 */
func (c *serverConnnection) preProcessRequest() (rawReqHeader bytes.Buffer, host, credential string, isHttps bool, err error) {
	tp := textproto.NewReader(c.brc)

	// http request first line: GET /index.html HTTP/1.0
	var requestLine string
	if requestLine, err = tp.ReadLine(); err != nil {
		return
	}

	method, requestURI, _, ok := parseRequestLine(requestLine)
	if !ok {
		err = &BadRequestError{"malformed http request"}
		return
	}

	// https request
	if method == "CONNECT" {
		isHttps = true
		requestURI = "http://" + requestURI
	}

	// get remote host
	uriInfo, err := url.ParseRequestURI(requestURI)
	if err != nil {
		return
	}

	// Subsequent lines: Key: value.
	mimeHeader, err := tp.ReadMIMEHeader()
	if err != nil {
		return
	}

	credential = mimeHeader.Get("Proxy-Authorization")

	if uriInfo.Host == "" {
		host = mimeHeader.Get("Host")
	} else {
		if strings.Index(uriInfo.Host, ":") == -1 {
			host = uriInfo.Host + ":80"
		} else {
			host = uriInfo.Host
		}
	}

	// rebuild http request header
	rawReqHeader.WriteString(requestLine + "\r\n")
	for k, vs := range mimeHeader {
		for _, v := range vs {
			rawReqHeader.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
		}
	}
	rawReqHeader.WriteString("\r\n")
	return
}

// auth provide basic authentication
func (c *serverConnnection) checkAuth(credential string) bool {
	if c.server.isNeedAuth() == false || c.server.validateCredential(c.logId, credential) {
		return true
	}
	// first send 407 to client
	_, err := c.rwc.Write(
		[]byte("HTTP/1.1 407 Proxy Authentication Required\r\nProxy-Authenticate: Basic realm=\"*\"\r\n\r\n"))
	if err != nil {
		log.Println(err)
	}
	return false
}

// tunnel http message between client and server
func (c *serverConnnection) tunnel(remoteConn net.Conn) {
	if remoteConn == nil {
		return
	}
	defer remoteConn.Close()
	defer c.rwc.Close()
	defer log.Printf("%d-end tunnel", c.logId)

	clientDoneCh, serverDoneCh := make(chan struct{}), make(chan struct{})
	go dataCopy(remoteConn, c.rwc, serverDoneCh)
	go dataCopy(c.rwc, remoteConn, clientDoneCh)
	wait(clientDoneCh, serverDoneCh)
}

func dataCopy(det net.Conn, src net.Conn, done chan struct{}) {
	buf := make([]byte, 1024*32)
	_, _ = io.CopyBuffer(det, src, buf)
	close(done)
}

func wait(clientDoneCh chan struct{}, serverDoneCh chan struct{}) {
	select {
	case <-clientDoneCh:
	case <-serverDoneCh:
	}
}

func parseRequestLine(line string) (method, requestURI, proto string, ok bool) {
	s1 := strings.Index(line, " ")
	s2 := strings.Index(line[s1+1:], " ")
	if s1 < 0 || s2 < 0 {
		return
	}
	s2 += s1 + 1
	return line[:s1], line[s1+1 : s2], line[s2+1:], true
}

type BadRequestError struct {
	what string
}

func (b *BadRequestError) Error() string {
	return b.what
}

func getIpAddr(conn net.Conn) (remoteIp string, remotePort int, localIp string, localPort int) {
	localAddr, ok1 := conn.LocalAddr().(*net.TCPAddr)
	tcpAddr, ok2 := conn.RemoteAddr().(*net.TCPAddr)
	if ok1 && ok2 {
		return localAddr.IP.String(), localAddr.Port, tcpAddr.IP.String(), tcpAddr.Port
	}
	return "-1", -1, "-1", -1
}
