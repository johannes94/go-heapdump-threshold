# go-heapdump-thresshold

If a container uses more memory than assigned to it, Kubernetes will kill it with a OOMKilled error. The problem with this is that after the container was killed, it is hard for engineers to debug what caused the unexpected memory usage. This package helps solve this problem.

With `go-heap-dump` you can define a memory limit and a memory thresshold. If the memory consupmtion of your Go application exceeds the thresshold, pprof will be called to dump a memory profile to the filesystem.

## Usage

For a entire code and Kubernetes pod example see the `example` directory of this project.

1. Use the HeapProfiler in your Go application
```go
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

// Start a background process that checks the thresshold every 30 seconds and dumps a heap profile if necessary
ctx := context.Background()
go heapProfiler.DumpHeapOnThreshhold(ctx, 30*time.Second)
```
2. Use [Kubernetes Downward API](https://kubernetes.io/docs/concepts/workloads/pods/downward-api/) to propagate the memory limit
```yaml
...
    env:
    - name: MEMORY_LIMIT
      valueFrom:
        resourceFieldRef:
          containerName: example-app
          resource: limits.memory
...
```
3. Use a volume to store the heap profiles between container restarts
```yaml
...
    volumeMounts:
      - mountPath: ./dump
        name: dump
  volumes:
  - name: dump
    emptyDir: {}
...
```

Once the allocated memory in the container process reaches the defined threshhold of `0.80*limit` it will dump a pprof memory profile to the `/dump/heapdump/<datetime>`. By using an emptyDir this dump is persisted throughout container restarts. If you'd like to persist it throughout different Pods use a PV instead.
