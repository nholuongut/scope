package kubernetes

import (
	"github.com/nholuongut/scope/report"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/labels"
)

const (
	// BetaStorageClassAnnotation is the annotation for default storage class
	BetaStorageClassAnnotation = "volume.beta.kubernetes.io/storage-class"
)

// PersistentVolumeClaim represents kubernetes PVC interface
type PersistentVolumeClaim interface {
	Meta
	Selector() (labels.Selector, error)
	GetNode(string) report.Node
	GetStorageClass() string
	GetCapacity() string
}

// persistentVolumeClaim represents kubernetes Persistent Volume Claims
type persistentVolumeClaim struct {
	*apiv1.PersistentVolumeClaim
	Meta
}

// NewPersistentVolumeClaim returns new Persistent Volume Claim type
func NewPersistentVolumeClaim(p *apiv1.PersistentVolumeClaim) PersistentVolumeClaim {
	return &persistentVolumeClaim{PersistentVolumeClaim: p, Meta: meta{p.ObjectMeta}}
}

// GetStorageClass will fetch storage class name from given PVC
func (p *persistentVolumeClaim) GetStorageClass() string {

	// Use Beta storage class annotation first
	storageClassName := p.Annotations[BetaStorageClassAnnotation]
	if storageClassName != "" {
		return storageClassName
	}
	if p.Spec.StorageClassName != nil {
		storageClassName = *p.Spec.StorageClassName
	}

	return storageClassName
}

// GetCapacity returns the storage size of PVC
func (p *persistentVolumeClaim) GetCapacity() string {
	capacity := p.Spec.Resources.Requests[apiv1.ResourceStorage]
	if capacity.String() != "" {
		return capacity.String()
	}
	return ""
}

// GetNode returns Persistent Volume Claim as Node
func (p *persistentVolumeClaim) GetNode(probeID string) report.Node {
	latests := map[string]string{
		NodeType:              "Persistent Volume Claim",
		Status:                string(p.Status.Phase),
		VolumeName:            p.Spec.VolumeName,
		StorageClassName:      p.GetStorageClass(),
		report.ControlProbeID: probeID,
	}

	if p.GetCapacity() != "" {
		latests[VolumeCapacity] = p.GetCapacity()
	}

	return p.MetaNode(report.MakePersistentVolumeClaimNodeID(p.UID())).
		WithLatests(latests).
		WithLatestActiveControls(CreateVolumeSnapshot, Describe)
}

// Selector returns all Persistent Volume Claim selector
func (p *persistentVolumeClaim) Selector() (labels.Selector, error) {
	selector, err := metav1.LabelSelectorAsSelector(p.Spec.Selector)
	if err != nil {
		return nil, err
	}
	return selector, nil
}
