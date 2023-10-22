package main

import (
  "net"
  "log"
  "fmt"
  "crypto/rand"
  "io"
  "time"
  "bytes"
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
  buf := new(bytes.Buffer)
  for {
    n, err := io.CopyN(buf, conn, int64(40000))
    if err != nil {
      log.Fatal(err)
    }
    fmt.Printf("recieved %d bytes over the network\n", n)
  }
}

func sendFile(size int) (err error) { // imitate a file
  file := make([]byte, size)
  _, err = io.ReadFull(rand.Reader, file)
  if err != nil {
    return err
  }
  conn, err := net.Dial("tcp", PORT)
  if err != nil {
    return
  }
  n, err := io.CopyN(conn, bytes.NewReader(file), int64(size))
  if err != nil {
    return
  }
  fmt.Printf("written %d bytes over the network\n", n)
  return nil
}

func main() {
  go func() {
    time.Sleep(3 * time.Second)
    sendFile(40000)
  }()
  server := &FileServer{}
  server.start()
}
