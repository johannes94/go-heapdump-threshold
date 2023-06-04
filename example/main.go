package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	heapprofiler "github.com/johannes94/go-heapdump-thresshold"
)

func main() {

	// In this example this env variable is set through K8s downward API
	// https://kubernetes.io/docs/concepts/workloads/pods/downward-api/
	strMemoryLimit := os.Getenv("MEMORY_LIMIT")
	limit, err := strconv.ParseUint(strMemoryLimit, 10, 64)
	if err != nil {
		fmt.Println(err)
	}

	// This create a NewHeapProfiler with following configuration
	// - Start to dump heap profiles when used heap is greater than or equal to 80% of the limit
	// - A limit defined by env variable MEMORY_LIMIT in bytes e.g. 100000000 for 100MB
	// - Create files at the "./dump" directory
	// - Wait for 1 Minute once a heap profile was written before writing another one
	heapProfiler := heapprofiler.NewHeapProfiler(0.80, limit, "./dump", 1*time.Minute)

	ctx := context.Background()

	// Start a background process that checks the thresshold every 30 seconds and dumps a heap profile if necessary
	go heapProfiler.DumpHeapOnThreshhold(ctx, 30*time.Second)

	// Start your other processes e. g. a HTTP listener
	// This example simulates a memory leak
	var data [][]int64
	fmt.Println("Starting example app")
	for {
		add := make([]int64, 700)
		data = append(data, add)
		time.Sleep(time.Microsecond * 100)
	}
}
