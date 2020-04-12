package client

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
)

var client *Client

func RunPrompter() {
	fmt.Println("\tWelcome to use client for MiniServer!")
	for {
		exit := BuildConnection()

		if exit {
			break
		}

		for {
			leave := SendRequest()
			if leave {
				break
			}
		}
	}

	fmt.Println("Goodbye.")
}

func BuildConnection() bool {

	promptSelect := promptui.Select{
		Label: "Do you want to Build a Connection?",
		Items: []string{"Yes", "No"},
	}

	_, conn, err := promptSelect.Run()

	if conn == "No" {
		return true
	}

	fmt.Println("Please enter the address which you want to connect to.")
	fmt.Println("ex: localhost:3000")
	prompt := promptui.Prompt{
		Label: "Server Address",
	}
	addr, err := prompt.Run()

	checkErr(err)

	go buildConn(addr)
	return false
}

func SendRequest() bool {
	fmt.Println("Please Select a Request Method.")
	prompt1 := promptui.Select{
		Label: "Request Method",
		Items: []string{"Get", "Quit"},
	}
	_, method, err := prompt1.Run()
	checkErr(err)

	if method == "Quit" {
		client.Close()
		return true
	}

	fmt.Println("Please enter a URI.")
	fmt.Println("ex: /")
	prompt2 := promptui.Prompt{
		Label: "URI",
	}
	uri, err := prompt2.Run()

	if client.Closed {
		fmt.Println("\t[Warning] : Connection was already closed!\n")
		return true
	}

	clientHandle(method, uri)

	return false
}

func clientHandle(method, uri string) {
	switch method {
	case "Get":
		r := NewRequest("Get", uri, client.Addr)
		sendRequest(r)
	}
}

func sendRequest(r *Request) {
	fmt.Printf("------Request %s------\n\n", r.Id)
	fmt.Printf("\t\t%s %s at %s\n", r.Method, r.URI, r.Timestamp)
	fmt.Printf("\n------Request %s------\n\n", r.Id)
	res := client.Send(r)
	fmt.Printf("------Response from %s------\n\n", res.Request.Id)
	fmt.Printf("\t\t%s\n", res.Body)
	fmt.Printf("\n------Response from %s------\n\n", res.Request.Id)
}

func buildConn(addr string) {
	ip, port := strings.Split(addr, ":")[0], strings.Split(addr, ":")[1]
	client = NewConn(ip, port)
	fmt.Println("Connection was built on", addr)
}

func checkErr(err error) {
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
}
