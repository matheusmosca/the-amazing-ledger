# Benchmark & Profilling

We are using GHZ (https://ghz.sh/) to run some perf tests with the Ledger.
For profilling we are using a Go tool a called pprof: go get -u github.com/google/pprof

## Running Benchmarks
To run the benchmarks you need to:

- generate an image/protoset:  
    ```$ buf build -o protoimage.bin --path=proto/ledger```  
    ```$ mv protoimage.bin perftests```  
- create scenarios (see perftests/scenarios)  
  
[local]  
  
- start database docker  
    ```$ docker-compose -f docker-compose-dev.yml up```  
- start the server  
    ```$ make build; ./build/server```  
- run  
    ```$ cd perftests```  
    ```> perftests > $ go run .```  
  
## Running the Profiler
[local]

- start database docker  
    ```$ docker-compose -f docker-compose-dev.yml up```  
- start the server  
    ```$ make compile; make build; ./build/server -cpuprofile cpu.profile -memprofile mem.profile```  
- run  
    ```$ cd perftests```  
    ```> perftests > $ go run .```  
- results  
    ```go tool pprof -http=:8080 cpu.profile```  
    ```go tool pprof -http=:8080 mem.profile```