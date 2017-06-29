package trigger

import (
	"os/exec"
	"time"
	"log"
)

type CmdTrigger struct {
	Command string
	Interval time.Duration
	trigger chan bool
	stop chan bool
	mask chan bool
}

func NewCmdTrigger(cmd string, interval time.Duration) *CmdTrigger {
	c := new(CmdTrigger)
	c.Command = cmd
	c.Interval = interval
	c.trigger = make(chan bool)
	c.mask = make(chan bool)
	c.stop = make(chan bool)
	return c
}

func (c *CmdTrigger) triggerLoop() {
	timer := time.NewTimer(c.Interval)
	for {
		select {
		case <-c.stop:
			timer.Stop()
			return
		case mask := <-c.mask:
			if mask {
				timer.Stop()
			} else {
				timer.Reset(c.Interval)
			}
		case <-timer.C:
			log.Printf("Running Trigger: %s\n", c.Command)
			cmd := exec.Command(c.Command)
			if err := cmd.Run(); err == nil {
				log.Printf("Trigger tripped")
				c.trigger <- true
			}
			timer.Reset(c.Interval)
		}
	}
}

func (c *CmdTrigger) Start() {
	log.Printf("Starting trigger wait loop")
	go c.triggerLoop()
}

func (c *CmdTrigger) Stop() {
	c.stop <- true
}

func (c *CmdTrigger) Mask() {
	c.mask <- true
}

func (c *CmdTrigger) Unmask() {
	c.mask <- false
}

func (c *CmdTrigger) Wait() {
	<- c.trigger
}
