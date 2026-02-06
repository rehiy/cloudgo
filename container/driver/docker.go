package drivers

import (
	"github.com/rehiy/cloudgo/container"
)

type DockerDriver struct {
	// add any necessary fields here
}

func (p *DockerDriver) ListContainers() ([]*container.ContainerInfo, error) {
	// TODO
	return nil, nil
}
