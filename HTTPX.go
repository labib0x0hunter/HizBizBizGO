package main 

import(
	"fmt"
	"os"
	"bufio"
	"sync"
	"net/http"
	"time"
	"strings"
)

// counter for random userAgent
// help Display Help , how to use it
var (
	counter int = 0
	help string = fmt.Sprintf("Usages : \n\tgo run crawler.go urls.txt")
	userAgents = []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.3 Safari/605.1.15",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:89.0) Gecko/20100101 Firefox/89.0",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (Linux; Android 10; SM-A505FN) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.120 Mobile Safari/537.36",
	}
)

// Help
func PrintHelp(){
	fmt.Println(help)
}

// Takes filename as input
// Gives Content of file, which is urls
func GetUrlFromFile(fileName string) []string {

	// Open file
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error in opening Files")
		os.Exit(1)
	}

	defer file.Close()

	// Read from files and store url in urls
	input := bufio.NewScanner(file)
	urls := []string{}
	for input.Scan() {
		url := input.Text()
		if !strings.HasPrefix(url, "https://"){
			continue
		}
		urls = append(urls, url)
	}
	return urls
}

// Make Request to urls
func makeRequest(url string, wg *sync.WaitGroup, mu *sync.Mutex){

	defer wg.Done()      // second execute

	// HTTP Client and Request
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	// Select Random userAgent
	mu.Lock()
	userAgent := userAgents[counter % 5]
	mu.Unlock()

	// Add Headers
	req.Header.Add("User-Agent", userAgent)

	// Make call to req
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close() // first exeecute

	// Print STATUS CODE
	mu.Lock()
	counter++
	fmt.Println("Url : ", url, " Status Code : ", resp.StatusCode)
	mu.Unlock()
}

func main() {

	// Check Argument , if file is provided
	args := os.Args[1:]
	if len(args) == 0 {
		PrintHelp()
		os.Exit(1)
	}

	// Filename
	// Read from file and store in urls
	fileName := args[0]
	urls := GetUrlFromFile(fileName)
	
	// For Race Condition & Goroutine Finish
	var (
		wg sync.WaitGroup
		mu sync.Mutex
	)

	// Rate Limiting , 5 request per second
	rate := time.Second / 5
	limiter := time.NewTicker(rate)
	defer limiter.Stop()

	// Make request to each url
	// Using goroutine as thread
	for _, url := range urls {
		<- limiter.C
		wg.Add(1)
		go makeRequest(url, &wg, &mu)
	}
	wg.Wait()
}
