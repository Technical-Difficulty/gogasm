package gogasm

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"sync"
	// "fmt"
	"io"
	"math"
)

var file string = "../server/big.txt"

func BenchmarkStreamRead(b *testing.B) {
	for n := 0; n < b.N; n++ {
		file, _ := os.Open(file)

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)

		out := make([]string, 0)
		for scanner.Scan() {
			out = append(out, scanner.Text())
		}

		file.Close()
	}
}

func BenchmarkStringsRead(b *testing.B) {
	for n := 0; n < b.N; n++ {
		data, _ := ioutil.ReadFile(file)
		lines := strings.Split(string(data), "\n")

		out := make([]string, 0)

		for i := range lines {
			out = append(out, lines[i])
		}
	}
}

func BenchmarkGoRead(b *testing.B) {
	for n := 0; n < b.N; n++ {
		f, _ := os.Open(file)
		defer f.Close()
		//sync pools to reuse the memory and decrease the preassure on //Garbage Collector
		linesPool := sync.Pool{New: func() interface{} {
			lines := make([]byte, 500*1024)
			return lines
		}}
		stringPool := sync.Pool{New: func() interface{} {
			lines := ""
			return lines
		}}
		// slicePool := sync.Pool{New: func() interface{} {
		// 	lines := make([]string, 100)
		// 	return lines
		// }}
		r := bufio.NewReader(f)
		var wg sync.WaitGroup //wait group to keep track off all threads
		for {

			buf := linesPool.Get().([]byte)
			n, err := r.Read(buf)
			buf = buf[:n]
			if n == 0 {
				if err != nil {
					// fmt.Println(err)
					break
				}
				if err == io.EOF {
					break
				}
				break
			}
			nextUntillNewline, err := r.ReadBytes('\n') //read entire line

			if err != io.EOF {
				buf = append(buf, nextUntillNewline...)
			}

			wg.Add(1)
			go func(chunk []byte, linesPool *sync.Pool, stringPool *sync.Pool) {

				//another wait group to process every chunk further
				var wg2 sync.WaitGroup
				logs := stringPool.Get().(string)
				logs = string(chunk)
				linesPool.Put(chunk) //put back the chunk in pool
				//split the string by "\n", so that we have slice of logs
				logsSlice := strings.Split(logs, "\n")
				stringPool.Put(logs) //put back the string pool
				chunkSize := 100     //process the bunch of 100 logs in thread
				n := len(logsSlice)
				noOfThread := n / chunkSize
				if n%chunkSize != 0 { //check for overflow
					noOfThread++
				}
				length := len(logsSlice)
				//traverse the chunk
				for i := 0; i < length; i += chunkSize {

					wg2.Add(1)
					//process each chunk in saperate chunk
					go func(s int, e int) {
						for i := s; i < e; i++ {
							text := logsSlice[i]
							if len(text) == 0 {
								continue
							}
							// fmt.Println(text)
						}
						wg2.Done()

					}(i*chunkSize, int(math.Min(float64((i+1)*chunkSize), float64(len(logsSlice)))))
					//passing the indexes for processing
				}
				wg2.Wait() //wait for a chunk to finish
				logsSlice = nil
				wg.Done()
			}(buf, &linesPool, &stringPool)
		}
		wg.Wait()
	}
}
