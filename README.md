*This is a proof-of-concept*

Planter
=======

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Latest Release](https://img.shields.io/github/release/tliron/planter.svg)](https://github.com/tliron/planter/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/tliron/planter)](https://goreportcard.com/report/github.com/tliron/planter)

Planter is a meta-scheduler for Kubernetes.

It is an operator that runs in management clusters and delegates workloads to one or more workload
clusters while also allowing the workloads to (re)configure those very workload clusters on which
they are to be deployed.

Planter is intended as a deliberately narrow solution to the ["bifurcation"](https://www.youtube.com/watch?v=6FULuWvXR84)
(or "chicken-and-egg") problem that arises from vertically-integrated workloads. It enables an
all-at-once fire-and-forget declarative frontend for the complex lifecycle management of the
workload-and-its-cluster interrelationships.

Note that we assume the existence of management clusters. These are Kubernetes clusters in which
we run the software used to manage the lifecycle of many workload clusters and even the workloads
themselves via policies, rules, and more direct operations. Which software, exactly, depends on
the Kubernetes platform provider. Running the management software in a Kubernetes cluster enables an
extensible, declarative approach to management, and allows Planter, as a Kubernetes operator, to
participate in this management layer. Indeed, Planter's scope is intentionally limited to solving
one specific challenge and it will depend on other provider-specific pieces (operators?) for handling
workload cluster installation and connectivity.

Planter is itself extensible by design: you can plug in your own heuristic logic to handle the custom
interralationships specific to your workload, cloud platform, and their unique resource types.

The name "planter" is a reference to
[this kind of farming equipment](https://en.wikipedia.org/wiki/Planter_(farm_implement)).

Rationale
---------

The reality of clouds, at least at the time of this writing, is that they are stacked platforms.
You have your hardware, which is configured via firmware, an operating system installed on top, and
then it's running various workload controllers, which orchestrate the local workloads.

But it's also a reality that many workloads, especially in the telco industry, are vertically
integrated. To deploy them you need to interact with all parts of the stack, from the firmware
through operating system kernels and drivers to the configuration of the cluster controllers.

Kubernetes adds a further separation between administrative configuration tasks (creating
custom resource definitions) and user tasks (creating custom resources), requiring us to order
our operations and break away from an all-at-once declarative approach.

How do we do it all? More specifically, how can we package a workload that can do it all?

One recent answer is GitOps. But GitOps, by itself, is a declarative illusion. Putting all your
various configurations in one place sure is convenient, but if each configuration unit is handled
by a different controller at a different time then the interrelationships are emergent rather than
declarative. GitOps doesn't "solve" the problem, if anything it reveals it. Planter is intended
as a straightforward solution to the next step after bundling together the configurations:
dispersing them.

How It Works
------------

Planter's core resource type is a "seed", which is a package (YAML file) combining all Kubernetes
resources relating to a workload at its Day 1 initialization. The seed includes all the familiar
Kubernetes resources—Deployments, DaemonSets, ServiceAccounts, Services, Persistent Volume Claims,
ConfigMaps, etc.—as well resources belonging to the management domain, such as cluster and node
configurations, for example to configure operating system kernels, RAM layouts (NUMA pages),
GPU/DPUs, SR-IOV, etc.

Planter then reconciles the seed with two goals:

1) **Placement**. Planter figures out which resource belongs to which cluster. The main differentiator
   is management vs. workload clusters, but multi-cluster workloads are supported, including
   deployments across multiple management and workload clusters.
2) **Dependency derivation**. Planter discovers the *hard temporal dependencies* of resources on other
   resources. The most obvious example is that resources cannot be deployed before the cluster in which
   they were placed exists. Other obvious ones is that namespaced resources cannot be declared before
   the namespace is created and custom resources cannot be declared before their custom resource definition
   is created. "Soft" dependencies (e.g. a Service not being functional until its selected Pods are ready)
   are ignored by Planter, as they are handled by Kubernetes controllers and custom operators.

Reconciliation employs built-in heuristics for common resource types, but your own logic can be plugged
in for custom resource types and custom dependencies. Plugins are simple executables that input the
seed, annotate resources as needed, and output the changes. They can be written in any language, e.g.
Python, Go, etc.

Planter will then schedule the seed, which simply means: declaring the resources at the right place
and time. It does not "deploy" these resources—that work is done by other controllers and operators.
It merely declares the intent. A successful planting, then, simply means that the resources have all
been declared where and when they are supposed to be declared. Whether or not the workload works at all
or properly is out of scope.

Note that "right time" is explicitly *not* about processing a classic
[directed acyclic graph](https://en.wikipedia.org/wiki/Directed_acyclic_graph). This is merely about
postpoing operations until a requirement is met. There are no success/failure forks and definitely
no conditional orchestration. Planter assumes and encourages cloud native, fully declarative design.
If you absolutely need DAG-like behavior, consider incorporating a Kubernetes workflow engine like
[Argo Workflows](https://argoproj.github.io/argo-workflows/).

Also note the emphasis on seeds being Day 1 designs—which is exactly why we are calling them
"seeds". The idea is that once planted the workload will take care of its own Day 2+ maintenance,
possibly via operators deployed by the seed itself. It is possible to use Planter to replant a seed,
but it should be understood as a kind of going back (soft reset) to initial conditions, if not to
initial state (hard reset). Ongoing management is out of scope: again, it's on you to design your
workloads based on cloud native principles.

FAQ
---

### Is Planter a replacement for Helm? Is a Planter seed like a Helm chart?

[Helm](https://helm.sh/) charts do not solve the chicken-and-egg problem. Some Helm-based solutions
indeed use multiple Helm charts that are meant to be deployed in order. But you might have to wait a
*very long time* between installing them, e.g. if one Helm chart causes a cluster to be installed it
can take hours before the cluster is ready. Unlike Planter, Helm is not a "fire and forget"
orchestration tool.

You could potentially use Helm to generate the Planter seed, via the `helm template` command.
