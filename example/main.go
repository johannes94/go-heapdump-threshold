package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	heapprofiler "github.com/johannes94/go-heapdump-thresshold"
)

func main() {
	var limit uint64 = os.Getenv("MEMORY_LIMIT")

	// This create a NewHeapProfiler with following configuration
	// - Start to dump heap profiles when used heap is greater than or equal to 80% of the limit
	// - A limit of 100MB
	// - Create files at the "./dump" directory
	// - Wait for 1 Minute once a heap profile was written before writing another one
	heapProfiler := heapprofiler.NewHeapProfiler(0.80, limit100MB, "./dump", 1*time.Minute)

	ctx := context.Background()

	// Start a background process that checks the thresshold every 30 seconds and dumps a heap profile if necessary
	go heapProfiler.DumpHeapOnThreshhold(ctx, 30*time.Second)

	// Start your other processes e. g. a HTTP listener
	if err := http.ListenAndServe(":8080", http.NotFoundHandler()); err != nil {
		fmt.Println(err)
	}
}
