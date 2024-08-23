package memcache

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type User struct {
	ID       int64
	Username string
	Email    string
}

func TestNewMap(t *testing.T) {
	userMap := NewMap[int64, User]()
	userMap.Set(1, User{ID: 1, Username: "user1", Email: "user1@gmail.com"})
	userMap.Set(2, User{ID: 2, Username: "user2", Email: "user2@gmail.com"})
	userMap.Set(3, User{ID: 3, Username: "user3", Email: "user3@gmail.com"})

	assert.Equalf(t, userMap.Len(), 3, "userMap.Len() = %d; want 3", userMap.Len())
	user, ok := userMap.Get(1)
	assert.Truef(t, ok, "userMap.Get(1) = %v; want true", ok)
	assert.Equalf(t, user.ID, int64(1), "user.ID = %d; want 1", user.ID)

	userMap.Delete(1)
	assert.Equalf(t, userMap.Len(), 2, "userMap.Len() = %d; want 2", userMap.Len())
	_, ok = userMap.Get(1)
	assert.Falsef(t, ok, "userMap.Get(1) = %v; want false", ok)

	userMap.ExpireKey(2, time.Second)
	time.Sleep(time.Second)
	_, ok = userMap.Get(2)
	assert.Falsef(t, ok, "userMap.Get(2) = %v; want false", ok)
}

func TestConcurrencyMap(t *testing.T) {
	userMap := NewMap[int, User]()

	start := time.Now()

	wg := &sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			userMap.Set(i, User{ID: int64(i), Username: "user", Email: fmt.Sprintf("user%d@gmail.com", i)})
			userMap.ExpireKey(i, time.Second)
			userMap.TTLKey(i)
			wg.Done()
		}(i)
	}

	wg.Wait()

	elapsed := time.Since(start)
	fmt.Printf("Set and Get 1000 keys in %s\n", elapsed)
	assert.Equalf(t, userMap.Len(), 1000, "userMap.Len() = %d; want 1000", userMap.Len())
	assert.Equalf(t, len(userMap.Keys()), 1000, "len(userMap.Keys()) = %d; want 1000", len(userMap.Keys()))
	assert.Equalf(t, len(userMap.Values()), 1000, "len(userMap.Values()) = %d; want 1000", len(userMap.Values()))

	time.Sleep(30 * time.Second)
	assert.Equalf(t, userMap.Len(), 0, "userMap.Len() = %d; want 0", userMap.Len())
}
