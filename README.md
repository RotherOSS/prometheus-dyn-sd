# prometheus-dyn-sd

A REST application written in Go that provides dynamic File-Based Service Discovery (File SD). It allows you to manage scrape targets on the fly using simple REST API endpoints.

Compatible with most telemetry backends that support that support Prometheus `file_sd_configs` protocol such as [Prometheus](https://prometheus.io/docs/guides/file-sd/), [Mimir](https://grafana.com/docs/mimir/latest/configure/about-configurations/), [Loki](https://grafana.com/docs/loki/latest/configure/) and [Grafana Alloy](https://grafana.com/docs/alloy/latest/).

## Getting started

The recommended way is to deploy the application using Docker Compose.

The example config below shows a simple configuration that exposes the api service and mounts the output sd file into the docker host.
```yaml
services:
    prometheus-dyn-sd:
        image: ghcr.io/splayfunityde/roonie:latest
        container_name: prometheus-dyn-sd
        restart: unless-stopped
        ports: 8010:8010
        environment:
          - PROMETHEUS_DYN_SD_FILEPATH=/opt/dynamic.json
        volumes:
          - ./data:/opt
```

> [!WARNING]  
> If you're opened port is exposed to the public it's recommended to secure your endpoints with some kind of authentication by using an external reverse proxy.

## Environment Variables

| Name  | Value | Required |
| ------------- | ------------- | ------------- |
| `PROMETHEUS_DYN_SD_FILEPATH` | Path to create the sd file | `Yes` |


## Exposed Endpoints

| URL  | Method | Description | Requires Body |
| ------------- | ------------- | ------------- | ------------- |
| `/hosts/{id}` | GET | Returns the body from the selected host | `No` |
| `/hosts/{id}` | PUT | Creates a new host | `Yes` |
| `/hosts/{id}` | PUT | Update an existing host | `Yes` |
| `/hosts/{id}` | DELETE | Delete an existing host | `No` |

The body for each required endpoint should have the following format.
```json
{
  "targets": [
    "host.local.yourdomain.de"
  ],
  "labels": {
    "webpage": "https://www.yourdomain.de",
    "mylabel": "custom"
  }
}
```
`hostname` represents the hostname which is used to ping the host. It can also be replaced with an IP address. Futhermore `labels` allows to add custom labels to the host for instance to monitor Web endpoints associated with the host.