package main

import (
	"flag"
	"bytes"
	"os/exec"
	"./semaphore"
	"path/filepath"
	"log"
	"time"
	"gopkg.in/yaml.v2"
	"io/ioutil"
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

func load_config(path string) *Config {
	c := new(Config)

	data,err := ioutil.ReadFile(path)
	log.Println(bytes.NewBuffer(data).String())

	err = yaml.Unmarshal(data, c)
	if err != nil {
				log.Fatalf("error: %v", err)
	}
	log.Println(c)
	c.Trigger,err = filepath.Abs(c.Trigger)
	if err != nil {
		log.Fatalf("Unable to resolve path: %s\n", err)
	}

	c.Action,err = filepath.Abs(c.Action)
	if err != nil {
		log.Fatalf("Unable to resolve path: %s\n", err)
	}

	c.Semaphore.Path,err = filepath.Abs(c.Semaphore.Path)
	if err != nil {
		log.Fatalf("Unable to resolve path: %s\n", err)
	}

	return c
}

func main() {
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
	default:
		sem = semaphore.New(cfg.Semaphore.Max)
	}

	for {
		log.Printf("Sleeping for %s\n", cfg.Interval)
		time.Sleep(cfg.Interval)

		trigger := exec.Command(cfg.Trigger)
		if err := trigger.Run(); err != nil {
			log.Printf("Trigger not tripped, continuing\n")
			continue
		} else {
			log.Printf("Trigger tripped, delaying for %s\n", cfg.Delay)
			time.Sleep(cfg.Delay)
		}

		log.Println("Acquiring semaphore")
		if ok, err := sem.Acquire(cfg.Semaphore.Id); !ok || err != nil {
			log.Printf("Failed to acquire semaphore: %s\n",err)
			continue
		}

		log.Printf("Executing action: %s\n", cfg.Action)
		action := exec.Command(cfg.Action)
		if err := action.Run(); err != nil {
			log.Fatalf("Action failed: %s\n", err)
			break
		}

		log.Println("Action completed, releasing semaphore")
		sem.Release(cfg.Semaphore.Id)
	}
}
