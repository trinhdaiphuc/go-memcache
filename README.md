# Go in-memory cache

## Overview
The memcache library provides a generic, type-safe in-memory key-value store with support for key expiration (TTL) and concurrency management. The library is designed to handle high concurrency scenarios while ensuring data consistency and safe access to the map data structure.

## Features

- Generic Map Interface: Supports keys and values of any comparable type.
- Concurrency-Safe Operations: Use of a command-based event loop to safely handle concurrent read/write operations.
- Key Expiration: Allows setting TTL (Time-To-Live) for individual keys or the entire map.
- Automatic Expiration Handling: Expired keys are automatically cleaned up.
- Standard Map Operations: Provides basic operations like set, get, delete, and retrieval of keys and values.

## Installation
To use this library, simply import it into your Go project:

```go
import "github.com/trinhdaiphuc/go-memcache"
```

## Usage

### Creating a New Map
You can create a new map instance by calling the NewMap function:

```go
m := memcache.NewMap[int, string]()
```
This creates a new map where keys are of type int and values are of type string.

### Basic Operations
#### Set a Key-Value Pair
To store a key-value pair in the map:

```go
m.Set(1, "value1")
```
This will store the value "value1" with the key 1.

#### Get a Value by Key
To retrieve a value by key:

```go
value, ok := m.Get(1)
if ok {
    fmt.Println("Value:", value)
} else {
    fmt.Println("Key not found or expired")
}
```
This will return the value associated with the key 1 if it exists and is not expired.

#### Delete a Key
To delete a key from the map:

```go
m.Delete(1)
```
This will remove the key 1 and its associated value from the map.

#### Get All Keys
To retrieve all keys from the map:

```go
keys := m.Keys()
fmt.Println("Keys:", keys)
```

#### Get All Values
To retrieve all values from the map:

```go
values := m.Values()
fmt.Println("Values:", values)
```

#### Get Map Length
To get the number of key-value pairs in the map:

```go
length := m.Len()
fmt.Println("Number of items in the map:", length)
```

### Key Expiration (TTL)

#### Set TTL for a Key

To set a TTL for a specific key:

```go
m.ExpireKey(1, time.Second*30)
```
This will set the key 1 to expire in 30 seconds.

#### Get TTL for a Key
To get the remaining TTL for a specific key:

```go
ttl := m.TTLKey(1)
fmt.Println("Remaining TTL:", ttl)
```

#### Set TTL for the Entire Map
To set a TTL for all keys in the map:

```go
m.Expire(time.Minute * 10)
```

This will set all keys in the map to expire in 10 minutes.

#### Check if the Map is Expired
To check if the entire map is expired:

```go
if m.IsExpired() {
    fmt.Println("Map has expired")
}
```

### Concurrency Considerations
The `go-memcache` library uses an event loop mechanism to handle concurrency. Each command (such as Set, Get, Delete) is executed sequentially through a command channel to ensure thread safety.

#### Internal Cleanup
The map automatically cleans up expired keys in the background using a ticker. This ensures that expired keys do not consume memory unnecessarily.

## Example
Here is a simple example of how to use the memcache library:

```go
package main

import (
    "fmt"
    "time"
	
    "github.com/trinhdaiphuc/go-memcache"
)

type User struct {
    ID       int64
    Username string
    Email    string
}

func main() {
    userMap := memcache.NewMap[int, User]()

    userMap.Set(1, User{ID: 1, Username: "user1", Email: "user1@gmail.com"})
    userMap.ExpireKey(1, time.Second*10)

    value, ok := userMap.Get(1)
    if ok {
        fmt.Println("Retrieved User:", value)
    } else {
        fmt.Println("User not found or expired")
    }

    time.Sleep(time.Second * 15)

    value, ok = userMap.Get(1)
    if !ok {
        fmt.Println("User has expired")
    }
}
```
In this example, the user with ID 1 is stored in the map and is set to expire in 10 seconds. After 15 seconds, the user is no longer available, demonstrating the TTL feature.

## Conclusion
The memcache library provides a powerful and flexible in-memory cache with support for concurrency and key expiration. It's suitable for applications that require high-performance caching with safe concurrent access.