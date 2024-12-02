#!/bin/bash

# To use staker scripts
# shellcheck disable=SC1091
. /etc/profile
# shellcheck disable=SC1091
. "${ASSETS_DIR}/variables.env"

echo "[INFO | configure] Calculating keyper configuration values..."

SUPPORTED_NETWORKS="gnosis"

export SHUTTER_P2P_ADVERTISEADDRESSES=["/ip4/${_DAPPNODE_GLOBAL_PUBLIC_IP}/tcp/${KEYPER_PORT}"]
export SHUTTER_BEACONAPIURL=http://${_DAPPNODE_GLOBAL_CONSENSUS_CLIENT_GNOSIS}
export SHUTTER_GNOSIS_NODE_CONTRACTSURL=ws://${_DAPPNODE_GLOBAL_EXECUTION_CLIENT_GNOSIS}
export SHUTTER_GNOSIS_NODE_ETHEREUMURL=ws://${_DAPPNODE_GLOBAL_EXECUTION_CLIENT_GNOSIS}
export VALIDATOR_PUBLIC_KEY=$(cat "${SHUTTER_CHAIN_DIR}/config/priv_validator_pubkey.hex")
export SHUTTER_METRICS_ENABLED=${SHUTTER_PUSH_METRICS_ENABLED}

echo "[INFO | configure] LISTEN: $SHUTTER_P2P_LISTENADDRESSES"

# Check if the keyper configuration file already exists
if [ -f "$KEYPER_CONFIG_FILE" ]; then
    echo "[INFO | configure] Keyper configuration file already exists"
else
    echo "[INFO | configure] Generating configuration files..."

    if [ ! -f "$KEYPER_GENERATED_CONFIG_FILE" ]; then
        echo "[ERROR | configure] Missing generated configuration file (${KEYPER_GENERATED_CONFIG_FILE})"
        exit 1
    fi

    cp "$KEYPER_GENERATED_CONFIG_FILE" "$KEYPER_CONFIG_FILE"
fi

go_shutter_settings --generated "$KEYPER_GENERATED_CONFIG_FILE" --config "$KEYPER_CONFIG_FILE" --output "$KEYPER_CONFIG_FILE" include-keyper-settings
