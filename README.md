[![Build](https://github.com/jidicula/cloudflare-cachebuster/actions/workflows/build.yml/badge.svg)](https://github.com/jidicula/cloudflare-cachebuster/actions/workflows/build.yml) [![Latest Release](https://github.com/jidicula/cloudflare-cachebuster/actions/workflows/release-draft.yml/badge.svg)](https://github.com/jidicula/cloudflare-cachebuster/actions/workflows/release-draft.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/jidicula/cloudflare-cachebuster)](https://goreportcard.com/report/github.com/jidicula/cloudflare-cachebuster) [![Go Reference](https://pkg.go.dev/badge/github.com/jidicula/cloudflare-cachebuster.svg)](https://pkg.go.dev/github.com/jidicula/cloudflare-cachebuster)

# cloudflare-cachebuster

Azure Function that calls Cloudflares [`purge_cache`](https://developers.cloudflare.com/api/operations/zone-purge) endpoint.

Requires 2 environment variables to be set:
* `CLOUDFLARE_ZONEID`: the Zone ID of the cache being purged.
* `CLOUDFLARE_PAT`: A Cloudflare PAT with `Zone.Cache Purge` permissions for the cache's zone.

## Deploying to Azure via GitHub Actions
Requires:
* `AZURE_RBAC_CREDENTIALS`: secret containing a [service principal's RBAC credentials](https://github.com/marketplace/actions/azure-functions-action#using-azure-service-principal-for-rbac-as-deployment-credential).
* `AZURE_FUNCTIONAPP_NAME`: variable containing the Function App's name.
