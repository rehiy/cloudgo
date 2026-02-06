package drivers

import (
	"github.com/rehiy/cloudgo/compute"
	"github.com/rehiy/cloudgo/provider"
	"github.com/rehiy/cloudgo/provider/alibaba"

	ecs "github.com/alibabacloud-go/ecs-20140526/v3/client"
	"github.com/alibabacloud-go/tea/tea"
)

type AlibabaEcsDriver struct {
	client *alibaba.Client
	ecs    *ecs.Client
	rq     *provider.ReqeustParam
}

func NewAlibabaEcsDriver(rq *provider.ReqeustParam) *AlibabaEcsDriver {

	client := alibaba.NewClient(rq)
	ecs, _ := client.Ecs()

	return &AlibabaEcsDriver{client, ecs, rq}

}

// List all instance
func (p *AlibabaEcsDriver) ListNodes() ([]*compute.Node, error) {

	resp, err := p.ecs.DescribeInstances(&ecs.DescribeInstancesRequest{})

	if err != nil {
		return nil, err
	}

	var nodes []*compute.Node

	for _, instance := range resp.Body.Instances.Instance {

		// TODO: fix mapping
		state := compute.NodeState(*instance.Status)
		osType := compute.OSType(*instance.OSType)

		node := &compute.Node{
			Id:        *instance.InstanceId,
			Name:      *instance.InstanceName,
			State:     state,
			PublicIp:  *instance.PublicIpAddress.IpAddress[0],
			PrivateIp: *instance.VpcAttributes.PrivateIpAddress.IpAddress[0],
			Size: &compute.NodeSize{
				Id: *instance.InstanceType,
			},
			Image: &compute.NodeImage{
				Id:     *instance.ImageId,
				Name:   *instance.OSName,
				OSType: osType,
			},
			Location: &compute.Location{
				Id: *instance.ZoneId,
			},
		}

		nodes = append(nodes, node)

	}

	return nodes, nil

}

// Detail instance by Id
func (p *AlibabaEcsDriver) DetailNode(id string) (*compute.Node, error) {

	resp, err := p.ecs.DescribeInstances(&ecs.DescribeInstancesRequest{
		InstanceIds: tea.String(`["` + id + `"]`),
	})

	if err != nil {
		return nil, err
	}

	instance := resp.Body.Instances.Instance[0]

	// TODO: fix mapping
	state := compute.NodeState(*instance.Status)
	osType := compute.OSType(*instance.OSType)

	node := &compute.Node{
		Id:        *instance.InstanceId,
		Name:      *instance.InstanceName,
		State:     state,
		PublicIp:  *instance.PublicIpAddress.IpAddress[0],
		PrivateIp: *instance.VpcAttributes.PrivateIpAddress.IpAddress[0],
		Size: &compute.NodeSize{
			Id: *instance.InstanceType,
		},
		Image: &compute.NodeImage{
			Id:     *instance.ImageId,
			Name:   *instance.OSName,
			OSType: osType,
		},
		Location: &compute.Location{
			Id: *instance.ZoneId,
		},
	}

	return node, nil

}

// Create new instance
func (p *AlibabaEcsDriver) CreateNode(opts *compute.NodeCreateOpts) (*compute.Node, error) {

	resp, err := p.ecs.CreateInstance(&ecs.CreateInstanceRequest{
		InstanceName: tea.String(opts.Name),
		InstanceType: tea.String(opts.Size.Id),
		ImageId:      tea.String(opts.Image.Id),
		ZoneId:       tea.String(opts.Location.Id),
	})

	if err != nil {
		return nil, err
	}

	instance := resp.Body.InstanceId

	node := &compute.Node{
		Id:   *instance,
		Name: opts.Name,
		Size: opts.Size,
	}

	return node, nil

}

// Destroy an existing instance
func (p *AlibabaEcsDriver) DestroyNode(node *compute.Node) error {

	_, err := p.ecs.DeleteInstance(&ecs.DeleteInstanceRequest{
		InstanceId: tea.String(node.Id),
	})

	return err

}

// Reboot instance
func (p *AlibabaEcsDriver) RebootNode(node *compute.Node) error {

	_, err := p.ecs.RebootInstance(&ecs.RebootInstanceRequest{
		InstanceId: tea.String(node.Id),
	})

	return err

}

