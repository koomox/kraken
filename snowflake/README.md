# Snow-flake           
Custom Epoch
By default this package uses the Twitter Epoch of 1230771600000 or Jan 01 2009 01:00:00. You can set your own epoch value by setting snowflake.Epoch to a time in milliseconds to use as the epoch.
```
+----------------------------------------------------------------------------------------------+
| 1 Bit Unused | 41 Bit Timestamp | 5 Bit Datacenter ID | 5 Bit Worker ID | 12 Bit Sequence ID |
+----------------------------------------------------------------------------------------------+
```
### Usage          
Example Code:       
```go
package main

import (
	"fmt"
	"sync"
	"github.com/koomox/kraken/snowflake"
)

func main() {
	datetime := "2009-01-01 01:00:00"
	loc, _ :=  time.LoadLocation("UTC")
	dt, _ := time.ParseInLocation("2006-01-02 15:04:05", datetime, loc)
	fmt.Printf("epoch: %v\n", dt.UnixNano()/1000000)

	var wg sync.WaitGroup
	s, err := snowflake.NewSnowflake(int64(0), int64(0))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Printf("ID: %v\n", s.NextID())
		}()
	}
	wg.Wait()
}
```
