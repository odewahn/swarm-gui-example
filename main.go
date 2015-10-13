package main

import (
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/andlabs/ui"
	"github.com/joho/godotenv"
	"github.com/samalba/dockerclient"
)

// Container shows metadata about containers running on the cluster
type Container struct {
	Image string
	Name  string
}

var (
	w      ui.Window
	docker *dockerclient.DockerClient
)

func gui() {
	connect()
	var c Container

	// Table of running containers
	table := ui.NewTable(reflect.TypeOf(c))
	go updateTable(table)
	killBtn := ui.NewButton("Kill")
	tableStack := ui.NewVerticalStack(table, killBtn)
	containerListGrp := ui.NewGroup("Containers", tableStack)
	containerListGrp.SetMargined(true)

	//Stack for the control
	l := ui.NewLabel("Image to start")
	imageName := ui.NewTextField()
	imageName.SetText("ipython/scipystack")
	startBtn := ui.NewButton("Launch")

	startBtn.OnClicked(func() {
		go Start(imageName.Text())
	})

	killBtn.OnClicked(func() {
		c := table.Selected()
		table.Lock()
		d := table.Data().(*[]Container)
		//this makes a shallow copy of the structure so that we can access elements per
		//   http://giantmachines.tumblr.com/post/51007535999/golang-struct-shallow-copy
		newC := *d
		table.Unlock()
		go Kill(newC[c].Name)
	})

	controlStack := ui.NewVerticalStack(l, imageName, startBtn)
	controlGrp := ui.NewGroup("Start Images", controlStack)
	controlGrp.SetMargined(true)

	// Now make a new 2 column stack
	mainStack := ui.NewHorizontalStack(controlGrp, containerListGrp)
	mainStack.SetStretchy(0)
	mainStack.SetStretchy(1)

	//stack := ui.NewVerticalStack(table)

	w = ui.NewWindow("Window", 600, 300, mainStack)
	w.SetMargined(true)

	w.OnClosing(func() bool {
		ui.Stop()
		return true
	})
	w.Show()
}

func updateTable(table ui.Table) {
	for {
		running := ps()
		table.Lock()
		d := table.Data().(*[]Container)
		*d = running
		table.Unlock()
		time.Sleep(1 * time.Second)
	}
}

func connect() {
	// Load the environment variables we need
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Set up connection to swarm
	tlsConfig, err := GetTLSConfig(os.Getenv("SWARM_CREDS_DIR"))
	if err != nil {
		log.Fatal("Could not find TLS certificate.")
	}
	// Setup the docker host
	docker, err = dockerclient.NewDockerClient(os.Getenv("DOCKER_HOST"), tlsConfig)
	if err != nil {
		log.Fatal("Error initializing docker: ", err)
	}
	log.Println("Swarm connection inialized", docker)

}

func ps() []Container {
	var out []Container
	// Get only running containers
	containers, err := docker.ListContainers(false, false, "")
	if err != nil {
		log.Fatal(err)
	}
	for _, c := range containers {
		//log.Println(c.Id)
		//container, _ := docker.InspectContainer(c.Id)
		cNew := Container{
			Image: c.Image,
			Name:  strings.Split(c.Names[0], "/")[2],
		}
		out = append(out, cNew)
	}
	return out
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