// Start instance
func (p *AlibabaEcsDriver) StartNode(node *compute.Node) error {

	_, err := p.ecs.StartInstance(&ecs.StartInstanceRequest{
		InstanceId: tea.String(node.Id),
	})

	return err

}

// Stop instance
func (p *AlibabaEcsDriver) StopNode(node *compute.Node) error {

	_, err := p.ecs.StopInstance(&ecs.StopInstanceRequest{
		InstanceId: tea.String(node.Id),
	})

	return err

}

// Get the current state of instance
func (p *AlibabaEcsDriver) GetNodeState(node *compute.Node) (compute.NodeState, error) {

	resp, err := p.ecs.DescribeInstances(&ecs.DescribeInstancesRequest{
		InstanceIds: tea.String(`["` + node.Id + `"]`),
	})

	if err != nil {
		return "", err
	}

	instance := resp.Body.Instances.Instance[0]

	// TODO: fix mapping
	state := compute.NodeState(*instance.Status)

	return state, nil

}

// Get the Console url for instance
func (p *AlibabaEcsDriver) GetNodeConsole(node *compute.Node) (string, error) {

	resp, err := p.ecs.DescribeInstanceVncUrl(&ecs.DescribeInstanceVncUrlRequest{
		InstanceId: tea.String(node.Id),
	})

	if err != nil {
		return "", err
	}

	return *resp.Body.VncUrl, nil

}

// Get the public IP address of instance
func (p *AlibabaEcsDriver) GetNodePublicIp(node *compute.Node) (string, error) {

	resp, err := p.ecs.DescribeInstances(&ecs.DescribeInstancesRequest{
		InstanceIds: tea.String(`["` + node.Id + `"]`),
	})

	if err != nil {
		return "", err
	}

	instance := resp.Body.Instances.Instance[0]
	ip := *instance.PublicIpAddress.IpAddress[0]

	return ip, nil

}

// Get the private IP address of instance
func (p *AlibabaEcsDriver) GetNodePrivateIp(node *compute.Node) (string, error) {

	resp, err := p.ecs.DescribeInstances(&ecs.DescribeInstancesRequest{
		InstanceIds: tea.String(`["` + node.Id + `"]`),
	})

	if err != nil {
		return "", err
	}

	instance := resp.Body.Instances.Instance[0]
	ip := *instance.VpcAttributes.PrivateIpAddress.IpAddress[0]

	return ip, nil

}

// List all available storage volumes for instance
func (p *AlibabaEcsDriver) ListVolumes(node *compute.Node) ([]*compute.StorageVolume, error) {

	resp, err := p.ecs.DescribeDisks(&ecs.DescribeDisksRequest{
		InstanceId: tea.String(node.Id),
	})

	if err != nil {
		return nil, err
	}

	var volumes []*compute.StorageVolume

	for _, disk := range resp.Body.Disks.Disk {

		volume := &compute.StorageVolume{
			Id:   *disk.DiskId,
			Name: *disk.DiskName,
			Size: int(*disk.Size),
		}

		volumes = append(volumes, volume)

	}

	return volumes, nil

}

// Attach volume to instance
func (p *AlibabaEcsDriver) AttachVolume(node *compute.Node, volume *compute.StorageVolume) error {

	_, err := p.ecs.AttachDisk(&ecs.AttachDiskRequest{
		InstanceId: tea.String(node.Id),
		DiskId:     tea.String(volume.Id),
	})

	return err

}

// Detach volume from instance
func (p *AlibabaEcsDriver) DetachVolume(node *compute.Node, volume *compute.StorageVolume) error {

	_, err := p.ecs.DetachDisk(&ecs.DetachDiskRequest{
		InstanceId: tea.String(node.Id),
		DiskId:     tea.String(volume.Id),
	})

	return err

}

