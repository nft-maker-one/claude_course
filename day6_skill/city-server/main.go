package main

import (
	"fmt"
	"log"
	"net/http"
)

func cityHandler(w http.ResponseWriter, r *http.Request) {
	cityName := "上海"
	fmt.Fprintf(w, "%s", cityName)
	log.Printf("Received request, responding with city: %s", cityName)

}

func main() {
	http.HandleFunc("/city", cityHandler)
	fmt.Println("Server is running on http://localhost:8080...")
	fmt.Println("Test it by running: curl http://localhost:8080/city")

	// 启动服务监听 8080 端口
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
