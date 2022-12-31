# kraken
Kraken project is golang extensions library
### Code Example      
```go
package main

import (
	"fmt"
	"time"
	"math/rand"
	"github.com/koomox/kraken/cache"
	"github.com/koomox/kraken/uuid"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	store := cache.NewWithStringComparator()
	var elements []string
	for i := 0; i < 10; i++ {
		elements = append(elements, uuid.NewUUID())
	}
	for k, v := range elements {
		store.Put(fmt.Sprintf("%v", k), v, 60 * time.Second)
	}
	store.Put("name", "kraken", 60 * time.Second)
	store.Remove("1")
	ids := []string{"name", "1", fmt.Sprintf("%v", rand.Intn(len(elements)))}
	for _, id := range ids {
		if v := store.Get(id); v != nil {
			fmt.Println(v.(string))
		}
	}
	store.CallbackFunc(func(v interface{}){
		if v != nil {
			fmt.Println(v.(string))
		}
	})
	store.CancelFunc(func(v interface{}) bool {
		if v != nil {
			if (v.(string) == "kraken") {
				fmt.Println(v.(string))
				return true
			}
		}
		return false
	})
	b, err := store.ToJSON()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(b))
}
```