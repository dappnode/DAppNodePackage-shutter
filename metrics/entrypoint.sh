#!/bin/sh

update_user_settings() {

    keys="KEYPER_NAME PUSHGATEWAY_URL PUSHGATEWAY_USERNAME PUSHGATEWAY_PASSWORD"

    if [ ! -f "$USER_SETTINGS_FILE" ]; then
        touch "$USER_SETTINGS_FILE"
    fi

    for key in $keys; do
        update_user_setting "$key" "$(eval echo \$"$key")"
    done
}

update_user_setting() {
    key=$1
    value=$(eval echo \$"$key")

    if [ -z "$value" ]; then
        echo "[INFO | metrics] Skipped updating $key in user settings file (empty value)"
        return 1
    fi

    if grep -q "^$key=" "$USER_SETTINGS_FILE"; then
        # Update the existing key
        sed -i "s|^$key=.*|$key=$value|" "$USER_SETTINGS_FILE"
        echo "[INFO | metrics] Updated $key in $USER_SETTINGS_FILE"
    else
        # Add the new key
        echo "$key=$value" >>"$USER_SETTINGS_FILE"
        echo "[INFO | metrics] Added $key to $USER_SETTINGS_FILE"
    fi
}

source_envs() {
    set -a # Export all variables

    # shellcheck disable=SC1091
    if [ -f "${ASSETS_DIR}/variables.env" ]; then
        . "${ASSETS_DIR}/variables.env"

    else
        echo "[ERROR | configure] Missing variables file (${ASSETS_DIR}/variables.env)"
        exit 1
    fi

    if [ -z "${_ASSETS_VERSION:-}" ]; then
        _ASSETS_VERSION="$(cat /assets/version)"
    fi

    # shellcheck disable=SC1090
    . "$USER_SETTINGS_FILE"

    set +a
}

replace_envs_in_yaml() {
    echo "[INFO | metrics] Replacing environment variables in the configuration file"
    sed "s|%{KEYPER_NAME}|$KEYPER_NAME|g; s|%{_ASSETS_VERSION}|$_ASSETS_VERSION|g" "$TEMPLATE_CONFIG_FILE" >"$CONFIG_FILE"
}

update_user_settings

if [ "${SHUTTER_PUSH_METRICS_ENABLED}" = "false" ]; then
    echo "[INFO | metrics] Metrics push is disabled"
    exit 0
fi

source_envs

replace_envs_in_yaml

exec /vmagent-prod \
    -promscrape.config="${CONFIG_FILE}" \
    -remoteWrite.url="${PUSHGATEWAY_URL}" \
    -remoteWrite.basicAuth.username="${PUSHGATEWAY_USERNAME}" \
    -remoteWrite.basicAuth.password="${PUSHGATEWAY_PASSWORD}"
