package drivers

import (
	"github.com/rehiy/cloudgo/compute"
	"github.com/rehiy/cloudgo/provider"
	"github.com/rehiy/cloudgo/provider/tencent"

	cbs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

type TencentCvmDriver struct {
	client *tencent.Client
	cbs    *cbs.Client
	cvm    *cvm.Client
	rq     *provider.ReqeustParam
}

func NewTencentCvmDriver(rq *provider.ReqeustParam) *TencentCvmDriver {

	client := tencent.NewClient(rq)
	cbs, _ := client.Cbs()
	cvm, _ := client.Cvm()

	return &TencentCvmDriver{client, cbs, cvm, rq}

}

// List all instance
func (p *TencentCvmDriver) ListNodes() ([]*compute.Node, error) {

	resp, err := p.cvm.DescribeInstances(&cvm.DescribeInstancesRequest{})

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
			PublicIp:  *instance.PublicIpAddresses[0],
			PrivateIp: *instance.PrivateIpAddresses[0],
			Size: &compute.NodeSize{
				Id: *instance.InstanceType,
			},
			Image: &compute.NodeImage{
				Id:     *instance.ImageId,
				Name:   *instance.OsName,
				OSType: osType,
			},
			Location: &compute.Location{
				Id: *instance.Placement.Zone,
			},
		}

		nodes = append(nodes, node)

	}

	return nodes, nil

}

