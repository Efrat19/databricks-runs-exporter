package main

import "github.com/prometheus/client_golang/prometheus"

// =============================
// Config config structure
// =============================
type Config struct {
	Metrics map[string]struct {
		Type        string
		Description string
		Labels      []string
		Value       string
		metricDesc  *prometheus.Desc
	}
}

// =============================
// QueryCollector exporter
// =============================
type QueryCollector struct{}


type rawRun struct {
	JobID                int `json:"job_id"`
	RunID                int `json:"run_id"`
	NumberInJob          int `json:"number_in_job"`
	OriginalAttemptRunID int `json:"original_attempt_run_id"`
	State                struct {
		LifeCycleState          string `json:"life_cycle_state"`
		ResultState             string `json:"result_state"`
		StateMessage            string `json:"state_message"`
		UserCancelledOrTimedout bool   `json:"user_cancelled_or_timedout"`
	} `json:"state"`
	Task struct {
		NotebookTask struct {
			NotebookPath   string `json:"notebook_path"`
			BaseParameters struct {
				SQLFilePath      string `json:"sql_file_path"`
				SecretPath       string `json:"secret_path"`
				CheckpointS3Path string `json:"checkpoint_s3_path"`
				S3BucketPath     string `json:"s3_bucket_path"`
				BootstrapServers string `json:"bootstrap_servers"`
				CategoryID       string `json:"category_id"`
				Topic            string `json:"topic"`
			} `json:"base_parameters"`
		} `json:"notebook_task"`
	} `json:"task"`
	ClusterSpec struct {
		ExistingClusterID string `json:"existing_cluster_id"`
	} `json:"cluster_spec"`
	ClusterInstance struct {
		ClusterID      string `json:"cluster_id"`
		SparkContextID string `json:"spark_context_id"`
	} `json:"cluster_instance"`
	StartTime         int64  `json:"start_time"`
	SetupDuration     int    `json:"setup_duration"`
	ExecutionDuration int    `json:"execution_duration"`
	CleanupDuration   int    `json:"cleanup_duration"`
	EndTime           int64  `json:"end_time"`
	Trigger           string `json:"trigger"`
	CreatorUserName   string `json:"creator_user_name"`
	RunName           string `json:"run_name"`
	RunPageURL        string `json:"run_page_url"`
	RunType           string `json:"run_type"`
	AttemptNumber     int    `json:"attempt_number"`
	Format            string `json:"format"`
}

type Run struct {
	JobID                   int    `json:"job_id"`
	RunID                   int    `json:"run_id"`
	NumberInJob             int    `json:"number_in_job"`
	OriginalAttemptRunID    int    `json:"original_attempt_run_id"`
	LifeCycleState          string `json:"life_cycle_state"`
	ResultState             string `json:"result_state"`
	StateMessage            string `json:"state_message"`
	UserCancelledOrTimedout bool   `json:"user_cancelled_or_timedout"`
	StartTime               int64  `json:"start_time"`
	SetupDuration           int    `json:"setup_duration"`
	ExecutionDuration       int    `json:"execution_duration"`
	CleanupDuration         int    `json:"cleanup_duration"`
	EndTime                 int64  `json:"end_time"`
	Trigger                 string `json:"trigger"`
	CreatorUserName         string `json:"creator_user_name"`
	RunName                 string `json:"run_name"`
	RunType                 string `json:"run_type"`
	AttemptNumber           int    `json:"attempt_number"`
	Format                  string `json:"format"`
}

type ApiResponse struct {
	Runs    []rawRun `json:"runs"`
	HasMore bool     `json:"has_more"`
}
