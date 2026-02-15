# Temporal Production Docker Compose Setup

## Architecture

```
                    ┌──────────────┐
                    │  Your Apps   │
                    │  (Workers)   │
                    └──────┬───────┘
                           │ gRPC :7233
                    ┌──────▼───────┐
             ┌──────┤   Temporal   ├──────┐
             │      │   Server     │      │
             │      └──────┬───────┘      │
             │             │              │
      ┌──────▼──────┐  ┌──▼───────┐  ┌───▼──────────┐
      │ Temporal UI │  │PostgreSQL│  │  Prometheus   │
      │   :8233     │  │  :5432   │  │    :9091      │
      └─────────────┘  └──────────┘  └───┬──────────┘
                                         │
                                    ┌────▼─────┐
                                    │ Grafana  │
                                    │  :3000   │
                                    └──────────┘
```

## Quick Start

```bash
# 1. Copy and configure environment variables
cp .env.example .env
# Edit .env and change all passwords!

# 2. Start all services
docker compose up -d

# 3. Verify
docker compose ps
docker compose logs temporal --tail 50

# 4. Access services
#    Temporal Web UI:  http://localhost:8233
#    Grafana:          http://localhost:3000
#    Prometheus:       http://localhost:9091
#    Temporal gRPC:    localhost:7233
```

## Service Ports

| Service             | Port  | Purpose                    |
|---------------------|-------|----------------------------|
| Temporal Server     | 7233  | gRPC (workers connect here)|
| Temporal UI         | 8233  | Web dashboard              |
| PostgreSQL          | 5432  | Database                   |
| Prometheus          | 9091  | Metrics UI                 |
| Grafana             | 3000  | Dashboards                 |
| Temporal Metrics    | 9090  | Prometheus scrape endpoint |

## Connecting Your Workers

### Go Worker

```go
c, err := client.Dial(client.Options{
    HostPort: "localhost:7233",
    Namespace: "default",
})
```

### Python Worker

```python
client = await Client.connect("localhost:7233", namespace="default")
```

## Managing Namespaces

```bash
# Enter admin tools container
docker exec -it temporal-admin-tools bash

# Create a new namespace
temporal operator namespace create \
  --namespace my-app \
  --retention 168h

# List namespaces
temporal operator namespace list

# Describe a namespace
temporal operator namespace describe --namespace my-app
```

## Monitoring

### Key Temporal Metrics to Watch

| Metric                                    | What it tells you                    |
|-------------------------------------------|--------------------------------------|
| `temporal_workflow_completed`              | Workflow completion rate             |
| `temporal_workflow_failed`                 | Workflow failure rate                |
| `temporal_workflow_task_schedule_to_start` | Worker backlog (latency)             |
| `temporal_activity_schedule_to_start`      | Activity queue wait time             |
| `temporal_persistence_latency`             | Database response time               |
| `temporal_service_requests`                | API request rate                     |

### Import Temporal Grafana Dashboards

Temporal provides official Grafana dashboards:
- Server General: https://github.com/temporalio/dashboards
- Import them via Grafana UI → Dashboards → Import

## Production Hardening Checklist

- [ ] Change all default passwords in `.env`
- [ ] Enable TLS/mTLS for gRPC connections
- [ ] Use a managed PostgreSQL (RDS, Cloud SQL) for HA
- [ ] Set up PostgreSQL backups and point-in-time recovery
- [ ] Configure proper resource limits in docker-compose
- [ ] Set up alerting rules in Grafana
- [ ] Enable Temporal archival for compliance
- [ ] Restrict network access (firewall rules)
- [ ] Set up log aggregation
- [ ] Test disaster recovery procedures

## Scaling Beyond Single Node

When you outgrow this single-node setup, split Temporal into
separate services for independent scaling:

```yaml
# Replace the single temporal service with:
temporal-frontend:
  image: temporalio/server:latest
  environment:
    - SERVICES=frontend
    # ... same DB config
  deploy:
    replicas: 2

temporal-history:
  image: temporalio/server:latest
  environment:
    - SERVICES=history
  deploy:
    replicas: 3

temporal-matching:
  image: temporalio/server:latest
  environment:
    - SERVICES=matching
  deploy:
    replicas: 2

temporal-worker:
  image: temporalio/server:latest
  environment:
    - SERVICES=worker
  deploy:
    replicas: 1
```

## Adding Elasticsearch (Optional)

If you need advanced search across millions of workflows:

```yaml
elasticsearch:
  image: elasticsearch:8.13.0
  environment:
    - discovery.type=single-node
    - xpack.security.enabled=false
    - ES_JAVA_OPTS=-Xms512m -Xmx512m
  ports:
    - "9200:9200"
  volumes:
    - es_data:/usr/share/elasticsearch/data
```

Then update the Temporal environment:
```yaml
- ENABLE_ES=true
- ES_SEEDS=elasticsearch
- ES_PORT=9200
- ES_VERSION=v8
```

## Backup & Recovery

```bash
# Backup PostgreSQL
docker exec temporal-postgresql \
  pg_dumpall -U temporal > temporal_backup_$(date +%Y%m%d).sql

# Restore
cat temporal_backup_20250215.sql | \
  docker exec -i temporal-postgresql psql -U temporal
```

## Troubleshooting

```bash
# Check Temporal server logs
docker compose logs temporal -f

# Check database connectivity
docker exec temporal-admin-tools \
  temporal operator cluster health

# Check workflow status
docker exec temporal-admin-tools \
  temporal workflow list --namespace default

# Check system metrics
curl http://localhost:9090/metrics | head -50
```
