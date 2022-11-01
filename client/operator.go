package client

import (
	"fmt"

	resources "github.com/tliron/planter/resources/planter.nephio.org/v1alpha1"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (self *Client) InstallOperator(sourceRegistryHost string, wait bool) error {
	var err error

	if _, err = self.CreateCustomResourceDefinition(&resources.SeedCustomResourceDefinition); err != nil {
		return err
	}

	if _, err = self.CreateCustomResourceDefinition(&resources.ClusterCustomResourceDefinition); err != nil {
		return err
	}

	if _, err = self.CreateCustomResourceDefinition(&resources.KernelConfigCustomResourceDefinition); err != nil {
		return err
	}

	if _, err = self.CreateCustomResourceDefinition(&resources.MemoryConfigCustomResourceDefinition); err != nil {
		return err
	}

	if _, err = self.CreateCustomResourceDefinition(&resources.WANCustomResourceDefinition); err != nil {
		return err
	}

	if _, err = self.CreateCustomResourceDefinition(&resources.NetworkCustomResourceDefinition); err != nil {
		return err
	}

	if _, err = self.createOperatorNamespace(); err != nil {
		return err
	}

	var serviceAccount *core.ServiceAccount
	if serviceAccount, err = self.createOperatorServiceAccount(); err != nil {
		return err
	}

	if _, err = self.createOperatorClusterRoleBinding(serviceAccount, self.ClusterRole); err != nil {
		return err
	}

	var operatorDeployment *apps.Deployment
	if operatorDeployment, err = self.createOperatorDeployment(sourceRegistryHost, serviceAccount, 1); err != nil {
		return err
	}

	if wait {
		if _, err := self.WaitForDeployment(self.Namespace, operatorDeployment.Name); err != nil {
			return err
		}
	}

	return nil
}

func (self *Client) UninstallOperator(wait bool) {
	var gracePeriodSeconds int64 = 0
	deleteOptions := meta.DeleteOptions{
		GracePeriodSeconds: &gracePeriodSeconds,
	}

	name := fmt.Sprintf("%s-operator", self.NamePrefix)

	// Deployment
	if err := self.Kubernetes.AppsV1().Deployments(self.Namespace).Delete(self.Context, name, deleteOptions); err != nil {
		self.Log.Warningf("%s", err)
	}

	// Cluster role binding
	if err := self.Kubernetes.RbacV1().ClusterRoleBindings().Delete(self.Context, self.NamePrefix, deleteOptions); err != nil {
		self.Log.Warningf("%s", err)
	}

	// Service account
	if err := self.Kubernetes.CoreV1().ServiceAccounts(self.Namespace).Delete(self.Context, self.NamePrefix, deleteOptions); err != nil {
		self.Log.Warningf("%s", err)
	}

	// Seed custom resource definition
	if err := self.APIExtensions.ApiextensionsV1().CustomResourceDefinitions().Delete(self.Context, resources.SeedCustomResourceDefinition.Name, deleteOptions); err != nil {
		self.Log.Warningf("%s", err)
	}

	// Cluster custom resource definition
	if err := self.APIExtensions.ApiextensionsV1().CustomResourceDefinitions().Delete(self.Context, resources.ClusterCustomResourceDefinition.Name, deleteOptions); err != nil {
		self.Log.Warningf("%s", err)
	}

	// KernelConfig custom resource definition
	if err := self.APIExtensions.ApiextensionsV1().CustomResourceDefinitions().Delete(self.Context, resources.KernelConfigCustomResourceDefinition.Name, deleteOptions); err != nil {
		self.Log.Warningf("%s", err)
	}

	// MemoryConfig custom resource definition
	if err := self.APIExtensions.ApiextensionsV1().CustomResourceDefinitions().Delete(self.Context, resources.MemoryConfigCustomResourceDefinition.Name, deleteOptions); err != nil {
		self.Log.Warningf("%s", err)
	}

	// WAN custom resource definition
	if err := self.APIExtensions.ApiextensionsV1().CustomResourceDefinitions().Delete(self.Context, resources.WANCustomResourceDefinition.Name, deleteOptions); err != nil {
		self.Log.Warningf("%s", err)
	}

	// Network custom resource definition
	if err := self.APIExtensions.ApiextensionsV1().CustomResourceDefinitions().Delete(self.Context, resources.NetworkCustomResourceDefinition.Name, deleteOptions); err != nil {
		self.Log.Warningf("%s", err)
	}

	if wait {
		getOptions := meta.GetOptions{}
		self.WaitForDeletion("operator deployment", func() bool {
			_, err := self.Kubernetes.AppsV1().Deployments(self.Namespace).Get(self.Context, name, getOptions)
			return err == nil
		})
		self.WaitForDeletion("cluster role binding", func() bool {
			_, err := self.Kubernetes.RbacV1().ClusterRoleBindings().Get(self.Context, self.NamePrefix, getOptions)
			return err == nil
		})
		self.WaitForDeletion("service account", func() bool {
			_, err := self.Kubernetes.CoreV1().ServiceAccounts(self.Namespace).Get(self.Context, self.NamePrefix, getOptions)
			return err == nil
		})
		self.WaitForDeletion("seed custom resource definition", func() bool {
			_, err := self.APIExtensions.ApiextensionsV1().CustomResourceDefinitions().Get(self.Context, resources.SeedCustomResourceDefinition.Name, getOptions)
			return err == nil
		})
		self.WaitForDeletion("cluster custom resource definition", func() bool {
			_, err := self.APIExtensions.ApiextensionsV1().CustomResourceDefinitions().Get(self.Context, resources.ClusterCustomResourceDefinition.Name, getOptions)
			return err == nil
		})
		self.WaitForDeletion("kernelconfig custom resource definition", func() bool {
			_, err := self.APIExtensions.ApiextensionsV1().CustomResourceDefinitions().Get(self.Context, resources.KernelConfigCustomResourceDefinition.Name, getOptions)
			return err == nil
		})
		self.WaitForDeletion("memoryconfig custom resource definition", func() bool {
			_, err := self.APIExtensions.ApiextensionsV1().CustomResourceDefinitions().Get(self.Context, resources.MemoryConfigCustomResourceDefinition.Name, getOptions)
			return err == nil
		})
		self.WaitForDeletion("wan custom resource definition", func() bool {
			_, err := self.APIExtensions.ApiextensionsV1().CustomResourceDefinitions().Get(self.Context, resources.WANCustomResourceDefinition.Name, getOptions)
			return err == nil
		})
		self.WaitForDeletion("network custom resource definition", func() bool {
			_, err := self.APIExtensions.ApiextensionsV1().CustomResourceDefinitions().Get(self.Context, resources.NetworkCustomResourceDefinition.Name, getOptions)
			return err == nil
		})
	}
}

func (self *Client) createOperatorNamespace() (*core.Namespace, error) {
	return self.CreateNamespace(&core.Namespace{
		ObjectMeta: meta.ObjectMeta{
			Name: self.Namespace,
		},
	})
}

func (self *Client) createOperatorServiceAccount() (*core.ServiceAccount, error) {
	return self.CreateServiceAccount(&core.ServiceAccount{
		ObjectMeta: meta.ObjectMeta{
			Name:   self.NamePrefix,
			Labels: self.Labels(fmt.Sprintf("%s-operator", self.NamePrefix), "operator", self.Namespace),
		},
	})
}

func (self *Client) createOperatorClusterRoleBinding(serviceAccount *core.ServiceAccount, role string) (*rbac.ClusterRoleBinding, error) {
	return self.CreateClusterRoleBinding(&rbac.ClusterRoleBinding{
		ObjectMeta: meta.ObjectMeta{
			Name:   self.NamePrefix,
			Labels: self.Labels(fmt.Sprintf("%s-operator", self.NamePrefix), "operator", self.Namespace),
		},
		Subjects: []rbac.Subject{
			{
				Kind:      rbac.ServiceAccountKind, // serviceAccount.Kind is empty
				Name:      serviceAccount.Name,
				Namespace: self.Namespace, // required
			},
		},
		RoleRef: rbac.RoleRef{
			APIGroup: rbac.GroupName,
			Kind:     "ClusterRole",
			Name:     role,
		},
	})
}

func (self *Client) createOperatorDeployment(sourceRegistryHost string, serviceAccount *core.ServiceAccount, replicas int32) (*apps.Deployment, error) {
	appName := fmt.Sprintf("%s-operator", self.NamePrefix)
	labels := self.Labels(appName, "operator", self.Namespace)

	deployment := &apps.Deployment{
		ObjectMeta: meta.ObjectMeta{
			Name:   appName,
			Labels: labels,
		},
		Spec: apps.DeploymentSpec{
			Replicas: &replicas,
			Selector: &meta.LabelSelector{
				MatchLabels: labels,
			},
			Template: core.PodTemplateSpec{
				ObjectMeta: meta.ObjectMeta{
					Labels: labels,
				},
				Spec: core.PodSpec{
					ServiceAccountName: serviceAccount.Name,
					Containers: []core.Container{
						{
							Name:            "operator",
							Image:           fmt.Sprintf("%s/%s", sourceRegistryHost, self.OperatorImageName),
							ImagePullPolicy: core.PullAlways,
							VolumeMounts: []core.VolumeMount{
								{
									Name:      "cache",
									MountPath: self.CachePath,
								},
							},
							Env: []core.EnvVar{
								{
									Name:  "PLANTER_OPERATOR_cache",
									Value: self.CachePath,
								},
								{
									Name:  "PLANTER_OPERATOR_concurrency",
									Value: "3",
								},
								{
									Name:  "PLANTER_OPERATOR_verbose",
									Value: "1",
								},
								{
									// For kutil's kubernetes.GetConfiguredNamespace
									Name: "KUBERNETES_NAMESPACE",
									ValueFrom: &core.EnvVarSource{
										FieldRef: &core.ObjectFieldSelector{
											FieldPath: "metadata.namespace",
										},
									},
								},
							},
							LivenessProbe: &core.Probe{
								ProbeHandler: core.ProbeHandler{
									HTTPGet: &core.HTTPGetAction{
										Port: intstr.FromInt(8086),
										Path: "/live",
									},
								},
							},
							ReadinessProbe: &core.Probe{
								ProbeHandler: core.ProbeHandler{
									HTTPGet: &core.HTTPGetAction{
										Port: intstr.FromInt(8086),
										Path: "/ready",
									},
								},
							},
							SecurityContext: self.DefaultSecurityContext(),
						},
					},
					Volumes: []core.Volume{
						{
							Name:         "cache",
							VolumeSource: self.VolumeSource("1Gi"),
						},
					},
				},
			},
		},
	}

	return self.CreateDeployment(deployment)
}
