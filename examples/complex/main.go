package main

import (
	"github.com/danryan/hal"
	_ "github.com/danryan/hal/adapter/irc"
	_ "github.com/danryan/hal/adapter/shell"
	_ "github.com/danryan/hal/adapter/slack"
	"github.com/danryan/hal/examples/complex/scripts"
	"log"
	"os"
)

// HAL is just another Go package, which means you are free to organize things
// however you deem best.

// You can define your handlers in the same file...
var openDoorsHandler = hal.Respond(`open the pod bay doors`, func(res *hal.Response) error {
	return res.Reply("I'm sorry, Dave. I can't do that.")
})

func Run() int {
	robot, err := hal.NewRobot()
	if err != nil {
		log.Println(err)
		return 1
	}

	// Or define them inside another function...
	var fooHandler = hal.Respond(`foo`, func(res *hal.Response) error {
		return res.Send("BAR")
	})

	// Or use the underlying hal.Listener struct...
	var tableFlipHandler = &hal.Listener{
		Method:  hal.HEAR,
		Pattern: `tableflip`,
		Handler: func(res *hal.Response) error {
			return res.Send(`(╯°□°）╯︵ ┻━┻`)
		},
	}

	// Or stick them in an entirely different package, and reference them
	// exactly in the ways you would expect.
	robot.Handle(
		scripts.PingHandler,
		scripts.SynHandler,
		openDoorsHandler,
		fooHandler,
		tableFlipHandler,
		// Or even inline!
		hal.Hear(`yo`, func(res *hal.Response) error {
			return res.Send("lo")
		}),
	)

	if err := robot.Run(); err != nil {
		log.Println(err)
		return 1
	}
	return 0
}

func main() {
	os.Exit(Run())
}
