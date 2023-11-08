# stuttgart-things/stageTime-server

gRPC Server for validating & producing revisionRuns (a collection of tekton pipelineRuns/stages)

## TASKS

```
task push # build & push image + chart to registry
```

## HELMFILE

<details><summary>SET VAULT CONNECTION</summary>

```bash
export VAULT_ADDR=https://${VAULT_FQDN}}
export VAULT_NAMESPACE=root

# APPROLE AUTH
export VAULT_AUTH_METHOD=approle
export VAULT_ROLE_ID=${VAULT_ROLE_ID}
export VAULT_SECRET_ID=${VAULT_SECRET_ID}

# TOKEN AUTH
export VAULT_AUTH_METHOD=token #default
export VAULT_TOKEN=${VAULT_TOKEN}
```

</details>

<details><summary>RENDER/APPLY</summary>

```bash
helmfile template --environment labul-pve-dev
helmfile sync --environment labul-pve-dev
```

</details>


## GRPC TEST CALL

```bash
# PORTWARD REDIS / EXAMPLE DEPLOYMENT!
kubectl port-forward --namespace stagetime-redis svc/redis-stack 28015:6379

# CHECK REDIS
redis-cli -h 127.0.0.1 -p 28015 -a <PASSWORD>

# SEND TEST DATA
go run tests/grpc-testClient-local.go
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
    dataType: stringData
    secretKVs:
      REDIS_SERVER: redis-stack-deployment-headless.redis-stack.svc.cluster.local
      REDIS_PORT: 6379
      REDIS_PASSWORD: <PASSWORD>

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
