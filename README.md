# databricks-runs-exporter
Prometheus Exporter for databricks runs.

## Installation
```bash
# Create the credentials secret 
kubectl create secret generic databricks-creds -n monitoring \ 
--from-literal=host=<databricks_host>  --from-literal=token=<databricks_token>
# Apply the exporter
kubectl apply -f https://github.com/Efrat19/databricks-runs-exporter/blob/main/k8s.yaml
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
| RUNS_SCRAPE_LIMIT             | 50      | Max runs to scrape on each interval. increase only if you have more then RUNS_SCRAPE_LIMIT runs created in a given RUNS_SCRAPE_TIMESPAN_SECONDS
| RUNS_SCRAPE_TIMESPAN_SECONDS  | 10      | How much seconds ago to scrape, match this number with your scrape interval (usually 10s). The exporter will also take a safety margin of 10 more seconds to ensure no runs are lost due to unexpected delays.

## Grafana Dahsboard

Coming soon...
## Credits
Awesome exporter template from [go-gywn/query-exporter-simple](https://github.com/go-gywn/query-exporter-simple.git)
