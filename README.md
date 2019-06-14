# rekonfig

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
        rekonfig.gitops.in/konfig-hash: c5f9708c6efxxx
```
