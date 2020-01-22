package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// VaultSyncerSpec defines the desired state of VaultSyncer
type VaultSyncerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Provider        string `json:"provider"`
	VaultName       string `json:"vaultname"`
	Consumer        string `json:"consumer"`
	SecretNamespace string `json:"secretnamespace"`
	SecretName      string `json:"secretname"`
	DeploymentList  string `json:"deploymentList"`
}

// VaultSyncerStatus defines the observed state of VaultSyncer
type VaultSyncerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	dateUpdated string `json:"string"`
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
