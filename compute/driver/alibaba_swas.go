package drivers

import (
	"github.com/rehiy/cloudgo/compute"
	"github.com/rehiy/cloudgo/provider"
	"github.com/rehiy/cloudgo/provider/alibaba"

	swas "github.com/alibabacloud-go/swas-open-20200601/client"
)

type AlibabaSwasDriver struct {
	client *alibaba.Client
	swas   *swas.Client
	rq     *provider.ReqeustParam
}

func NewAlibabaSwasDriver(rq *provider.ReqeustParam) *AlibabaSwasDriver {

	client := alibaba.NewClient(rq)
	swas, _ := client.Swas()

	return &AlibabaSwasDriver{client, swas, rq}

}

// List all instance
func (p *AlibabaSwasDriver) ListNodes() ([]*compute.Node, error) {
	return nil, nil
}

// Detail instance by Id
func (p *AlibabaSwasDriver) DetailNode(id string) (*compute.Node, error) {
	return nil, nil
}

// Create new instance
func (p *AlibabaSwasDriver) CreateNode(opts *compute.NodeCreateOpts) (*compute.Node, error) {
	return nil, nil
}

// Destroy an existing instance
func (p *AlibabaSwasDriver) DestroyNode(node *compute.Node) error {
	return nil
}

// Reboot instance
func (p *AlibabaSwasDriver) RebootNode(node *compute.Node) error {
	return nil
}

// Start instance
func (p *AlibabaSwasDriver) StartNode(node *compute.Node) error {
	return nil
}

// Stop instance
func (p *AlibabaSwasDriver) StopNode(node *compute.Node) error {
	return nil
}

// Get the current state of instance
func (p *AlibabaSwasDriver) GetNodeState(node *compute.Node) (compute.NodeState, error) {
	return compute.NodeState(""), nil
}

// Get the Console url for instance
func (p *AlibabaSwasDriver) GetNodeConsole(node *compute.Node) (string, error) {
	return "", nil
}

// Get the public IP address of instance
func (p *AlibabaSwasDriver) GetNodePublicIp(node *compute.Node) (string, error) {
	return "", nil
}

// Get the private IP address of instance
func (p *AlibabaSwasDriver) GetNodePrivateIp(node *compute.Node) (string, error) {
	return "", nil
}

// List all available storage volumes for instance
func (p *AlibabaSwasDriver) ListVolumes(node *compute.Node) ([]*compute.StorageVolume, error) {
	return nil, nil
}

// Attach volume to instance
func (p *AlibabaSwasDriver) AttachVolume(node *compute.Node, volume *compute.StorageVolume) error {
	return nil
}

// Detach volume from instance
func (p *AlibabaSwasDriver) DetachVolume(node *compute.Node, volume *compute.StorageVolume) error {
	return nil
}

// List all snapshots for instance
func (p *AlibabaSwasDriver) ListSnapshots(node *compute.Node) ([]*compute.VolumeSnapshot, error) {
	return nil, nil
}

// Create snapshot for instance
func (p *AlibabaSwasDriver) CreateSnapshot(node *compute.Node, name string) (*compute.VolumeSnapshot, error) {
	return nil, nil
}

// Destroy snapshot for instance
func (p *AlibabaSwasDriver) DestroySnapshot(node *compute.Node, snapshot *compute.VolumeSnapshot) error {
	return nil
}

// Apply snapshot to instance
func (p *AlibabaSwasDriver) ApplySnapshot(node *compute.Node, snapshot *compute.VolumeSnapshot) error {
	return nil
}

// List all available images for instance
func (p *AlibabaSwasDriver) ListImages() ([]*compute.NodeImage, error) {
	return nil, nil
}

// Apply Image to instance
func (p *AlibabaSwasDriver) ApplyImage(node *compute.Node, image *compute.NodeImage) error {
	return nil
}

// List all available sizes for instance
func (p *AlibabaSwasDriver) ListSizes() ([]*compute.NodeSize, error) {
	return nil, nil
}

// Resize instance
func (p *AlibabaSwasDriver) ResizeNode(node *compute.Node, opts *compute.NodeResizeOpts) error {
	return nil
}

// List all available locations for instance
func (p *AlibabaSwasDriver) ListLocations() ([]*compute.Location, error) {
	return nil, nil
}
