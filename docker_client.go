package main

import (
	"crypto/rand"
	"log"
	"os"

	"github.com/samalba/dockerclient"
)

// From https://www.socketloop.com/tutorials/golang-how-to-generate-random-string
func getHostName() string {
	dictionary := "0123456789abcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, 12)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(bytes)
}

// Start a container
func Start(imageName string) {

	hostName := getHostName() //Get a random name to use as a hostname
	log.Printf("Starting image %s as name %s", imageName, hostName)

	// Create the container
	containerConfig := &dockerclient.ContainerConfig{
		Image: imageName,
		Cmd:   []string{"/bin/sh", "-c", "ipython notebook --no-browser --port 8888 --ip=* --NotebookApp.allow_origin=*"},
		ExposedPorts: map[string]struct{}{
			"8888/tcp": {},
		},
		Hostname:   hostName,
		Domainname: os.Getenv("DOMAIN_NAME"),
	}
	containerID, err := docker.CreateContainer(containerConfig, hostName)
	if err != nil {
		log.Println(err)
	}

	hostConfig := &dockerclient.HostConfig{
		PublishAllPorts: true,
	}
	err = docker.StartContainer(containerID, hostConfig)
	if err != nil {
		log.Println(err)
	}

	log.Printf("Started container %s", containerID)

}

// Kill a container
func Kill(containerID string) {
	log.Println("Stopping container ", containerID)
	err := docker.StopContainer(containerID, 5)
	if err != nil {
		log.Println("Could not kill container ", containerID, err)
	}
	docker.RemoveContainer(containerID, true, true)
	if err != nil {
		log.Println("Could not remove container ", containerID, err)
	}
	log.Println("Removed container ", containerID)
}
