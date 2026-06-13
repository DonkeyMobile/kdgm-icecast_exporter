# Icecast exporter for Prometheus

A [Prometheus](https://prometheus.io/) exporter that scrapes stats from
[Icecast](http://icecast.org/) streaming media server via its JSON API
(`/status-json.xsl`, requires Icecast 2.4.0+).

By default icecast_exporter listens on port 9146 for HTTP requests.

## Metrics

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `icecast_up` | gauge | | 1 if Icecast is reachable |
| `icecast_server_start` | gauge | | Timestamp of server startup |
| `icecast_listeners` | gauge | listenurl, server_type | Currently connected listeners |
| `icecast_stream_start` | gauge | listenurl, server_type | Timestamp of active source connection |
| `icecast_exporter_scrape_errors_total` | counter | | Errors scraping Icecast |

## Running

```
docker run --rm -p 9146:9146 kdgm/icecast_exporter \
  -icecast.scrape-uri http://icecast:8000/status-json.xsl
```

### Flags

```
  -web.listen-address string
        Address to listen on for web interface and telemetry. (default ":9146")
  -web.telemetry-path string
        Path under which to expose metrics. (default "/metrics")
  -icecast.scrape-uri string
        URI on which to scrape Icecast. (default "http://localhost:8000/status-json.xsl")
  -icecast.timeout duration
        Timeout for trying to get stats from Icecast. (default 5s)
```
