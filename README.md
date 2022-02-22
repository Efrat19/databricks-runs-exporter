# databricks-runs-exporter
Prometheus Exporter for databricks runs.

## Installation
```bash

# Create the credentials secret 
kubectl create secret generic databricks-creds -n monitoring \ 
--from-literal=host=<databricks_host>  --from-literal=token=<databricks_token>

# Apply the exporter
kubectl apply -f https://raw.githubusercontent.com/Efrat19/databricks-runs-exporter/main/k8s.yaml
```
## Metrics
### `runs_total`
- **Type:** Count
- **Namespace:** databricks_exporter
- **Description:** runs info from databricks
- **Labels:** job_id, run_id, number_in_job, original_attempt_run_id, life_cycle_state, result_state, user_cancelled_or_timedout, start_time, setup_duration, execution_duration, cleanup_duration, end_time, trigger, creator_user_name, run_name, run_type, attempt_number, format

## Env Variables

| Variable                      | Default | Description |
| ----------------------------- | ------- | -----------------------------
| DATABRICKS_HOST               | ""      | Databricks host URL, required
| DATABRICKS_TOKEN              | ""      | Databricks token with permissions to list jobs, required
| RUNS_SCRAPE_LIMIT             | 50      | Max runs to scrape on each interval. Increase if you have more then `RUNS_SCRAPE_LIMIT` runs created in a given `RUNS_SCRAPE_TIMESPAN_SECONDS`
| RUNS_SCRAPE_TIMESPAN_SECONDS  | 10      | How much seconds ago to look for created jobs. Match this value to your scrape interval (which is usually 10s)

## Grafana Dahsboard

[dashboard.json](https://raw.githubusercontent.com/Efrat19/databricks-runs-exporter/main/dashboard.json)

## Credits
Awesome exporter template from [go-gywn/query-exporter-simple](https://github.com/go-gywn/query-exporter-simple.git)
