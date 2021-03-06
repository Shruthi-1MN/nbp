// Copyright (c) 2018 Huawei Technologies Co., Ltd. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package opensds

import (
	"fmt"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"

	csi "github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/golang/glog"
	sdscontroller "github.com/opensds/nbp/client/opensds"
	"github.com/opensds/nbp/csi/util"
	c "github.com/opensds/opensds/client"
	"github.com/opensds/opensds/contrib/connector"
	"github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/pkg/utils"
	"github.com/opensds/opensds/pkg/utils/constants"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

////////////////////////////////////////////////////////////////////////////////
//                            Controller Service                              //
////////////////////////////////////////////////////////////////////////////////

var (
	// Client opensds client
	Client *c.Client
)

func init() {
	Client = sdscontroller.GetClient("", "")
}

// GetDefaultProfile implementation
func GetDefaultProfile() (*model.ProfileSpec, error) {
	profiles, err := Client.ListProfiles()
	if err != nil {
		glog.Error("Get default profile failed: ", err)
		return nil, err
	}

	for _, profile := range profiles {
		if profile.Name == "default" {
			return profile, nil
		}
	}

	return nil, status.Error(codes.FailedPrecondition, "No default profile")
}

// FindVolume implementation
func FindVolume(req *model.VolumeSpec) (bool, bool, *model.VolumeSpec, error) {
	isExist := false
	volumes, err := Client.ListVolumes()

	if err != nil {
		glog.Error("List volumes failed: ", err)

		return false, false, nil, err
	}

	for _, volume := range volumes {
		if volume.Name == req.Name {
			isExist = true

			if (volume.Size == req.Size) && (volume.ProfileId == req.ProfileId) &&
				(volume.AvailabilityZone == req.AvailabilityZone) &&
				(volume.SnapshotId == req.SnapshotId) {
				glog.V(5).Infof("Volume already exists and is compatible")

				return true, true, volume, nil
			}
		}
	}

	return isExist, false, nil, nil
}

// CreateVolume implementation
func (p *Plugin) CreateVolume(
	ctx context.Context,
	req *csi.CreateVolumeRequest) (
	*csi.CreateVolumeResponse, error) {

	glog.V(5).Info("start to CreateVolume")
	defer glog.V(5).Info("end to CreateVolume")

	// build volume body
	volumebody := &model.VolumeSpec{}
	volumebody.Name = req.Name
	allocationUnitBytes := util.GiB
	if req.CapacityRange != nil {
		volumeSizeBytes := int64(req.CapacityRange.RequiredBytes)
		volumebody.Size = (volumeSizeBytes + allocationUnitBytes - 1) / allocationUnitBytes
		if volumebody.Size < 1 {
			//Using default volume size
			volumebody.Size = 1
		}
	} else {
		//Using default volume size
		volumebody.Size = 1
	}
	var secondaryAZ = util.OpensdsDefaultSecondaryAZ
	var enableReplication = false
	for k, v := range req.GetParameters() {
		switch strings.ToLower(k) {
		case KParamProfile:
			volumebody.ProfileId = v
		case KParamAZ:
			volumebody.AvailabilityZone = v
		case KParamEnableReplication:
			if strings.ToLower(v) == "true" {
				enableReplication = true
			}
		case KParamSecondaryAZ:
			secondaryAZ = v
		}
	}

	contentSource := req.GetVolumeContentSource()
	if nil != contentSource {
		snapshot := contentSource.GetSnapshot()
		if snapshot != nil {
			volumebody.SnapshotId = snapshot.GetSnapshotId()
		}
	}

	if "" == volumebody.ProfileId {
		defaultRrf, err := GetDefaultProfile()
		if err != nil {
			return nil, err
		}

		volumebody.ProfileId = defaultRrf.Id
	}

	if "" == volumebody.AvailabilityZone {
		volumebody.AvailabilityZone = "default"
	}

	glog.V(5).Infof("CreateVolume volumebody: %v", volumebody)

	isExist, isCompatible, findVolume, err := FindVolume(volumebody)
	if err != nil {
		return nil, err
	}

	var v *model.VolumeSpec

	if isExist {
		if isCompatible {
			v = findVolume
		} else {
			return nil, status.Error(codes.AlreadyExists,
				"Volume already exists but is incompatible")
		}
	} else {
		createVolume, err := Client.CreateVolume(volumebody)
		if err != nil {
			isExist, isCompatible, findV, findErr := FindVolume(volumebody)
			if findErr != nil {
				return nil, findErr
			}

			if !(isExist && isCompatible) {
				glog.Error("failed to CreateVolume", err)
				return nil, err
			}

			v = findV
			glog.V(5).Infof("Although the return failed, it was actually successful. volume = %v", findV)
		} else {
			v = createVolume
		}
	}

	glog.V(5).Infof("opensds volume = %v", v)
	// return volume info
	volumeinfo := &csi.Volume{
		CapacityBytes: v.Size * allocationUnitBytes,
		VolumeId:      v.Id,
		VolumeContext: map[string]string{
			KVolumeName:      v.Name,
			KVolumeStatus:    v.Status,
			KVolumeAZ:        v.AvailabilityZone,
			KVolumePoolId:    v.PoolId,
			KVolumeProfileId: v.ProfileId,
			KVolumeLvPath:    v.Metadata["lvPath"],
		},
	}

	glog.V(5).Infof("resp volumeinfo = %v", volumeinfo)
	if enableReplication && !isExist {
		volumebody.AvailabilityZone = secondaryAZ
		volumebody.Name = SecondaryPrefix + req.Name
		sVol, err := Client.CreateVolume(volumebody)
		if err != nil {
			glog.Errorf("failed to create secondar volume: %v", err)
			return nil, err
		}
		replicaBody := &model.ReplicationSpec{
			Name:              req.Name,
			PrimaryVolumeId:   v.Id,
			SecondaryVolumeId: sVol.Id,
			ReplicationMode:   model.ReplicationModeSync,
			ReplicationPeriod: 0,
		}
		replicaResp, err := Client.CreateReplication(replicaBody)
		if err != nil {
			glog.Errorf("Create replication failed: %v", err)
			return nil, err
		}
		volumeinfo.VolumeContext[KVolumeReplicationId] = replicaResp.Id
	}

	return &csi.CreateVolumeResponse{
		Volume: volumeinfo,
	}, nil
}

func getReplicationByVolume(volId string) *model.ReplicationSpec {
	replications, _ := Client.ListReplications()
	for _, r := range replications {
		if volId == r.PrimaryVolumeId || volId == r.SecondaryVolumeId {
			return r
		}
	}
	return nil
}

// DeleteVolume implementation
func (p *Plugin) DeleteVolume(
	ctx context.Context,
	req *csi.DeleteVolumeRequest) (
	*csi.DeleteVolumeResponse, error) {
	glog.V(5).Info("start to DeleteVolume")
	defer glog.V(5).Info("end to DeleteVolume")
	volId := req.VolumeId

	r := getReplicationByVolume(volId)
	if r != nil {
		if err := Client.DeleteReplication(r.Id, nil); err != nil {
			return nil, err
		}
		if err := Client.DeleteVolume(r.PrimaryVolumeId, &model.VolumeSpec{}); err != nil {
			return nil, err
		}
		if err := Client.DeleteVolume(r.SecondaryVolumeId, &model.VolumeSpec{}); err != nil {
			return nil, err
		}
	} else {
		if err := Client.DeleteVolume(volId, &model.VolumeSpec{}); err != nil {
			return nil, err
		}
	}

	return &csi.DeleteVolumeResponse{}, nil
}

// isStringMapEqual implementation
func isStringMapEqual(metadataA, metadataB map[string]string) bool {
	glog.V(5).Infof("start to isStringMapEqual, metadataA = %v, metadataB = %v!",
		metadataA, metadataB)
	if len(metadataA) != len(metadataB) {
		glog.V(5).Infof("len(metadataA)(%v) != len(metadataB)(%v) ",
			len(metadataA), len(metadataB))
		return false
	}

	for key, valueA := range metadataA {
		valueB, ok := metadataB[key]
		if !ok || (valueA != valueB) {
			glog.V(5).Infof("ok = %v, key = %v, valueA = %v, valueB = %v!",
				ok, key, valueA, valueB)
			return false
		}
	}

	return true
}

// isVolumePublished Check if the volume is published and compatible
func isVolumePublished(canAtMultiNode bool, attachReq *model.VolumeAttachmentSpec,
	metadata map[string]string) (*model.VolumeAttachmentSpec, error) {
	glog.V(5).Infof("start to isVolumePublished, canAtMultiNode = %v, attachReq = %v",
		canAtMultiNode, attachReq)

	attachments, err := Client.ListVolumeAttachments()
	if err != nil {
		glog.V(5).Info("ListVolumeAttachments failed: " + err.Error())
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	for _, attachSpec := range attachments {
		if attachSpec.VolumeId == attachReq.VolumeId {
			if attachSpec.Host != attachReq.Host {
				if !canAtMultiNode {
					msg := fmt.Sprintf("the volume %s has been published to another node and does not have MULTI_NODE volume capability",
						attachReq.VolumeId)
					return nil, status.Error(codes.FailedPrecondition, msg)
				}
			} else {
				// Opensds does not have volume_capability and readonly parameters,
				// but needs to check other parameters to determine compatibility?
				if attachSpec.Platform == attachReq.Platform &&
					attachSpec.OsType == attachReq.OsType &&
					attachSpec.Initiator == attachReq.Initiator &&
					isStringMapEqual(attachSpec.Metadata, metadata) &&
					attachSpec.AccessProtocol == attachReq.AccessProtocol {
					glog.V(5).Info("Volume published and is compatible")

					return attachSpec, nil
				}

				glog.Error("Volume published but is incompatible, incompatible attachement Id = " + attachSpec.Id)
				return nil, status.Error(codes.AlreadyExists, "Volume published but is incompatible")
			}
		}
	}

	glog.V(5).Info("Need to create a new attachment")
	return nil, nil
}

// ControllerPublishVolume implementation
func (p *Plugin) ControllerPublishVolume(
	ctx context.Context,
	req *csi.ControllerPublishVolumeRequest) (
	*csi.ControllerPublishVolumeResponse, error) {

	glog.V(5).Info("start to ControllerPublishVolume")
	defer glog.V(5).Info("end to ControllerPublishVolume")

	//check volume is exist
	volSpec, errVol := Client.GetVolume(req.VolumeId)
	if errVol != nil || volSpec == nil {
		msg := fmt.Sprintf("the volume %s is not exist", req.VolumeId)
		return nil, status.Error(codes.NotFound, msg)
	}

	pool, err := Client.GetPool(volSpec.PoolId)
	if err != nil || pool == nil {
		msg := fmt.Sprintf("the pool %s is not sxist", volSpec.PoolId)
		glog.Error(msg)
		return nil, status.Error(codes.NotFound, msg)
	}

	var protocol = strings.ToLower(pool.Extras.IOConnectivity.AccessProtocol)
	if protocol == "" {
		// Default protocol is iscsi
		protocol = "iscsi"
	}

	var initator string
	hostName, wwpns, _, iqns := extractInfoFromNodeId(req.NodeId)

	switch protocol {
	case connector.FcDriver:
		if len(wwpns) <= 0 {
			msg := fmt.Sprintf("protocol is %v, but no wwpn", protocol)
			glog.Error(msg)
			return nil, status.Error(codes.FailedPrecondition, msg)
		}

		initator = strings.Join(wwpns, ",")
		break
	case connector.IscsiDriver:
		if len(iqns) <= 0 {
			msg := fmt.Sprintf("protocol is %v, but no iqn", protocol)
			glog.Error(msg)
			return nil, status.Error(codes.FailedPrecondition, msg)
		}

		initator = iqns[0]
		break
	case connector.RbdDriver:
		break
	default:
		msg := fmt.Sprintf("protocol cannot be %v", protocol)
		glog.Error(msg)
		return nil, status.Error(codes.InvalidArgument, msg)
	}

	attachReq := &model.VolumeAttachmentSpec{
		VolumeId: req.VolumeId,
		HostInfo: model.HostInfo{
			Host:      hostName,
			Platform:  runtime.GOARCH,
			OsType:    runtime.GOOS,
			Initiator: initator,
		},
		Metadata:       req.VolumeContext,
		AccessProtocol: protocol,
	}

	mode := req.VolumeCapability.AccessMode.Mode
	canAtMultiNode := false

	if csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER == mode ||
		csi.VolumeCapability_AccessMode_MULTI_NODE_READER_ONLY == mode ||
		csi.VolumeCapability_AccessMode_MULTI_NODE_SINGLE_WRITER == mode {
		canAtMultiNode = true
	}

	expectedMetadata := utils.MergeStringMaps(attachReq.Metadata, volSpec.Metadata)
	existAttachment, err := isVolumePublished(canAtMultiNode, attachReq, expectedMetadata)
	if err != nil {
		return nil, err
	}

	var attachSpec *model.VolumeAttachmentSpec

	if nil == existAttachment {
		newAttachment, errAttach := Client.CreateVolumeAttachment(attachReq)
		if errAttach != nil {
			msg := fmt.Sprintf("the volume %s failed to publish to node %s.", req.VolumeId, req.NodeId)
			glog.Errorf("failed to ControllerPublishVolume: %v", attachReq)
			return nil, status.Error(codes.FailedPrecondition, msg)
		}

		attachSpec = newAttachment
	} else {
		attachSpec = existAttachment
	}

	resp := &csi.ControllerPublishVolumeResponse{
		PublishContext: map[string]string{
			KPublishHostIp:       attachSpec.Ip,
			KPublishHostName:     attachSpec.Host,
			KPublishAttachId:     attachSpec.Id,
			KPublishAttachStatus: attachSpec.Status,
		},
	}

	if replicationId, ok := req.VolumeContext[KVolumeReplicationId]; ok {
		r, err := Client.GetReplication(replicationId)
		if err != nil {
			return nil, status.Error(codes.FailedPrecondition, "Get replication failed")
		}

		attachReq.VolumeId = r.SecondaryVolumeId
		existAttachment, err := isVolumePublished(canAtMultiNode, attachReq, expectedMetadata)
		if err != nil {
			return nil, err
		}

		if nil == existAttachment {
			newAttachment, errAttach := Client.CreateVolumeAttachment(attachReq)
			if errAttach != nil {
				msg := fmt.Sprintf("the volume %s failed to publish to node %s.", req.VolumeId, req.NodeId)
				glog.Errorf("failed to ControllerPublishVolume: %v", attachReq)
				return nil, status.Error(codes.FailedPrecondition, msg)
			}

			attachSpec = newAttachment
		} else {
			attachSpec = existAttachment
		}

		resp.PublishContext[KPublishSecondaryAttachId] = attachSpec.Id
	}
	return resp, nil
}

func extractInfoFromNodeId(nodeId string) (string, []string, []string, []string) {
	var hostName string
	var wwpns []string
	var wwnns []string
	var iqns []string

	glog.V(5).Info("nodeId: " + nodeId)
	hostNameAndInitor := strings.Split(nodeId, ",")
	hostNameAndInitorLen := len(hostNameAndInitor)

	if hostNameAndInitorLen >= 1 {
		hostName = hostNameAndInitor[0]
	}

	var previousParameter string
	previousIndex := 0

	for i := 1; i < hostNameAndInitorLen; i++ {
		if strings.HasPrefix(hostNameAndInitor[i], connector.Wwpn+":") {
			wwpns = append(wwpns, strings.Split(hostNameAndInitor[i], connector.Wwpn+":")[1])
			previousParameter = connector.Wwpn
			previousIndex = len(wwpns) - 1
		} else {
			if strings.HasPrefix(hostNameAndInitor[i], connector.Wwnn+":") {
				wwnns = append(wwnns, strings.Split(hostNameAndInitor[i], connector.Wwnn+":")[1])
				previousParameter = connector.Wwnn
				previousIndex = len(wwnns) - 1
			} else {
				if strings.HasPrefix(hostNameAndInitor[i], connector.Iqn+":") {
					iqns = append(iqns, strings.Split(hostNameAndInitor[i], connector.Iqn+":")[1])
					previousParameter = connector.Iqn
					previousIndex = len(iqns) - 1
				} else {
					switch previousParameter {
					case connector.Wwpn:
						wwpns[previousIndex] = wwpns[previousIndex] + "," + hostNameAndInitor[i]
						break
					case connector.Wwnn:
						wwnns[previousIndex] = wwnns[previousIndex] + "," + hostNameAndInitor[i]
						break
					case connector.Iqn:
						iqns[previousIndex] = iqns[previousIndex] + "," + hostNameAndInitor[i]
						break
					default:
						glog.Error("The format of nodeId is incorrect")

					}
				}

			}
		}
	}

	return hostName, wwpns, wwnns, iqns
}

// ControllerUnpublishVolume implementation
func (p *Plugin) ControllerUnpublishVolume(
	ctx context.Context,
	req *csi.ControllerUnpublishVolumeRequest) (
	*csi.ControllerUnpublishVolumeResponse, error) {

	glog.V(5).Infof("start to ControllerUnpublishVolume, req VolumeId = %v, NodeId = %v, ControllerUnpublishSecrets =%v",
		req.VolumeId, req.NodeId, req.Secrets)
	defer glog.V(5).Info("end to ControllerUnpublishVolume")

	//check volume is exist
	volSpec, errVol := Client.GetVolume(req.VolumeId)
	if errVol != nil || volSpec == nil {
		msg := fmt.Sprintf("the volume %s is not exist", req.VolumeId)
		return nil, status.Error(codes.NotFound, msg)
	}

	attachments, err := Client.ListVolumeAttachments()
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, "Failed to unpublish volume.")
	}

	hostName, _, _, _ := extractInfoFromNodeId(req.NodeId)
	var acts []*model.VolumeAttachmentSpec

	for _, attachSpec := range attachments {
		if attachSpec.VolumeId == req.VolumeId && (req.NodeId == "" || attachSpec.Host == hostName) {
			acts = append(acts, attachSpec)
		}
	}

	if r := getReplicationByVolume(req.VolumeId); r != nil {
		for _, attachSpec := range attachments {
			if attachSpec.VolumeId == r.SecondaryVolumeId && (req.NodeId == "" || attachSpec.Host == hostName) {
				acts = append(acts, attachSpec)
			}
		}
	}

	for _, act := range acts {
		err = Client.DeleteVolumeAttachment(act.Id, act)
		if err != nil {
			msg := fmt.Sprintf("the volume %s failed to unpublish from node %s.", req.VolumeId, req.NodeId)
			glog.Errorf("failed to ControllerUnpublishVolume: %v", err)
			return nil, status.Error(codes.FailedPrecondition, msg)
		}

		glog.V(5).Infof("attachment %v has been successfully deleted", act.Id)
	}

	return &csi.ControllerUnpublishVolumeResponse{}, nil
}

