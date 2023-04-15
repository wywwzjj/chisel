// chisel end-to-end testSocks
// ======================
//
//                    (direct)
//         .--------------->----------------.
//        /    chisel         chisel         \
// request--->client:2001--->server:2002---->fileserver:3000
//        \                                  /
//         '--> crowbar:4001--->crowbar:4002'
//              client           server
//
// crowbar and chisel binaries should be in your PATH

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"time"
)

func run() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("go run main.go [testSocks] or [bench]")
	}
	for _, a := range args {
		switch a {
		case "testSocks":
			testSocks()
		case "bench":
			// bench()
		}
	}
}

func main() {
	startServer()

	time.Sleep(100 * time.Millisecond)

	startClient()

	time.Sleep(100 * time.Millisecond)

	defer func() {
		if r := recover(); r != nil {
			log.Print(r)
		}
	}()

	testSocks()
}

func startServer() {
	server := exec.Command("./chisel-darwin_amd64", "server",
		"--v",
		"--key", "foobar",
		"--reverse",
		"--socks5",
		"--port", "8001",
		"--backend", "https://qq.com")
	server.Stdout = os.Stdout
	if err := server.Start(); err != nil {
		fmt.Println(err)
	}
}

func startClient() {
	client := exec.Command("./chisel-darwin_amd64", "client",
		"--v",
		"--fingerprint", "OHclTPr2X1+S7CdRWW7dLFP7SwgtZy6jub2UmnpbTXw=",
		"http://127.0.0.1:8001",
		"R:127.0.0.1:12347:socks")
	client.Stdout = os.Stdout
	if err := client.Start(); err != nil {
		fmt.Println(err)
	}
}

func testSocks() {
	proxyURL, _ := url.Parse("socks5://rdd1:yyds@127.0.0.1:12347")

	// create an HTTP client with the SOCKS5 proxy transport
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}

	// use the client to make an HTTP request
	resp, err := client.Get("https://example.com")
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		panic("failed to get HTTP response using socks5")
	} else {
		fmt.Println("Success")
	}
}
