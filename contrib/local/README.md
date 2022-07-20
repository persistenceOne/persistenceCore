# Local net

Scripts for running local net locally

## In system
* First build the binaries with `make build` in root directory.
* Change directory into `cd contrib/local`
* `make`: Runs cleanup, `setup.sh` script, start node.
  * Initiate the gensis file
  * Override the config params
* Above command will be blocking, hence advised to run in a separate terminal
* Run test commands
  * `make run-gov-contract`: Create proposal via proposal, vote on the proposal, and initiate the contract, run test commands
* `make clean` cleanup

## In Docker
* `make docker-setup`: Pull the core image, and run `make` command, run in background
* `make docker-exec`: Open and bash shell into the background container, can run further test commands there
  * `make run-gov-contract` inside the container
* `make docker-clean`: Cleanup the containers after testing

## Permissionless wasm
By default the chain runs with wasm as a permissioned module where all the contracts
are uploaded via a gov proposal. Inorder to start the chain as permissionless use
following
```bash
# system
WASM_PERMISSIONLESS=true make clean setup
## or
WASM_PERMISSIONLESS=true make

# for docker
WASM_PERMISSIONLESS=true make docker-setup
```

With this we start the container such that the wasm module runns in permissionless fashion. For testing,
run `make run-contract` for permissionless contract testing.