// ValidateVolumeCapabilities implementation
func (p *Plugin) ValidateVolumeCapabilities(
	ctx context.Context,
	req *csi.ValidateVolumeCapabilitiesRequest) (
	*csi.ValidateVolumeCapabilitiesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// ListVolumes implementation
func (p *Plugin) ListVolumes(
	ctx context.Context,
	req *csi.ListVolumesRequest) (
	*csi.ListVolumesResponse, error) {

	glog.V(5).Info("start to ListVolumes")
	defer glog.V(5).Info("end to ListVolumes")

	// only support list all the volumes at present
	volumes, err := Client.ListVolumes()
	if err != nil {
		return nil, err
	}

	ens := []*csi.ListVolumesResponse_Entry{}
	for _, v := range volumes {
		if v != nil {

			volumeinfo := &csi.Volume{
				CapacityBytes: v.Size,
				VolumeId:      v.Id,
				VolumeContext: map[string]string{
					"Name":             v.Name,
					"Status":           v.Status,
					"AvailabilityZone": v.AvailabilityZone,
					"PoolId":           v.PoolId,
					"ProfileId":        v.ProfileId,
				},
			}

			ens = append(ens, &csi.ListVolumesResponse_Entry{
				Volume: volumeinfo,
			})
		}
	}

	return &csi.ListVolumesResponse{
		Entries: ens,
	}, nil
}

// GetCapacity implementation
func (p *Plugin) GetCapacity(
	ctx context.Context,
	req *csi.GetCapacityRequest) (
	*csi.GetCapacityResponse, error) {

	glog.V(5).Info("start to GetCapacity")
	defer glog.V(5).Info("end to GetCapacity")

	pools, err := Client.ListPools()
	if err != nil {
		return nil, err
	}

	// calculate all the free capacity of pools
	freecapacity := int64(0)
	for _, p := range pools {
		if p != nil {
			freecapacity += int64(p.FreeCapacity)
		}
	}

	return &csi.GetCapacityResponse{
		AvailableCapacity: freecapacity,
	}, nil
}

// ControllerGetCapabilities implementation
func (p *Plugin) ControllerGetCapabilities(
	ctx context.Context,
	req *csi.ControllerGetCapabilitiesRequest) (
	*csi.ControllerGetCapabilitiesResponse, error) {

	glog.V(5).Info("start to ControllerGetCapabilities")
	defer glog.V(5).Info("end to ControllerGetCapabilities")

	return &csi.ControllerGetCapabilitiesResponse{
		Capabilities: []*csi.ControllerServiceCapability{
			&csi.ControllerServiceCapability{
				Type: &csi.ControllerServiceCapability_Rpc{
					Rpc: &csi.ControllerServiceCapability_RPC{
						Type: csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
					},
				},
			},
			&csi.ControllerServiceCapability{
				Type: &csi.ControllerServiceCapability_Rpc{
					Rpc: &csi.ControllerServiceCapability_RPC{
						Type: csi.ControllerServiceCapability_RPC_PUBLISH_UNPUBLISH_VOLUME,
					},
				},
			},
			&csi.ControllerServiceCapability{
				Type: &csi.ControllerServiceCapability_Rpc{
					Rpc: &csi.ControllerServiceCapability_RPC{
						Type: csi.ControllerServiceCapability_RPC_LIST_VOLUMES,
					},
				},
			},
			&csi.ControllerServiceCapability{
				Type: &csi.ControllerServiceCapability_Rpc{
					Rpc: &csi.ControllerServiceCapability_RPC{
						Type: csi.ControllerServiceCapability_RPC_GET_CAPACITY,
					},
				},
			},
			&csi.ControllerServiceCapability{
				Type: &csi.ControllerServiceCapability_Rpc{
					Rpc: &csi.ControllerServiceCapability_RPC{
						Type: csi.ControllerServiceCapability_RPC_CREATE_DELETE_SNAPSHOT,
					},
				},
			},
			&csi.ControllerServiceCapability{
				Type: &csi.ControllerServiceCapability_Rpc{
					Rpc: &csi.ControllerServiceCapability_RPC{
						Type: csi.ControllerServiceCapability_RPC_LIST_SNAPSHOTS,
					},
				},
			},
		},
	}, nil
}

// FindSnapshot implementation
func FindSnapshot(req *model.VolumeSnapshotSpec) (bool, bool, *model.VolumeSnapshotSpec, error) {
	isExist := false
	snapshots, err := Client.ListVolumeSnapshots()

	if err != nil {
		glog.Error("List volume snapshots failed: ", err)

		return false, false, nil, err
	}

	for _, snapshot := range snapshots {
		if snapshot.Name == req.Name {
			isExist = true

			if (snapshot.VolumeId == req.VolumeId) && (snapshot.ProfileId == req.ProfileId) {
				glog.V(5).Infof("snapshot already exists and is compatible")

				return true, true, snapshot, nil
			}
		}
	}

	return isExist, false, nil, nil
}

// CreateSnapshot implementation
func (p *Plugin) CreateSnapshot(
	ctx context.Context,
	req *csi.CreateSnapshotRequest) (
	*csi.CreateSnapshotResponse, error) {

	defer glog.V(5).Info("end to CreateSnapshot")
	glog.V(5).Infof("start to CreateSnapshot, Name: %v, SourceVolumeId: %v, CreateSnapshotSecrets: %v, parameters: %v!",
		req.Name, req.SourceVolumeId, req.Secrets, req.Parameters)

	if 0 == len(req.Name) {
		return nil, status.Error(codes.InvalidArgument, "Snapshot Name cannot be empty")
	}

	if 0 == len(req.SourceVolumeId) {
		return nil, status.Error(codes.InvalidArgument, "Source Volume ID cannot be empty")
	}

	snapReq := &model.VolumeSnapshotSpec{
		Name:     req.Name,
		VolumeId: req.SourceVolumeId,
	}

	for k, v := range req.GetParameters() {
		switch strings.ToLower(k) {
		// TODO: support profile name
		case KParamProfile:
			snapReq.ProfileId = v
		}
	}

	glog.Infof("opensds CreateVolumeSnapshot request body: %v", snapReq)
	var snapshot *model.VolumeSnapshotSpec
	isExist, isCompatible, findSnapshot, err := FindSnapshot(snapReq)

	if err != nil {
		return nil, err
	}

	if isExist {
		if isCompatible {
			snapshot = findSnapshot
		} else {
			return nil, status.Error(codes.AlreadyExists,
				"Snapshot already exists but is incompatible")
		}
	} else {
		createSnapshot, err := Client.CreateVolumeSnapshot(snapReq)
		if err != nil {
			glog.Error("failed to CreateVolumeSnapshot", err)
			return nil, err
		}

		snapshot = createSnapshot
	}

	glog.V(5).Infof("opensds snapshot = %v", snapshot)
	creationTime, err := p.convertStringToPtypesTimestamp(snapshot.CreatedAt)
	if nil != err {
		return nil, err
	}

	return &csi.CreateSnapshotResponse{
		Snapshot: &csi.Snapshot{
			SizeBytes:      snapshot.Size * util.GiB,
			SnapshotId:     snapshot.Id,
			SourceVolumeId: snapshot.VolumeId,
			CreationTime:   creationTime,
			ReadyToUse:     true,
		},
	}, nil
}

func (p *Plugin) convertStringToPtypesTimestamp(timeStr string) (*timestamp.Timestamp, error) {
	timeAt, err := time.Parse(constants.TimeFormat, timeStr)
	if nil != err {
		return nil, status.Error(codes.Internal, err.Error())
	}
	ptypesTime, err := ptypes.TimestampProto(timeAt)
	if err != nil {
		return nil, err
	}
	return ptypesTime, nil
}

// DeleteSnapshot implementation
func (p *Plugin) DeleteSnapshot(
	ctx context.Context,
	req *csi.DeleteSnapshotRequest) (
	*csi.DeleteSnapshotResponse, error) {

	defer glog.V(5).Info("end to DeleteSnapshot")
	glog.V(5).Infof("start to DeleteSnapshot, SnapshotId: %v, DeleteSnapshotSecrets: %v!",
		req.SnapshotId, req.Secrets)

	if 0 == len(req.SnapshotId) {
		return nil, status.Error(codes.InvalidArgument, "Snapshot ID cannot be empty")
	}

	err := Client.DeleteVolumeSnapshot(req.SnapshotId, nil)

	if nil != err {
		return nil, err
	}

	return &csi.DeleteSnapshotResponse{}, nil
}

// ListSnapshots implementation
func (p *Plugin) ListSnapshots(
	ctx context.Context,
	req *csi.ListSnapshotsRequest) (
	*csi.ListSnapshotsResponse, error) {

	defer glog.V(5).Info("end to ListSnapshots")
	glog.V(5).Infof("start to ListSnapshots, MaxEntries: %v, StartingToken: %v, SourceVolumeId: %v, SnapshotId: %v!",
		req.MaxEntries, req.StartingToken, req.SourceVolumeId, req.SnapshotId)

	var opts map[string]string
	allSnapshots, err := Client.ListVolumeSnapshots(opts)
	if nil != err {
		return nil, err
	}

	snapshotId := req.GetSnapshotId()
	snapshotIDLen := len(snapshotId)
	sourceVolumeId := req.GetSourceVolumeId()
	sourceVolumeIdLen := len(sourceVolumeId)
	var snapshotsFilterByVolumeId []*model.VolumeSnapshotSpec
	var snapshotsFilterById []*model.VolumeSnapshotSpec
	var filterResult []*model.VolumeSnapshotSpec

	for _, snapshot := range allSnapshots {
		if snapshot.VolumeId == sourceVolumeId {
			snapshotsFilterByVolumeId = append(snapshotsFilterByVolumeId, snapshot)
		}

		if snapshot.Id == snapshotId {
			snapshotsFilterById = append(snapshotsFilterById, snapshot)
		}
	}

	switch {
	case (0 == snapshotIDLen) && (0 == sourceVolumeIdLen):
		if len(allSnapshots) <= 0 {
			glog.V(5).Info("len(allSnapshots) <= 0")
			return &csi.ListSnapshotsResponse{}, nil
		}

		filterResult = allSnapshots
		break
	case (0 == snapshotIDLen) && (0 != sourceVolumeIdLen):
		if len(snapshotsFilterByVolumeId) <= 0 {
			return nil, status.Error(codes.NotFound, fmt.Sprintf("no snapshot with source volume id %s", sourceVolumeId))
		}

		filterResult = snapshotsFilterByVolumeId
		break
	case (0 != snapshotIDLen) && (0 == sourceVolumeIdLen):
		if len(snapshotsFilterById) <= 0 {
			return nil, status.Error(codes.NotFound, fmt.Sprintf("no snapshot with id %s", snapshotId))
		}

		filterResult = snapshotsFilterById
		break
	case (0 != snapshotIDLen) && (0 != sourceVolumeIdLen):
		for _, snapshot := range snapshotsFilterById {
			if snapshot.VolumeId == sourceVolumeId {
				filterResult = append(filterResult, snapshot)
			}
		}

		if len(filterResult) <= 0 {
			return nil, status.Error(codes.NotFound,
				fmt.Sprintf("no snapshot with id %v and source volume id %v", snapshotId, sourceVolumeId))
		}

		break
	}

	glog.V(5).Infof("filterResult=%v.", filterResult)
	var sortedKeys []string
	snapshotsMap := make(map[string]*model.VolumeSnapshotSpec)

	for _, snapshot := range filterResult {
		sortedKeys = append(sortedKeys, snapshot.Id)
		snapshotsMap[snapshot.Id] = snapshot
	}
	sort.Strings(sortedKeys)

	var sortResult []*model.VolumeSnapshotSpec
	for _, key := range sortedKeys {
		sortResult = append(sortResult, snapshotsMap[key])
	}

	var (
		ulenSnapshots = int32(len(sortResult))
		maxEntries    = req.MaxEntries
		startingToken int32
	)

	if v := req.StartingToken; v != "" {
		i, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			return nil, status.Error(codes.Aborted, "parsing the startingToken failed")
		}
		startingToken = int32(i)
	}

	if startingToken >= ulenSnapshots {
		return nil, status.Errorf(
			codes.Aborted,
			"startingToken=%d >= len(snapshots)=%d",
			startingToken, ulenSnapshots)
	}

	// If maxEntries is 0 or greater than the number of remaining entries then
	// set maxEntries to the number of remaining entries.
	var sliceResult []*model.VolumeSnapshotSpec
	var nextToken string
	nextTokenIndex := startingToken + maxEntries

	if maxEntries == 0 || nextTokenIndex >= ulenSnapshots {
		sliceResult = sortResult[startingToken:]
	} else {
		sliceResult = sortResult[startingToken:nextTokenIndex]
		nextToken = fmt.Sprintf("%d", nextTokenIndex)
	}

	glog.V(5).Infof("sliceResult=%v, nextToken=%v.", sliceResult, nextToken)
	if len(sliceResult) <= 0 {
		return &csi.ListSnapshotsResponse{NextToken: nextToken}, nil
	}

	entries := []*csi.ListSnapshotsResponse_Entry{}
	for _, snapshot := range sliceResult {
		creationTime, err := p.convertStringToPtypesTimestamp(snapshot.CreatedAt)
		if nil != err {
			return nil, err
		}
		entries = append(entries, &csi.ListSnapshotsResponse_Entry{
			Snapshot: &csi.Snapshot{
				SizeBytes:      snapshot.Size * util.GiB,
				SnapshotId:     snapshot.Id,
				SourceVolumeId: snapshot.VolumeId,
				CreationTime:   creationTime,
				ReadyToUse:     true,
			},
		})
	}

	glog.V(5).Infof("entries=%v.", entries)
	return &csi.ListSnapshotsResponse{
		Entries:   entries,
		NextToken: nextToken,
	}, nil
}
