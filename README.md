# stuttgart-things/stageTime-server

gRPC Server for validating & producing revisionRuns (a collection of tekton pipelineRuns/stages)

## EXAMPLE DEPLOYMENT (DRAFT)

<details><summary>values</summary>

```
cat <<EOF > stageTime-server.yaml
---
secrets:
  redis-connection:
    name: redis-connection
    labels:
      app: stagetime-server
    dataType: data
    secretKVs:
      REDIS_SERVER: MTAuMzEuMTAxLjEzOA==
      REDIS_PORT: NjM3OQ==
      REDIS_PASSWORD: QXRsYW43aXM=
      REDIS_QUEUE: cmVkaXNxdWV1ZTp5YWNodC1yZXZpc2lvbnJ1bnM=

customresources:
  yas-ingress-certificate:
    apiVersion: cert-manager.io/v1
    kind: Certificate
    metadata:
      name: stagetime-server-ingress
      labels:
        app: stagetime-server
    spec:
      commonName: yas.app.sthings.tiab.ssc.sva.de
      dnsNames:
        - yas.app.sthings.tiab.ssc.sva.de
      issuerRef:
        name: cluster-issuer-ssc #cluster-issuer-approle
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
    hostname: yas
    # clusterName: dev
    domain: app.sthings.tiab.ssc.sva.de
    tls:
      secretName: stagetime-server-ingress-tls
      host: yas.app.sthings.tiab.ssc.sva.de
```
EOF

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
