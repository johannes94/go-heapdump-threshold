# go-heapdump-thresshold

If a container uses more memory than assigned to it, Kubernetes will kill it with a OOMKilled error. The problem with this is that after the container was killed, it is hard for engineers to debug what caused the unexpected memory usage. This package