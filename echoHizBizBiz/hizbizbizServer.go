package main 

import(
    "fmt"
    "net"
    "strings"
    "time"
)

var users = map[string]string{
    "user1": "pass1",
    "user2": "pass2",
}

func currentTime() string {
    return time.Now().Format("2006-01-02 15:04:05")
}

func getUser(conn net.Conn) (string, string, error) {
    buffer := make([]byte, 1024)
    conn.Write([]byte("username: "))
    n, err := conn.Read(buffer)
    if err != nil {
        return "", "", fmt.Errorf("Error reading username: %v", err)
    }
    username := strings.TrimSpace(string(buffer[:n]))

    conn.Write([]byte("password: "))
    n, err = conn.Read(buffer)
    if err != nil {
        return "", "", fmt.Errorf("Error reading password: %v", err)
    }
    password := strings.TrimSpace(string(buffer[:n]))
    return username, password, nil
}

func handleConn(conn net.Conn) {
    defer conn.Close()
    conn.Write([]byte("Welcome to HizBizBiz\n"))
    buffer := make([]byte, 1024)

    var username string
    for {
        user, pass, err := getUser(conn)
        if err != nil {
            fmt.Println("Auth failed")
            return
        }
        checkPass, ok := users[user]
        if ok && checkPass == pass {
            conn.Write([]byte("Login successful\n"))
            username = user
            fmt.Println(username + " logged in " + currentTime())
            break
        }
        conn.Write([]byte("Wrong Credentials. Try Again\n"))
    }

    for {
        conn.Write([]byte(" >>> "))
        n, err := conn.Read(buffer)
        if err != nil {
            fmt.Println("Client disconnected")
            return
        }
        input := strings.TrimSpace(string(buffer[:n]))
       
        cmd := input
        switch cmd {
        case "whoami":
            conn.Write([]byte(username + "\n"))
        case "get":
            conn.Write([]byte(currentTime() + "\n"))
        case "exit":
            conn.Write([]byte("Goodbye " + username + "!\n"))
            fmt.Println(username + " logged out " + currentTime())
            return
        default:
            conn.Write([]byte("Unknown Command\n"))
        }
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
