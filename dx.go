package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/dazeus/dazeus-go"
)

func handleCommand(dz *dazeus.DaZeus, ev dazeus.Event) {
	if len(ev.Params) > 1 {
		// The full message after the command sits in ev.Params[0]. Each word of ev.Params[0] is put into ev.Params[1..n].
		maybeDice := ev.Params[0]
		fmt.Printf("maybeDice: %v\n", maybeDice)
		if num, err := strconv.Atoi(maybeDice); err == nil && num > 0 {
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			random := r.Intn(num) + 1
			ev.Reply(fmt.Sprintf("Throwing a d%d... It is: %d", num, random), true)
			return
		}
	}
	ev.Reply("Sorry, cannot interpret this command as a dice roll!", true)
}

func main() {
	connStr := "unix:/tmp/dazeus.sock"
	if len(os.Args) > 1 {
		connStr = os.Args[1]
	}
	defer func() {
		if p := recover(); p != nil {
			fmt.Printf("Paniek! %v\n", p)
			debug.PrintStack()
		}
	}()

	dz, err := dazeus.ConnectWithLoggingToStdErr(connStr)
	if err != nil {
		panic(err)
	}

	if _, hlerr := dz.HighlightCharacter(); hlerr != nil {
		panic(hlerr)
	}

	_, err = dz.SubscribeCommand("throw", dazeus.NewUniversalScope(), func(ev dazeus.Event) {
		handleCommand(dz, ev)
	})
	if err != nil {
		panic(err)
	}

	listenerr := dz.Listen()
	panic(listenerr)
}
