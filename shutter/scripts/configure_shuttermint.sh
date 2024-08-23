#!/bin/bash

# In case the config file has been imported
get_values_from_config_file() {

    #Check if the config file exists
    if [ ! -f "$SHUTTER_CHAIN_CONFIG_FILE" ]; then
        echo "[INFO | configure] Previous Shutter chain configuration file not found"
        return
    fi

    echo "[INFO | configure] Getting values from an existing shutter chain configuration file..."

    SHUTTERMINT_MONIKER=$(get_imported_value_if_undefined.sh "moniker" "$SHUTTER_CHAIN_CONFIG_FILE")
}

# shellcheck disable=SC1091
if [[ -f ${ASSETS_DIR}/variables.env ]]; then
    . "${ASSETS_DIR}/variables.env"
else
    echo "[ERROR | configure] Missing variables file (${ASSETS_DIR}/variables)"
    exit 1
fi

rm "${SHUTTER_CHAIN_DIR}/config/genesis.json"
ln -s /assets/genesis.json "${SHUTTER_CHAIN_DIR}/config/genesis.json"

SHUTTERMINT_MONIKER=${KEYPER_NAME:-$(openssl rand -hex 8)}

sed -i "/^seeds =/c\seeds = \"${_ASSETS_SHUTTERMINT_SEED_NODES}\"" "$SHUTTER_CHAIN_CONFIG_FILE"
sed -i "/^moniker =/c\moniker = \"${SHUTTERMINT_MONIKER}\"" "$SHUTTER_CHAIN_CONFIG_FILE"
sed -i "/^genesis_file =/c\genesis_file = \"/assets/genesis.json\"" "$SHUTTER_CHAIN_CONFIG_FILE"
sed -i "/^external_address =/c\external_address = \"${_DAPPNODE_GLOBAL_PUBLIC_IP}:${CHAIN_PORT}\"" "$SHUTTER_CHAIN_CONFIG_FILE"
sed -i "/^addr_book_strict =/c\addr_book_strict = true" "$SHUTTER_CHAIN_CONFIG_FILE"
sed -i "/^pex =/c\pex = true" "$SHUTTER_CHAIN_CONFIG_FILE"
if [ "$SHUTTER_PUSH_METRICS_ENABLED" = "true" ]; then
    sed -i "/^prometheus =/c\prometheus = true" "$SHUTTER_CHAIN_CONFIG_FILE"
    sed -i "/^prometheus_listen_addr =/c\prometheus_listen_addr = \"0.0.0.0:${PROMETHEUS_LISTEN_PORT}\"" "$SHUTTER_CHAIN_CONFIG_FILE"
fi
