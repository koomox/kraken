# kraken
Kraken project is golang extensions library
### Code Example      
```go
package main

import (
	"fmt"
	"time"
	"math/rand"
	"github.com/koomox/kraken/memory"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	store := memory.NewStore()
	elements := []string{
		"1daf5910-dc70-4c19-baab-609c727f6cde",
		"5d0f8663-e8c2-4fa0-8cab-539918d79ebb",
		"b4a35d44-2a68-42b4-a54d-8e6a9e02f5be",
		"c1f2de02-d1e9-44fd-a779-1ee8d28bd152",
		"c2ae5949-75c1-4f43-8793-81a35599ddfc",
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
	b, err := store.ToJSON()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(b))
}
```