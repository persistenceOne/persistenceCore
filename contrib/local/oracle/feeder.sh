#!/bin/bash
set -e errexit

DIR="$HOME"
mkdir -p $DIR

# Move into the repo directory
cd $DIR

## Clone the oracle-feeder repo.
git clone git@github.com:persistenceOne/oracle-feeder.git

# Move into the oracle feeder directory
cd $DIR/oracle-feeder

# Checkout specific branch and build the binary
git checkout tikaryan/add-mock-provider && go install

# initialize oracle-feeder configuration
touch $DIR/oracle-feeder/config.toml

# setup price-feeder configuration
tee $DIR/oracle-feeder/config.toml <<EOF
gas_adjustment = 1.5
fees = "100uxprt"

[server]
listen_addr = "0.0.0.0:7171"
read_timeout = "20s"
verbose_cors = true
write_timeout = "20s"

[[deviation_thresholds]]
base = "DUMMY"
threshold = "1.5"

[[currency_pairs]]
base = "DUMMY"
providers = [
  "mock",
]
quote = "USD"

[account]
address = "persistence1fk57xyyxz89krc3nn5law2xekan84gufzagszw"
chain_id = "testing"
validator = "persistencevaloper1fk57xyyxz89krc3nn5law2xekan84guftegdth"

[keyring]
backend = "test"
dir = "/tmp/trash/.persistenceCore"

[rpc]
grpc_endpoint = "localhost:9090"
rpc_timeout = "100ms"
tmrpc_endpoint = "http://localhost:26657"
EOF

# start price-feeder
echo "###Start the oracle price feeder"
export PRICE_FEEDER_KEY_PASS="test"

# start oracle-feeder in background process and redirect output to log file in current directory.
oracle-feeder $DIR/oracle-feeder/config.toml > /dev/null