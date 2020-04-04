# opencensus-redis
opencensus hook for go-redis

## Usage

```go
import (
    "github.com/kit-x/opencensus-redis/ochook"
    "github.com/go-redis/redis/v7"
)

redisOpts := redis.Options{
    Addr: "127.0.0.1:6379",
}
traceOptions := []ochook.TraceOption{
    ochook.WithAllowRoot(true),
}

client := redis.NewClient(&redisOpts)
client.AddHook(ochook.New(traceOptions...))
client.Ping()
```