package hosts

import (
	"context"
	"fmt"
	"github.com/Hargeek/kube-tools/client"
	"github.com/hargeek/knode-hosts-inject/pkg/kube"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog"
	"os"
	"path/filepath"
	"time"
)

var NodeHosts nodeHosts

type nodeHosts struct{}

func (n *nodeHosts) StartInjectAndRefreshNodeHostsTask() {
	klog.Info("start inject and refresh node hosts task")
	nodeHostsFilePath := os.Getenv("KUBE_NODE_HOSTS_PATH")
	if nodeHostsFilePath == "" {
		klog.Fatalf("the env KUBE_NODE_HOSTS_PATH is empty, but it is required")
	}
	clientSet := client.KubeClientSetM.GetClient()

	startTime := time.Now()
	delayDuration := 60 * time.Second

	// check if the delay period has passed
	isDelayPassed := func() bool {
		return time.Since(startTime) > delayDuration
	}

	wrappedUpdateHostsFile := func() {
		if isDelayPassed() {
			n.updateHostsFile(clientSet, nodeHostsFilePath)
		} else {
			klog.Info("skipping first hosts file update during initial delay period")
		}
	}

	onNodeAdd := func() {
		wrappedUpdateHostsFile()
	}

	onNodeDelete := func() {
		wrappedUpdateHostsFile()
	}

	go kube.Node.NodesWalkWithInformer(clientSet, onNodeAdd, onNodeDelete)

	klog.Info("initial hosts file update after 60 seconds delay after startup")
	time.Sleep(delayDuration)
	n.updateHostsFile(clientSet, nodeHostsFilePath)
}

func (n *nodeHosts) updateHostsFile(clientSet *kubernetes.Clientset, nodeHostsFilePath string) {
	nodes, err := clientSet.CoreV1().Nodes().List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		klog.Fatalf("failed to list nodes: %v", err)
		return
	}

	// hosts file content of head
	hostsContent := `127.0.0.1   localhost localhost.localdomain localhost4 localhost4.localdomain4
::1         localhost localhost.localdomain localhost6 localhost6.localdomain6

`

	// gen hosts content
	for _, node := range nodes.Items {
		var internalIP, hostname string
		for _, addr := range node.Status.Addresses {
			if addr.Type == v1.NodeInternalIP {
				internalIP = addr.Address
			} else if addr.Type == v1.NodeHostName {
				hostname = addr.Address
			}
		}
		if internalIP != "" && hostname != "" {
			hostsContent += fmt.Sprintf("%s %s\n", internalIP, hostname)
		}
	}

	// write hosts file
	err = os.WriteFile(filepath.Clean(nodeHostsFilePath), []byte(hostsContent), 0644)
	if err != nil {
		klog.Fatalf("failed to write hosts file: %v", err)
	}
	klog.Info("hosts file refresh successfully")
}
