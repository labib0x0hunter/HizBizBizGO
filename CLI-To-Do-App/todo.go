package main 

import(
    "fmt"
    "net"
    "strings"
    "strconv"
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

// To-Do
// type Task struct {
//     id      int
//     description string
//     status  string
//     createdAt   time.Time
//     updatedAt   time.Time
// }

// # Adding a new task
// task-cli add "Buy groceries"
// # Output: Task added successfully (ID: 1)

// # Updating and deleting tasks
// task-cli update 1 "Buy groceries and cook dinner"
// task-cli delete 1

// # Marking a task as in progress or done
// task-cli mark-in-progress 1
// task-cli mark-done 1

// # Listing all tasks
// task-cli list

// # Listing tasks by status
// task-cli list done
// task-cli list todo
// task-cli list in-progress

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
    conn.Write([]byte("Welcome to To-Do CLI\n"))

    var tasks []string
    var completed []bool

    buffer := make([]byte, 1024)
    for{
        conn.Write([]byte(" -> "))
        n, err := conn.Read(buffer)
        if err != nil {
            fmt.Println("Client Disconnected\n")
            return
        }
        input := buffer[:n]
        parts := strings.Fields(strings.TrimSpace(string(input)))

        // No Input
        if len(parts) == 0 {
            continue
        }

        cmd := parts[0]      // Cmd type
        switch cmd {
        case "add" :
            parts = parts[1: ]
            task := strings.Join(parts, " ")
            if len(task) == 0 {
                conn.Write([]byte("No task is given\n"))
                continue
            }
            exist := false
            for _, t := range tasks {
                if t == task {
                    conn.Write([]byte("Task already exits\n"))
                    exist = true
                    break
                }
            }
            if exist == false {
                completed = append(completed, false)
                tasks = append(tasks, task)
                conn.Write([]byte("Task Added\n"))
            }
        case "len" :
            conn.Write([]byte(strconv.Itoa(len(tasks)) + "\n"))
        case "list" :
            if len(tasks) == 0 {
                conn.Write([]byte("Task list is empty\n"))
                continue
            }
            taskCount := len(tasks)
            if len(parts) == 1 {       // All Tasks
                for i := 0; i < taskCount; i++ {
                    var task string
                    if completed[i] {
                        task = fmt.Sprintf("[+] %d : %s", i + 1, tasks[i])
                    } else {
                        task = fmt.Sprintf("[-] %d : %s", i + 1, tasks[i])
                    }
                    conn.Write([]byte(task + "\n"))
                }
            } else {
                if parts[1] == "done" {
                    for i := 0; i < taskCount; i++ {
                        if completed[i] {
                            conn.Write([]byte("[+] " + tasks[i] + "\n"))
                        }
                    }
                } else if parts[1] == "undone" {
                    for i := 0; i < taskCount; i++ {
                        if !completed[i] {
                            conn.Write([]byte("[-] " + tasks[i] + "\n"))
                        }
                    }
                }
            }
            
        case "delete":
            taskId, err := strconv.Atoi(parts[1])
            if err != nil {
                conn.Write([]byte("Error in getting task id\n"))
                continue
            }
            taskCount := len(tasks)
            if taskId > taskCount || taskId < 1{
                conn.Write([]byte("Task is not in list\n"))
                continue
            }
            taskId--
            tasks = append(tasks[ : taskId], tasks[taskId + 1 : ]...)
            completed = append(completed[ : taskId], completed[taskId + 1 : ]...)
            conn.Write([]byte("Task is Removed\n"))
        case "done":
            if len(parts) == 1 {
                conn.Write([]byte("Expected an task id\n"))
                continue
            }
            taskId, err := strconv.Atoi(parts[1])
            if err != nil {
                conn.Write([]byte("Error in getting task id\n"))
                continue
            }
            taskCount := len(tasks)
            if taskId > taskCount || taskId < 1 {
                conn.Write([]byte("Task is not in list\n"))
                continue
            }
            taskId--
            completed[taskId] = true
            conn.Write([]byte("Marked as completed\n"))
        case "exit":
            conn.Write([]byte("Goodbye !\n"))
            return
        default:
            conn.Write([]byte("Unknown Command\n"))
        }
    }
}

// To-Do
func validateAddTask(){}
func validateTaskId(){}

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
