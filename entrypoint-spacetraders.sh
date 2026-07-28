#!/usr/bin/bash

set -eu pipefail

SPACETRADERS_ACCOUNT_TOKEN=$(aws secretsmanager get-secret-value --secret-id spacetraders/account-token --query SecretString --output text)
#SPACETRADERS_AGENT_TOKEN=$(aws secretsmanager get-secret-value --secret-id MyTestSecret --query SecretString --output text)

GRAFANA_CLOUD_TOKEN=$(aws secretsmanager get-secret-value --secret-id spacetraders/grafana-cloud-write-token --query SecretString --output text)

git pull
./pull-state.sh
go run . TVRJ