// List all snapshots for instance
func (p *AlibabaEcsDriver) ListSnapshots(node *compute.Node) ([]*compute.VolumeSnapshot, error) {

	resp, err := p.ecs.DescribeDisks(&ecs.DescribeDisksRequest{
		InstanceId: tea.String(node.Id),
	})

	if err != nil {
		return nil, err
	}

	var snapshots []*compute.VolumeSnapshot

	for _, disk := range resp.Body.Disks.Disk {

		snapshot := &compute.VolumeSnapshot{
			Id:   *disk.DiskId,
			Name: *disk.DiskName,
			Size: int(*disk.Size),
		}

		snapshots = append(snapshots, snapshot)

	}

	return snapshots, nil

}

// Create snapshot for instance
func (p *AlibabaEcsDriver) CreateSnapshot(node *compute.Node, name string) (*compute.VolumeSnapshot, error) {

	resp, err := p.ecs.CreateSnapshot(&ecs.CreateSnapshotRequest{
		DiskId:       tea.String(node.Id),
		SnapshotName: tea.String(name),
	})

	if err != nil {
		return nil, err
	}

	snapshot := &compute.VolumeSnapshot{
		Id:    *resp.Body.SnapshotId,
		State: compute.StorageVolumeStateCREATING,
		Name:  name,
	}

	return snapshot, nil

}

// Destroy snapshot for instance
func (p *AlibabaEcsDriver) DestroySnapshot(node *compute.Node, snapshot *compute.VolumeSnapshot) error {

	_, err := p.ecs.DeleteSnapshot(&ecs.DeleteSnapshotRequest{
		SnapshotId: tea.String(snapshot.Id),
	})

	return err

}

// Apply snapshot to instance
func (p *AlibabaEcsDriver) ApplySnapshot(node *compute.Node, snapshot *compute.VolumeSnapshot) error {

	_, err := p.ecs.DescribeDisks(&ecs.DescribeDisksRequest{
		InstanceId: tea.String(node.Id),
	})

	return err

}

// List all available images for instance
func (p *AlibabaEcsDriver) ListImages() ([]*compute.NodeImage, error) {

	resp, err := p.ecs.DescribeImages(&ecs.DescribeImagesRequest{
		ImageOwnerAlias: tea.String("self"),
	})

	if err != nil {
		return nil, err
	}

	var images []*compute.NodeImage

	for _, image := range resp.Body.Images.Image {

		images = append(images, &compute.NodeImage{
			Id:   *image.ImageId,
			Name: *image.ImageName,
		})

	}

	return images, nil

}

// Apply Image to instance
func (p *AlibabaEcsDriver) ApplyImage(node *compute.Node, image *compute.NodeImage) error {

	return nil //TODO: implement

}

// List all available sizes for instance
func (p *AlibabaEcsDriver) ListSizes() ([]*compute.NodeSize, error) {

	resp, err := p.ecs.DescribeInstanceTypes(&ecs.DescribeInstanceTypesRequest{})

	if err != nil {
		return nil, err
	}

	var sizes []*compute.NodeSize

	for _, instanceType := range resp.Body.InstanceTypes.InstanceType {

		arch := compute.Architecture(*instanceType.CpuArchitecture)

		sizes = append(sizes, &compute.NodeSize{
			Id:           *instanceType.InstanceTypeId,
			Name:         *instanceType.InstanceTypeFamily,
			Architecture: arch,
			Gpu:          int(*instanceType.GPUAmount),
			Cpu:          int(*instanceType.CpuCoreCount),
			Ram:          int(*instanceType.MemorySize),
			Disk:         int(*instanceType.DiskQuantity),
		})

	}

	return sizes, nil

}

// Resize instance
func (p *AlibabaEcsDriver) ResizeNode(node *compute.Node, opts *compute.NodeResizeOpts) error {

	_, err := p.ecs.ModifyInstanceSpec(&ecs.ModifyInstanceSpecRequest{
		InstanceId:   tea.String(node.Id),
		InstanceType: tea.String(opts.Size.Id),
	})

	return err

}

// List all available locations for instance
func (p *AlibabaEcsDriver) ListLocations() ([]*compute.Location, error) {

	resp, err := p.ecs.DescribeRegions(&ecs.DescribeRegionsRequest{})

	if err != nil {
		return nil, err
	}

	var locations []*compute.Location

	for _, region := range resp.Body.Regions.Region {

		locations = append(locations, &compute.Location{
			Id:   *region.RegionId,
			Name: *region.LocalName,
		})

	}

	return locations, nil

}
