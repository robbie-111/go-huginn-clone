package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"go-huginn-clone/router"
)

func findAvailablePort(startPort int) int {
	for port := startPort; port < 65535; port++ {
		addr := fmt.Sprintf(":%d", port)
		ln, err := net.Listen("tcp", addr)
		if err == nil {
			ln.Close()
			return port
		}
	}
	log.Fatal("사용 가능한 포트를 찾을 수 없습니다")
	return -1
}

func main() {
	r := router.New()

	startPort := 3001
	if p := os.Getenv("PORT"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil {
			startPort = parsed
		}
	}

	port := findAvailablePort(startPort)
	addr := fmt.Sprintf(":%d", port)

	if port != startPort {
		log.Printf("포트 %d 사용 중 → 포트 %d 로 자동 전환", startPort, port)
	}
	log.Printf("Huginn (Go) starting on http://localhost%s", addr)
	log.Printf("Login with any username/password combination")

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
