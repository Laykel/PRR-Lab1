package clock

import (
	"fmt"
	"time"
)

func Clock() {
    t1 := time.Now()

    t2 := time.Now()

    diff := t2.Sub(t1)

    fmt.Println(t1)
    fmt.Println(t2)
    fmt.Println(diff)
}
