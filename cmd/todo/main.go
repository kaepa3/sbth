package main

import (
	"context"
	"time"

	"github.com/kaepa3/sbth"
)

func main() {
	addr := ""
	ctx, _ := context.WithTimeout(context.Background(), 7*time.Second)
	sbth.Scan(addr, ctx)

}
