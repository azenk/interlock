package main

import (
	"flag"
	"os/exec"
	"github.com/azenk/interlock/semaphore"
	"github.com/azenk/interlock/trigger"
	"path/filepath"
	"log"
	"time"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"github.com/hashicorp/consul/api"
)

type Config struct {
	Trigger	string  `yaml:trigger`
	Action string  `yaml:action`
	Interval time.Duration `yaml:interval`
	Delay time.Duration  `yaml:delay`
	Semaphore struct{
			Type string `yaml:type`
			Path string `yaml:path`
			Max int  `yaml:max`
			Id string `yaml:id`
			Token string `yaml:token`
		}
}

const (
	Loading = iota
	TriggerWait
	PerformAction
	ActionComplete
	PreCheck
	HoldsSemaphore
	AcquireSemaphore
	ReleaseSemaphore
	PostCheck
	Error
)

func load_config(path string) *Config {
	c := new(Config)

	data,err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Unable to read config file: %s\n", err)
	}

	err = yaml.Unmarshal(data, c)
	if err != nil {
				log.Fatalf("error: %v", err)
	}
	c.Trigger,err = filepath.Abs(c.Trigger)
	if err != nil {
		log.Fatalf("Unable to resolve path: %s\n", err)
	}

	c.Action,err = filepath.Abs(c.Action)
	if err != nil {
		log.Fatalf("Unable to resolve path: %s\n", err)
	}

	if c.Semaphore.Type == "file" {
		c.Semaphore.Path,err = filepath.Abs(c.Semaphore.Path)
		if err != nil {
			log.Fatalf("Unable to resolve path: %s\n", err)
		}
	}

	return c
}

func main() {
	state := Loading
	log.Println("Loading")
	config_path := new(string)
	flag.StringVar(config_path, "config", "config.yml", "path to config file")
	flag.Parse()
	path,_ := filepath.Abs(*config_path)
	log.Println(path)
	cfg := load_config(path)
	log.Printf("trigger: %s", cfg.Trigger)
	// Set up semaphore
	var sem semaphore.Semaphore
	switch cfg.Semaphore.Type {
	case "file":
		log.Printf("Creating file semaphore at: %s\n", cfg.Semaphore.Path)
		sem = semaphore.NewSemaphoreFile(cfg.Semaphore.Path, cfg.Semaphore.Max)
	case "consul":
		log.Printf("Creating new consul semaphore.")
		client, err := api.NewClient(api.DefaultConfig())
		if err != nil {
		    log.Fatalf("Unable to connect to consul")
		}
		sem = semaphore.NewSemaphoreConsul(client.KV(), cfg.Semaphore.Path, cfg.Semaphore.Max)
	default:
		sem = semaphore.New(cfg.Semaphore.Max)
	}

	log.Println("Starting Trigger")
	t := trigger.NewCmdTrigger(cfg.Trigger, cfg.Interval)
	t.Start()
	defer t.Stop()
	state = HoldsSemaphore

	for {
		switch state {
		case TriggerWait:
			log.Printf("Waiting for trigger\n")
			t.Unmask()
			t.Wait()
			log.Printf("Trigger received!\n")
			t.Mask()
			state = PreCheck

		case HoldsSemaphore:
			log.Println("Checking to see if we already hold the semaphore")
			if hold, err := sem.Holds(cfg.Semaphore.Id); err != nil {
				state = Error
			} else if hold {
				state = ActionComplete
			} else {
				state = TriggerWait
			}

		case PreCheck:
			log.Println("Running pre checks")
			state = AcquireSemaphore

		case AcquireSemaphore:
			log.Println("Acquiring semaphore")
			if ok, err := sem.Acquire(cfg.Semaphore.Id); !ok || err != nil {
				log.Printf("Failed to acquire semaphore: %s\n",err)
				state = Error
			}
			state = PerformAction

		case PerformAction:
			log.Printf("Executing action: %s\n", cfg.Action)
			action := exec.Command(cfg.Action)
			if err := action.Run(); err != nil {
				log.Fatalf("Action failed: %s\n", err)
				state = Error
			}
			state = ActionComplete

		case ActionComplete:
			log.Println("Test if action was completed")
			state = PostCheck

		case PostCheck:
			log.Println("Running post checks")
			state = ReleaseSemaphore

		case ReleaseSemaphore:
			log.Println("Action completed, releasing semaphore")
			sem.Release(cfg.Semaphore.Id)
			state = TriggerWait

		case Error:
			log.Println("Error")
			break
		}
	}
}
