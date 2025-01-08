package cmd

import "github.com/hargeek/knode-hosts-inject/pkg/hosts"

func RunServer() {
	hosts.NodeHosts.StartInjectAndRefreshNodeHostsTask()
	//make main goroutine run forever
	select {}
}
