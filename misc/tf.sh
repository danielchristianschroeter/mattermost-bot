#!/bin/bash
# Terraform wrapper script
set -e -o pipefail

ACTION="$1"
ENV="$2"
ZONE="$3"
TF_DIR="/opt/terraform/environments/$ENV/$ZONE"
LOG_DIR="/opt/terraform/log/"
DATE=$(date +'%Y%m%d-%H:%M:%S')
LOG_FILE="$LOG_DIR/${DATE}_terraform_apply.log"

if [ -z "$ACTION" ] || [ -z "$ENV" ] || [ -z "$ZONE" ]; then
    echo "Usage: tf <plan or apply> <env> <zone>"
    exit 1
fi

if [[ ! -d "$TF_DIR" ]]; then
    echo "$TF_DIR not found."
    exit 1
fi

if [[ "$ACTION" == "plan" ]]; then
  terraform -chdir="$TF_DIR" plan -no-color >/dev/null -input=false -out=tfplan \
    && terraform -chdir="$TF_DIR" show tfplan -no-color \
    && echo "To apply this plan, run: !tf apply $ENV $ZONE"
fi

if [[ "$ACTION" == "apply" ]]; then
  if [ -d "$LOG_DIR" ]; then
    if [ -f "$TF_DIR/tfplan" ]; then
      terraform -chdir="$TF_DIR" show -no-color tfplan | tee -a "$LOG_FILE" \
        && terraform -chdir="$TF_DIR" apply -no-color -input=false -auto-approve tfplan | tee -a "$LOG_FILE"
    else
      echo "Error: tfplan not found in $TF_DIR"
      exit 1
    fi
  else
    echo "Error: log directory $LOG_DIR not found"
    exit 1
  fi
fi
