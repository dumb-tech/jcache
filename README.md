# Just Cache

![jcache illustration](assets/images/logo.jpg)

**Just Cache** — It is a simple in-memory cache with support for item expiration and automatic cleanup.

## Description

The **jcache** package provides the ability to store data as a key-value with a specified time to live (TTL). There are two strategies for cleaning up expired items:
- **On The Fly** — removes expired items as they are found;
- **Collect** — collects expired items and removes them in bulk.

## Installation

Install the package using the command:

```bash
go get github.com/dumb-tech/jcache
```

## Usage
Below is an example of how to use the package:

```go
package main

import (
	"fmt"
	"time"

	"github.com/dumb-tech/jcache"
)

func main() {
	// Create a new cache with a cleanup interval of 1 minute and a capacity of 10000 items.
	cache := jcache.New(1*time.Minute, 10000)
	defer cache.Close()

	// Set a value with a TTL of 10 seconds.
	err := cache.Set("key", "value", 10*time.Minute)
	if err != nil {
		fmt.Println("Error setting value:", err)
		return
	}

	// Retrieve the value.
	val := cache.Get("key")
	fmt.Println("Value:", val)
}
```

## Documentation

+ `New(interval time.Duration, capacity int64) *JustCache`
Creates a new cache instance with the specified cleanup interval and capacity.

+ `Default() *JustCache`
Returns a cache instance with default settings.

+ `WithStrategy(strategy CleanupStrategy) *JustCache`
Sets the cleanup strategy and returns the modified cache.

+ `WithCleanupInterval(interval time.Duration) *JustCache`
Sets the cleanup interval and returns the modified cache.

+ `WithCapacity(capacity int64) *JustCache`
Sets the cache capacity and returns the modified cache.

+ `Get(key string) any`
Retrieves the value associated with the specified key.

+ `Has(key string) bool`
Checks whether the specified key exists in the cache.

+ `Item(key string) Item`
Returns a cache item (key and value) for the specified key.

+ `Set(key string, value any, ttl time.Duration) error`
Stores a key-value pair in the cache with a time-to-live duration.
Returns an error if the cache is full.

+ `Del(key string)`
Deletes the specified key from the cache.

+ `Keys() []string`
Returns a slice of all keys in the cache.

+ `Items() []Item`
Returns a slice of all cache items.

+ `Clear()`
Removes all items from the cache.

+ `Clean(now time.Time)`
Removes expired items from the cache based on the provided time.

+ `Close() error`
Stops the cleanup process and clears the cache.

## Testing
```bash
go test -v ./...
```