package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var count int = 0

func hello(w http.ResponseWriter, req *http.Request) {
	count++
    time.Sleep(2 * time.Second)
	fmt.Println(fmt.Sprintf("ID %d: %s", count, req.URL.Path))
	fmt.Fprintf(w, "you found me mother fucker\n")
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	data, _ := ioutil.ReadFile("big.txt")
	lines := strings.Split(string(data), "\n")
	lineCount := len(lines)
	offset := rand.Intn(lineCount) % 501

	for _, line := range lines[offset:500] {
		http.HandleFunc(fmt.Sprintf("/%s", line), hello)
	}

	http.HandleFunc("/", http.NotFound)

	fmt.Println("Offset: ", offset)
	http.ListenAndServe(":3000", nil)
}
