![gogasm-demo](https://i.imgur.com/mjnn11K.gif)

# gogasm
Directory scanner written in GO and orgasmicly fast

# usage
We are current still in development but feel free to check it out
```
git clone https://github.com/Technical-Difficulty/gogasm.git
```
Build the project
```
cd gogasm && go build
```
Run it!
```
./gogasm -w /path/to/wordlist -a localhost -p 80
```
# Testing
run the test server
```
cd server && go run server.go
```
run http client benchmarks to get mem and cpu profile 
```
go test -cpuprofile cpu.prof -memprofile mem.prof -bench .
```
then
```
go tool pprof mem.prof
```
or
```
go tool pprof cpu.prof
```
Test src benchmarks
```
go test ./src -bench=. -benchmem -benchtime=1000x
```
