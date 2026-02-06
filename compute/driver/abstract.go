package drivers

import (
	"github.com/rehiy/cloudgo/compute"
)

type AbstractDriver struct{}

func NewAbstractDriver() *AbstractDriver {
	return &AbstractDriver{}
}

// List all instance
func (p *AbstractDriver) ListNodes() ([]*compute.Node, error) {
	return nil, nil
}

// Detail instance by Id
func (p *AbstractDriver) DetailNode(id string) (*compute.Node, error) {
	return nil, nil
}

// Create new instance
func (p *AbstractDriver) CreateNode(opts *compute.NodeCreateOpts) (*compute.Node, error) {
	return nil, nil
}

// Destroy an existing instance
func (p *AbstractDriver) DestroyNode(node *compute.Node) error {
	return nil
}

// Reboot instance
func (p *AbstractDriver) RebootNode(node *compute.Node) error {
	return nil
}

// Start instance
func (p *AbstractDriver) StartNode(node *compute.Node) error {
	return nil
}

// Stop instance
func (p *AbstractDriver) StopNode(node *compute.Node) error {
	return nil
}

// Get the current state of instance
func (p *AbstractDriver) GetNodeState(node *compute.Node) (compute.NodeState, error) {
	return compute.NodeState(""), nil
}

// Get the Console url for instance
func (p *AbstractDriver) GetNodeConsole(node *compute.Node) (string, error) {
	return "", nil
}

// Get the public IP address of instance
func (p *AbstractDriver) GetNodePublicIp(node *compute.Node) (string, error) {
	return "", nil
}

// Get the private IP address of instance
func (p *AbstractDriver) GetNodePrivateIp(node *compute.Node) (string, error) {
	return "", nil
}

// List all available storage volumes for instance
func (p *AbstractDriver) ListVolumes(node *compute.Node) ([]*compute.StorageVolume, error) {
	return nil, nil
}

// Attach volume to instance
func (p *AbstractDriver) AttachVolume(node *compute.Node, volume *compute.StorageVolume) error {
	return nil
}

// Detach volume from instance
func (p *AbstractDriver) DetachVolume(node *compute.Node, volume *compute.StorageVolume) error {
	return nil
}

// List all snapshots for instance
func (p *AbstractDriver) ListSnapshots(node *compute.Node) ([]*compute.VolumeSnapshot, error) {
	return nil, nil
}

// Create snapshot for instance
func (p *AbstractDriver) CreateSnapshot(node *compute.Node, name string) (*compute.VolumeSnapshot, error) {
	return nil, nil
}

// Destroy snapshot for instance
func (p *AbstractDriver) DestroySnapshot(node *compute.Node, snapshot *compute.VolumeSnapshot) error {
	return nil
}

// Apply snapshot to instance
func (p *AbstractDriver) ApplySnapshot(node *compute.Node, snapshot *compute.VolumeSnapshot) error {
	return nil
}

// List all available images for instance
func (p *AbstractDriver) ListImages() ([]*compute.NodeImage, error) {
	return nil, nil
}

// Apply Image to instance
func (p *AbstractDriver) ApplyImage(node *compute.Node, image *compute.NodeImage) error {
	return nil
}

// List all available sizes for instance
func (p *AbstractDriver) ListSizes() ([]*compute.NodeSize, error) {
	return nil, nil
}

// Resize instance
func (p *AbstractDriver) ResizeNode(node *compute.Node, opts *compute.NodeResizeOpts) error {
	return nil
}

// List all available locations for instance
func (p *AbstractDriver) ListLocations() ([]*compute.Location, error) {
	return nil, nil
}
