#!/usr/bin/bash

set -eu pipefail

export CALLSIGN=TVRJ

export AWS_REGION=eu-west-2

export SPACETRADERS_ACCOUNT_TOKEN=$(aws secretsmanager get-secret-value --secret-id spacetraders/account-token --query SecretString --output text)
export SPACETRADERS_AGENT_TOKEN=$(aws secretsmanager get-secret-value --secret-id spacetraders/agent-token --query SecretString --output text)
export GRAFANA_CLOUD_TOKEN=$(aws secretsmanager get-secret-value --secret-id spacetraders/grafana-cloud-write-token --query SecretString --output text)
export GRAFANA_METRICS_ENDPOINT=https://prometheus-prod-55-prod-gb-south-1.grafana.net/api/prom/push
export GRAFANA_METRICS_UERNAME=2920303
export GRAFANA_LOGS_ENDPOINT=https://logs-prod-035.grafana.net
export GRAFANA_LOGS_USERNAME=1455824

git pull
./pull-state.sh
go run .
