package controller

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"log"
)

func CreateNewContainer(image, user string) (string, error) {
	cli, err := client.NewEnvClient()
	ctx := context.Background()
	if err != nil {
		log.Println("error creating New docker client")
	}

	//container port
	containerPort, err := nat.NewPort("tcp", "80")
	if err != nil {
		panic("unable to get the port")
	}

	//port binding
	hostBinding := nat.PortBinding{
		HostIP:   "0.0.0.0",
		HostPort: "8080",
	}

	//binding the port with the host ip
	portBinding := nat.PortMap{containerPort: []nat.PortBinding{hostBinding}}
	cont, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			User:  user,
			Image: image,
		},
		&container.HostConfig{
			PortBindings: portBinding,
		},
		nil, "")
	if err != nil {
		panic(err)
	}

	//start container
	cli.ContainerStart(ctx, cont.ID, types.ContainerStartOptions{})
	fmt.Printf("Container ID: %s started\n", cont.ID)

	return cont.ID, err

}

func PauseContainer(contID string) error {
	cli, err := client.NewEnvClient()
	ctx := context.Background()
	if err != nil {
		panic("Error creating container environment %s ")
	}
	err = cli.ContainerStop(ctx, contID, nil)

	fmt.Println("Success")
	return err
}