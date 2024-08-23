#!/bin/bash

get_imported_value_if_undefined() {
    local key=$1
    local file=$1

    # Eval key value to check if it is already defined
    local value=${!key}

    # If value is defined, return it
    if [ -n "$value" ]; then
        echo "$value"
        return
    fi

    grep -oP "^$key\s*=\s*\"\K[^\"]+" "$file"
}

get_imported_value_if_undefined "$1" "$2"
