## Cache Manger Plugin for Jeoga Go Framework

### Usage
```go
import "github.com/jeogago/fm/pkg/cache"
import "github.com/jeogago/fm/pkg/cacherdbms"

configStore := cacherdbms.Config{
  DB: db,
}
rdbms := cacherdbms.New(&configStore)
config := cache.Config{
  Store:     rdbms,
  GlobalKey: "jcache", // optional
}
cache := cache.New(&config)
```

## Usage
### Put
```go
cache.Put("name", "jeoga", 60) // Expire in 60 seconds
```

### Get
```go
cache.Get("name")
```

### Del
```go
cache.Del("name")
```

### Flush
```go
cache.Flush()
```

