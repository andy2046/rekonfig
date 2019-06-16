# Rekonfig

Rekonfig monitors `Deployment`s and `ConfigMap`s/`Secret`s mounted, ensures that each `Deployment`'s `Pod`s always have up to date configuration.

Whenever a `ConfigMap` or `Secret` is updated, Rekonfig can trigger a `Rolling Update` of the `Deployment`.

## Installation

To deploy Reconfig to a Kubernetes cluster.

```bash
$ kubectl apply -f ./deploy
```

## Configuration

Rekonfig watches all `Deployment`s within a Kubernetes cluster but only processes those with the annotation `rekonfig.gitops.in/update-on-konfig-change: "true"`.

As shown below, once enabled, Rekonfig will set the configuration hash as an annotation `rekonfig.gitops.in/konfig-hash` on the `Deployment`'s `PodTemplate`.

```yaml
apiVersion: apps/v1beta1
kind: StatefulSet
metadata:
  annotations:
    rekonfig.gitops.in/update-on-konfig-change: "true"
  labels:
    app: my-app
  name: my-app
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      app: my-app
  serviceName: my-app-headless
  template:
    metadata:
      annotations:
        rekonfig.gitops.in/konfig-hash: "<SHA256_HASH>"
# ...
```

## Test

```bash
$ make test
```

## Build

```bash
$ make build
```
