#!/bin/bash
set -e

echo "Initiating connection handshake..."
hermes -c ./config.toml create connection test-1 test-2

sleep 2
