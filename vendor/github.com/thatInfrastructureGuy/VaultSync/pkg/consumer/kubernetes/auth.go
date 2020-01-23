package kubernetes

import (
	"log"
	"os"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Config contains the kubernetes clientset and configs.
type Config struct {
	clientset        *kubernetes.Clientset
	secretObject     *apiv1.Secret
	SecretName       string
	Namespace        string
	KubeSecretExists bool
}

// authenticate authenticates and authorizes the client.
func (k *Config) authenticate() (err error) {
	var config *rest.Config
	kubeconfigPath := os.Getenv("KUBECONFIG")
	if len(kubeconfigPath) > 0 {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	} else {
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		log.Printf("Cannot connect to Kubernetes: %v\n", err)
		return err
	}
	k.clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
