package main 

import(
    "fmt"
    "net"
    "strings"
)

// Server Configuration
type Config struct {
    address string
    port    string
}

// Server
type Server struct {
    listener net.Listener
    cfg     Config
}

// Starts the server
func (s *Server) Start() {
    addr := net.JoinHostPort(s.cfg.address, s.cfg.port)
    fmt.Println(addr)
    var err error
    s.listener, err = net.Listen("tcp", addr)
    if err != nil {
        fmt.Println("Error in Listening")
        return
    }

    defer s.listener.Close()
    fmt.Println("Server is listening at ", addr)

    // Handle connection
    for{
        conn, err := s.listener.Accept()
        if err != nil {
            fmt.Println("Error in Incoming request")
            continue
        }

        // Start Goroutine
        go s.HandleConnection(conn)
    }
}

// Handle Connection for each connection
func (s *Server) HandleConnection(conn net.Conn) {
    defer conn.Close()
    conn.Write([]byte("Welcome to Server\n"))

    buffer := make([]byte, 1024)
    for{
        conn.Write([]byte(" -> "))
        n, err := conn.Read(buffer)
        if err != nil {
            fmt.Println("Client Disconnected\n")
            return
        }
        input := buffer[:n]
        response := strings.TrimSpace(string(input))
        conn.Write([]byte(response + "\n"))
    }
}

func main() {
    cfg := Config{
        address : "127.0.0.1",
        port :  "8080",
    }

    server := Server{
        cfg : cfg,
    }

    server.Start()
}
