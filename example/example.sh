#!/usr/bin/env bash
#env GOOS=linux GOARCH=amd64 go build #for linux amd64 systems

generateAlerts() {

    envsubst "$DATA" < monitors/monitoring.yaml >  generated.yaml

cat generated.yaml
  EXCLUDE=$(echo $EXCLUDE | tr "|" "\n")

  for ITEM in $EXCLUDE
  do
    export ITEM=$ITEM
    envsubst '$ITEM' < monitors/exclude.yaml >>  excluded.yaml
    echo '' >> excluded.yaml
  done
  if [  -f "excluded.yaml" ]; then
    GENERATED_TEXT=`cat generated.yaml`
    EXCLUDE_TEXT=`cat excluded.yaml`
    echo "${GENERATED_TEXT//##EXCLUDE##/$EXCLUDE_TEXT}" > generated.yaml
    rm excluded.yaml
  fi
   ## for debug purpose only
    #cat generated.yaml > $QUERY.yaml
}


temp_file="/es-alert-cli"
install_dir="/usr/local/bin"

apk add curl
curl -sL "https://github.com/Trendyol/es-alert-cli/releases/download/0.5.0/es-alert-cli" -o "/usr/local/bin/es-alert-cli"
mv "$temp_file" "$install_dir/"
chmod +x "${install_dir}/$(basename "$temp_file")"

LOGGING_CLUSTER_URL="my-logging-01-kibana-url.com"
LOGGING_CLUSTER_IP="http://XX.XX.XXX.XXX:9200" #elastic ip
LOGGING_CLUSTER_INDEX_PATTERN_ID="a3959da0-19c8-11bd-8c51-6f32cab50115"

while IFS=, read -r QUERY COUNT MONITOR ALERT LEVEL SEVERITY EXCLUDE CLUSTER APPNAME KIBANA_TRIGGER_NAME; do
  if [[ "$QUERY" == "QUERY" ]]; then continue; fi
  export QUERY=$QUERY
  export COUNT=$COUNT
  export MONITOR=$MONITOR
  export ALERT=$ALERT
  export LEVEL=$LEVEL
  export SEVERITY=$SEVERITY
  export CLUSTER=$CLUSTER
  export APPNAME=$APPNAME
  export KIBANA_TRIGGER_NAME=$KIBANA_TRIGGER_NAME
  export GENERATED_KIBANA_URL=$LOGGING_CLUSTER_URL
  export INDEX_PATTERN_ID
  export DC="DC1"


  DATA='$QUERY:$COUNT:$MONITOR:$ALERT:$LEVEL:$SEVERITY:$APPNAME:$CLUSTER_URL'

if [[ $CLUSTER == *"dc1"* ]]; then
    CLUSTER=$LOGGING_CLUSTER_URL
    INDEX_PATTERN_ID=$LOGGING_CLUSTER_INDEX_PATTERN_ID
    DATA='$QUERY:$COUNT:$MONITOR:$ALERT:$LEVEL:$SEVERITY:$APPNAME:$CLUSTER:$INDEX_PATTERN_ID:$DC'
    generateAlerts

    es-alert-cli upsert -c $LOGGING_CLUSTER_IP -n generated.yaml
fi

done < projects.csv
