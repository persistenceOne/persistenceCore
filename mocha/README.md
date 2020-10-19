## Pre-requisites

```
Install npm and nodeJS
cd ~
curl -sL https://deb.nodesource.com/setup_10.x -o nodesource_setup.sh
sudo bash nodesource_setup.sh
sudo apt install nodejs

To check which version of Node.js you have installed after these initial steps, type:
node -v

For more information, visit https://www.digitalocean.com/community/tutorials/how-to-install-node-js-on-ubuntu-18-04.
```

* * *

## Installation

```

cd /persistenceOne/assetMantle/mocha

npm install

npm run test:awesome

NOTE: If any error comes which says: Error: Cannot find module 'xxx'
then run "npm install xxx --save"
```

* * *

## Documentation

For more information, visit https://autom8able.com.

* * *

## Testing

To test, go to the  go/src/github.com/persistenceOne/assetMantle/mocha folder and type :

    $ npm run test:awesome
   
* * *
 
## Report

[mochawesome] Report JSON saved to go/src/github.com/persistenceOne/assetMantle/mocha/mochawesome-report/mochawesome.json

[mochawesome] Report HTML saved to go/src/github.com/persistenceOne/assetMantle/mocha/mochawesome-report/mochawesome.html


* * *