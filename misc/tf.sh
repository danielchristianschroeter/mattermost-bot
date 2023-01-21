#!/bin/bash
# Terraform wrapper script
set -e -o pipefail
ACTION=$1
ENV=$2
ZONE=$3

if [ -z "$ACTION" ] || [ -z "$ENV" ] || [ -z "$ZONE" ]; then
    echo "Usage: tf <plan or apply> <env> <zone>"
    exit 1
fi

if [[ ! -d "/opt/terraform/environments/$ENV/$ZONE" ]]; then
    echo "/opt/terraform/environments/$ENV/$ZONE not found."
    exit 1
fi

if [[ "$ACTION" == "plan" ]]; then
  terraform -chdir=/opt/terraform/environments/$ENV/$ZONE plan -no-color >/dev/null -input=false -out=tfplan
  terraform -chdir=/opt/terraform/environments/$ENV/$ZONE show tfplan -no-color
  echo "To apply this plan, run: !tf apply $ENV $ZONE"
fi

if [[ "$ACTION" == "apply" ]]; then
  terraform -chdir=/opt/terraform/environments/$ENV/$ZONE apply -no-color -input=false -auto-approve tfplan
fi

