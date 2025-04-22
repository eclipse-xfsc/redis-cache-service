# cache

## Prerequisites

- Kubernetes 1.19+
- Helm 3+

## Install Helm Chart

```console
helm install [RELEASE_NAME] -f values.yaml .
```

_See [configuration](#configuration) below._

_See [helm install](https://helm.sh/docs/helm/helm_install/) for command documentation._

## Dependencies

By default, dependencies are not included in the application/service's Helm chart. Please install dependencies  separately using their respective vendor Helm charts. The dependencies that have to be installed manually are:

- [nats](https://nats-io.github.io/k8s/helm/charts/)
- [redis](https://github.com/bitnami/charts/tree/main/bitnami/redis)
To disable dependencies during installation, see [multiple releases](#multiple-releases) below.

_See [helm dependency](https://helm.sh/docs/helm/helm_dependency/) for command documentation._

## Uninstall Helm Chart

```console
helm uninstall [RELEASE_NAME]
```

This removes all the Kubernetes components associated with the chart and deletes the release.

_See [helm uninstall](https://helm.sh/docs/helm/helm_uninstall/) for command documentation._

## Upgrading Chart

```console
# Upgrade chart to development environment values
helm upgrade -f values-dev.yaml [RELEASE_NAME] .
```

## Configuration

See [Customizing the Chart Before Installing](https://helm.sh/docs/intro/using_helm/#customizing-the-chart-before-installing). To see all configurable options with detailed comments:

```console
helm show values .
```

### Secrets

The following Helm template snippet demonstrates how to use Kubernetes secrets in your Helm charts:

```yaml
{{- range $key, $value := .Values.secretEnv }}
- name: "{{ $key }}"
  valueFrom:
    secretRef:
      name: "{{ $value.name }}"
{{- end }}
```
This template iterates over a list of secret environment variables specified in the Helm chart values file (values.yaml). For each secret environment variable defined, it retrieves the corresponding value from a Kubernetes Secret and injects it into the container as an environment variable.

#### How to use

To utilize this template in your Helm charts, follow these steps:

    Define Secret Environment Variables: In your Helm chart values.yaml file, define a list of secret environment variables along with the names of the corresponding Kubernetes Secret objects. For example:

```yaml
secretEnv:
  DB_PASSWORD: 
    name: my-db-secret
  API_KEY:
    name: api-key-secret
```
In this example, DB_PASSWORD and API_KEY are the environment variable names, and my-db-secret and api-key-secret are the names of the Kubernetes Secret objects containing the respective values.

Update Helm Template: Update your Helm chart template file (e.g., deployment.yaml) to include the provided template snippet. This will dynamically fetch the values of the secret environment variables from the specified Kubernetes Secrets and inject them into your container.

Deploy Helm Chart: Deploy your Helm chart to your Kubernetes cluster using the helm install command.

Edit Helm Chart Values: Open the values.yaml file in your Helm chart and configure the following values:

    useConfigMap: Set to true to use a ConfigMap, or false otherwise.
    configMapName: Name of the ConfigMap containing environment variables.
    useSecretRef: Set to true to use a Secret, or false otherwise.
    secretRefName: Name of the Secret containing sensitive environment variables.
    
### Istio intergration

Optional Istio integration is done the following way in the `values.yaml` file:

```yaml
istio:
  injection:
    pod: true
```

You may also `helm show values` on this chart's [dependencies](#dependencies) for additional options
## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| autoscaling.enabled | bool | `false` | Enable autoscaling |
| autoscaling.maxReplicas | int | `3` | Maximum replicas |
| autoscaling.minReplicas | int | `1` | Minimum replicas |
| autoscaling.targetCPUUtilizationPercentage | int | `70` | CPU target for autoscaling trigger |
| autoscaling.targetMemoryUtilizationPercentage | int | `70` | Memory target for autoscaling trigger |
| cache.http.host | string | `""` | Host for the cache HTTP service |
| cache.http.port | int | `8080` | Port for the cache HTTP service |
| cache.http.timeout.idle | string | `"120s"` | Timeout duration for idle connections in the cache HTTP service |
| cache.http.timeout.read | string | `"10s"` | Timeout duration for read operations in the cache HTTP service |
| cache.http.timeout.write | string | `"10s"` | Timeout duration for write operations in the cache HTTP service |
| cache.nats.subject | string | `"external_cache_events"` | NATS subject for cache events |
| cache.nats.url | string | `"nats:4222"` | URL for connecting to NATS server |
| image.name | string | `"gaiax/cache"` | Name of the cache image |
| image.pullPolicy | string | `"IfNotPresent"` | Pull policy for the cache image |
| image.pullSecrets | string | `"deployment-key-light"` | Secret used to pull the cache image |
| image.repository | string | `"eu.gcr.io/vrgn-infra-prj"` | Repository for the cache image |
| image.sha | string | `""` | SHA checksum of the cache image |
| image.tag | string | `""` | Tag of the cache image |
| ingress.annotations."kubernetes.io/ingress.class" | string | `"nginx"` | Ingress class annotation for cache |
| ingress.annotations."nginx.ingress.kubernetes.io/rewrite-target" | string | `"/$2"` | Rewrite target annotation for cache |
| ingress.enabled | bool | `true` | Enable ingress for cache |
| ingress.frontendDomain | string | `"tsa.xfsc.dev"` | Domain for cache frontend |
| ingress.frontendTlsSecretName | string | `"cert-manager-tls"` | Secret name for TLS certificate |
| ingress.tlsEnabled | bool | `true` | Enable TLS for cache ingress |
| log.encoding | string | `"json"` | Encoding format for logs |
| log.level | string | `"debug"` | Log level for cache |
| metrics.enabled | bool | `true` | Enable Prometheus metrics for cache |
| metrics.port | int | `2112` | Port for Prometheus metrics |
| name | string | `"cache"` | Name of the cache application |
| nameOverride | string | `""` | Override for application name |
| podAnnotations | object | `{}` | Annotations for cache pods |
| redis.addr | string | `"redis-master:6379"` | Address of the Redis server |
| redis.db | int | `0` | Redis database number |
| redis.expiration | string | `"1h"` | Expiration duration for cache entries |
| redis.pass.secretName | string | `"redis-pass"` | Secret name for Redis password |
| redis.user.secretName | string | `"redis-user"` | Secret name for Redis username |
| replicaCount | int | `1` | Number of cache instances to deploy |
| resources.limits.cpu | string | `"150m"` | CPU limit for cache pods |
| resources.limits.memory | string | `"128Mi"` | Memory limit for cache pods |
| resources.requests.cpu | string | `"25m"` | CPU request for cache pods |
| resources.requests.memory | string | `"64Mi"` | Memory request for cache pods |
| security.runAsGid | int | `0` | Group ID used by the cache application |
| security.runAsNonRoot | bool | `false` | Run cache application as non-root |
| security.runAsUid | int | `0` | User ID used by the cache application |
| service.port | int | `8080` | Port for cache service |

----------------------------------------------
Autogenerated from chart metadata using [helm-docs v1.12.0](https://github.com/norwoodj/helm-docs/releases/v1.12.0)
