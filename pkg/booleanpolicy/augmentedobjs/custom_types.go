package augmentedobjs

import (
	"github.com/stackrox/rox/generated/storage"
)

// This block enumerates custom tags.
const (
	ComponentAndVersionCustomTag  = "Component And Version"
	ContainerNameCustomTag        = "Container Name"
	DockerfileLineCustomTag       = "Dockerfile Line"
	EnvironmentVarCustomTag       = "Environment Variable"
	ImageScanCustomTag            = "Image Scan"
	NetworkFlowSrcNameCustomTag   = "Network Flow Source Name"
	NetworkFlowDstNameCustomTag   = "Network Flow Destination Name"
	NetworkFlowDstPortCustomTag   = "Network Flow Destination Port"
	NetworkFlowL4Protocol         = "Network Flow L4 Protocol"
	NotInNetworkBaselineCustomTag = "Not In Network Baseline"
	NotInProcessBaselineCustomTag = "Not In Baseline"
	KubernetesAPIVerbCustomTag    = "Kubernetes API Verb"
	KubernetesResourceCustomTag   = "Kubernetes Resource"
)

type dockerfileLine struct {
	Line string `search:"Dockerfile Line"`
}

type componentAndVersion struct {
	ComponentAndVersion string `search:"Component And Version"`
}

type baselineResult struct {
	NotInBaseline bool `search:"Not In Baseline"`
}

type networkFlowDetails struct {
	SrcEntityName        string             `search:"Network Flow Source Name"`
	DstEntityName        string             `search:"Network Flow Destination Name"`
	DstPort              uint32             `search:"Network Flow Destination Port"`
	L4Protocol           storage.L4Protocol `search:"Network Flow L4 Protocol"`
	NotInNetworkBaseline bool               `search:"Not In Network Baseline"`
}

type envVar struct {
	EnvVar string `search:"Environment Variable"`
}
