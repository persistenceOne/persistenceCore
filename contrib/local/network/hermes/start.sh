#!/bin/bash

# Load shell variables
. ./network/hermes/variables.sh

# Start the hermes relayer in multi-paths mode
echo "Starting hermes relayer..."
hermes -c ./config.toml start
