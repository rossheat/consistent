package main

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	key       string
	debug     bool
	model     string
	instances int
	delay     int
}

func NewConfig() Config {

	key := flag.String("key", "", "Your Anthropic API key")
	debug := flag.Bool("debug", false, "Start the progrom in debug mode")
	model := flag.String("model", "claude-3-5-sonnet-20240620", "The name of the Anthropic model you'd like to question")
	instances := flag.Int("instances", 50, "The number times your question is sent to the model API")
	delay := flag.Int("delay", 500, "Milliseconds of delay between calling the Anthropic API")

	flag.Parse()

	if *key == "" {
		fmt.Println("Please provide your Anthropic API key")
		os.Exit(1)
	}
	return Config{key: *key, debug: *debug, model: *model, instances: *instances, delay: *delay}
}
