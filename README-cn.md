[English](./README.md) | 简体中文

# knode-hosts-inject

![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.22-%23007d9c)
[![GoDoc](https://godoc.org/github.com/hargeek/knode-hosts-inject?status.svg)](https://pkg.go.dev/github.com/hargeek/knode-hosts-inject)
[![Contributors](https://img.shields.io/github/contributors/hargeek/knode-hosts-inject)](https://github.com/hargeek/knode-hosts-inject/graphs/contributors)
[![License](https://img.shields.io/github/license/hargeek/knode-hosts-inject)](./LICENSE)

## 核心功能

集群节点扩缩容时（有机器新加入集群或者有机器从集群中删除），自动确保集群中现有节点的host文件始终包含所有现有节点的“IP-主机名”映射关系

## 快速部署

```bash
~ kubectl apply -k https://github.com/Hargeek/knode-hosts-inject.git//deploy\?ref\=master
```

## 部署清单

```bash
~ tree ./deploy
./deploy
├── cluster-role-binding.yaml # ClusterRoleBinding
├── cluster-role.yaml         # ClusterRole
├── daemonset.yaml            # DaemonSet
├── kustomization.yaml        # Kustomize File
└── service-account.yaml      # ServiceAccount
```
