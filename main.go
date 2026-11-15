package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const API_URL = "http://localhost:11434/api/chat"

// define the user message , using struct tag so we can produce json structure .
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"` // slice that contain all messages from user and the model .
	Stream   bool      `json:"stream"`
}

type ChatResponse struct {
	Message Message `json:"message"`
}

func printBanner() {
	fmt.Println(`
   ██████╗██╗     ██╗██╗     ██╗     ███╗   ███╗
  ██╔════╝██║     ██║██║     ██║     ████╗ ████║
  ██║     ██║     ██║██║     ██║     ██╔████╔██║
  ██║     ██║     ██║██║     ██║     ██║╚██╔╝██║
  ╚██████╗███████╗██║███████╗███████╗██║ ╚═╝ ██║
   ╚═════╝╚══════╝╚═╝╚══════╝╚══════╝╚═╝     ╚═╝

             LOCAL LLM TERMINAL
`)
	fmt.Println("  ----------------------------------------")
}

// in the case the user have more then llm (or none)
func getModels() ([]string, error) {

	out, err := exec.Command("ollama", "list").Output()
	if err != nil {
		return nil, err
	}

	var models []string

	lines := strings.Split(string(out), "\n")

	for i, line := range lines {

		// Skip the header
		if i == 0 {
			continue
		}

		fields := strings.Fields(line)

		// First column contains model name
		if len(fields) > 0 {
			models = append(models, fields[0])
		}
	}

	return models, nil
}

func selectModel(sc *bufio.Scanner) (string, error) {

	models, err := getModels()
	if err != nil {
		return "", err
	}

	if len(models) == 0 {
		return "", fmt.Errorf("no installed models found")
	}

	fmt.Println("\n  Available models:\n")

	for i, model := range models {
		fmt.Printf("  [%d] %s\n", i+1, model)
	}

	fmt.Printf("\n  Select model [1-%d] > ", len(models))

	if !sc.Scan() {
		return "", fmt.Errorf("failed to read input: %v", sc.Err())
	}

	choice, err := strconv.Atoi(strings.TrimSpace(sc.Text()))
	if err != nil {
		return "", fmt.Errorf("please enter a number")
	}

	if choice < 1 || choice > len(models) {
		return "", fmt.Errorf("invalid model selection")
	}

	return models[choice-1], nil
}

func CatchSignal() {

	// create channel to catch the signals :
	sigchan := make(chan os.Signal, 1)

	//listen for SIGINIT:
	signal.Notify(sigchan, syscall.SIGINT)

	// now waiting for the signal :
	sig := <-sigchan

	fmt.Println(" bye bye .....", sig)

	os.Exit(0)

}

// loadDocu checks if document exists and passes it to python script
func loadDocu(path string) (string, error) {
	// check if the file exist :
	info, err := os.Stat(path)
	if err != nil {
		fmt.Println("the file doesnt exist ...")
		return "", err
	}
	// check if its not directory :
	if info.IsDir() {
		return "", fmt.Errorf("%s it directory not a file ...", info.Name())
	}

	cmd := exec.Command("python3", "readoc.py", path)
	content, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func main() {

	sc := bufio.NewScanner(os.Stdin)

	printBanner()
	model, err := selectModel(sc)
	if err != nil {
		fmt.Println("error: ", err)
		return
	}

	// lunch go routine to catch signal , if user press ctrl + c
	go CatchSignal()

	for {
		//get the input from user :
		fmt.Print("\nYou > ")
		if !sc.Scan() {
			fmt.Println("Input error:", sc.Err())
			return
		}

		userInput := strings.TrimSpace(sc.Text())

		// handle upload check
		if userInput == "/upload" || strings.HasPrefix(userInput, "/upload ") {
			rest := strings.TrimSpace(strings.TrimPrefix(userInput, "/upload"))
			if rest == "" {
				fmt.Println("error: please provide a file path after /upload")
				continue
			}

			parts := strings.SplitN(rest, " ", 2)
			path := parts[0]

			fileContent, err := loadDocu(path)
			if err != nil {
				fmt.Println("error:", err)
				continue
			}

			fmt.Println("file uploaded!")

			if len(parts) == 2 && strings.TrimSpace(parts[1]) != "" {
				userInput = fmt.Sprintf("Here is the text content extracted from the document:\n\"\"\"\n%s\n\"\"\"\n\nQuestion: %s", strings.TrimSpace(fileContent), strings.TrimSpace(parts[1]))
			} else {
				userInput = fmt.Sprintf("Summarize the following document text:\n\"\"\"\n%s\n\"\"\"", strings.TrimSpace(fileContent))
			}
		}
		// 2. Construct the HTTP request DATA :
		requestData := ChatRequest{
			Model: model,
			Messages: []Message{
				{
					Role:    "user",
					Content: userInput,
				},
			},
			Stream: false,
		}

		// 3. Convert Go structure into JSON
		jsonData, err := json.Marshal(requestData)
		if err != nil {
			fmt.Println("JSON error:", err)
			return
		}

		// 4. Create HTTP request
		req, err := http.NewRequest(
			"POST",
			API_URL,
			bytes.NewBuffer(jsonData),
		)
		if err != nil {
			fmt.Println("Request error:", err)
			return
		}

		// set the HTTP header (for json structure)
		req.Header.Set("Content-Type", "application/json")

		// 5. Send request
		client := &http.Client{
			Timeout: 5 * time.Minute,
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Connection error:", err)
			return
		}
		defer resp.Body.Close()

		// 6. Check HTTP status (if its ok or not ...)
		if resp.StatusCode != http.StatusOK {
			fmt.Println("Server returned:", resp.Status)
			return
		}

		// 7. Decode JSON response
		var result ChatResponse

		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			fmt.Println("JSON decode error:", err)
			return
		}

		// 8. Display response
		fmt.Println("\nAi>", result.Message.Content)
	}
}
