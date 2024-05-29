package bean

import (
	"github.com/devtron-labs/common-lib/utils/k8s"
	"github.com/devtron-labs/common-lib/utils/remoteConnection/bean"
	remoteConnectionBean "github.com/devtron-labs/common-lib/utils/remoteConnection/bean"
)

type ClusterBean struct {
	Id            int               `json:"id" validate:"number"`
	ClusterName   string            `json:"cluster_name,omitempty" validate:"required"`
	Description   string            `json:"description"`
	ServerUrl     string            `json:"server_url,omitempty" validate:"url,required"`
	PrometheusUrl string            `json:"prometheus_url,omitempty" validate:"validate-non-empty-url"`
	Active        bool              `json:"active"`
	ProxyUrl      string            `json:"proxyUrl,omitempty"`
	Config        map[string]string `json:"config,omitempty"`
	//PrometheusAuth          *PrometheusAuth                  `json:"prometheusAuth,omitempty"`
	DefaultClusterComponent []*DefaultClusterComponent       `json:"defaultClusterComponent"`
	AgentInstallationStage  int                              `json:"agentInstallationStage,notnull"` // -1=external, 0=not triggered, 1=progressing, 2=success, 3=fails
	K8sVersion              string                           `json:"k8sVersion"`
	HasConfigOrUrlChanged   bool                             `json:"-"`
	UserName                string                           `json:"userName,omitempty"`
	InsecureSkipTLSVerify   bool                             `json:"insecureSkipTlsVerify"`
	ErrorInConnecting       string                           `json:"errorInConnecting"`
	IsCdArgoSetup           bool                             `json:"isCdArgoSetup"`
	IsVirtualCluster        bool                             `json:"isVirtualCluster"`
	isClusterNameEmpty      bool                             `json:"-"`
	ClusterUpdated          bool                             `json:"clusterUpdated"`
	ToConnectWithSSHTunnel  bool                             `json:"toConnectWithSSHTunnel,omitempty"`
	SSHTunnelConfig         *bean.SSHTunnelConfig            `json:"sshTunnelConfig,omitempty"`
	RemoteConnectionConfig  *bean.RemoteConnectionConfigBean `json:"remoteConnectionConfig"`
}
type DefaultClusterComponent struct {
	ComponentName  string `json:"name"`
	AppId          int    `json:"appId"`
	InstalledAppId int    `json:"installedAppId,omitempty"`
	EnvId          int    `json:"envId"`
	EnvName        string `json:"envName"`
	Status         string `json:"status"`
}

func (bean ClusterBean) GetClusterConfig() *k8s.ClusterConfig {
	//bean = *adapter.ConvertClusterBeanToNewClusterBean(&bean)
	configMap := bean.Config
	bearerToken := configMap[k8s.BearerToken]
	clusterCfg := &k8s.ClusterConfig{
		ClusterId:             bean.Id,
		ClusterName:           bean.ClusterName,
		Host:                  bean.ServerUrl,
		BearerToken:           bearerToken,
		InsecureSkipTLSVerify: bean.InsecureSkipTLSVerify,
	}
	if bean.InsecureSkipTLSVerify == false {
		clusterCfg.KeyData = configMap[k8s.TlsKey]
		clusterCfg.CertData = configMap[k8s.CertData]
		clusterCfg.CAData = configMap[k8s.CertificateAuthorityData]
	}
	if bean.RemoteConnectionConfig != nil {
		clusterCfg.RemoteConnectionConfig = &remoteConnectionBean.RemoteConnectionConfigBean{
			RemoteConnectionConfigId: bean.RemoteConnectionConfig.RemoteConnectionConfigId,
			ConnectionMethod:         remoteConnectionBean.RemoteConnectionMethod(bean.RemoteConnectionConfig.ConnectionMethod),
			ProxyConfig:              (*remoteConnectionBean.ProxyConfig)(bean.RemoteConnectionConfig.ProxyConfig),
			SSHTunnelConfig:          (*remoteConnectionBean.SSHTunnelConfig)(bean.RemoteConnectionConfig.SSHTunnelConfig),
		}
	}
	return clusterCfg
}
