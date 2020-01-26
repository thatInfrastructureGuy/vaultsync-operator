package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// VaultSyncerSpec defines the desired state of VaultSyncer
type VaultSyncerSpec struct {
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	// +kubebuilder:validation:MaxLength=30
	// +kubebuilder:validation:MinLength=1
	Provider string `json:"provider"`
	// +kubebuilder:validation:MaxLength=30
	// +kubebuilder:validation:MinLength=1
	VaultName string `json:"vaultName"`
	// +kubebuilder:validation:MaxLength=30
	// +kubebuilder:validation:MinLength=1
	Consumer string `json:"consumer"`
	// +kubebuilder:validation:MaxLength=30
	// +kubebuilder:validation:MinLength=1
	SecretNamespace string `json:"secretNamespace"`
	// +kubebuilder:validation:MaxLength=30
	SecretName string `json:"secretName,omitempty"`

	// +kubebuilder:validation:MaxItems=500
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:UniqueItems=true
	DeploymentList              []string `json:"deploymentList,omitempty"`
	StatefulsetList             []string `json:"statefulsetList,omitempty"`
	RefreshRate                 int      `json:"refreshRate,omitempty"`
	ConvertHyphensToUnderscores bool     `json:"convertHyphensToUnderscores,omitempty"`
}

// VaultSyncerStatus defines the observed state of VaultSyncer
type VaultSyncerStatus struct {
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	dateUpdated string `json:"dateUpdated,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VaultSyncer is the Schema for the vaultsyncers API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=vaultsyncers,scope=Namespaced
type VaultSyncer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VaultSyncerSpec   `json:"spec,omitempty"`
	Status VaultSyncerStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VaultSyncerList contains a list of VaultSyncer
type VaultSyncerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VaultSyncer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VaultSyncer{}, &VaultSyncerList{})
}
