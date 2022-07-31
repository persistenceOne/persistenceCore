# Local net

Scripts for running local net locally (tested for mac and linux).

## Prerequiest

For any of the types bellow, need to build the local clients

* From root dir `persistenceCore/`
* Run `make build`: This will create the binary at `./build/persistenceCore`

## In system

* Change directory into `cd contrib/local`
* `make`: Runs cleanup, `setup.sh` script, start node.
    * Initiate the gensis file
    * Override the config params
* Above command will be blocking, hence advised to run in a separate terminal
* Run test commands
    * `make run-cw20-govinit`: Create cw20 contract in wasm via a proposal, instantiate the contract via proposal as
      well, run test commands
* `make clean` cleanup

## In Docker

* Change directory into `cd contrib/local`
* `make docker-setup`: Pull the core image, and run `make` command, run in background, export ports to local host
* Then you can connect to the chain running in docker at localhost
* Run `make run-gov-contract` from `contrib/local` file on local system
* `make docker-clean`: Cleanup the containers after testing
* Optionally: `make docker-exec`: Open and bash shell into the background container, can run further test commands there

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

Commands we can run with permissionless wasm are

* `make run-contract`
* `make run-gov-contract`
* `make run-cw20-base`
