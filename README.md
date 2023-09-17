# stuttgart-things/stageTime-server

gRPC Server for validating & producing revisionRuns (a collection of tekton pipelineRuns/stages)

## TASKS

```
task push # build image + chart
```


## EXAMPLE DEPLOYMENT (DRAFT)

<details><summary>values</summary>

```yaml
cat <<EOF > stageTime-server.yaml
---
secrets:
  redis-connection:
    name: redis-connection
    labels:
      app: stagetime-server
    dataType: data
    secretKVs:
      REDIS_SERVER: cmVkaXMtc3RhY2stZGVwbG95bWVudC1oZWFkbGVzcy5yZWRpcy1zdGFjay5zdmMuY2x1c3Rlci5sb2NhbAo=
      REDIS_PORT: NjM3OQ==
      REDIS_PASSWORD: d2Vhaw==

customresources:
  stagetime-ingress-certificate:
    apiVersion: cert-manager.io/v1
    kind: Certificate
    metadata:
      name: stagetime-server-ingress
      labels:
        app: stagetime-server
    spec:
      commonName: stagetime.app.sthings.tiab.ssc.sva.de
      dnsNames:
        - stagetime.app.sthings.tiab.ssc.sva.de
      issuerRef:
        name: cluster-issuer-approle
        kind: ClusterIssuer
      secretName: stagetime-server-ingress-tls

ingress:
  stagetime-server:
    labels:
      app: stagetime-server
    name: stagetime-server
    ingressClassName: nginx
    annotations:
      nginx.ingress.kubernetes.io/ssl-redirect: "true"
      nginx.ingress.kubernetes.io/backend-protocol: "GRPC"
    service:
      name: stagetime-server-service
      port: 80
      path: /
      pathType: Prefix
    hostname: stagetime
    domain: dev21.sthings-vsphere.labul.sva.de
    tls:
      secretName: stagetime-server-ingress-tls
      host: stagetime.dev21.sthings-vsphere.labul.sva.de
EOF
```

```bash
helm upgrade --install server helm/
stagetime-server/ -n stagetime --create-namespace --values stageTime-server.yaml
```

</details>

## LICENSE

<details><summary><b>APACHE 2.0</b></summary>

Copyright 2023 patrick hermann.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

</details>

Author Information
------------------
Patrick Hermann, stuttgart-things 04/2023
