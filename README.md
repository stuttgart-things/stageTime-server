# stuttgart-things/stageTime-server

gRPC Server for validating & producing revisionRuns (a collection of tekton pipelineRuns/stages)

## EXAMPLE DEPLOYMENT (DRAFT)

<details><summary>values</summary>
  
```
cat <<EOF > stageTime-server.yaml
---
configmaps:
  yas-configuration:
    PIPELINE_WORKSPACE: tekton-cd
ingress:
  yacht-application-server:
    hostname: yas
    clusterName: dev11
    domain: 4sthings.tiab.ssc.sva.de
    tls:
      host: yas.dev11.4sthings.tiab.ssc.sva.de
customresources:
  yas-ingress-certificate:
    spec:
      commonName: yas.dev11.4sthings.tiab.ssc.sva.de
      dnsNames:
      - yas.dev11.4sthings.tiab.ssc.sva.de
      issuerRef:
        name: cluster-issuer-approle
secrets:
  redis-connection:
    secretKVs:
      REDIS_SERVER: cmVkaXMtZGVwbG95bWVudC1oZWFkbGVzcy55YWNodC5zdmMuY2x1c3Rlci5sb2NhbA==
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

