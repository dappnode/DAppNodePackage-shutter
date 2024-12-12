#!/bin/bash

# To use staker scripts
# shellcheck disable=SC1091
. /etc/profile

SUPPORTED_NETWORKS="gnosis"

# Set environment variables before running the health checks.
echo "[INFO | chain] Setting environment variables for beacon API and execution API..."
export SHUTTER_BEACONAPIURL=$(get_beacon_api_url_from_global_env "$NETWORK" "$SUPPORTED_NETWORKS")
export SHUTTER_GNOSIS_NODE_CONTRACTSURL=http://execution.gnosis.dncore.dappnode:8545

perform_node_healthcheck() {
    echo "[INFO | chain] Starting health check for beacon API and execution API..."

    while true; do
       # Check the syncing status using JSON-RPC
        echo "[INFO | chain] Checking execution API syncing status at: $SHUTTER_GNOSIS_NODE_CONTRACTSURL/eth/v1/node/syncing"
        sync_status=$(curl -s -X POST "$SHUTTER_GNOSIS_NODE_CONTRACTSURL" \
            -H "Content-Type: application/json" \
            -d '{"jsonrpc":"2.0","method":"eth_syncing","params":[],"id":1}' | jq -r '.result')

        if [ "$sync_status" != "false" ]; then
            echo "[WARN | chain] Execution API is syncing. Sync status: $sync_status. Retrying in 30 seconds..."
            sleep 30
            continue
        else
            echo "[INFO | chain] Execution API is not syncing."
        fi
        
        # Check the syncing status using JSON-RPC
        echo "[INFO | chain] Checking execution API syncing status at: $SHUTTER_GNOSIS_NODE_CONTRACTSURL/eth/v1/node/syncing"
        
        # Perform the request and capture the response
        response=$(curl -s -w "%{http_code}" -o /tmp/sync_response.json -X POST "$SHUTTER_GNOSIS_NODE_CONTRACTSURL" \
            -H "Content-Type: application/json" \
            -d '{"jsonrpc":"2.0","method":"eth_syncing","params":[],"id":1}')
        
        # Extract HTTP status code and response body
        http_code="${response: -3}"
        response_body=$(cat /tmp/sync_response.json)

        # Check if the response body is empty or the HTTP status code is not 200
        if [ -z "$response_body" ] || [ "$http_code" -ne 200 ]; then
            echo "[ERROR | chain] No response or error from execution API. Not healthy"
            sleep 30
            continue
        fi

        # Parse the JSON response to check the syncing status
        syncing_status=$(echo "$response_body" | jq -r '.result.is_syncing')

        # If syncing is true or an error occurs, treat it as not synced
        if [ "$syncing_status" = "true" ]; then
            echo "[WARN | chain] Execution API is syncing. Sync status: $syncing_status. Retrying in 30 seconds..."
            sleep 30
            continue
        else
            echo "[INFO | chain] Execution API is healthy."
        fi

        # If we reach this point, all checks have passed.
        echo "[INFO | chain] All services are healthy. Exiting health check loop."
        break
    done
}

run_chain() {
    echo "[INFO | chain] Starting chain process..."
    echo "[DEBUG | chain] Command: $SHUTTER_BIN chain --config \"$SHUTTER_CHAIN_CONFIG_FILE\""
    $SHUTTER_BIN chain --config "$SHUTTER_CHAIN_CONFIG_FILE"
}

# Run health checks first
perform_node_healthcheck

# If everything is healthy, run the chain
run_chain
