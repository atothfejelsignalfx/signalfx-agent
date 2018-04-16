# Kubernetes Deployments

The agent runs in Kubernetes and has some monitors and observers specific to
that environment.  

The resources in this directory can be used to deploy the agent to K8s.  They
are generated from our [Helm](https://github.com/kubernetes/helm) chart,
which is available in the main [Helm Charts
repo](https://github.com/kubernetes/charts/tree/master/stable/signalfx-agent).

Make sure you change the `kubernetes_cluster` global dimension to something
specific to your cluster in the `configmap.yaml` resource before deploying.

## Development

These resources can be refreshed from the Helm chart by using the
`generate-from-helm` script in this dir.
