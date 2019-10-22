# muxprom

> Add prometheus metrics to a [Mux](https://github.com/gorilla/mux) server easily.

## Usage

```go
package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/AverageMarcus/muxprom"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/healthz", healthz).Methods("GET")

	mp := muxprom.MuxProm{
		MetricName:  "requests", // Default if not provided
		MetricsPath: "/metrics", // Default if not provided
	}
	mp.RegisterPrometheus(router)

	log.Println("Listening on port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
```
