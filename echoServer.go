package main 

import(
    "fmt"
    "net"
    "strings"
)

func handleConn(conn net.Conn) {
    defer conn.Close()
    conn.Write([]byte("Welcome to HizBizBiz\n"))
    buffer := make([]byte, 1024)
    for{
        conn.Write([]byte(" >>> "))
        n, err := conn.Read(buffer)
        if err != nil {
            fmt.Println("Client disconnected")
            return
        }
        input := strings.TrimSpace(string(buffer[:n]))
        response := fmt.Sprintf("%s\n", input)
        conn.Write([]byte(response))
    }
}

func main() {
        
    listener, err := net.Listen("tcp", ":8080")
    if err != nil{
        panic("Error - 01")
    }
    defer listener.Close()

    // Handle multiple client
    for {
        conn, err := listener.Accept()
        if err != nil{
            fmt.Println("Error accepting client : ", err)
            continue
        }
        go handleConn(conn)
    }
}
