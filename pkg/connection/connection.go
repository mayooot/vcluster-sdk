package connection

import (
	"context"
	"fmt"

	"github.com/loft-sh/log"
	"github.com/loft-sh/vcluster/pkg/util/clihelper"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	defaultClientQPS   = 100
	defaultClientBurst = 200
)

func GetVClusterClientset(ctx context.Context, client *kubernetes.Clientset, vclusterName string, vclusterNamespace string, logger log.Logger) (*kubernetes.Clientset, error) {
	apiConfig, err := clihelper.GetKubeConfig(ctx, client, vclusterName, vclusterNamespace, logger)
	if err != nil {
		return nil, fmt.Errorf("could not get api server config: %w", err)
	}

	clientConfig, err := clientcmd.NewNonInteractiveClientConfig(*apiConfig, "", &clientcmd.ConfigOverrides{}, nil).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("could not get api server config: %w", err)
	}

	apply(clientConfig,
		withInsecure(),
		withHost(vclusterName, vclusterNamespace),
		withDefaultQPSAndBurst(),
	)

	return kubernetes.NewForConfig(clientConfig)
}

type patcherOption func(*rest.Config)

func apply(base *rest.Config, patchers ...patcherOption) *rest.Config {
	for _, patcher := range patchers {
		patcher(base)
	}
	return base
}

func withInsecure() patcherOption {
	return func(config *rest.Config) {
		config.Insecure = true
		config.TLSClientConfig.CertFile = ""
		config.TLSClientConfig.CAData = []byte{}
	}
}

func withHost(vclusterName, vclusterNamespace string) patcherOption {
	return func(config *rest.Config) {
		config.Host = fmt.Sprintf("https://%s.%s.svc.cluster.local:8443", vclusterName, vclusterNamespace)
	}
}

func withDefaultQPSAndBurst() patcherOption {
	return func(config *rest.Config) {
		config.QPS = defaultClientQPS
		config.Burst = defaultClientBurst
	}
}
