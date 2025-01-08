English | [简体中文](./README-cn.md)

# knode-hosts-inject

![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.22-%23007d9c)
[![GoDoc](https://godoc.org/github.com/hargeek/gin-auto-redoc?status.svg)](https://pkg.go.dev/github.com/hargeek/gin-auto-redoc)
[![Contributors](https://img.shields.io/github/contributors/hargeek/gin-auto-redoc)](https://github.com/hargeek/gin-auto-redoc/graphs/contributors)
[![License](https://img.shields.io/github/license/hargeek/gin-auto-redoc)](./LICENSE)

## What is knode-hosts-inject

When the cluster nodes are scaled up or down (new machines join the cluster or machines are removed from the cluster), knode-hosts-inject ensures that the host file of the existing nodes in the cluster always contains the "IP-hostname" mapping of all existing nodes

## Quick Start

```bash
~ kubectl apply -k https://github.com/Hargeek/knode-hosts-inject.git//deploy\?ref\=master
```

## Deploy Manifests

```bash
~ tree ./deploy
./deploy
├── cluster-role-binding.yaml # ClusterRoleBinding
├── cluster-role.yaml         # ClusterRole
├── daemonset.yaml            # DaemonSet
├── kustomization.yaml        # Kustomize File
└── service-account.yaml      # ServiceAccount
```
