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
	"github.com/koomox/kraken/snowflake"
	"time"
)

func main() {
	datetime := "2009-01-01 01:00:00"
	loc, _ := time.LoadLocation("UTC")
	dt, _ := time.ParseInLocation("2006-01-02 15:04:05", datetime, loc)
	fmt.Printf("epoch: %v\n", dt.UnixNano()/1000000)

	ch := make(chan bool, 1)
	snowflake.WithBackground(snowflake.NewSnowflake(0, 0))
	length := 10
	for i := 0; i < length; i++ {
		go func(i int, id int64) {
			fmt.Printf("%d: %v\n", i, id)
			ch <- true
		}(i, snowflake.NextID())
	}

	for i := 0; i < length; i++{
		select {
		case <-ch:
		case <-time.After(time.Second):
			fmt.Println("timeout")
		}
	}
}
```
