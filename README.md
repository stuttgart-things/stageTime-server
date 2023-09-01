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


Author Information
------------------
Patrick Hermann, stuttgart-things 04/2023
