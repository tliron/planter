#!/bin/bash
set -e

SEED=$(</dev/stdin)

SEED="$SEED

---

apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition

metadata:
  name: networks.planter.nephio.org
  annotations:
    planter.nephio.org/cluster: network-service|edge1 # hard dependency

spec:
  group: planter.nephio.org
  names:
    singular: network
    plural: networks
    kind: Network
    listKind: NetworkList
    categories:
    - all # will appear in \"kubectl get all\"
  scope: Namespaced
  versions:
  - name: v1alpha1
    served: true
    storage: true # one and only one version must be marked with storage=true
    schema:
      openAPIV3Schema:
        description: >-
          Planter network
        type: object
        required: [ spec ]
        properties:
          spec:
            type: object
            properties:
              provider:
                description: >-
                  Network provider
                type: string
"

echo "$SEED"
