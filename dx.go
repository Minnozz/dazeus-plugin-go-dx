package main

import (
	"runtime/debug"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"strconv"
	"time"

	"github.com/dazeus/dazeus-go"
)

var myCommand string

func handlePrivmsg(dz *dazeus.DaZeus, ev dazeus.Event) {
	if len(ev.Params) > 1 {
		// The message sits in evt.Params[0]
		if idx := strings.Index(ev.Params[0], myCommand); idx > -1 {
			if len(ev.Params[0]) > 2 {
				maybeDice := ev.Params[0][2:]
				fmt.Printf("maybeDice: %v\n", maybeDice)
				if num, err := strconv.Atoi(maybeDice); err == nil && num > 0 {
					r := rand.New(rand.NewSource(time.Now().UnixNano()))
					random := r.Intn(num+1)
					ev.Reply(fmt.Sprintf("Throwing a d%d... It is: %d", num, random), true)
					return
				}
			}
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

	if hl, hlerr := dz.HighlightCharacter(); hlerr != nil {
		panic(hlerr)
	} else {
		myCommand = hl + "d"
	}

	_, err = dz.Subscribe(dazeus.EventPrivMsg, func(ev dazeus.Event) {
		handlePrivmsg(dz, ev)
	})
	if err != nil {
		panic(err)
	}

	listenerr := dz.Listen()
	panic(listenerr)
}
