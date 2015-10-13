package main

import (
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/andlabs/ui"
	"github.com/joho/godotenv"
	"github.com/odewahn/swarm-manager/db"
	"github.com/odewahn/swarm-manager/manager"
	"github.com/odewahn/swarm-manager/models"
	"github.com/samalba/dockerclient"
)

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
	container_list_grp := ui.NewGroup("Containers", table)
	container_list_grp.SetMargined(true)

	//Stack for the control
	l := ui.NewLabel("Image to start")
	imageName := ui.NewTextField()
	imageName.SetText("ipython/scipystack")
	start_btn := ui.NewButton("Launch")

	start_btn.OnClicked(func() {
		m := &models.Container{
			Image:      imageName.Text(),
			User:       "odewahn",
			Domainname: "i3.odewahn.com",
		}
		status := make(chan string)

		//ui.NewForeignEvent(status, func() {})
		go manager.Start(m, status)

		//<-status //block until we get a message back that the status record is ready

	})

	control_stack := ui.NewVerticalStack(l, imageName, start_btn)
	control_grp := ui.NewGroup("Start Images", control_stack)
	control_grp.SetMargined(true)

	// Now make a new 2 column stack
	main_stack := ui.NewHorizontalStack(control_grp, container_list_grp)
	main_stack.SetStretchy(0)
	main_stack.SetStretchy(1)

	//stack := ui.NewVerticalStack(table)

	w = ui.NewWindow("Window", 600, 300, main_stack)
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

	manager.Init()
	db.Init()

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
		c_new := Container{
			Image: c.Image,
			Name:  strings.Split(c.Names[0], "/")[2],
		}
		out = append(out, c_new)
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
