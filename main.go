package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/andlabs/ui"
	"github.com/joho/godotenv"
	"github.com/odewahn/swarm-manager/manager"
	"github.com/samalba/dockerclient"
)

var (
	w      ui.Window
	docker *dockerclient.DockerClient
)

func gui() {
	connect()
	b := ui.NewButton("Click me")
	stack := ui.NewVerticalStack(b)

	b.OnClicked(func() {
		ps()
	})

	w = ui.NewWindow("Window", 400, 500, stack)
	w.OnClosing(func() bool {
		ui.Stop()
		return true
	})
	w.Show()
}

func connect() {
	// Load the environment variables we need
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Set up connection to swarm
	tlsConfig, err := manager.GetTLSConfig(os.Getenv("SWARM_CREDS_DIR"))
	if err != nil {
		log.Fatal("Could not create TLS certificate.")
	}
	// Setup the docker host
	docker, err = dockerclient.NewDockerClient(os.Getenv("DOCKER_HOST"), tlsConfig)
	if err != nil {
		log.Fatal("Error initializing docker: ", err)
	}
	log.Println("Swarm connection inialized", docker)

}

func ps() {
	// Get only running containers
	containers, err := docker.ListContainers(false, false, "")
	if err != nil {
		log.Fatal(err)
	}
	for _, c := range containers {
		//log.Println(c.Id)
		container, _ := docker.InspectContainer(c.Id)
		fmt.Printf("%s \t %s \t %s\n",
			c.Image,
			container.Config.Hostname+"."+container.Config.Domainname,
			strings.Split(c.Names[0], "/")[2])
	}
}

func main() {

	// This runs the code that displays our GUI.
	// All code that interfaces with package ui (except event handlers) must be run from within a ui.Do() call.
	go ui.Do(gui)

	err := ui.Go()
	if err != nil {
		log.Print(err)
	}

}
