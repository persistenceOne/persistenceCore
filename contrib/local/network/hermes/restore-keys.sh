#!/bin/bash
set -e

# Load shell variables

### Sleep is needed otherwise the relayer crashes when trying to init
sleep 1s
### Restore Keys
$HERMES_BINARY -c ./hermes/config.toml keys restore test-1 -m "marble allow december print trial know resource cry next segment twice nose because steel omit confirm hair extend shrimp seminar one minor phone deputy"
sleep 5s

$HERMES_BINARY -c ./hermes/config.toml keys restore test-2 -m "record gift you once hip style during joke field prize dust unique length more pencil transfer quit train device arrive energy sort steak upset"
sleep 5s
