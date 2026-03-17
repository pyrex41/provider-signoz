# Provider SigNoz

`provider-signoz` is a [Crossplane](https://crossplane.io/) provider for
[SigNoz](https://signoz.io), built using [Upjet](https://github.com/crossplane/upjet) code
generation tools. It wraps [`terraform-provider-signoz`](https://github.com/pyrex41/terraform-provider-signoz)
and exposes XRM-conformant managed resources for managing SigNoz dashboards, alerts, and notification channels via Kubernetes.

## Compatibility

| Provider Version | Terraform Provider | SigNoz Version | Notes |
|-----------------|-------------------|----------------|-------|
| v0.2.29 | 0.0.12-rc26 | 0.104.0 – 0.115.x | Current release; full resource support |
| v0.2.23 | 0.0.12-rc20 | 0.104.0 – 0.110.x | Added notification channels |
| < v0.2.10 | < 0.0.12-rc12 | < 0.104.0 | Legacy single-object API format only |

**Tested against:** SigNoz v0.110.1 (production), v0.115.0 (development)

## Managed Resources

Resources are available in both cluster-scoped and namespace-scoped variants:

| Resource | API Group | Description |
|----------|-----------|-------------|
| `Dashboard` | `dashboard.signoz.crossplane.io` | Manage SigNoz dashboards (v4 and v5 query formats) |
| `Alert` | `alert.signoz.crossplane.io` | Manage SigNoz alert rules |
| `NotificationChannel` | `notificationchannel.signoz.crossplane.io` | Manage notification channels (Slack, webhook) |

Namespace-scoped variants use the `.signoz.m.crossplane.io` API group suffix.

### v5 Dashboard Mutation Shielding

SigNoz's API mutates dashboard widget JSON during v5 migration (template variable rewriting, orderBy injection, operator changes). The provider uses `LateInitializer.IgnoredFields` to prevent Upjet from copying mutated values back into `spec.forProvider`, keeping your manifests as the source of truth. Shielded fields: `widgets`, `variables`, `layout`, `panel_map`, `version`.

## Getting Started

1. Install the provider in your Crossplane cluster
2. Create a `ProviderConfig` with your SigNoz endpoint and credentials
3. Apply managed resources to create dashboards, alerts, and notification channels

## Developing

Run code-generation pipeline:
```console
go run cmd/generator/main.go "$PWD"
```

Run against a Kubernetes cluster:

```console
make run
```

Build, push, and install:

```console
make all
```

Build binary:

```console
make build
```

## Report a Bug

For filing bugs, suggesting improvements, or requesting new features, please
open an [issue](https://github.com/facilitygrid/provider-signoz/issues).
