#!/usr/bin/env python3

import sys
from ruamel.yaml import YAML
yaml=YAML(typ='safe')
yaml.default_flow_style = False


ANNOTATION_CLUSTER = 'planter.nephio.org/cluster'
ANNOTATION_USER = 'planter.nephio.org/user'
ANNOTATION_DEPENDENCIES = 'planter.nephio.org/dependencies'

KIND_NAMESPACE = 'v1|Namespace'
KIND_SERVICE_ACCOUNT = 'v1|ServiceAccount'
KIND_ROLE = 'rbac.authorization.k8s.io/v1|Role'
KIND_ROLE_BINDING = 'rbac.authorization.k8s.io/v1|RoleBinding'
KIND_CRD = 'apiextensions.k8s.io/v1|CustomResourceDefinition'
KIND_CLUSTER = 'planter.nephio.org/v1alpha1|Cluster'


processors = {}
resources = {}


def process_resource(resource):
    metadata = resource.get('metadata', {})
    annotations = metadata.get('annotations', {})

    if 'namespace' in metadata:
        namespace = get_namespace(metadata['namespace'])
        if namespace is not None:
            add_dependency(annotations, namespace)
            if not annotations.get(ANNOTATION_CLUSTER):
                cluster = get_cluster_for(namespace)
                if cluster:
                    annotations[ANNOTATION_CLUSTER] = get_resource_fullname(cluster)

    if not annotations.get(ANNOTATION_CLUSTER):
        annotations[ANNOTATION_CLUSTER] = 'SELF'

    crd = get_resource_crd(resource)
    if crd:
        add_dependency(annotations, crd)

    metadata['annotations'] = annotations
    resource['metadata'] = metadata


def process_admin(resource):
    process_resource(resource)

    metadata = resource.get('metadata', {})
    annotations = metadata.get('annotations', {})

    if not annotations.get(ANNOTATION_USER):
        annotations[ANNOTATION_USER] = 'admin'

    metadata['annotations'] = annotations
    resource['metadata'] = metadata


processors[''] = process_resource
processors[KIND_NAMESPACE] = process_admin
processors[KIND_SERVICE_ACCOUNT] = process_admin
processors[KIND_ROLE] = process_admin
processors[KIND_ROLE_BINDING] = process_admin
processors[KIND_CRD] = process_admin


def get_resource_name(resource):
    return resource.get('metadata', {}).get('name', '')


def get_resource_namepsace(resource):
    return resource.get('metadata', {}).get('namespace', '')


def get_resource_fullname(resource):
    return get_resource_namepsace(resource) + '|' + get_resource_name(resource)


def get_resource_kind(resource):
    return resource.get('apiVersion', '') + '|' + resource.get('kind', '')


def get_resource_id(resource):
    return get_resource_kind(resource) + '|' + get_resource_fullname(resource)


def get_resource_crd(resource):
    kind = get_resource_kind(resource)
    for crd in get_crds():
        if kind in get_crd_kinds(crd):
            return crd
    return None


def get_crds():
    for id, resource in resources.items():
        if id.startswith(KIND_CRD + '||'):
            yield resource


def get_crd_kinds(crd):
    spec = crd.get('spec', {})
    group = spec.get('group', '')
    kind = spec.get('names', {}).get('kind', '')
    versions = spec.get('versions', [])
    for version in versions:
        yield group + '/' + version.get('name', '') + '|' + kind


def get_namespace(name):
    return resources.get(KIND_NAMESPACE + '||' + name)


def get_cluster(fullname):
    return resources.get(KIND_CLUSTER + '|' + fullname)


def get_cluster_for(resource):
    annotations = resource.get('metadata', {}).get('annotations', {})
    fullname = annotations.get(ANNOTATION_CLUSTER)
    if fullname:
        return get_cluster(fullname)
    return None


def add_dependency(annotations, dependency):
    dependencies = annotations.get(ANNOTATION_DEPENDENCIES, '')
    if dependencies:
        dependencies += ','
    dependencies += get_resource_id(dependency)
    annotations[ANNOTATION_DEPENDENCIES] = dependencies


def main():
    seed = list(yaml.load_all(sys.stdin)) or {}

    for resource in seed:
        resources[get_resource_id(resource)] = resource

    for resource in seed:
        kind = get_resource_kind(resource)
        if kind:
            processor = processors.get(kind, processors.get(''))
            if processor:
                processor(resource)

    yaml.dump_all(seed, sys.stdout)


if __name__ == '__main__':
    main()
