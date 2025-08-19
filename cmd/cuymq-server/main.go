package main

import (
  "fmt"
  "io"
  "log"
  "net"
)

const (
  DefaultListenAddr = ":5678"
)

// Message Broker server
type Server struct {
  listenAddr  string
  listener    net.Listener
  quitCh      chan struct{}
}

func NewServer(listenAddr string) *Server {
  return &Server{
    listenAddr: listenAddr,
    quitCh: make(chan struct{}),
  }
}

func (s *Server) Start() error {
  ln, err := net.Listen("tcp", s.listenAddr)
  if err != nil {
    return fmt.Errorf("Failed to listen on %s: %w", s.listenAddr, err)
  }

  s.listener = ln

  log.Printf("CuyMQ server started on %s", s.listenAddr)

  go s.acceptLoop()

  <-s.quitCh
  return s.listener.Close()
}

func (s *Server) acceptLoop()  {
  for {
    conn, err := s.listener.Accept()
    if err != nil {
      log.Println("accept error:", err)
      break
    }

    log.Printf("new connection from: %s", conn.RemoteAddr())

    go s.hanndleConnection(conn)
  }
}


func (s *Server) hanndleConnection(conn net.Conn)  {
  defer conn.Close()
  buf := make([]byte, 8)
  for {
    n, err := conn.Read(buf)
    if err != nil {
      if err != io.EOF {
        log.Printf("read error: %v", err)
      }
      break
    }
    msg := string(buf[:n])
    log.Println(n, "this is n")
    log.Printf("Received from %s: %s", conn.RemoteAddr(), msg)

    _, err = conn.Write([]byte("server received: " + msg))
    if err != nil {
      log.Printf("write error: %v", err)
      break
    }
  }
  log.Printf("connection closed: %s", conn.RemoteAddr())
}

func main()  {
  server := NewServer(DefaultListenAddr)

  log.Fatal(server.Start())
}
