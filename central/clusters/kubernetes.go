package clusters

import (
	"encoding/base64"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/renderer"
	"github.com/stackrox/rox/pkg/zip"
)

func init() {
	deployers[storage.ClusterType_KUBERNETES_CLUSTER] = newKubernetes()
}

type kubernetes struct{}

func newKubernetes() Deployer {
	return &kubernetes{}
}

var monitoringFilenames = []string{
	"kubernetes/kubectl/telegraf.conf",
}

var admissionController = "kubernetes/kubectl/admission-controller.yaml"

func (k *kubernetes) Render(c Wrap, ca []byte) ([]*zip.File, error) {
	fields, err := fieldsFromWrap(c)
	if err != nil {
		return nil, err
	}

	fields["K8sCommand"] = "kubectl"

	filenames := renderer.FileNameMap{
		"kubernetes/common/ca-setup.sh":  "ca-setup-sensor.sh",
		"kubernetes/common/delete-ca.sh": "delete-ca-sensor.sh",
	}
	filenames.Add(
		"kubernetes/kubectl/sensor.sh",
		"kubernetes/kubectl/sensor.yaml",
		"kubernetes/kubectl/sensor-rbac.yaml",
		"kubernetes/kubectl/sensor-netpol.yaml",
		"kubernetes/kubectl/delete-sensor.sh",
		"kubernetes/kubectl/sensor-pod-security.yaml",
	)

	if c.MonitoringEndpoint != "" {
		filenames.Add(monitoringFilenames...)
	}

	if c.AdmissionController {
		fields["CABundle"] = base64.StdEncoding.EncodeToString(ca)
		filenames.Add(admissionController)
	}

	allFiles, err := renderer.RenderFiles(filenames, fields)
	if err != nil {
		return nil, err
	}

	assetFiles, err := renderer.LoadAssets(renderer.NewFileNameMap(dockerAuthAssetFile))
	if err != nil {
		return nil, err
	}

	allFiles = append(allFiles, assetFiles...)
	return allFiles, err
}
