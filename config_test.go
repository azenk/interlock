package main

import (
	"testing"
	"path/filepath"
)

func TestLoadConfig(t *testing.T) {
	path, err := filepath.Abs("testing/sample_config.yml")
	if err != nil {
		t.Errorf("Unable to build path to config file: %s\n", err)
	}
	c := load_config(path)
	if c.Semaphore.Max != 1 {
		t.Errorf("Config file not loading properly")
	}

}