// Detail instance by Id
func (p *TencentCvmDriver) DetailNode(id string) (*compute.Node, error) {

	resp, err := p.cvm.DescribeInstances(&cvm.DescribeInstancesRequest{
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
		PublicIp:  *instance.PublicIpAddresses[0],
		PrivateIp: *instance.PrivateIpAddresses[0],
		Size: &compute.NodeSize{
			Id: *instance.InstanceType,
		},
		Image: &compute.NodeImage{
			Id:     *instance.ImageId,
			Name:   *instance.OsName,
			OSType: osType,
		},
		Location: &compute.Location{
			Id: *instance.Placement.Zone,
		},
	}

	return node, nil

}

// Create new instance
func (p *TencentCvmDriver) CreateNode(opts *compute.NodeCreateOpts) (*compute.Node, error) {

	// TODO: create instance
	resp, err := p.cvm.RunInstances(&cvm.RunInstancesRequest{
		// InstanceChargeType: &opts.InstanceChargeType,
		// InstanceType:       &opts.InstanceType,
		// Placement: &cvm.Placement{
		// 	Zone: &opts.Zone,
		// },
		// ImageId:       &opts.ImageId,
		// InstanceCount: &opts.InstanceCount,
		// InstanceName:  &opts.InstanceName,
		// LoginSettings: &cvm.LoginSettings{
		// 	Password: &opts.Password,
		// },
		// SecurityGroupIds: []*string{&opts.SecurityGroupId},
		// SystemDisk: &cvm.SystemDisk{
		// 	DiskType: &opts.DiskType,
		// 	DiskSize: &opts.DiskSize,
		// },
		// VirtualPrivateCloud: &cvm.VirtualPrivateCloud{
		// 	VpcId:    &opts.VpcId,
		// 	SubnetId: &opts.SubnetId,
		// },
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
func (p *TencentCvmDriver) DestroyNode(node *compute.Node) error {

	_, err := p.cvm.TerminateInstances(&cvm.TerminateInstancesRequest{
		InstanceIds: []*string{&node.Id},
	})

	return err

}

// Reboot instance
func (p *TencentCvmDriver) RebootNode(node *compute.Node) error {

	_, err := p.cvm.RebootInstances(&cvm.RebootInstancesRequest{
		InstanceIds: []*string{&node.Id},
	})

	return err

}

// Start instance
func (p *TencentCvmDriver) StartNode(node *compute.Node) error {

	_, err := p.cvm.StartInstances(&cvm.StartInstancesRequest{
		InstanceIds: []*string{&node.Id},
	})

	return err

}

// Stop instance
func (p *TencentCvmDriver) StopNode(node *compute.Node) error {

	_, err := p.cvm.StopInstances(&cvm.StopInstancesRequest{
		InstanceIds: []*string{&node.Id},
	})

	return err

}

// Get the current state of instance
func (p *TencentCvmDriver) GetNodeState(node *compute.Node) (compute.NodeState, error) {

	resp, err := p.cvm.DescribeInstancesStatus(&cvm.DescribeInstancesStatusRequest{
		InstanceIds: []*string{&node.Id},
	})

	if err != nil {
		return "", err
	}

	if len(resp.Response.InstanceStatusSet) == 0 {
		return "", nil
	}

	instanceStatus := resp.Response.InstanceStatusSet[0]

	// TODO: fix mapping
	state := compute.NodeState(*instanceStatus.InstanceState)

	return state, nil

}

// Get the Console url for instance
func (p *TencentCvmDriver) GetNodeConsole(node *compute.Node) (string, error) {

	resp, err := p.cvm.DescribeInstanceVncUrl(&cvm.DescribeInstanceVncUrlRequest{
		InstanceId: &node.Id,
	})

	if err != nil {
		return "", err
	}

	return *resp.Response.InstanceVncUrl, nil

}

// Get the public IP address of instance
func (p *TencentCvmDriver) GetNodePublicIp(node *compute.Node) (string, error) {

	resp, err := p.cvm.DescribeInstances(&cvm.DescribeInstancesRequest{
		InstanceIds: []*string{&node.Id},
	})

	if err != nil {
		return "", err
	}

	if len(resp.Response.InstanceSet) == 0 {
		return "", nil
	}

	instance := resp.Response.InstanceSet[0]
	ip := *instance.PublicIpAddresses[0]

	return ip, nil

}

// Get the private IP address of instance
func (p *TencentCvmDriver) GetNodePrivateIp(node *compute.Node) (string, error) {

	resp, err := p.cvm.DescribeInstances(&cvm.DescribeInstancesRequest{
		InstanceIds: []*string{&node.Id},
	})

	if err != nil {
		return "", err
	}

	if len(resp.Response.InstanceSet) == 0 {
		return "", nil
	}

	instance := resp.Response.InstanceSet[0]
	ip := *instance.PrivateIpAddresses[0]

	return ip, nil

}

// List all available storage volumes for instance
func (p *TencentCvmDriver) ListVolumes(node *compute.Node) ([]*compute.StorageVolume, error) {

	resp, err := p.cvm.DescribeInstances(&cvm.DescribeInstancesRequest{
		InstanceIds: []*string{&node.Id},
	})

	if err != nil {
		return nil, err
	}

	instance := resp.Response.InstanceSet[0]

	volumes := []*compute.StorageVolume{
		{
			Id:   *instance.SystemDisk.DiskId,
			Type: *instance.SystemDisk.DiskType,
			Size: int(*instance.SystemDisk.DiskSize),
		},
	}

	for _, disk := range instance.DataDisks {

		volumes = append(volumes, &compute.StorageVolume{
			Id:   *disk.DiskId,
			Type: *disk.DiskType,
			Size: int(*disk.DiskSize),
		})

	}

	return volumes, nil

}

// Attach volume to instance
func (p *TencentCvmDriver) AttachVolume(node *compute.Node, volume *compute.StorageVolume) error {

	_, err := p.cbs.AttachDisks(&cbs.AttachDisksRequest{
		InstanceId: &node.Id,
		DiskIds:    []*string{&volume.Id},
	})

	return err

}

// Detach volume from instance
func (p *TencentCvmDriver) DetachVolume(node *compute.Node, volume *compute.StorageVolume) error {

	_, err := p.cbs.DetachDisks(&cbs.DetachDisksRequest{
		InstanceId: &node.Id,
		DiskIds:    []*string{&volume.Id},
	})

	return err

}

// List all snapshots for instance
func (p *TencentCvmDriver) ListSnapshots(node *compute.Node) ([]*compute.VolumeSnapshot, error) {

	volumes, err := p.ListVolumes(node)

	if err != nil {
		return nil, err
	}

	volumeIds := []*string{}
	for _, volume := range volumes {
		volumeIds = append(volumeIds, &volume.Id)
	}

	filterName := "disk-id"

	resp, err := p.cbs.DescribeSnapshots(&cbs.DescribeSnapshotsRequest{
		Filters: []*cbs.Filter{
			{
				Name:   &filterName,
				Values: volumeIds,
			},
		},
	})

	if err != nil {
		return nil, err
	}

	snapshots := []*compute.VolumeSnapshot{}

	for _, snapshot := range resp.Response.SnapshotSet {

		snapshots = append(snapshots, &compute.VolumeSnapshot{
			Id:   *snapshot.SnapshotId,
			Name: *snapshot.SnapshotName,
			Size: int(*snapshot.DiskSize),
		})

	}

	return snapshots, nil

}

// Create snapshot for instance
func (p *TencentCvmDriver) CreateSnapshot(node *compute.Node, name string) (*compute.VolumeSnapshot, error) {

	resp, err := p.cbs.CreateSnapshot(&cbs.CreateSnapshotRequest{
		DiskId:       &node.Id,
		SnapshotName: &name,
	})

	if err != nil {
		return nil, err
	}

	snapshot := &compute.VolumeSnapshot{
		Id: *resp.Response.SnapshotId,
	}

	return snapshot, nil

}

// Destroy snapshot for instance
func (p *TencentCvmDriver) DestroySnapshot(node *compute.Node, snapshot *compute.VolumeSnapshot) error {

	_, err := p.cbs.DeleteSnapshots(&cbs.DeleteSnapshotsRequest{
		SnapshotIds: []*string{&snapshot.Id},
	})

	return err

}

// Apply snapshot to instance
func (p *TencentCvmDriver) ApplySnapshot(node *compute.Node, snapshot *compute.VolumeSnapshot) error {

	_, err := p.cbs.ApplySnapshot(&cbs.ApplySnapshotRequest{
		DiskId:     &snapshot.Id,
		SnapshotId: &snapshot.Id,
	})

	return err

}

// List all available images for instance
func (p *TencentCvmDriver) ListImages() ([]*compute.NodeImage, error) {

	resp, err := p.cvm.DescribeImages(&cvm.DescribeImagesRequest{
		ImageIds: []*string{},
	})

	if err != nil {
		return nil, err
	}

	images := []*compute.NodeImage{}

	for _, image := range resp.Response.ImageSet {

		images = append(images, &compute.NodeImage{
			Id:   *image.ImageId,
			Name: *image.ImageName,
			// Description: *image.ImageDescription,
			// Os:          *image.OsName,
			// OsVersion:   *image.OsVersion,
		})

	}

	return images, nil

}

// Apply Image to instance
func (p *TencentCvmDriver) ApplyImage(node *compute.Node, image *compute.NodeImage) error {

	_, err := p.cvm.ResetInstance(&cvm.ResetInstanceRequest{
		InstanceId: &node.Id,
		ImageId:    &image.Id,
	})

	return err

}

// List all available sizes for instance
func (p *TencentCvmDriver) ListSizes() ([]*compute.NodeSize, error) {

	resp, err := p.cvm.DescribeZoneInstanceConfigInfos(&cvm.DescribeZoneInstanceConfigInfosRequest{})

	if err != nil {
		return nil, err
	}

	sizes := []*compute.NodeSize{}

	for _, instanceType := range resp.Response.InstanceTypeQuotaSet {

		disk := instanceType.LocalDiskTypeList[0].MinSize
		arch := compute.Architecture(*instanceType.InstanceFamily)

		sizes = append(sizes, &compute.NodeSize{
			Id:           *instanceType.InstanceType,
			Name:         *instanceType.InstanceType,
			Architecture: arch,
			Gpu:          int(*instanceType.Gpu),
			Cpu:          int(*instanceType.Cpu),
			Ram:          int(*instanceType.Memory),
			Disk:         int(*disk),
		})

	}

	return sizes, nil

}

// Resize instance
func (p *TencentCvmDriver) ResizeNode(node *compute.Node, opts *compute.NodeResizeOpts) error {

	_, err := p.cvm.ResetInstancesType(&cvm.ResetInstancesTypeRequest{
		InstanceIds:  []*string{&node.Id},
		InstanceType: &opts.Size.Id,
	})
	return err

}

// List all instance
func (p *TencentCvmDriver) ListLocations() ([]*compute.Location, error) {

	resp, err := p.cvm.DescribeZones(&cvm.DescribeZonesRequest{})

	if err != nil {
		return nil, err
	}

	var locations []*compute.Location

	for _, zone := range resp.Response.ZoneSet {

		location := &compute.Location{
			Id:   *zone.Zone,
			Name: *zone.ZoneName,
		}

		locations = append(locations, location)

	}

	return locations, nil

}
