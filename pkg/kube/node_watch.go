package kube

import (
	"context"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog"
	"time"
)

var Node node

type node struct{}

func (n *node) NodesWalkWithWatch(clientSet *kubernetes.Clientset) {
	watcher, err := clientSet.CoreV1().Nodes().Watch(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		klog.Error("Failed to start node watcher: ", err)
		return
	}

	klog.Info("Start node watcher successfully")
	ch := watcher.ResultChan()

	for event := range ch {
		switch event.Type {
		case watch.Added:
			node := event.Object.(*coreV1.Node)
			klog.Info("node added: ", node.Name)
		case watch.Deleted:
			node := event.Object.(*coreV1.Node)
			klog.Info("node deleted: ", node.Name)
			//case watch.Modified:
			//	node := event.Object.(*coreV1.Node)
			//	klog.Info("node modified: ", node.Name)
		}
	}
}

func (n *node) NodesWalkWithInformer(clientSet *kubernetes.Clientset, onNodeAdd func(), onNodeDelete func()) {
	klog.Info("start node status listen")
	factory := informers.NewSharedInformerFactory(clientSet, 0)
	nodeInformer := factory.Core().V1().Nodes().Informer()
	_, err := nodeInformer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				node := obj.(*coreV1.Node)
				nodeName, nodeIP := n.parseNodeNameAndInternalIP(node)
				klog.Infof("new node added, node name: %s, node ip: %s", nodeName, nodeIP)
				onNodeAdd()
			},
			DeleteFunc: func(obj interface{}) {
				node := obj.(*coreV1.Node)
				nodeName, nodeIP := n.parseNodeNameAndInternalIP(node)
				klog.Infof("node deleted, node name: %s, node ip: %s", nodeName, nodeIP)
				onNodeDelete()
			},
			//UpdateFunc: func(oldObj, newObj interface{}) {
			//	node := newObj.(*coreV1.Node)
			//	klog.Info("node modified: ", node.Name)
			//},
		})
	if err != nil {
		klog.Error("failed to add event handler: ", err)
		return
	}

	stopCh := make(chan struct{})
	defer close(stopCh)

	factory.Start(stopCh)

	if !cache.WaitForCacheSync(stopCh, nodeInformer.HasSynced) {
		klog.Error("failed to sync node informer")
		return
	}

	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				klog.Info("node status listen process running...")
			case <-stopCh:
				ticker.Stop()
				klog.Info("node status listen process stopped")
				return
			}
		}
	}()

	<-stopCh

	return
}

func (n *node) parseNodeNameAndInternalIP(node *coreV1.Node) (nodeName, nodeIP string) {
	for _, addr := range node.Status.Addresses {
		switch addr.Type {
		case coreV1.NodeHostName:
			nodeName = addr.Address
		case coreV1.NodeInternalIP:
			nodeIP = addr.Address
		}
	}
	return
}
