#!/bin/bash

if [ -z "$1" ]; then
    echo "Error: Please provide the k8s context name as first parameter."
    echo "Usage: $0 CONTEXT_NAME"
    exit 1
fi

K8S_CONTEXT=$1
OUTPUT_FILE="magic_infra_task_k8s_report.json"

echo "{" > $OUTPUT_FILE

NAMESPACES=$(kubectl --context $K8S_CONTEXT get ns -o jsonpath="{.items[*].metadata.name}")

for namespace in $NAMESPACES
do
    if [ "$namespace" == "kube-system" ]; then
        continue
    fi

    echo "\"$namespace\": " >> $OUTPUT_FILE
    kubectl --context $K8S_CONTEXT get all -o json -n $namespace >> $OUTPUT_FILE
    echo "," >> $OUTPUT_FILE
done

if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i '' '$ s/,$//' $OUTPUT_FILE
else
    sed -i '$ s/,$//' $OUTPUT_FILE
fi

echo "}" >> $OUTPUT_FILE
