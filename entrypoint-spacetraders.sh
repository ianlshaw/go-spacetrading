#!/usr/bin/bash

set -eu pipefail

SPACETRADERS_ACCOUNT_TOKEN=$(aws secretsmanager get-secret-value --secret-id spacetraders/account-token --query SecretString --output text)
#SPACETRADERS_AGENT_TOKEN=$(aws secretsmanager get-secret-value --secret-id MyTestSecret --query SecretString --output text)

GRAFANA_CLOUD_TOKEN=$(aws secretsmanager get-secret-value --secret-id spacetraders/grafana-cloud-write-token --query SecretString --output text)
GRAFANA_METRICS_ENDPOINT=https://prometheus-prod-55-prod-gb-south-1.grafana.net/api/prom/push
GRAFANA_METRICS_UERNAME=2920303
GRAFANA_LOGS_ENDPOINT=https://logs-prod-035.grafana.net
GRAFANA_LOGS_USERNAME=1455824

git pull
./pull-state.sh
go run . TVRJ
