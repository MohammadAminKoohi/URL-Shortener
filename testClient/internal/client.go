package internal

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func RequestSender(stopCh <-chan bool) {
	for {
		select {
		case <-stopCh:
			fmt.Println("Stopping request sender")
			return
		default:
			randomUrl := GenerateRandomUrl()
			fmt.Printf("Sending URL: %s\n", randomUrl)
			SendPostRequest(randomUrl)
			time.Sleep(250 * time.Millisecond)
		}
	}
}

func RunClient() {
	stopCh := make(chan bool)
	for i := 0; i < 10; i++ {
		time.Sleep(200 * time.Millisecond)
		go RequestSender(stopCh)
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter 0 to stop: ")
		input, _ := reader.ReadString('\n')
		if input == "0\n" {
			close(stopCh)
			break
		}
	}
}
