package main

import (
  "net"
  "log"
  "fmt"
  "crypto/rand"
  "io"
  "time"
)

const (
  PORT = ":8090"
)

type FileServer struct{}

func (fs *FileServer) start() {
  ln, err := net.Listen("tcp", PORT)
  if err != nil {
    log.Fatal(err)
  }
  for {
    conn, err := ln.Accept()
    if err != nil {
      log.Fatal(err)
    }
    go fs.readConn(conn)
  }
}

func (fs FileServer) readConn(conn net.Conn) {
  buf := make([]byte, 2048)
  for {
    n, err := conn.Read(buf)
    if err != nil {
      log.Fatal(err)
    }
    fmt.Printf("recieved %d bytes over the network\n", n)
  }
}

func sendFile(size int) error { // imitate a file
  file := make([]byte, size)
  _, err := io.ReadFull(rand.Reader, file)
  if err != nil {
    return err
  }
  conn, err := net.Dial("tcp", PORT)
  if err != nil {
    return err
  }
  var n int
  n, err = conn.Write(file)
  fmt.Printf("written %d bytes over the network\n", n)
  return nil
}

func main() {
  go func() {
    time.Sleep(3 * time.Second)
    sendFile(1000)
  }()
  server := &FileServer{}
  server.start()
}
