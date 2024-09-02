#!/bin/bash

# shellcheck disable=SC1091
. "${ASSETS_DIR}/variables.env"

ASSETS_GENESIS_FILE="/assets/genesis.json"
CHAIN_GENESIS_FILE="${SHUTTER_CHAIN_DIR}/config/genesis.json"

rm "$CHAIN_GENESIS_FILE"
ln -s "$ASSETS_GENESIS_FILE" "$CHAIN_GENESIS_FILE"

SHUTTERMINT_MONIKER=${KEYPER_NAME:-$(openssl rand -hex 8)}

sed -i "/^seeds =/c\seeds = \"${_ASSETS_SHUTTERMINT_SEED_NODES}\"" "$SHUTTER_CHAIN_CONFIG_FILE"
sed -i "/^moniker =/c\moniker = \"${SHUTTERMINT_MONIKER}\"" "$SHUTTER_CHAIN_CONFIG_FILE"
sed -i "/^genesis_file =/c\genesis_file = \"${ASSETS_GENESIS_FILE}\"" "$SHUTTER_CHAIN_CONFIG_FILE"
sed -i "/^external_address =/c\external_address = \"${_DAPPNODE_GLOBAL_PUBLIC_IP}:${CHAIN_PORT}\"" "$SHUTTER_CHAIN_CONFIG_FILE"
sed -i "/^addr_book_strict =/c\addr_book_strict = true" "$SHUTTER_CHAIN_CONFIG_FILE"
sed -i "/^pex =/c\pex = true" "$SHUTTER_CHAIN_CONFIG_FILE"
if [ "$SHUTTER_PUSH_METRICS_ENABLED" = "true" ]; then
    sed -i "/^prometheus =/c\prometheus = true" "$SHUTTER_CHAIN_CONFIG_FILE"
    sed -i "/^prometheus_listen_addr =/c\prometheus_listen_addr = \"0.0.0.0:${PROMETHEUS_LISTEN_PORT}\"" "$SHUTTER_CHAIN_CONFIG_FILE"
fi
