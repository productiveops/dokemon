package dockerapi

import "github.com/productiveops/dokemon/pkg/server/store"

// Containers

type Port struct {
	IP string `json:"ip"`
	PrivatePort uint16 `json:"privatePort"`	// Container
	PublicPort uint16 `json:"publicPort"`	// Host
	Type string `json:"type"`
}

type Container struct {
	Id     	string `json:"id"`
	Name   	string `json:"name"`
	Image   string `json:"image"`
	Status 	string `json:"status"`
	State  	string `json:"state"`
	Ports  	[]Port `json:"ports"`
	Stale 	string `json:"stale"`	// yes, no, error, processing
}

type DockerContainerList struct {
	All bool `json:"all"`
}

type DockerContainerListResponse struct {
	Items []Container `json:"items"`
}

type DockerContainerStart struct {
	Id string `json:"id"`
}

type DockerContainerStop struct {
	Id string `json:"id"`
}

type DockerContainerRestart struct {
	Id string `json:"id"`
}

type DockerContainerRemove struct {
	Id    string `json:"id"`
	Force bool   `json:"force"`
}

type DockerContainerLogs struct {
	Id string `json:"id"`
}

type DockerContainerTerminal struct {
	Id string `json:"id"`
}

// Images

type Image struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Tag     string `json:"tag"`
	Size    int64  `json:"size"`
	Dangling bool  `json:"dangling"`
	Created int64  `json:"created"`
	InUse   bool   `json:"inUse"`
}

type DockerImageList struct {
	All bool `json:"all"`
}

type DockerImageListResponse struct {
	Items []Image `json:"items"`
}

type DockerImagePull struct {
	Image string `json:"image"`
	Tag   string `json:"tag"`
}

type DockerImageRemove struct {
	Id    string `json:"id"`
	Force bool   `json:"force"`
}

type DockerImagesPrune struct {
	All bool `json:"all"` // Remove all unused images, not just dangling
}

type DockerImagesPruneDeletedItem struct {
	Deleted  string `json:"deleted"`
	Untagged string `json:"untagged"`
}

type DockerImagesPruneResponse struct {
	ImagesDeleted  []DockerImagesPruneDeletedItem `json:"imagesDeleted"`
	SpaceReclaimed uint64                        `json:"spaceReclaimed"`
}

type DockerVolumeList struct {
}

type Volume struct {
	Driver string `json:"driver"`
	Name string `json:"name"`
}

type DockerVolumeListResponse struct {
	Items []Volume `json:"items"`
}

type DockerVolumeRemove struct {
	Name    string `json:"name"`
}

type DockerVolumesPrune struct {
	All bool `json:"all"` // Remove all unused volumes, not just anonymous ones
}

type DockerVolumesPruneResponse struct {
	VolumesDeleted  []string `json:"volumesDeleted"`
	SpaceReclaimed uint64    `json:"spaceReclaimed"`
}

type DockerNetworkList struct {
}

type Network struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Driver string `json:"driver"`
	Scope string `json:"scope"`
}

type DockerNetworkListResponse struct {
	Items []Network `json:"items"`
}

type DockerNetworkRemove struct {
	Id    string `json:"id"`
}

type DockerNetworksPrune struct {
}

type DockerNetworksPruneResponse struct {
	NetworksDeleted  []string `json:"networksDeleted"`
}

type DockerComposeList struct {
}

type DockerComposeGet struct {
	ProjectName string `json:"projectName"`
}

type ComposeItemInternal struct {
	Name string `json:"Name"`
	Status string `json:"Status"`
	ConfigFiles string `json:"ConfigFiles"`
}

type ComposeItem struct {
	Name string `json:"name"`
	Status string `json:"status"`
	ConfigFiles string `json:"configFiles"`
	Stale string `json:"stale"`
}

type DockerComposeListResponse struct {
	Items []ComposeItem `json:"items"`
}

type DockerComposeContainerList struct {
	ProjectName string `json:"projectName"`
}

type ComposeContainerInternal struct {
	Id string `json:"ID"`
	Name string `json:"Name"`
	Image string `json:"Image"`
	Service string `json:"Service"`
	Status string `json:"Status"`
	State string `json:"State"`
	Ports string `json:"Ports"`
}

type ComposeContainer struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Image string `json:"image"`
	Service string `json:"service"`
	Status string `json:"status"`
	State string `json:"state"`
	Ports string `json:"ports"`
	Stale string `json:"stale"`
}

type DockerComposeContainerListResponse struct {
	Items []ComposeContainer `json:"items"`
}

type DockerComposeLogs struct {
	ProjectName string `json:"projectName"`
}

type DockerComposeDeploy struct {
	ProjectName string `json:"projectName"`
	Definition string `json:"definition"`
	Variables map[string]store.VariableValue `json:"variables"`
}

type DockerComposePull struct {
	ProjectName string `json:"projectName"`
	Definition string `json:"definition"`
	Variables map[string]store.VariableValue `json:"variables"`
}

type DockerComposeUp struct {
	ProjectName string `json:"projectName"`
	Definition string `json:"definition"`
	Variables map[string]store.VariableValue `json:"variables"`
}

type DockerComposeDown struct {
	ProjectName string `json:"projectName"`
}

type DockerComposeDownNoStreaming struct {
	ProjectName string `json:"projectName"`
}

type DockerComposeProjectUnique struct {
	ProjectName string `json:"projectName"`
}

type DockerComposeProjectCreate struct {
	ProjectName string `json:"projectName"`
	Definition string `json:"definition"`
}

type DockerComposeProjectUpdate struct {
	ProjectName string `json:"projectName"`
	Definition string `json:"definition"`
}

type DockerComposeProjectDelete struct {
	ProjectName string `json:"projectName"`
}

type DockerComposeProjectDefinition struct {
	ProjectName string `json:"projectName"`
}

type DockerComposeProjectDefinitionResponse struct {
	ProjectName string `json:"projectName"`
	Definition string `json:"definition"`
}