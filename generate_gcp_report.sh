#!/bin/bash

if [ -z "$1" ]; then
    echo "Error: Please provide a project ID as first parameter."
    echo "Usage: $0 PROJECT_ID"
    exit 1
fi

PROJECT_ID="$1"

echo "Generating report for GCP project $PROJECT_ID..."

#######
# GKE #
#######
GKE_CLUSTERS=$(gcloud container clusters list --format=json --project $PROJECT_ID)

#######
# GCS #
#######
GCS_BUCKETS_INFO=()
GCS_BUCKETS_IAM=()

BUCKETS=$(gcloud storage buckets list --project $PROJECT_ID --format=json)

for BUCKET_NAME in $(echo "${BUCKETS}" | jq -r '.[].name')
do
    BUCKET_URL="gs://${BUCKET_NAME}"
    BUCKET_INFO=$(gcloud storage buckets describe $BUCKET_URL --format=json)
    BUCKET_IAM=$(gcloud storage buckets get-iam-policy $BUCKET_URL --format=json)
    GCS_BUCKETS_INFO+=("$BUCKET_INFO")
    GCS_BUCKETS_IAM+=("{\"bucket\":\"$BUCKET_NAME\",\"iam_policy\":$BUCKET_IAM}")
done

GCS_BUCKETS_INFO_JOINED=$(IFS=,; echo "${GCS_BUCKETS_INFO[*]}")
GCS_BUCKETS_IAM_JOINED=$(IFS=,; echo "${GCS_BUCKETS_IAM[*]}")

#######
# IAM #
#######
IAM_ROLES=$(gcloud iam roles list --format=json --project $PROJECT_ID)
IAM_SERVICE_ACCOUNTS=$(gcloud iam service-accounts list --format=json --project $PROJECT_ID)
IAM_BINDINGS=$(gcloud projects get-iam-policy $PROJECT_ID --format=json)

#######
# GCR #
#######
ARTIFACT_REGISTRIES=$(gcloud artifacts repositories list --format=json --project $PROJECT_ID)

#######
# GCS #
#######
VM_INSTANCES=$(gcloud compute instances list --format=json --project $PROJECT_ID)

#######
# VPC #
#######
VPC_NETWORKS=$(gcloud compute networks list --format=json --project $PROJECT_ID)
VPC_SUBNETS=$(gcloud compute networks subnets list --format=json --project $PROJECT_ID)
VPC_FIREWALL_RULES=$(gcloud compute firewall-rules list --format=json --project $PROJECT_ID)

##########
# Output #
##########
cat <<EOF > magic_infra_task_gcp_report.json
{
  "gke": ${GKE_CLUSTERS:-null},
  "vm_instances": ${VM_INSTANCES:-null},
  "vpc_networks": ${VPC_NETWORKS:-null},
  "vpc_subnets": ${VPC_SUBNETS:-null},
  "vpc_firewall_rules": ${VPC_FIREWALL_RULES:-null},
  "gcs": {
    "buckets_info": [$GCS_BUCKETS_INFO_JOINED],
    "buckets_iam": [$GCS_BUCKETS_IAM_JOINED]
  },
  "iam_roles": ${IAM_ROLES:-null},
  "iam_service_accounts": ${IAM_SERVICE_ACCOUNTS:-null},
  "iam_bindings": ${IAM_BINDINGS:-null},
  "artifact_registries": ${ARTIFACT_REGISTRIES:-null}
}
EOF

echo "Report saved to 'magic_infra_task_gcp_report.json'"
