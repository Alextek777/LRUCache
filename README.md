# Thread Safe LRUCache Implementation in Go

[![Go Reference](https://pkg.go.dev/badge/github.com/yourusername/yourrepo.svg)](https://pkg.go.dev/github.com/yourusername/yourrepo)
[![Tests](https://github.com/yourusername/yourrepo/actions/workflows/tests.yml/badge.svg)](https://github.com/yourusername/yourrepo/actions/workflows/tests.yml)

A thread-safe LRU (Least Recently Used) cache implementation in Go with TTL support.

## Features

- 🚀 Thread-safe operations using sync.RWMutex
- ⏳ Time-based eviction (TTL)
- 📏 Size-based eviction (LRU)
- 🔄 Automatic cleanup of expired items
- 📊 Metrics and statistics (optional)

## Installation
```sh
go get github.com/Alextek777/LRUCache
````


## Usage
```golang
package main

import (
	"fmt"
	"time"

	"github.com/Alextek777/LRUCache"
)

func main() {
	cache := LRUCache.New[int, string](10, time.Second)

	cache.Push(1, "one")
	cache.Push(2, "two")
	cache.Push(3, "three")

	time.Sleep(2 * time.Second)

	cache.Push(4, "four") // only this value will be left
	v, ok := cache.Get(4)
	if !ok {
		fmt.Println("value not found")
	}

	fmt.Println("found value: ", v)
	cache.Remove(4)

	cache.Display()
}
```

### Output
```text
LRUCache
Key:  2  Value:  two
Key:  3  Value:  three

LRUCache TTL List
Key:  2  TTL :  2025-04-25 10:28:23.04043588 +0300 MSK m=+60.000050486
Key:  3  TTL :  2025-04-25 10:28:23.040436962 +0300 MSK m=+60.000051568
```

```golang
package main

import (
	"time"

	"github.com/Alextek777/LRUCache"
)

func main() {

	cache := LRUCache.New[int, string](10, time.Second)

	cache.Push(1, "one")
	cache.Push(2, "two")
	cache.Push(3, "three")

	time.Sleep(2 * time.Second)

	cache.Push(4, "four") // only this value will be left

	cache.Display()
}
```

### Output
```text
LRUCache
Key:  4  Value:  four

LRUCache TTL List
Key:  4  TTL :  2025-04-25 10:47:19.679317282 +0300 MSK m=+3.001068451
```
