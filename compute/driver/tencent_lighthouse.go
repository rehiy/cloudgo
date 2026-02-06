package drivers

import (
	"github.com/rehiy/cloudgo/compute"
	"github.com/rehiy/cloudgo/provider"
	"github.com/rehiy/cloudgo/provider/tencent"

	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
)

type TencentLighthouseDriver struct {
	client     *tencent.Client
	lighthouse *lighthouse.Client
	rq         *provider.ReqeustParam
}

func NewTencentLighthouseDriver(rq *provider.ReqeustParam) *TencentLighthouseDriver {

	client := tencent.NewClient(rq)
	lighthouse, _ := client.Lighthouse()

	return &TencentLighthouseDriver{client, lighthouse, rq}

}

// List all instance
func (p *TencentLighthouseDriver) ListNodes() ([]*compute.Node, error) {

	resp, err := p.lighthouse.DescribeInstances(&lighthouse.DescribeInstancesRequest{})

	if err != nil {
		return nil, err
	}

	var nodes []*compute.Node

	for _, instance := range resp.Response.InstanceSet {

		// TODO: fix mapping
		state := compute.NodeState(*instance.InstanceState)
		osType := compute.OSType(*instance.OsName)

		node := &compute.Node{
			Id:        *instance.InstanceId,
			Name:      *instance.InstanceName,
			State:     state,
			PublicIp:  *instance.PublicAddresses[0],
			PrivateIp: *instance.PrivateAddresses[0],
			Size: &compute.NodeSize{
				Id: *instance.BundleId,
			},
			Image: &compute.NodeImage{
				Id:     "",
				Name:   *instance.OsName,
				OSType: osType,
			},
			Location: &compute.Location{
				Id: *instance.Zone,
			},
		}

		nodes = append(nodes, node)

	}

	return nodes, nil

}

// Detail instance by Id
func (p *TencentLighthouseDriver) DetailNode(id string) (*compute.Node, error) {

	resp, err := p.lighthouse.DescribeInstances(&lighthouse.DescribeInstancesRequest{
		InstanceIds: []*string{&id},
	})

	if err != nil {
		return nil, err
	}

	if len(resp.Response.InstanceSet) == 0 {
		return nil, nil
	}

	instance := resp.Response.InstanceSet[0]

	// TODO: fix mapping
	state := compute.NodeState(*instance.InstanceState)
	osType := compute.OSType(*instance.OsName)

	node := &compute.Node{
		Id:        *instance.InstanceId,
		Name:      *instance.InstanceName,
		State:     state,
		PublicIp:  *instance.PublicAddresses[0],
		PrivateIp: *instance.PrivateAddresses[0],
		Size: &compute.NodeSize{
			Id: *instance.BundleId,
		},
		Image: &compute.NodeImage{
			Id:     "",
			Name:   *instance.OsName,
			OSType: osType,
		},
		Location: &compute.Location{
			Id: *instance.Zone,
		},
	}

	return node, nil

}

// Create new instance
func (p *TencentLighthouseDriver) CreateNode(opts *compute.NodeCreateOpts) (*compute.Node, error) {

	resp, err := p.lighthouse.CreateInstances(&lighthouse.CreateInstancesRequest{
		// TODO: create instance
	})

	if err != nil {
		return nil, err
	}

	if len(resp.Response.InstanceIdSet) == 0 {
		return nil, nil
	}

	instanceId := resp.Response.InstanceIdSet[0]

	node, err := p.DetailNode(*instanceId)

	if err != nil {
		return nil, err
	}

	return node, nil

}

// Destroy an existing instance
func (p *TencentLighthouseDriver) DestroyNode(node *compute.Node) error {

	_, err := p.lighthouse.TerminateInstances(&lighthouse.TerminateInstancesRequest{
		InstanceIds: []*string{&node.Id},
	})

	return err

}

// Reboot instance
func (p *TencentLighthouseDriver) RebootNode(node *compute.Node) error {

	_, err := p.lighthouse.RebootInstances(&lighthouse.RebootInstancesRequest{
		InstanceIds: []*string{&node.Id},
	})

	return err

}

// Start instance
func (p *TencentLighthouseDriver) StartNode(node *compute.Node) error {

	_, err := p.lighthouse.StartInstances(&lighthouse.StartInstancesRequest{
		InstanceIds: []*string{&node.Id},
	})

	return err

}

// Stop instance
func (p *TencentLighthouseDriver) StopNode(node *compute.Node) error {

	_, err := p.lighthouse.StopInstances(&lighthouse.StopInstancesRequest{
		InstanceIds: []*string{&node.Id},
	})

	return err

}

// Get the current state of instance
func (p *TencentLighthouseDriver) GetNodeState(node *compute.Node) (compute.NodeState, error) {

	resp, err := p.lighthouse.DescribeInstances(&lighthouse.DescribeInstancesRequest{
		InstanceIds: []*string{&node.Id},
	})

	if err != nil {
		return "", err
	}

	if len(resp.Response.InstanceSet) == 0 {
		return "", nil
	}

	instance := resp.Response.InstanceSet[0]

	// TODO: fix mapping
	state := compute.NodeState(*instance.InstanceState)

	return state, nil

}

// Get the Console url for instance
func (p *TencentLighthouseDriver) GetNodeConsole(node *compute.Node) (string, error) {
	return "", nil
}

// Get the public IP address of instance
func (p *TencentLighthouseDriver) GetNodePublicIp(node *compute.Node) (string, error) {
	return "", nil
}

// Get the private IP address of instance
func (p *TencentLighthouseDriver) GetNodePrivateIp(node *compute.Node) (string, error) {
	return "", nil
}

// List all available storage volumes for instance
func (p *TencentLighthouseDriver) ListVolumes(node *compute.Node) ([]*compute.StorageVolume, error) {
	return nil, nil
}

// Attach volume to instance
func (p *TencentLighthouseDriver) AttachVolume(node *compute.Node, volume *compute.StorageVolume) error {
	return nil
}

// Detach volume from instance
func (p *TencentLighthouseDriver) DetachVolume(node *compute.Node, volume *compute.StorageVolume) error {
	return nil
}

// List all snapshots for instance
func (p *TencentLighthouseDriver) ListSnapshots(node *compute.Node) ([]*compute.VolumeSnapshot, error) {
	return nil, nil
}

// Create snapshot for instance
func (p *TencentLighthouseDriver) CreateSnapshot(node *compute.Node, name string) (*compute.VolumeSnapshot, error) {
	return nil, nil
}

// Destroy snapshot for instance
func (p *TencentLighthouseDriver) DestroySnapshot(node *compute.Node, snapshot *compute.VolumeSnapshot) error {
	return nil
}

// Apply snapshot to instance
func (p *TencentLighthouseDriver) ApplySnapshot(node *compute.Node, snapshot *compute.VolumeSnapshot) error {
	return nil
}

// List all available images for instance
func (p *TencentLighthouseDriver) ListImages() ([]*compute.NodeImage, error) {
	return nil, nil
}

// Apply Image to instance
func (p *TencentLighthouseDriver) ApplyImage(node *compute.Node, image *compute.NodeImage) error {
	return nil
}

// List all available sizes for instance
func (p *TencentLighthouseDriver) ListSizes() ([]*compute.NodeSize, error) {
	return nil, nil
}

// Resize instance
func (p *TencentLighthouseDriver) ResizeNode(node *compute.Node, opts *compute.NodeResizeOpts) error {
	return nil
}

// List all available locations for instance
func (p *TencentLighthouseDriver) ListLocations() ([]*compute.Location, error) {
	return nil, nil
}
