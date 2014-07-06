package hal

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// Robot receives messages from an adapter and sends them to listeners
type Robot struct {
	Name    string
	Alias   string
	Adapter Adapter

	handlers   []Handler
	signalChan chan os.Signal
}

// Handlers returns the robot's handlers
func (robot *Robot) Handlers() []Handler {
	return robot.handlers
}

// NewRobot returns a new Robot instance
func NewRobot() (*Robot, error) {
	robot := &Robot{
		Name:       Config.Name,
		signalChan: make(chan os.Signal, 1),
	}

	adapter, err := NewAdapter(robot)
	if err != nil {
		Logger.Error(err)
		return nil, err
	}
	robot.SetAdapter(adapter)

	return robot, nil
}

// Handle registers a new handler with the robot
func (robot *Robot) Handle(handlers ...Handler) {
	robot.handlers = append(robot.handlers, handlers...)
}

// Receive dispatches messages to our handlers
func (robot *Robot) Receive(msg *Message) error {
	Logger.Debugf("%s - robot received message", Config.AdapterName)
	for _, handler := range robot.handlers {
		response := NewResponse(robot, msg)
		err := handler.Handle(response)
		if err != nil {
			Logger.Error(err)
			return err
		}
	}
	return nil
}

// Run initiates the startup process
func (robot *Robot) Run() error {

	Logger.Info("starting robot")
	Logger.Infof("starting %s adapter", Config.AdapterName)
	go robot.Adapter.Run()
	// Start the HTTP server after the adapter, as adapter.Run() adds additional
	// handlers to the router.
	Logger.Debug("starting HTTP server")
	go http.ListenAndServe(`:`+string(Config.Port), Router)

	signal.Notify(robot.signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	stop := false
	for !stop {
		select {
		case sig := <-robot.signalChan:
			switch sig {
			case syscall.SIGINT, syscall.SIGTERM:
				stop = true
			}
		}
	}
	// Stop listening for new signals
	signal.Stop(robot.signalChan)

	// Initiate the shutdown process for our robot
	robot.Stop()

	return nil
}

// Stop initiates the shutdown process
func (robot *Robot) Stop() error {
	fmt.Println() // so we don't break up the log formatting when running interactively ;)
	Logger.Infof("stopping %s adapter", Config.AdapterName)

	robot.Adapter.Stop()
	Logger.Info("stopping robot")

	return nil
}

func (robot *Robot) respondRegex(pattern string) string {
	str := `^(?:`
	if robot.Alias != "" {
		str += `(?:` + robot.Alias + `|` + robot.Name + `)`
	} else {
		str += robot.Name
	}
	str += `[:,]?)\s+(?:` + pattern + `)`
	return str
}

// SetName sets robot's name
func (robot *Robot) SetName(name string) {
	robot.Name = name
}

// SetAdapter sets robot's adapter
func (robot *Robot) SetAdapter(adapter Adapter) {
	robot.Adapter = adapter
}
