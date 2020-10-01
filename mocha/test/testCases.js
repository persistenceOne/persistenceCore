"use strict";

//During the test the env variable is set to test
process.env.NODE_ENV = 'test';
let chai = require('chai');
let chaiHttp = require('chai-http');
let async = require('async');
let should = chai.should();
let expect = chai.expect
const crypto = require("crypto");
const { request } = require('chai')
var assert = chai.assert

var config = require('./config.json');
var temp = require('./helpers/helpers');
var helper = new temp()
const { type } = require('os');
const retry = require('async-retry')
const fetch = require('node-fetch')

var server = config.ip + config.port
var retry_count = config.tries_threshold

chai.use(chaiHttp);

const name1 = crypto.randomBytes(16).toString("hex");
const name2 = crypto.randomBytes(16).toString("hex");

var names1 = {
    address: '',
    typekey: '',
    value: '',
    signature: ''
}

var names2 = {
    address: '',
    typekey: '',
    value: '',
    signature: ''
}

var txHash = ''

////////////////////////////////////////////////////////////////////////////////////////////////
/*
Account Creation for User1
*/
////////////////////////////////////////////////////////////////////////////////////////////////

describe('Account Creation, Sign and Broadcast for User1', async () => {

    beforeEach(function (done) {
        this.timeout(4000)
        setTimeout(function () {
            done()
        }, 3000)
    })

    it('it should create an account for: ' + name1, async () => {

        let name = {
            name: name1
        }

        var err, res = await chai.request(server)
            .post(config.keysAdd)
            .send(name)

        res.should.have.status(200);
        res.body.should.be.a('object');

        expect(res.body.result.success).to.be.true
        expect(res.body.result.keyOutput.address).to.not.equal(null)
        expect(res.body.result.keyOutput.address).to.not.equal('')

        names1.address = res.body.result.keyOutput.address

    });

    it(name1 + ' should be able to signTx ', async () => {

        let obj = {
            "baseReq": {
                "from": config.testAccountAddress,
                "chain_id": config.chain_id
            },
            "type": "cosmos-sdk/StdTx",
            "value": {
                "msg": [
                    {
                        "type": "cosmos-sdk/MsgSend",
                        "value": {
                            "from_address": config.testAccountAddress,
                            "to_address": names1.address,
                            "amount": [
                                {
                                    "denom": "stake",
                                    "amount": "1000"
                                }
                            ]
                        }
                    }
                ],
                "fee": {
                    "amount": [],
                    "gas": "200000"
                },
                "signatures": null,
                "memo": ""
            }
        }

        var err, res = await chai.request(server)
            .post(config.signTx)
            .send(obj)

        res.should.have.status(200);
        res.body.should.be.a('object');
        expect(res.body.result.success).to.be.true

        names1.typekey = res.body.result.tx.signatures[0].pub_key.type
        names1.value = res.body.result.tx.signatures[0].pub_key.value
        names1.signature = res.body.result.tx.signatures[0].signature

    });

    it(name1 + ' should be able to broadcastTx ', async () => {

        let obj = {
            "tx": {
                "msg": [
                    {
                        "type": "cosmos-sdk/MsgSend",
                        "value": {
                            "from_address": config.testAccountAddress,
                            "to_address": names1.address,
                            "amount": [
                                {
                                    "denom": "stake",
                                    "amount": "1000"
                                }
                            ]
                        }
                    }
                ],
                "fee": {
                    "amount": [],
                    "gas": "200000"
                },
                "signatures": [
                    {
                        "pub_key": {
                            "type": names1.typekey,
                            "value": names1.value
                        },
                        "signature": names1.signature
                    }
                ],
                "memo": ""
            },
            "mode": "sync"
        }


        var err, res = await chai.request(server)
            .post(config.broadcastTx)
            .send(obj)



        res.should.have.status(200);
        res.body.should.be.a('object');
        expect(res.body.txhash).to.not.equal(null)
        expect(res.body.txhash).to.not.equal('')

        var hash = res.body.txhash

        var err, res = await chai.request(server)
            .get('/txs/' + hash)
    });
});

////////////////////////////////////////////////////////////////////////////////////////////////
// /*
// Account Creation for User2
// */
////////////////////////////////////////////////////////////////////////////////////////////////

describe('Account Creation, Sign and Broadcast for User2', async () => {

    beforeEach(function (done) {
        this.timeout(4000)
        setTimeout(function () {
            done()
        }, 3000)
    })

    it('it should create an account for: ' + name2, async () => {

        let name = {
            name: name2
        }

        var err, res = await chai.request('http://localhost:1317')
            .post(config.keysAdd)
            .send(name)



        res.should.have.status(200);
        res.body.should.be.a('object');

        expect(res.body.result.success).to.be.true
        expect(res.body.result.keyOutput.address).to.not.equal(null)
        expect(res.body.result.keyOutput.address).to.not.equal('')

        names2.address = res.body.result.keyOutput.address
    });

    it(name2 + ' should be able to signTx ', async () => {

        let obj = {
            "baseReq": {
                "from": config.testAccountAddress,
                "chain_id": config.chain_id
            },
            "type": "cosmos-sdk/StdTx",
            "value": {
                "msg": [
                    {
                        "type": "cosmos-sdk/MsgSend",
                        "value": {
                            "from_address": config.testAccountAddress,
                            "to_address": names2.address,
                            "amount": [
                                {
                                    "denom": "stake",
                                    "amount": "1000"
                                }
                            ]
                        }
                    }
                ],
                "fee": {
                    "amount": [],
                    "gas": "200000"
                },
                "signatures": null,
                "memo": ""
            }
        }

        var err, res = await chai.request(server)
            .post(config.signTx)
            .send(obj)



        res.should.have.status(200);
        res.body.should.be.a('object');
        expect(res.body.result.success).to.be.true

        names2.typekey = res.body.result.tx.signatures[0].pub_key.type
        names2.value = res.body.result.tx.signatures[0].pub_key.value
        names2.signature = res.body.result.tx.signatures[0].signature

    });

    it(name2 + ' should be able to broadcastTx ', async () => {

        let obj = {
            "tx": {
                "msg": [
                    {
                        "type": "cosmos-sdk/MsgSend",
                        "value": {
                            "from_address": config.testAccountAddress,
                            "to_address": names2.address,
                            "amount": [
                                {
                                    "denom": "stake",
                                    "amount": "1000"
                                }
                            ]
                        }
                    }
                ],
                "fee": {
                    "amount": [],
                    "gas": "200000"
                },
                "signatures": [
                    {
                        "pub_key": {
                            "type": names2.typekey,
                            "value": names2.value
                        },
                        "signature": names2.signature
                    }
                ],
                "memo": ""
            },
            "mode": "sync"
        }


        var err, res = await chai.request(server)
            .post(config.broadcastTx)
            .send(obj)



        res.should.have.status(200);
        res.body.should.be.a('object');
        expect(res.body.txhash).to.not.equal(null)
        expect(res.body.txhash).to.not.equal('')

        var hash = res.body.txhash

        var err, res = await chai.request(server)
            .get('/txs/' + hash)


    });
});

////////////////////////////////////////////////////////////////////////////////////////////////
/*
Nub TX
*/
////////////////////////////////////////////////////////////////////////////////////////////////

describe('Nub Tx', async () => {

    beforeEach(function (done) {
        this.timeout(4000)
        setTimeout(function () {
            done()
        }, 3000)
    })

    it('nubTx: ', async () => {

        let obj = {
            "type": config.nubType,
            "value": {
                "baseReq": {
                    "from": config.testAccountAddress,
                    "chain_id": config.chain_id
                },
                "nubID": config.nubID
            }
        }

        var err, res = await chai.request(server)
            .post(config.nubPath)
            .send(obj)

        res.should.have.status(200);
        res.body.should.be.a('object');
        expect(res.body.txhash).to.not.equal(null)
        expect(res.body.txhash).to.not.equal('')

        txHash = res.body.txhash

    });

    it('Query Tx: ', async () => {
        var err, res = await chai.request(server)
            .get('/txs/' + txHash)
        var data1 = JSON.stringify(res.body)

        async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

            if ((data1.indexOf('error') != -1)) {
                callbackretry('failed')
            } else {
                //continue
            }
        }, function (err, response) {
            if (err) {
                //do nothing
            } else {
                expect(res.body.raw_log).to.not.contain('failed')
                expect(res.body).to.not.contain('error')
            }
        })
    });
});

////////////////////////////////////////////////////////////////////////////////////////////////
/*
Issue Identity
*/
////////////////////////////////////////////////////////////////////////////////////////////////

describe('Identity', async () => {
    describe('Issue Identity 1', async () => {

        beforeEach(function (done) {
            this.timeout(4000)
            setTimeout(function () {
                done()
            }, 3000)
        })


        var data = {
            'clasificationID': '',
            'hashID': '',
            'chainID': '',
            'clshashID': ''
        }


        it('Query Identity: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qIdentity)

            var list = res.body.result.value.identities.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == config.nubID) {
                    data.clasificationID = value.value.id.value.classificationID.value.idString
                    data.hashID = value.value.id.value.hashID.value.idString

                }

            });

        });

        it('Define Identity: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.defineIdentityType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": data.clasificationID + '|' + data.hashID,
                    "mutableTraits": "mutableTraits1:S|num1",
                    "immutableTraits": "immutableTraits1:S|",
                    "mutableMetaTraits": "mutableMetaTraits1:S|num3",
                    "immutableMetaTraits": "immutableMetaTraits1:S|num4"
                }
            }


            var err, res = await chai.request(server)
                .post(config.defineIdentityPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });

        it('Query Classification: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qClassification)

            var list = res.body.result.value.classifications.value.list
            list.forEach(function (value) {

                if (value.value.immutableTraits.value.properties.value.propertyList[0].value.id.value.idString == "immutableMetaTraits1") {
                    data.chainID = value.value.id.value.chainID.value.idString
                    data.clshashID = value.value.id.value.hashID.value.idString
                }

            });

        });

        it('Issue Identity: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.issuedentityType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "to": config.testAccountAddress,
                    "fromID": data.clasificationID + '|' + data.hashID,
                    "classificationID": data.chainID + '.' + data.clshashID,
                    "mutableProperties": "mutableTraits1:S|num1",
                    "immutableProperties": "immutableTraits1:S|",
                    "mutableMetaProperties": "mutableMetaTraits1:S|num3",
                    "immutableMetaProperties": "immutableMetaTraits1:S|num4"
                }
            }


            var err, res = await chai.request(server)
                .post(config.issueIdentityPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });
    });

    describe('Provision Key', async () => {

        beforeEach(function (done) {
            this.timeout(4000)
            setTimeout(function () {
                done()
            }, 3000)
        })


        var data = {
            'clasificationID': '',
            'hashID': '',
            'chainID': '',
            'clshashID': ''
        }


        it('Query Identity: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qIdentity)

            var list = res.body.result.value.identities.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == config.nubID) {
                    data.clasificationID = value.value.id.value.classificationID.value.idString
                    data.hashID = value.value.id.value.hashID.value.idString

                }

            });

        });

        it('Provision Key: ', async () => {

            let obj = {
                "type": config.provisionKeyType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "to": names1.address,
                    "identityID": data.clasificationID + '|' + data.hashID
                }
            }


            var err, res = await chai.request(server)
                .post(config.provisionKeyPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });
    });

    describe('Unprovision Key', async () => {

        beforeEach(function (done) {
            this.timeout(4000)
            setTimeout(function () {
                done()
            }, 3000)
        })


        var data = {
            'clasificationID': '',
            'hashID': '',
            'chainID': '',
            'clshashID': ''
        }


        it('Query Identity: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qIdentity)

            var list = res.body.result.value.identities.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == config.nubID) {
                    data.clasificationID = value.value.id.value.classificationID.value.idString
                    data.hashID = value.value.id.value.hashID.value.idString

                }

            });

        });

        it('Unprovision Key: ', async () => {

            let obj = {
                "type": config.unprovisionKeyType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "to": names1.address,
                    "identityID": data.clasificationID + '|' + data.hashID
                }
            }


            var err, res = await chai.request(server)
                .post(config.unprovisionKeyPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });
    });

    describe('Provision an unprovision Key', async () => {

        beforeEach(function (done) {
            this.timeout(4000)
            setTimeout(function () {
                done()
            }, 3000)
        })


        var data = {
            'clasificationID': '',
            'hashID': '',
            'chainID': '',
            'clshashID': ''
        }

        it('Query Identity: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qIdentity)

            var list = res.body.result.value.identities.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == config.nubID) {
                    data.clasificationID = value.value.id.value.classificationID.value.idString
                    data.hashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Provision an unprovision Key: ', async () => {

            let obj = {
                "type": config.provisionKeyType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "to": names1.address,
                    "identityID": data.clasificationID + '|' + data.hashID
                }
            }


            var err, res = await chai.request(server)
                .post(config.provisionKeyPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });
    });

    describe('Issue Identity 2', async () => {

        beforeEach(function (done) {
            this.timeout(4000)
            setTimeout(function () {
                done()
            }, 3000)
        })


        var data = {
            'clasificationID': '',
            'hashID': '',
            'chainID': '',
            'clshashID': ''
        }


        it('Query Identity: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qIdentity)

            var list = res.body.result.value.identities.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == config.nubID) {
                    data.clasificationID = value.value.id.value.classificationID.value.idString
                    data.hashID = value.value.id.value.hashID.value.idString

                }

            });

        });

        it('Define Identity: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.defineIdentityType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": data.clasificationID + '|' + data.hashID,
                    "mutableTraits": "mutableTraits2:S|num1",
                    "immutableTraits": "immutableTraits2:S|",
                    "mutableMetaTraits": "mutableMetaTraits2:S|num3",
                    "immutableMetaTraits": "immutableMetaTraits2:S|num4"
                }
            }


            var err, res = await chai.request(server)
                .post(config.defineIdentityPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });

        it('Query Classification: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qClassification)

            var list = res.body.result.value.classifications.value.list
            list.forEach(function (value) {

                if (value.value.immutableTraits.value.properties.value.propertyList[0].value.id.value.idString == "immutableMetaTraits2") {
                    data.chainID = value.value.id.value.chainID.value.idString
                    data.clshashID = value.value.id.value.hashID.value.idString
                }

            });

        });

        it('Issue Identity 2: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.issuedentityType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "to": config.testAccountAddress,
                    "fromID": data.clasificationID + '|' + data.hashID,
                    "classificationID": data.chainID + '.' + data.clshashID,
                    "mutableProperties": "mutableTraits2:S|num1",
                    "immutableProperties": "immutableTraits2:S|",
                    "mutableMetaProperties": "mutableMetaTraits2:S|num3",
                    "immutableMetaProperties": "immutableMetaTraits2:S|num4"
                }
            }


            var err, res = await chai.request(server)
                .post(config.issueIdentityPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });
    });

    describe('Issue Identity 3', async () => {

        beforeEach(function (done) {
            this.timeout(4000)
            setTimeout(function () {
                done()
            }, 3000)
        })


        var data = {
            'clasificationID': '',
            'hashID': '',
            'chainID': '',
            'clshashID': ''
        }


        it('Query Identity: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qIdentity)

            var list = res.body.result.value.identities.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == config.nubID) {
                    data.clasificationID = value.value.id.value.classificationID.value.idString
                    data.hashID = value.value.id.value.hashID.value.idString

                }

            });

        });

        it('Define Identity: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.defineIdentityType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": data.clasificationID + '|' + data.hashID,
                    "mutableTraits": "mutableTraits3:S|num1",
                    "immutableTraits": "immutableTraits3:S|",
                    "mutableMetaTraits": "mutableMetaTraits3:S|num3",
                    "immutableMetaTraits": "immutableMetaTraits3:S|num4"
                }
            }


            var err, res = await chai.request(server)
                .post(config.defineIdentityPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });

        it('Query Classification: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qClassification)

            var list = res.body.result.value.classifications.value.list
            list.forEach(function (value) {

                if (value.value.immutableTraits.value.properties.value.propertyList[0].value.id.value.idString == "immutableMetaTraits3") {
                    data.chainID = value.value.id.value.chainID.value.idString
                    data.clshashID = value.value.id.value.hashID.value.idString
                }

            });

        });

        it('Issue Identity 3: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.issuedentityType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "to": config.testAccountAddress,
                    "fromID": data.clasificationID + '|' + data.hashID,
                    "classificationID": data.chainID + '.' + data.clshashID,
                    "mutableProperties": "mutableTraits3:S|num1",
                    "immutableProperties": "immutableTraits3:S|",
                    "mutableMetaProperties": "mutableMetaTraits3:S|num3",
                    "immutableMetaProperties": "immutableMetaTraits3:S|num4"
                }
            }


            var err, res = await chai.request(server)
                .post(config.issueIdentityPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });
    });

})

////////////////////////////////////////////////////////////////////////////////////////////////
/*
Mint Asset
*/
////////////////////////////////////////////////////////////////////////////////////////////////


describe('Assets', async () => {

    describe('Mint Asset', async () => {

        beforeEach(function (done) {
            this.timeout(4000)
            setTimeout(function () {
                done()
            }, 3000)
        })


        var data = {
            'clasificationID': '',
            'hashID': '',
            'chainID': '',
            'clshashID': ''
        }

        it('Query Identity: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qIdentity)

            var list = res.body.result.value.identities.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "immutableMetaTraits1") {
                    data.clasificationID = value.value.id.value.classificationID.value.idString
                    data.hashID = value.value.id.value.hashID.value.idString
                }

            });

        });

        it('Define Asset: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.defineAssetType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": data.clasificationID + '|' + data.hashID,
                    "mutableTraits": "ASSET1:S|num1" + ",burn:H|1",
                    "immutableTraits": "ASSET2:S|",
                    "mutableMetaTraits": "ASSET3:S|num3",
                    "immutableMetaTraits": "ASSET4:S|num4"
                }
            }


            var err, res = await chai.request(server)
                .post(config.defineAssetPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });

        it('Query Classification: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qClassification)

            var list = res.body.result.value.classifications.value.list
            list.forEach(function (value) {

                if (value.value.immutableTraits.value.properties.value.propertyList[0].value.id.value.idString == "ASSET4") {
                    data.chainID = value.value.id.value.chainID.value.idString
                    data.clshashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Mint Asset: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.mintAssetType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "toID": data.clasificationID + '|' + data.hashID,
                    "fromID": data.clasificationID + '|' + data.hashID,
                    "classificationID": data.chainID + '.' + data.clshashID,
                    "mutableProperties": "ASSET1:S|num1" + ",burn:H|1",
                    "immutableProperties": "ASSET2:S|num2",
                    "mutableMetaProperties": "ASSET3:S|num3",
                    "immutableMetaProperties": "ASSET4:S|num4"

                }
            }


            var err, res = await chai.request(server)
                .post(config.mintAssetPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });
    });

    describe('Mutate Asset', async () => {

        beforeEach(function (done) {
            this.timeout(4000)
            setTimeout(function () {
                done()
            }, 3000)
        })


        var data = {
            'iclasificationID': '',
            'ihashID': '',
            'aclasificationID': '',
            'ahashID': ''
        }

        it('Query Identity: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qIdentity)

            var list = res.body.result.value.identities.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "immutableMetaTraits1") {
                    data.iclasificationID = value.value.id.value.classificationID.value.idString
                    data.ihashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Query Asset: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qAsset)

            var list = res.body.result.value.assets.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "ASSET4") {
                    data.aclasificationID = value.value.id.value.classificationID.value.idString
                    data.ahashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Mutate Asset: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.mutateAssetType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": data.iclasificationID + '|' + data.ihashID,
                    "assetID": data.aclasificationID + '|' + data.ahashID,
                    "mutableProperties": "ASSET1:S|",
                    "mutableMetaProperties": "ASSET3:S|num3"
                }
            }


            var err, res = await chai.request(server)
                .post(config.mutateAssetPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });
    });

    describe('Mint Asset with meta properties', async () => {

        beforeEach(function (done) {
            this.timeout(4000)
            setTimeout(function () {
                done()
            }, 3000)
        })


        var data = {
            'clasificationID': '',
            'hashID': '',
            'chainID': '',
            'clshashID': ''
        }

        it('Query Identity: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qIdentity)

            var list = res.body.result.value.identities.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "immutableMetaTraits1") {
                    data.clasificationID = value.value.id.value.classificationID.value.idString
                    data.hashID = value.value.id.value.hashID.value.idString
                }

            });

        });

        it('Query Classification: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qClassification)

            var list = res.body.result.value.classifications.value.list
            list.forEach(function (value) {

                if (value.value.immutableTraits.value.properties.value.propertyList[0].value.id.value.idString == "ASSET4") {
                    data.chainID = value.value.id.value.chainID.value.idString
                    data.clshashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Mint Asset with meta properties: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.mintAssetType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "toID": data.clasificationID + '|' + data.hashID,
                    "fromID": data.clasificationID + '|' + data.hashID,
                    "classificationID": data.chainID + '.' + data.clshashID,
                    "mutableProperties": "ASSET1:S|num5" + ",burn:H|1",
                    "immutableProperties": "ASSET2:S|num6",
                    "mutableMetaProperties": "ASSET3:S|num7",
                    "immutableMetaProperties": "ASSET4:S|num8"

                }
            }


            var err, res = await chai.request(server)
                .post(config.mintAssetPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });
    });

    describe('Mutate asset non meta properteies to meta properties', async () => {

        beforeEach(function (done) {
            this.timeout(4000)
            setTimeout(function () {
                done()
            }, 3000)
        })


        var data = {
            'iclasificationID': '',
            'ihashID': '',
            'chainID': ',',
            'clshashID': '',
            'aclasificationID': '',
            'ahashID': ''
        }

        it('Query Identity: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qIdentity)

            var list = res.body.result.value.identities.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "immutableMetaTraits1") {
                    data.iclasificationID = value.value.id.value.classificationID.value.idString
                    data.ihashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Define Asset: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.defineAssetType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": data.clasificationID + '|' + data.hashID,
                    "mutableTraits": "ASSETS1:S|num1" + ",burn:H|1",
                    "immutableTraits": "ASSETS2:S|",
                    "mutableMetaTraits": "ASSETS3:S|num3",
                    "immutableMetaTraits": "ASSETS4:S|num4"
                }
            }


            var err, res = await chai.request(server)
                .post(config.defineAssetPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });

        it('Query Classification: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qClassification)

            var list = res.body.result.value.classifications.value.list
            list.forEach(function (value) {

                if (value.value.immutableTraits.value.properties.value.propertyList[0].value.id.value.idString == "ASSETS4") {
                    data.chainID = value.value.id.value.chainID.value.idString
                    data.clshashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Mint Asset: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.mintAssetType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "toID": data.clasificationID + '|' + data.hashID,
                    "fromID": data.clasificationID + '|' + data.hashID,
                    "classificationID": data.chainID + '.' + data.clshashID,
                    "mutableProperties": "ASSETS1:S|num5" + ",burn:H|1",
                    "immutableProperties": "ASSETS2:S|num6",
                    "mutableMetaProperties": "ASSETS3:S|num7",
                    "immutableMetaProperties": "ASSETS4:S|num8"

                }
            }


            var err, res = await chai.request(server)
                .post(config.mintAssetPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });

        it('Query Asset: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qAsset)

            var list = res.body.result.value.assets.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "ASSETS4") {
                    data.aclasificationID = value.value.id.value.classificationID.value.idString
                    data.ahashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Meta Reveal: ', async () => {

            let obj = {
                "type": config.metaRevealType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "metaFact": "S|num5"
                }
            }

            var err, res = await chai.request(server)
                .post(config.metaRevealPath)
                .send(obj)

            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash
        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });

        it('Mutate Asset non meta properties to meta properties: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.mutateAssetType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": data.iclasificationID + '|' + data.ihashID,
                    "assetID": data.aclasificationID + '|' + data.ahashID,
                    "mutableProperties": "ASSETS1:S|",
                    "mutableMetaProperties": "ASSETS3:S|num3"
                }
            }


            var err, res = await chai.request(server)
                .post(config.mutateAssetPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash
        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });
    });

    describe('Mint Asset with 22 properties', async () => {

        beforeEach(function (done) {
            this.timeout(4000)
            setTimeout(function () {
                done()
            }, 3000)
        })


        var data = {
            'clasificationID': '',
            'hashID': '',
            'chainID': '',
            'clshashID': ''
        }

        it('Query Identity: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qIdentity)

            var list = res.body.result.value.identities.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "immutableMetaTraits1") {
                    data.clasificationID = value.value.id.value.classificationID.value.idString
                    data.hashID = value.value.id.value.hashID.value.idString
                }

            });

        });

        it('Define Asset: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.defineAssetType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": data.clasificationID + '|' + data.hashID,
                    "mutableTraits": "ASSETP1:S|A,ASSETP11:S|B,ASSETP12:S|C,ASSETP13:S|D,ASSETP14:S|E,burn:H|1",
                    "immutableTraits": "ASSETP2:S|G,ASSETP21:S|H,ASSETP22:S|I,ASSETP23:S|J,ASSETP24:S|K",
                    "mutableMetaTraits": "ASSETP3:S|L,ASSETP31:S|M,ASSETP32:S|N,ASSETP33:S|O,ASSETP34:S|P",
                    "immutableMetaTraits": "ASSETP4:S|Q,ASSETP41:S|R,ASSETP42:S|S,ASSETP43:S|T,ASSETP44:S|U,ASSETP45:S|V"
                }
            }


            var err, res = await chai.request(server)
                .post(config.defineAssetPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });

        it('Query Classification: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qClassification)

            var list = res.body.result.value.classifications.value.list
            list.forEach(function (value) {

                if (value.value.immutableTraits.value.properties.value.propertyList[0].value.id.value.idString == "ASSETP4") {
                    data.chainID = value.value.id.value.chainID.value.idString
                    data.clshashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Mint Asset: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.mintAssetType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "toID": data.clasificationID + '|' + data.hashID,
                    "fromID": data.clasificationID + '|' + data.hashID,
                    "classificationID": data.chainID + '.' + data.clshashID,
                    "mutableProperties": "ASSETP1:S|A,ASSETP11:S|B,ASSETP12:S|C,ASSETP13:S|D,ASSETP14:S|E,burn:H|1",
                    "immutableProperties": "ASSETP2:S|G,ASSETP21:S|H,ASSETP22:S|I,ASSETP23:S|J,ASSETP24:S|K",
                    "mutableMetaProperties": "ASSETP3:S|L,ASSETP31:S|M,ASSETP32:S|N,ASSETP33:S|O,ASSETP34:S|P",
                    "immutableMetaProperties": "ASSETP4:S|Q,ASSETP41:S|R,ASSETP42:S|S,ASSETP43:S|T,ASSETP44:S|U,ASSETP45:S|V"

                }
            }
            var err, res = await chai.request(server)
                .post(config.mintAssetPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });
    });

    describe('Add asset properties on mutation', async () => {

        beforeEach(function (done) {
            this.timeout(4000)
            setTimeout(function () {
                done()
            }, 3000)
        })


        var data = {
            'iclasificationID': '',
            'ihashID': '',
            'chainID': ',',
            'clshashID': '',
            'aclasificationID': '',
            'ahashID': ''
        }

        it('Query Identity: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qIdentity)

            var list = res.body.result.value.identities.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "immutableMetaTraits1") {
                    data.iclasificationID = value.value.id.value.classificationID.value.idString
                    data.ihashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Define Asset: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.defineAssetType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": data.clasificationID + '|' + data.hashID,
                    "mutableTraits": "ASSET_A5:S|, burn:H|1",
                    "immutableTraits": "ASSET_A6:S|",
                    "mutableMetaTraits": "ASSET_A7:S|",
                    "immutableMetaTraits": "ASSET_A8:S|"
                }
            }


            var err, res = await chai.request(server)
                .post(config.defineAssetPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });

        it('Query Classification: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qClassification)

            var list = res.body.result.value.classifications.value.list
            list.forEach(function (value) {

                if (value.value.immutableTraits.value.properties.value.propertyList[0].value.id.value.idString == "ASSET_A8") {
                    data.chainID = value.value.id.value.chainID.value.idString
                    data.clshashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Mint Asset: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.mintAssetType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "toID": data.clasificationID + '|' + data.hashID,
                    "fromID": data.clasificationID + '|' + data.hashID,
                    "classificationID": data.chainID + '.' + data.clshashID,
                    "mutableProperties": "ASSET_A5:S|A, burn:H|1",
                    "immutableProperties": "ASSET_A6:S|B",
                    "mutableMetaProperties": "ASSET_A7:S|C",
                    "immutableMetaProperties": "ASSET_A8:S|D"

                }
            }


            var err, res = await chai.request(server)
                .post(config.mintAssetPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });

        it('Query Asset: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qAsset)

            var list = res.body.result.value.assets.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "ASSET_A8") {
                    data.aclasificationID = value.value.id.value.classificationID.value.idString
                    data.ahashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Add asset properties on mutation: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.mutateAssetType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": data.iclasificationID + '|' + data.ihashID,
                    "assetID": data.aclasificationID + '|' + data.ahashID,
                    "mutableProperties": "ASSET_A5:S|A",
                    "mutableMetaProperties": "ASSET_A7:S|C"
                }
            }


            var err, res = await chai.request(server)
                .post(config.mutateAssetPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash
        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });
    });

    describe('Mint Asset with more than 22 properties', async () => {

        beforeEach(function (done) {
            this.timeout(4000)
            setTimeout(function () {
                done()
            }, 3000)
        })


        var data = {
            'clasificationID': '',
            'hashID': '',
            'chainID': '',
            'clshashID': ''
        }

        it('Query Identity: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qIdentity)

            var list = res.body.result.value.identities.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "immutableMetaTraits1") {
                    data.clasificationID = value.value.id.value.classificationID.value.idString
                    data.hashID = value.value.id.value.hashID.value.idString
                }

            });

        });

        it('Define Asset: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.defineAssetType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": data.clasificationID + '|' + data.hashID,
                    "mutableTraits": "P1:S|A,P11:S|B,P12:S|C,P13:S|D,P14:S|E,P14:S|F,burn:H|1",
                    "immutableTraits": "P2:S|G,P21:S|H,P22:S|I,P23:S|J,P24:S|K",
                    "mutableMetaTraits": "P3:S|L,P31:S|M,P32:S|N,P33:S|O,P34:S|P",
                    "immutableMetaTraits": "P4:S|Q,P41:S|R,P42:S|S,P43:S|T,P44:S|U,P45:S|V"
                }
            }


            var err, res = await chai.request(server)
                .post(config.defineAssetPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });

        it('Query Classification: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qClassification)

            var list = res.body.result.value.classifications.value.list
            list.forEach(function (value) {

                if (value.value.immutableTraits.value.properties.value.propertyList[0].value.id.value.idString == "P4") {
                    data.chainID = value.value.id.value.chainID.value.idString
                    data.clshashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Mint Asset with more than 22 properties: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.mintAssetType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "toID": data.clasificationID + '|' + data.hashID,
                    "fromID": data.clasificationID + '|' + data.hashID,
                    "classificationID": data.chainID + '.' + data.clshashID,
                    "mutableProperties": "P1:S|A,P11:S|B,P12:S|C,P13:S|D,P14:S|E,P14:S|F,burn:H|1",
                    "immutableProperties": "P2:S|G,P21:S|H,P22:S|I,P23:S|J,P24:S|K",
                    "mutableMetaProperties": "P3:S|L,P31:S|M,P32:S|N,P33:S|O,P34:S|P",
                    "immutableMetaProperties": "P4:S|Q,P41:S|R,P42:S|S,P43:S|T,P44:S|U,P45:S|V"

                }
            }
            var err, res = await chai.request(server)
                .post(config.mintAssetPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });
    });

    describe('Mutate Asset to add more that 22 properties', async () => {

        beforeEach(function (done) {
            this.timeout(4000)
            setTimeout(function () {
                done()
            }, 3000)
        })


        var data = {
            'iclasificationID': '',
            'ihashID': '',
            'aclasificationID': '',
            'ahashID': ''
        }

        it('Query Identity: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qIdentity)

            var list = res.body.result.value.identities.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "immutableMetaTraits1") {
                    data.iclasificationID = value.value.id.value.classificationID.value.idString
                    data.ihashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Query Asset: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qAsset)

            var list = res.body.result.value.assets.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "ASSET_A8") {
                    data.aclasificationID = value.value.id.value.classificationID.value.idString
                    data.ahashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Mutate Asset to add more that 22 properties: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.mutateAssetType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": data.iclasificationID + '|' + data.ihashID,
                    "assetID": data.aclasificationID + '|' + data.ahashID,
                    "mutableProperties": "ASSET_A5:S|A,ASSET_A5:S|B,ASSET_A5:S|C,ASSET_A5:S|D,ASSET_A5:S|E,ASSET_A5:S|F,ASSET_A5:S|G,ASSET_A5:S|H,ASSET_A5:S|I,ASSET_A5:S|J,ASSET_A5:S|K,ASSET_A5:S|L,ASSET_A5:S|M,ASSET_A5:S|N",
                    "mutableMetaProperties": "ASSET_A7:S|O,ASSET_A7:S|P,ASSET_A7:S|Q,ASSET_A7:S|R,ASSET_A7:S|S,ASSET_A7:S|T,ASSET_A7:S|U,ASSET_A7:S|V,ASSET_A7:S|W"
                }
            }


            var err, res = await chai.request(server)
                .post(config.mutateAssetPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });
    });

    describe('Burn Asset', async () => {

        beforeEach(function (done) {
            this.timeout(4000)
            setTimeout(function () {
                done()
            }, 3000)
        })


        var data = {
            'iclasificationID': '',
            'ihashID': '',
            'aclasificationID': '',
            'ahashID': ''
        }

        it('Query Identity: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qIdentity)

            var list = res.body.result.value.identities.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "immutableMetaTraits1") {
                    data.iclasificationID = value.value.id.value.classificationID.value.idString
                    data.ihashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Query Asset: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qAsset)

            var list = res.body.result.value.assets.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "ASSET_P4") {
                    data.aclasificationID = value.value.id.value.classificationID.value.idString
                    data.ahashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Meta Reveal: ', async () => {

            let obj = {
                "type": config.metaRevealType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "metaFact": "H|1"
                }
            }

            var err, res = await chai.request(server)
                .post(config.metaRevealPath)
                .send(obj)

            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash
        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });

        it('Burn Asset', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.burnAssetType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": data.iclasificationID + '|' + data.ihashID,
                    "assetID": data.aclasificationID + '|' + data.ahashID
                }
            }


            var err, res = await chai.request(server)
                .post(config.burnAssetPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });
    });

    describe('Mint Asset with burn greater than forseeable block height', async () => {

        beforeEach(function (done) {
            this.timeout(4000)
            setTimeout(function () {
                done()
            }, 3000)
        })


        var data = {
            'iclasificationID': '',
            'ihashID': '',
            'chainID': '',
            'clshashID': '',
            'aclasificationID': '',
            'ahashID': ''
        }

        it('Query Identity: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qIdentity)

            var list = res.body.result.value.identities.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "immutableMetaTraits1") {
                    data.iclasificationID = value.value.id.value.classificationID.value.idString
                    data.ihashID = value.value.id.value.hashID.value.idString
                }

            });

        });

        it('Define Asset: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.defineAssetType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": data.iclasificationID + '|' + data.ihashID,
                    "mutableTraits": "AssetA:S|num1" + ",burn:H|1",
                    "immutableTraits": "AssetB:S|",
                    "mutableMetaTraits": "AssetC:S|num3",
                    "immutableMetaTraits": "AssetD:S|num4"
                }
            }


            var err, res = await chai.request(server)
                .post(config.defineAssetPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });

        it('Query Classification: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qClassification)

            var list = res.body.result.value.classifications.value.list
            list.forEach(function (value) {

                if (value.value.immutableTraits.value.properties.value.propertyList[0].value.id.value.idString == "AssetD") {
                    data.chainID = value.value.id.value.chainID.value.idString
                    data.clshashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Mint Asset: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.mintAssetType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "toID": data.iclasificationID + '|' + data.ihashID,
                    "fromID": data.iclasificationID + '|' + data.ihashID,
                    "classificationID": data.chainID + '.' + data.clshashID,
                    "mutableProperties": "AssetA:S|num1" + ",burn:H|100000000",
                    "immutableProperties": "AssetB:S|num2",
                    "mutableMetaProperties": "AssetC:S|num3",
                    "immutableMetaProperties": "AssetD:S|num4"

                }
            }


            var err, res = await chai.request(server)
                .post(config.mintAssetPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });

        it('Query Asset: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qAsset)

            var list = res.body.result.value.assets.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "AssetD") {
                    data.aclasificationID = value.value.id.value.classificationID.value.idString
                    data.ahashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Meta Reveal: ', async () => {

            let obj = {
                "type": config.metaRevealType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "metaFact": "H|100000000"
                }
            }

            var err, res = await chai.request(server)
                .post(config.metaRevealPath)
                .send(obj)

            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash
        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });
        it('Burn Asset', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.burnAssetType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": data.iclasificationID + '|' + data.ihashID,
                    "assetID": data.aclasificationID + '|' + data.ahashID
                }
            }


            var err, res = await chai.request(server)
                .post(config.burnAssetPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });

        it('Meta Reveal: ', async () => {

            let obj = {
                "type": config.metaRevealType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "metaFact": "H|100"
                }
            }

            var err, res = await chai.request(server)
                .post(config.metaRevealPath)
                .send(obj)

            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash
        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });

        it('Mutate Asset: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.mutateAssetType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": data.iclasificationID + '|' + data.ihashID,
                    "assetID": data.aclasificationID + '|' + data.ahashID,
                    "mutableProperties": "AssetA:S|ABCd,burn:H|100",
                    "mutableMetaProperties": "AssetC:S|num3"
                }
            }


            var err, res = await chai.request(server)
                .post(config.mutateAssetPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });

        it('Burn Asset', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.burnAssetType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": data.iclasificationID + '|' + data.ihashID,
                    "assetID": data.aclasificationID + '|' + data.ahashID
                }
            }


            var err, res = await chai.request(server)
                .post(config.burnAssetPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });
    });

    describe('Mint Asset with extra properties when mutable trait is not defined', async () => {

        beforeEach(function (done) {
            this.timeout(4000)
            setTimeout(function () {
                done()
            }, 3000)
        })


        var data = {
            'clasificationID': '',
            'hashID': '',
            'chainID': '',
            'clshashID': ''
        }

        it('Query Identity: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qIdentity)

            var list = res.body.result.value.identities.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "immutableMetaTraits1") {
                    data.clasificationID = value.value.id.value.classificationID.value.idString
                    data.hashID = value.value.id.value.hashID.value.idString
                }

            });

        });

        it('Define Asset: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.defineAssetType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": data.clasificationID + '|' + data.hashID,
                    "mutableTraits": "ASSETA1:S|,burn:H|1",
                    "immutableTraits": "ASSETA2:S|G",
                    "mutableMetaTraits": "ASSETA3:S|L",
                    "immutableMetaTraits": "ASSETA4:S|Q"
                }
            }


            var err, res = await chai.request(server)
                .post(config.defineAssetPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });

        it('Query Classification: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qClassification)

            var list = res.body.result.value.classifications.value.list
            list.forEach(function (value) {

                if (value.value.immutableTraits.value.properties.value.propertyList[0].value.id.value.idString == "AssetD") {
                    data.chainID = value.value.id.value.chainID.value.idString
                    data.clshashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Mint Asset: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.mintAssetType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "toID": data.clasificationID + '|' + data.hashID,
                    "fromID": data.clasificationID + '|' + data.hashID,
                    "classificationID": data.chainID + '.' + data.clshashID,
                    "mutableProperties": "ASSETA1:S|A,burn:H|1,ASSETA1:S|B,ASSETA1:S|C",
                    "immutableProperties": "ASSETA2:S|G",
                    "mutableMetaProperties": "ASSETA3:S|L",
                    "immutableMetaProperties": "ASSETA4:S|Q"
                }
            }


            var err, res = await chai.request(server)
                .post(config.mintAssetPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });
    });
});


////////////////////////////////////////////////////////////////////////////////////////////////
/*
Splits
*/
////////////////////////////////////////////////////////////////////////////////////////////////

describe('Splits', async () => {

    describe('send split of an asset', async () => {

        beforeEach(function (done) {
            this.timeout(4000)
            setTimeout(function () {
                done()
            }, 3000)
        })

        var data = {
            'iclasificationID1': '',
            'ihashID1': '',
            'iclasificationID2': '',
            'ihashI2': '',
            'iclasificationID3': '',
            'ihashID3': '',
            'chainID': '',
            'clshashID': '',
            'aclasificationID': '',
            'ahashID': ''
        }

        it('Query Identity 1: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qIdentity)

            var list = res.body.result.value.identities.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "immutableMetaTraits1") {
                    data.iclasificationID1 = value.value.id.value.classificationID.value.idString
                    data.ihashID1 = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Query Identity 2: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qIdentity)

            var list = res.body.result.value.identities.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "immutableMetaTraits2") {
                    data.iclasificationID2 = value.value.id.value.classificationID.value.idString
                    data.ihashID2 = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Query Asset: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qAsset)

            var list = res.body.result.value.assets.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == " ") {
                    data.aclasificationID = value.value.id.value.classificationID.value.idString
                    data.ahashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Send split of an asset: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.sendSplitType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": data.iclasificationID1 + '|' + data.ihashID1,
                    "toID": data.iclasificationID2 + '|' + data.ihashID2,
                    "ownableID": data.aclasificationID + '.' + data.ahashID,
                    "split": config.splitval
                }
            }

            var err, res = await chai.request(server)
                .post(config.sendSplitPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });
    })

    describe('Unwrap a coin', async () => {

        beforeEach(function (done) {
            this.timeout(4000)
            setTimeout(function () {
                done()
            }, 3000)
        })

        var data = {
            'iclasificationID1': '',
            'ihashID1': '',
            'aclasificationID': '',
            'ahashID': ''
        }

        it('Query Identity: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qIdentity)

            var list = res.body.result.value.identities.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "immutableMetaTraits1") {
                    data.iclasificationID1 = value.value.id.value.classificationID.value.idString
                    data.ihashID1 = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Query Asset: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qAsset)

            var list = res.body.result.value.assets.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "AssetD") {
                    data.aclasificationID = value.value.id.value.classificationID.value.idString
                    data.ahashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Unwrap a coin: ', async () => {

            let obj = {
                "type": config.unwrapCoinType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": data.iclasificationID1 + '|' + data.ihashID1,
                    "ownableID": data.aclasificationID + '.' + data.ahashID,
                    "split": config.splitval
                }
            }

            var err, res = await chai.request(server)
                .post(config.unwrapCoinPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });
    })

    describe('Wrap a coin', async () => {

        beforeEach(function (done) {
            this.timeout(4000)
            setTimeout(function () {
                done()
            }, 3000)
        })

        var data = {
            'iclasificationID1': '',
            'ihashID1': ''
        }

        it('Query Identity: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qIdentity)

            var list = res.body.result.value.identities.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "immutableMetaTraits1") {
                    data.iclasificationID1 = value.value.id.value.classificationID.value.idString
                    data.ihashID1 = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Wrap a coin: ', async () => {

            let obj = {
                "type": config.wrapCoinType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": data.iclasificationID1 + '|' + data.ihashID1,
                    "coins": config.coins
                }
            }


            var err, res = await chai.request(server)
                .post(config.wrapCoinPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });
    })

    describe('send split of an coin', async () => {

        beforeEach(function (done) {
            this.timeout(4000)
            setTimeout(function () {
                done()
            }, 3000)
        })

        var data = {
            'iclasificationID1': '',
            'ihashID1': '',
            'iclasificationID2': '',
            'ihashI2': '',
            'iclasificationID3': '',
            'ihashID3': '',
            'chainID': '',
            'clshashID': '',
            'aclasificationID': '',
            'ahashID': ''
        }

        it('Query Identity 1: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qIdentity)

            var list = res.body.result.value.identities.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "immutableMetaTraits1") {
                    data.iclasificationID1 = value.value.id.value.classificationID.value.idString
                    data.ihashID1 = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Query Identity 2: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qIdentity)

            var list = res.body.result.value.identities.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "immutableMetaTraits2") {
                    data.iclasificationID2 = value.value.id.value.classificationID.value.idString
                    data.ihashID2 = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Query Asset: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qAsset)

            var list = res.body.result.value.assets.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "AssetD") {
                    data.aclasificationID = value.value.id.value.classificationID.value.idString
                    data.ahashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Send split of an coin: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.sendSplitType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": data.iclasificationID1 + '|' + data.ihashID1,
                    "toID": data.iclasificationID2 + '|' + data.ihashID2,
                    "ownableID": data.aclasificationID + '.' + data.ahashID,
                    "split": config.splitval
                }
            }

            var err, res = await chai.request(server)
                .post(config.sendSplitPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });
    })
})

////////////////////////////////////////////////////////////////////////////////////////////////
/*
  Metas
*/
////////////////////////////////////////////////////////////////////////////////////////////////

describe('Metas', async () => {

    describe('Reveal a meta', async () => {

        beforeEach(function (done) {
            this.timeout(4000)
            setTimeout(function () {
                done()
            }, 3000)
        })

        var data = {
            'iclasificationID1': '',
            'ihashID1': '',
            'chainID': '',
            'clshashID': '',
            'aclasificationID': '',
            'ahashID': ''
        }

        it('Query Identity: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qIdentity)

            var list = res.body.result.value.identities.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "immutableMetaTraits1") {
                    data.iclasificationID1 = value.value.id.value.classificationID.value.idString
                    data.ihashID1 = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Define Asset: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.defineAssetType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": data.iclasificationID1 + '|' + data.ihashID1,
                    "mutableTraits": "ASSET_PA:S|AAA, burn:H|4",
                    "immutableTraits": "ASSET_PB:D|0.344",
                    "mutableMetaTraits": "ASSET_PC:I|ID",
                    "immutableMetaTraits": "ASSET_PD:S|A"
                }
            }


            var err, res = await chai.request(server)
                .post(config.defineAssetPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });

        it('Query Classification: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qClassification)

            var list = res.body.result.value.classifications.value.list
            list.forEach(function (value) {

                if (value.value.immutableTraits.value.properties.value.propertyList[0].value.id.value.idString == "ASSET_PD") {
                    data.chainID = value.value.id.value.chainID.value.idString
                    data.clshashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Mint Asset: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.mintAssetType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "toID": data.iclasificationID1 + '|' + data.ihashID1,
                    "fromID": data.iclasificationID1 + '|' + data.ihashID1,
                    "classificationID": data.chainID + '.' + data.clshashID,
                    "mutableProperties": "ASSET_PA:S|AAA, burn:H|4",
                    "immutableProperties": "ASSET_PB:D|0.344",
                    "mutableMetaProperties": "ASSET_PC:I|ID",
                    "immutableMetaProperties": "ASSET_PD:S|A"

                }
            }


            var err, res = await chai.request(server)
                .post(config.mintAssetPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });

        it('Query Asset: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qAsset)

            var list = res.body.result.value.assets.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "ASSET_PD") {
                    data.aclasificationID = value.value.id.value.classificationID.value.idString
                    data.ahashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Meta Reveal: ', async () => {

            let obj = {
                "type": config.metaRevealType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "metaFact": "S|AAA"
                }
            }

            var err, res = await chai.request(server)
                .post(config.metaRevealPath)
                .send(obj)

            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash
        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });
    })

    describe('Reveal a meta of id type', async () => {

        beforeEach(function (done) {
            this.timeout(4000)
            setTimeout(function () {
                done()
            }, 3000)
        })

        it('Meta Reveal: ', async () => {

            let obj = {
                "type": config.metaRevealType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "metaFact": "I|ID"
                }
            }

            var err, res = await chai.request(server)
                .post(config.metaRevealPath)
                .send(obj)

            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash
        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });
    })

    describe('Reveal a meta of string type', async () => {

        beforeEach(function (done) {
            this.timeout(4000)
            setTimeout(function () {
                done()
            }, 3000)
        })

        it('Meta Reveal: ', async () => {

            let obj = {
                "type": config.metaRevealType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "metaFact": "S|AAA"
                }
            }

            var err, res = await chai.request(server)
                .post(config.metaRevealPath)
                .send(obj)

            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash
        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });
    })

    describe('Reveal a meta of dec type', async () => {

        beforeEach(function (done) {
            this.timeout(4000)
            setTimeout(function () {
                done()
            }, 3000)
        })

        it('Meta Reveal: ', async () => {

            let obj = {
                "type": config.metaRevealType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "metaFact": "D|0.344"
                }
            }

            var err, res = await chai.request(server)
                .post(config.metaRevealPath)
                .send(obj)

            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash
        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });
    })

    describe('Reveal a meta of height type', async () => {

        beforeEach(function (done) {
            this.timeout(4000)
            setTimeout(function () {
                done()
            }, 3000)
        })

        it('Meta Reveal: ', async () => {

            let obj = {
                "type": config.metaRevealType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "metaFact": "H|4"
                }
            }

            var err, res = await chai.request(server)
                .post(config.metaRevealPath)
                .send(obj)

            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash
        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });
    })

    describe('Reveal an already revealed meta', async () => {

        beforeEach(function (done) {
            this.timeout(4000)
            setTimeout(function () {
                done()
            }, 3000)
        })

        it('Meta Reveal: ', async () => {

            let obj = {
                "type": config.metaRevealType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "metaFact": "H|4"
                }
            }

            var err, res = await chai.request(server)
                .post(config.metaRevealPath)
                .send(obj)

            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash
        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });
    })
})

////////////////////////////////////////////////////////////////////////////////////////////////
/*
  Orders
*/
////////////////////////////////////////////////////////////////////////////////////////////////

describe('Orders', async () => {

    describe('Create an asset make order', async () => {

        beforeEach(function (done) {
            this.timeout(4000)
            setTimeout(function () {
                done()
            }, 3000)
        })

        var data = {
            'clasificationID': '',
            'hashID': '',
            'chainID': '',
            'clshashID': '',
            'aclasificationID': '',
            'ahashID': ''
        }


        it('Query Identity: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qIdentity)

            var list = res.body.result.value.identities.value.list
            list.forEach(function (value) {
                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "immutableMetaTraits1") {
                    data.clasificationID = value.value.id.value.classificationID.value.idString
                    data.hashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Define Order: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.defineOrderType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": data.clasificationID + '|' + data.hashID,
                    "mutableTraits": "ORDER_MUTABLE2:S|ORDER_CLASSIFICATION_MUTABLE_1",
                    "immutableTraits": "ORDER_IMMUTABLE2:S|ORDER_CLASSIFICATION_IMMUTABLE_1",
                    "mutableMetaTraits": "ORDER_MUTABLE_META2:S|ORDER_CLASSIFICATION_MUTABLE_META_1,exchangeRate:D|" + config.exchangeRate + ",makerOwnableSplit:D|" + config.makerownablesplit + ",expiry:H|" + config.expiry + ",makerSplit:D|" + config.makerownablesplit + ",takerID:S|",
                    "immutableMetaTraits": "ORDER_IMMUTABLE_META2:S|num4"
                }
            }

            var err, res = await chai.request(server)
                .post(config.defineOrderPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash
        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });

        it('Query Classification: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qClassification)

            var list = res.body.result.value.classifications.value.list
            list.forEach(function (value) {

                if (value.value.immutableTraits.value.properties.value.propertyList[0].value.id.value.idString == "ORDER_IMMUTABLE_META2") {
                    data.chainID = value.value.id.value.chainID.value.idString
                    data.clshashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Query Asset: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qAsset)

            var list = res.body.result.value.assets.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "ASSET4") {
                    data.aclasificationID = value.value.id.value.classificationID.value.idString
                    data.ahashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Asset Make Order: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.makeOrderType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": data.clasificationID + '|' + data.hashID,
                    "classificationID": data.chainID + '.' + data.clshashID,
                    "makerOwnableID": data.aclasificationID + '|' + data.ahashID,
                    "takerOwnableID": config.makerownableid,
                    "expiresIn": config.expiry,
                    "makerOwnableSplit": config.makerownablesplit,
                    "mutableProperties": "ORDER_MUTABLE2:S|ORDER_CLASSIFICATION_MUTABLE_1",
                    "immutableProperties": "ORDER_IMMUTABLE2:S|ORDER_CLASSIFICATION_IMMUTABLE_1",
                    "mutableMetaProperties": "ORDER_MUTABLE_META2:S|ORDER_CLASSIFICATION_MUTABLE_META_1,exchangeRate:D|" + config.exchangeRate + ",makerSplit:D|" + config.makerownablesplit,
                    "immutableMetaProperties": "ORDER_IMMUTABLE_META2:S|num4"
                }
            }

            var err, res = await chai.request(server)
                .post(config.makeOrderPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });
    })

    describe('Take an asset take order', async () => {

        beforeEach(function (done) {
            this.timeout(4000)
            setTimeout(function () {
                done()
            }, 3000)
        })

        var data = {
            'clasificationID': '',
            'hashID': '',
            'chainID': '',
            'clshashID': '',
            'aclasificationID': '',
            'ahashID': '',
            'oclassificationID': '',
            'omakerownableid': '',
            'otakerownableid': '',
            'omakerid': '',
            'ohashid': ''
        }


        it('Query Identity: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qIdentity)

            var list = res.body.result.value.identities.value.list
            list.forEach(function (value) {
                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "immutableMetaTraits1") {
                    data.clasificationID = value.value.id.value.classificationID.value.idString
                    data.hashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Query Orders: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qOrder)

            var list = res.body.result.value.orders.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "ORDER_IMMUTABLE_META2") {
                    data.oclassificationID = value.value.id.value.classificationID.value.idString
                    data.omakerownableid = value.value.id.value.makerOwnableID.value.idString
                    data.otakerownableid = value.value.id.value.takerOwnableID.value.idString
                    data.omakerid = value.value.id.value.makerID.value.idString
                    data.ohashid = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Take Order: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.takeOrderType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": data.clasificationID + '|' + data.hashID,
                    "takerOwnableSplit": config.makerownablesplit,
                    "orderID": data.oclassificationID + '*' + data.omakerownableid + '*' + data.otakerownableid + '*' + data.omakerid + '*' + data.ohashid
                }
            }

            var err, res = await chai.request(server)
                .post(config.takeOrderPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });
    })

    describe('Cancel an asset order', async () => {

        beforeEach(function (done) {
            this.timeout(4000)
            setTimeout(function () {
                done()
            }, 3000)
        })

        var data = {
            'clasificationID': '',
            'hashID': '',
            'chainID': '',
            'clshashID': '',
            'aclasificationID': '',
            'ahashID': '',
            'oclassificationID': '',
            'omakerownableid': '',
            'otakerownableid': '',
            'omakerid': '',
            'ohashid': ''
        }


        it('Query Identity: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qIdentity)

            var list = res.body.result.value.identities.value.list
            list.forEach(function (value) {
                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "immutableMetaTraits1") {
                    data.clasificationID = value.value.id.value.classificationID.value.idString
                    data.hashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Define Order: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.defineOrderType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": data.clasificationID + '|' + data.hashID,
                    "mutableTraits": "ORDER_MUTABLE3:S|ORDER_CLASSIFICATION_MUTABLE_1",
                    "immutableTraits": "ORDER_IMMUTABLE3:S|ORDER_CLASSIFICATION_IMMUTABLE_1",
                    "mutableMetaTraits": "ORDER_MUTABLE_META3:S|ORDER_CLASSIFICATION_MUTABLE_META_1,exchangeRate:D|" + config.exchangeRate + ",makerOwnableSplit:D|" + config.makerownablesplit + ",expiry:H|" + config.expiry + ",makerSplit:D|" + config.makerownablesplit + ",takerID:S|ID",
                    "immutableMetaTraits": "ORDER_IMMUTABLE_META3:S|num4"
                }
            }

            var err, res = await chai.request(server)
                .post(config.defineOrderPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash
        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });

        it('Query Classification: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qClassification)

            var list = res.body.result.value.classifications.value.list
            list.forEach(function (value) {

                if (value.value.immutableTraits.value.properties.value.propertyList[0].value.id.value.idString == "ORDER_IMMUTABLE_META3") {
                    data.chainID = value.value.id.value.chainID.value.idString
                    data.clshashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Query Asset: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qAsset)

            var list = res.body.result.value.assets.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "AssetD") {
                    data.aclasificationID = value.value.id.value.classificationID.value.idString
                    data.ahashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Make Order: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.makeOrderType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": data.clasificationID + '|' + data.hashID,
                    "classificationID": data.chainID + '.' + data.clshashID,
                    "makerOwnableID": data.aclasificationID + '|' + data.ahashID,
                    "takerOwnableID": config.makerownableid,
                    "expiresIn": config.expiry,
                    "makerOwnableSplit": config.makerownablesplit,
                    "mutableProperties": "ORDER_MUTABLE3:S|ORDER_CLASSIFICATION_MUTABLE_1",
                    "immutableProperties": "ORDER_IMMUTABLE3:S|ORDER_CLASSIFICATION_IMMUTABLE_1",
                    "mutableMetaProperties": "ORDER_MUTABLE_META3:S|ORDER_CLASSIFICATION_MUTABLE_META_1,exchangeRate:D|" + config.exchangeRate + ",makerSplit:D|" + config.makerownablesplit,
                    "immutableMetaProperties": "ORDER_IMMUTABLE_META3:S|num4"
                }
            }

            var err, res = await chai.request(server)
                .post(config.makeOrderPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });

        it('Query Orders: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qOrder)

            var list = res.body.result.value.orders.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "ORDER_IMMUTABLE_META3") {
                    data.oclassificationID = value.value.id.value.classificationID.value.idString
                    data.omakerownableid = value.value.id.value.makerOwnableID.value.idString
                    data.otakerownableid = value.value.id.value.takerOwnableID.value.idString
                    data.omakerid = value.value.id.value.makerID.value.idString
                    data.ohashid = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Cancel Order: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.cancelOrderType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": data.clasificationID + '|' + data.hashID,
                    "orderID": data.oclassificationID + '*' + data.omakerownableid + '*' + data.otakerownableid + '*' + data.omakerid + '*' + data.ohashid
                }
            }

            var err, res = await chai.request(server)
                .post(config.cancelOrderPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });
    })

    describe('Create a coin make order', async () => {

        beforeEach(function (done) {
            this.timeout(4000)
            setTimeout(function () {
                done()
            }, 3000)
        })

        var data = {
            'clasificationID': '',
            'hashID': '',
            'chainID': '',
            'clshashID': '',
            'aclasificationID': '',
            'ahashID': ''
        }


        it('Query Identity: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qIdentity)

            var list = res.body.result.value.identities.value.list
            list.forEach(function (value) {
                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "immutableMetaTraits1") {
                    data.clasificationID = value.value.id.value.classificationID.value.idString
                    data.hashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Define Order: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.defineOrderType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": data.clasificationID + '|' + data.hashID,
                    "mutableTraits": "ORDER_MUTABLE4:S|ORDER_CLASSIFICATION_MUTABLE_1",
                    "immutableTraits": "ORDER_IMMUTABLE4:S|ORDER_CLASSIFICATION_IMMUTABLE_1",
                    "mutableMetaTraits": "ORDER_MUTABLE_META4:S|ORDER_CLASSIFICATION_MUTABLE_META_1,exchangeRate:D|" + config.exchangeRate + ",makerOwnableSplit:D|" + config.makerownablesplit + ",expiry:H|" + config.expiry + ",makerSplit:D|" + config.makerownablesplit + ",takerID:S|ID",
                    "immutableMetaTraits": "ORDER_IMMUTABLE_META4:S|num4"
                }
            }

            var err, res = await chai.request(server)
                .post(config.defineOrderPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash
        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });

        it('Query Classification: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qClassification)

            var list = res.body.result.value.classifications.value.list
            list.forEach(function (value) {

                if (value.value.immutableTraits.value.properties.value.propertyList[0].value.id.value.idString == "ORDER_IMMUTABLE_META4") {
                    data.chainID = value.value.id.value.chainID.value.idString
                    data.clshashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Wrap a coin: ', async () => {

            let obj = {
                "type": config.wrapCoinType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": data.clasificationID + '|' + data.hashID,
                    "coins": config.coins
                }
            }


            var err, res = await chai.request(server)
                .post(config.wrapCoinPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });

        it('Coin Make Order: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.makeOrderType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": data.clasificationID + '|' + data.hashID,
                    "classificationID": data.chainID + '.' + data.clshashID,
                    "makerOwnableID": config.makerownableid,
                    "takerOwnableID": config.takerownableid,
                    "expiresIn": config.expiry,
                    "makerOwnableSplit": config.makerownablesplit,
                    "mutableProperties": "ORDER_MUTABLE4:S|ORDER_CLASSIFICATION_MUTABLE_1",
                    "immutableProperties": "ORDER_IMMUTABLE4:S|ORDER_CLASSIFICATION_IMMUTABLE_1",
                    "mutableMetaProperties": "ORDER_MUTABLE_META4:S|ORDER_CLASSIFICATION_MUTABLE_META_1,exchangeRate:D|" + config.exchangeRate,
                    "immutableMetaProperties": "ORDER_IMMUTABLE_META4:S|num4"
                }
            }


            var err, res = await chai.request(server)
                .post(config.makeOrderPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });
    })

    describe('Take an coin take order', async () => {

        beforeEach(function (done) {
            this.timeout(4000)
            setTimeout(function () {
                done()
            }, 3000)
        })

        var data = {
            'clasificationID': '',
            'hashID': '',
            'chainID': '',
            'clshashID': '',
            'aclasificationID': '',
            'ahashID': '',
            'oclassificationID': '',
            'omakerownableid': '',
            'otakerownableid': '',
            'omakerid': '',
            'ohashid': ''
        }


        it('Query Identity: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qIdentity)

            var list = res.body.result.value.identities.value.list
            list.forEach(function (value) {
                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "immutableMetaTraits1") {
                    data.clasificationID = value.value.id.value.classificationID.value.idString
                    data.hashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Query Orders: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qOrder)

            var list = res.body.result.value.orders.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "ORDER_IMMUTABLE_META4") {
                    data.oclassificationID = value.value.id.value.classificationID.value.idString
                    data.omakerownableid = value.value.id.value.makerOwnableID.value.idString
                    data.takerownableid = value.value.id.value.takerOwnableID.value.idString
                    data.omakerid = value.value.id.value.makerID.value.idString
                    data.ohashid = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Take Order: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.takeOrderType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": data.clasificationID + '|' + data.hashID,
                    "takerOwnableSplit": config.makerownablesplit,
                    "orderID": data.oclassificationID + '*' + data.omakerownableid + '*' + data.takerownableid + '*' + data.omakerid + '*' + data.ohashid
                }
            }


            var err, res = await chai.request(server)
                .post(config.takeOrderPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });
    })

    describe('Cancel a coin order', async () => {

        beforeEach(function (done) {
            this.timeout(4000)
            setTimeout(function () {
                done()
            }, 3000)
        })

        var data = {
            'clasificationID': '',
            'hashID': '',
            'chainID': '',
            'clshashID': '',
            'aclasificationID': '',
            'ahashID': '',
            'oclassificationID': '',
            'omakerownableid': '',
            'otakerownableid': '',
            'omakerid': '',
            'ohashid': ''
        }


        it('Query Identity: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qIdentity)

            var list = res.body.result.value.identities.value.list
            list.forEach(function (value) {
                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "immutableMetaTraits1") {
                    data.clasificationID = value.value.id.value.classificationID.value.idString
                    data.hashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Define Order: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.defineOrderType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": data.clasificationID + '|' + data.hashID,
                    "mutableTraits": "ORDER_MUTABLE5:S|ORDER_CLASSIFICATION_MUTABLE_1",
                    "immutableTraits": "ORDER_IMMUTABLE5:S|ORDER_CLASSIFICATION_IMMUTABLE_1",
                    "mutableMetaTraits": "ORDER_MUTABLE_META5:S|ORDER_CLASSIFICATION_MUTABLE_META_1,exchangeRate:D|" + config.exchangeRate + ",makerOwnableSplit:D|" + config.makerownablesplit + ",expiry:H|" + config.expiry + ",makerSplit:D|" + config.makerownablesplit + ",takerID:S|ID",
                    "immutableMetaTraits": "ORDER_IMMUTABLE_META5:S|num4"
                }
            }

            var err, res = await chai.request(server)
                .post(config.defineOrderPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash
        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });

        it('Query Classification: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qClassification)

            var list = res.body.result.value.classifications.value.list
            list.forEach(function (value) {

                if (value.value.immutableTraits.value.properties.value.propertyList[0].value.id.value.idString == "ORDER_IMMUTABLE_META5") {
                    data.chainID = value.value.id.value.chainID.value.idString
                    data.clshashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Query Asset: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qAsset)

            var list = res.body.result.value.assets.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "AssetD") {
                    data.aclasificationID = value.value.id.value.classificationID.value.idString
                    data.ahashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Make Coin Order: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.makeOrderType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": data.clasificationID + '|' + data.hashID,
                    "classificationID": data.chainID + '.' + data.clshashID,
                    "makerOwnableID": config.makerownableid,
                    "takerOwnableID": config.makerownableid,
                    "expiresIn": config.expiry,
                    "makerOwnableSplit": config.makerownablesplit,
                    "mutableProperties": "ORDER_MUTABLE5:S|ORDER_CLASSIFICATION_MUTABLE_1",
                    "immutableProperties": "ORDER_IMMUTABLE5:S|ORDER_CLASSIFICATION_IMMUTABLE_1",
                    "mutableMetaProperties": "ORDER_MUTABLE_META5:S|ORDER_CLASSIFICATION_MUTABLE_META_1,exchangeRate:D|" + config.exchangeRate + ",makerSplit:D|" + config.makerownablesplit,
                    "immutableMetaProperties": "ORDER_IMMUTABLE_META5:S|num4"
                }
            }


            var err, res = await chai.request(server)
                .post(config.makeOrderPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });

        it('Query Orders: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qOrder)

            var list = res.body.result.value.orders.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "ORDER_IMMUTABLE_META5") {
                    data.oclassificationID = value.value.id.value.classificationID.value.idString
                    data.omakerownableid = value.value.id.value.makerOwnableID.value.idString
                    data.takerownableid = value.value.id.value.takerOwnableID.value.idString
                    data.omakerid = value.value.id.value.makerID.value.idString
                    data.ohashid = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Cancel Order: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.cancelOrderType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": data.clasificationID + '|' + data.hashID,
                    "orderID": data.oclassificationID + '*' + data.omakerownableid + '*' + data.takerownableid + '*' + data.omakerid + '*' + data.ohashid
                }
            }

            var err, res = await chai.request(server)
                .post(config.cancelOrderPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });
    })

    describe('Create an order with correct takerID', async () => {

        beforeEach(function (done) {
            this.timeout(4000)
            setTimeout(function () {
                done()
            }, 3000)
        })

        var data = {
            'clasificationID1': '',
            'hashID1': '',
            'clasificationID2': '',
            'hashID2': '',
            'chainID': '',
            'clshashID': '',
            'aclasificationID': '',
            'ahashID': ''
        }


        it('Query Identity 1: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qIdentity)

            var list = res.body.result.value.identities.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "immutableMetaTraits1") {
                    data.clasificationID1 = value.value.id.value.classificationID.value.idString
                    data.hashID1 = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Query Identity 2: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qIdentity)

            var list = res.body.result.value.identities.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "immutableMetaTraits2") {
                    data.clasificationID2 = value.value.id.value.classificationID.value.idString
                    data.hashID2 = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Define Order: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.defineOrderType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": data.clasificationID1 + '|' + data.hashID1,
                    "mutableTraits": "ORDER_MUTABLE6:S|ORDER_CLASSIFICATION_MUTABLE_1",
                    "immutableTraits": "ORDER_IMMUTABLE6:S|ORDER_CLASSIFICATION_IMMUTABLE_1",
                    "mutableMetaTraits": "ORDER_MUTABLE_META6:S|ORDER_CLASSIFICATION_MUTABLE_META_1,exchangeRate:D|" + config.exchangeRate + ",makerOwnableSplit:D|" + config.makerownablesplit + ",expiry:H|" + config.expiry + ",makerSplit:D|" + config.makerownablesplit + ",takerID:S|ID",
                    "immutableMetaTraits": "ORDER_IMMUTABLE_META6:S|num4"
                }
            }

            var err, res = await chai.request(server)
                .post(config.defineOrderPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash
        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });

        it('Query Classification: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qClassification)

            var list = res.body.result.value.classifications.value.list
            list.forEach(function (value) {

                if (value.value.immutableTraits.value.properties.value.propertyList[0].value.id.value.idString == "ORDER_IMMUTABLE_META6") {
                    data.chainID = value.value.id.value.chainID.value.idString
                    data.clshashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Query Asset: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qAsset)

            var list = res.body.result.value.assets.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "ASSET4") {
                    data.aclasificationID = value.value.id.value.classificationID.value.idString
                    data.ahashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Asset Make Order: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.makeOrderType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": data.clasificationID1 + '|' + data.hashID1,
                    "classificationID": data.chainID + '.' + data.clshashID,
                    "makerOwnableID": data.aclasificationID + '|' + data.ahashID,
                    "takerOwnableID": config.makerownableid,
                    "expiresIn": config.expiry,
                    "makerOwnableSplit": config.makerownablesplit,
                    "mutableProperties": "ORDER_MUTABLE6:S|ORDER_CLASSIFICATION_MUTABLE_1",
                    "immutableProperties": "ORDER_IMMUTABLE6:S|ORDER_CLASSIFICATION_IMMUTABLE_1",
                    "mutableMetaProperties": "ORDER_MUTABLE_META6:S|ORDER_CLASSIFICATION_MUTABLE_META_1,exchangeRate:D|" + config.exchangeRate + ",makerSplit:D|" + config.makerownablesplit,
                    "immutableMetaProperties": "ORDER_IMMUTABLE_META6:S|num4"
                }
            }


            var err, res = await chai.request(server)
                .post(config.makeOrderPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });
    })

    describe('Create an order with other takerID', async () => {

        beforeEach(function (done) {
            this.timeout(4000)
            setTimeout(function () {
                done()
            }, 3000)
        })

        var data = {
            'clasificationID1': '',
            'hashID1': '',
            'clasificationID2': '',
            'hashID2': '',
            'clasificationID3': '',
            'hashID3': '',
            'chainID': '',
            'clshashID': '',
            'aclasificationID': '',
            'ahashID': ''
        }


        it('Query Identity 1: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qIdentity)

            var list = res.body.result.value.identities.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "immutableMetaTraits1") {
                    data.clasificationID1 = value.value.id.value.classificationID.value.idString
                    data.hashID1 = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Query Identity 2: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qIdentity)

            var list = res.body.result.value.identities.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "immutableMetaTraits2") {
                    data.clasificationID2 = value.value.id.value.classificationID.value.idString
                    data.hashID2 = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Query Identity 3: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qIdentity)

            var list = res.body.result.value.identities.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "immutableMetaTraits3") {
                    data.clasificationID3 = value.value.id.value.classificationID.value.idString
                    data.hashID3 = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Query Classification: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qClassification)

            var list = res.body.result.value.classifications.value.list
            list.forEach(function (value) {

                if (value.value.immutableTraits.value.properties.value.propertyList[0].value.id.value.idString == "ORDER_IMMUTABLE_META2") {
                    data.chainID = value.value.id.value.chainID.value.idString
                    data.clshashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Query Asset: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qAsset)

            var list = res.body.result.value.assets.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "ASSET4") {
                    data.aclasificationID = value.value.id.value.classificationID.value.idString
                    data.ahashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Asset Make Order: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.makeOrderType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": data.clasificationID + '|' + data.hashID,
                    "classificationID": data.chainID + '.' + data.clshashID,
                    "makerOwnableID": data.aclasificationID + '|' + data.ahashID,
                    "takerOwnableID": config.makerownableid,
                    "expiresIn": config.expiry,
                    "makerOwnableSplit": config.makerownablesplit,
                    "mutableProperties": "ORDER_MUTABLE6:S|ORDER_CLASSIFICATION_MUTABLE_1",
                    "immutableProperties": "ORDER_IMMUTABLE6:S|ORDER_CLASSIFICATION_IMMUTABLE_1",
                    "mutableMetaProperties": "ORDER_MUTABLE_META6:S|ORDER_CLASSIFICATION_MUTABLE_META_1,exchangeRate:D|" + config.exchangeRate + ",makerSplit:D|" + config.makerownablesplit + ",takerID:S|" + data.clasificationID3 * data.hashID3,
                    "immutableMetaProperties": "ORDER_IMMUTABLE_META6:S|num4"
                }
            }


            var err, res = await chai.request(server)
                .post(config.makeOrderPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });
    })

    describe('Create order with takerID', async () => {

        beforeEach(function (done) {
            this.timeout(4000)
            setTimeout(function () {
                done()
            }, 3000)
        })

        var data = {
            'clasificationID': '',
            'hashID': '',
            'chainID': '',
            'clshashID': '',
            'aclasificationID': '',
            'ahashID': ''
        }


        it('Query Identity: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qIdentity)

            var list = res.body.result.value.identities.value.list
            list.forEach(function (value) {
                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "immutableMetaTraits1") {
                    data.clasificationID = value.value.id.value.classificationID.value.idString
                    data.hashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Define Order: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.defineOrderType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": data.clasificationID + '|' + data.hashID,
                    "mutableTraits": "ORDER_MUTABLE7:S|ORDER_CLASSIFICATION_MUTABLE_1",
                    "immutableTraits": "ORDER_IMMUTABLE7:S|ORDER_CLASSIFICATION_IMMUTABLE_1",
                    "mutableMetaTraits": "ORDER_MUTABLE_META7:S|ORDER_CLASSIFICATION_MUTABLE_META_1,exchangeRate:D|" + config.exchangeRate + ",makerOwnableSplit:D|" + config.makerownablesplit + ",expiry:H|" + config.expiry + ",makerSplit:D|" + config.makerownablesplit + ",takerID:S|",
                    "immutableMetaTraits": "ORDER_IMMUTABLE_META7:S|num4"
                }
            }

            var err, res = await chai.request(server)
                .post(config.defineOrderPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash
        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });

        it('Query Classification: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qClassification)

            var list = res.body.result.value.classifications.value.list
            list.forEach(function (value) {

                if (value.value.immutableTraits.value.properties.value.propertyList[0].value.id.value.idString == "ORDER_IMMUTABLE_META7") {
                    data.chainID = value.value.id.value.chainID.value.idString
                    data.clshashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Query Asset: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qAsset)

            var list = res.body.result.value.assets.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "ASSET4") {
                    data.aclasificationID = value.value.id.value.classificationID.value.idString
                    data.ahashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Asset Make Order: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.makeOrderType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": data.clasificationID + '|' + data.hashID,
                    "classificationID": data.chainID + '.' + data.clshashID,
                    "makerOwnableID": data.aclasificationID + '|' + data.ahashID,
                    "takerOwnableID": config.makerownableid,
                    "expiresIn": config.expiry,
                    "makerOwnableSplit": config.makerownablesplit,
                    "mutableProperties": "ORDER_MUTABLE7:S|ORDER_CLASSIFICATION_MUTABLE_1",
                    "immutableProperties": "ORDER_IMMUTABLE7:S|ORDER_CLASSIFICATION_IMMUTABLE_1",
                    "mutableMetaProperties": "ORDER_MUTABLE_META7:S|ORDER_CLASSIFICATION_MUTABLE_META_1,exchangeRate:D|" + config.exchangeRate + ",makerSplit:D|" + config.makerownablesplit + ",takerID:S|" + data.clasificationID * data.hashID,
                    "immutableMetaProperties": "ORDER_IMMUTABLE_META7:S|num4"
                }
            }

            var err, res = await chai.request(server)
                .post(config.makeOrderPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });
    })

    describe('Sell assets with splits, where taker gives more splits than he is supposed to', async () => {

        beforeEach(function (done) {
            this.timeout(4000)
            setTimeout(function () {
                done()
            }, 3000)
        })

        var data = {
            'clasificationID1': '',
            'hashID1': '',
            'clasificationID2': '',
            'hashID2': '',
            'clasificationID3': '',
            'hashID3': '',
            'chainID': '',
            'clshashID': '',
            'aclasificationID': '',
            'ahashID': ''
        }


        it('Query Identity 1: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qIdentity)

            var list = res.body.result.value.identities.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "immutableMetaTraits1") {
                    data.clasificationID1 = value.value.id.value.classificationID.value.idString
                    data.hashID1 = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Query Identity 2: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qIdentity)

            var list = res.body.result.value.identities.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "immutableMetaTraits2") {
                    data.clasificationID2 = value.value.id.value.classificationID.value.idString
                    data.hashID2 = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Query Identity 3: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qIdentity)

            var list = res.body.result.value.identities.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "immutableMetaTraits3") {
                    data.clasificationID3 = value.value.id.value.classificationID.value.idString
                    data.hashID3 = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Define Asset: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.defineAssetType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": data.clasificationID + '|' + data.hashID,
                    "mutableTraits": "ASSETS10:S|num1" + ",burn:H|1",
                    "immutableTraits": "ASSETS11:S|",
                    "mutableMetaTraits": "ASSETS12:S|num3",
                    "immutableMetaTraits": "ASSETS13:S|num4"
                }
            }


            var err, res = await chai.request(server)
                .post(config.defineAssetPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });



        it('Query Classification: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qClassification)

            var list = res.body.result.value.classifications.value.list
            list.forEach(function (value) {

                if (value.value.immutableTraits.value.properties.value.propertyList[0].value.id.value.idString == "ASSETS13") {
                    data.chainID = value.value.id.value.chainID.value.idString
                    data.clshashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Mint Asset: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.mintAssetType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "toID": data.clasificationID1 + '|' + data.hashID1,
                    "fromID": data.clasificationID1 + '|' + data.hashID1,
                    "classificationID": data.chainID + '.' + data.clshashID,
                    "mutableProperties": "ASSETS10:S|num1" + ",burn:H|1",
                    "immutableProperties": "ASSETS11:S|abc",
                    "mutableMetaProperties": "ASSETS12:S|num3",
                    "immutableMetaProperties": "ASSETS13:S|num4"
                }
            }

            var err, res = await chai.request(server)
                .post(config.mintAssetPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });

        it('Query Asset: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qAsset)

            var list = res.body.result.value.assets.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "ASSETS13") {
                    data.aclasificationID = value.value.id.value.classificationID.value.idString
                    data.ahashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Define Order: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.defineOrderType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": data.clasificationID1 + '|' + data.hashID1,
                    "mutableTraits": "ORDER_MUTABLE21:S|ORDER_CLASSIFICATION_MUTABLE_1",
                    "immutableTraits": "ORDER_IMMUTABLE22:S|ORDER_CLASSIFICATION_IMMUTABLE_1",
                    "mutableMetaTraits": "ORDER_MUTABLE_META23:S|ORDER_CLASSIFICATION_MUTABLE_META_1,exchangeRate:D|" + config.exchangeRate + ",makerOwnableSplit:D|" + config.makerownablesplit + ",expiry:H|" + config.expiry + ",makerSplit:D|" + config.makerownablesplit + ",takerSplit:D|,takerID:S|",
                    "immutableMetaTraits": "ORDER_IMMUTABLE_META24:S|num4"
                }
            }

            var err, res = await chai.request(server)
                .post(config.defineOrderPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash
        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });

        it('Query Classification: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qClassification)

            var list = res.body.result.value.classifications.value.list
            list.forEach(function (value) {

                if (value.value.immutableTraits.value.properties.value.propertyList[0].value.id.value.idString == "ORDER_IMMUTABLE_META24") {
                    data.chainID = value.value.id.value.chainID.value.idString
                    data.clshashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Asset Make Order: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.makeOrderType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": data.clasificationID1 + '|' + data.hashID1,
                    "classificationID": data.chainID + '.' + data.clshashID,
                    "makerOwnableID": data.aclasificationID + '|' + data.ahashID,
                    "takerOwnableID": config.makerownableid,
                    "expiresIn": config.expiry,
                    "makerOwnableSplit": config.makerownablesplit,
                    "mutableProperties": "ORDER_MUTABLE21:S|ORDER_CLASSIFICATION_MUTABLE_1",
                    "immutableProperties": "ORDER_IMMUTABLE22:S|ORDER_CLASSIFICATION_IMMUTABLE_1",
                    "mutableMetaProperties": "ORDER_MUTABLE_META23:S|ORDER_CLASSIFICATION_MUTABLE_META_1,exchangeRate:D|" + config.exchangeRate + ",makerSplit:D|" + config.makerownablesplit + ",takerID:S|,makerID:S|" + data.clasificationID3 * data.hashID3,
                    "immutableMetaProperties": "ORDER_IMMUTABLE_META24:S|num4"
                }
            }


            var err, res = await chai.request(server)
                .post(config.makeOrderPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });
    })

    describe('Orders (splits) with exchange value other than smallest dec', async () => {

        beforeEach(function (done) {
            this.timeout(4000)
            setTimeout(function () {
                done()
            }, 3000)
        })

        var data = {
            'clasificationID1': '',
            'hashID1': '',
            'clasificationID2': '',
            'hashID2': '',
            'clasificationID3': '',
            'hashID3': '',
            'chainID': '',
            'clshashID': '',
            'aclasificationID': '',
            'ahashID': ''
        }


        it('Query Identity 1: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qIdentity)

            var list = res.body.result.value.identities.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "immutableMetaTraits1") {
                    data.clasificationID1 = value.value.id.value.classificationID.value.idString
                    data.hashID1 = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Query Identity 2: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qIdentity)

            var list = res.body.result.value.identities.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "immutableMetaTraits2") {
                    data.clasificationID2 = value.value.id.value.classificationID.value.idString
                    data.hashID2 = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Query Identity 3: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qIdentity)

            var list = res.body.result.value.identities.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "immutableMetaTraits3") {
                    data.clasificationID3 = value.value.id.value.classificationID.value.idString
                    data.hashID3 = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Define Asset: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.defineAssetType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": data.clasificationID + '|' + data.hashID,
                    "mutableTraits": "ASSETS101:S|num1" + ",burn:H|1",
                    "immutableTraits": "ASSETS111:S|",
                    "mutableMetaTraits": "ASSETS121:S|num3",
                    "immutableMetaTraits": "ASSETS131:S|num4"
                }
            }


            var err, res = await chai.request(server)
                .post(config.defineAssetPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });



        it('Query Classification: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qClassification)

            var list = res.body.result.value.classifications.value.list
            list.forEach(function (value) {

                if (value.value.immutableTraits.value.properties.value.propertyList[0].value.id.value.idString == "ASSETS131") {
                    data.chainID = value.value.id.value.chainID.value.idString
                    data.clshashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Mint Asset: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.mintAssetType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "toID": data.clasificationID1 + '|' + data.hashID1,
                    "fromID": data.clasificationID1 + '|' + data.hashID1,
                    "classificationID": data.chainID + '.' + data.clshashID,
                    "mutableProperties": "ASSETS101:S|num1" + ",burn:H|1",
                    "immutableProperties": "ASSETS111:S|abc",
                    "mutableMetaProperties": "ASSETS121:S|num3",
                    "immutableMetaProperties": "ASSETS131:S|num4"
                }
            }

            var err, res = await chai.request(server)
                .post(config.mintAssetPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });

        it('Query Asset: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qAsset)

            var list = res.body.result.value.assets.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "ASSETS131") {
                    data.aclasificationID = value.value.id.value.classificationID.value.idString
                    data.ahashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Define Order: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.defineOrderType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": data.clasificationID1 + '|' + data.hashID1,
                    "mutableTraits": "ORDER_MUTABLE211:S|ORDER_CLASSIFICATION_MUTABLE_1",
                    "immutableTraits": "ORDER_IMMUTABLE222:S|ORDER_CLASSIFICATION_IMMUTABLE_1",
                    "mutableMetaTraits": "ORDER_MUTABLE_META233:S|ORDER_CLASSIFICATION_MUTABLE_META_1,exchangeRate:D|" + config.exchangeRate + ",makerOwnableSplit:D|100" + ",expiry:H|" + config.expiry + ",makerSplit:D|" + config.makerownablesplit + ",takerSplit:D|,takerID:S|",
                    "immutableMetaTraits": "ORDER_IMMUTABLE_META244:S|num4"
                }
            }

            var err, res = await chai.request(server)
                .post(config.defineOrderPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash
        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });

        it('Query Classification: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qClassification)

            var list = res.body.result.value.classifications.value.list
            list.forEach(function (value) {

                if (value.value.immutableTraits.value.properties.value.propertyList[0].value.id.value.idString == "ORDER_IMMUTABLE_META244") {
                    data.chainID = value.value.id.value.chainID.value.idString
                    data.clshashID = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Asset Make Order: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.makeOrderType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": data.clasificationID1 + '|' + data.hashID1,
                    "classificationID": data.chainID + '.' + data.clshashID,
                    "makerOwnableID": data.aclasificationID + '|' + data.ahashID,
                    "takerOwnableID": config.makerownableid,
                    "expiresIn": config.expiry,
                    "makerOwnableSplit": config.makerownablesplit,
                    "mutableProperties": "ORDER_MUTABLE21:S|ORDER_CLASSIFICATION_MUTABLE_1",
                    "immutableProperties": "ORDER_IMMUTABLE22:S|ORDER_CLASSIFICATION_IMMUTABLE_1",
                    "mutableMetaProperties": "ORDER_MUTABLE_META23:S|ORDER_CLASSIFICATION_MUTABLE_META_1,exchangeRate:D|" + config.exchangeRate + ",makerSplit:D|100",
                    "immutableMetaProperties": "ORDER_IMMUTABLE_META24:S|num4"
                }
            }


            var err, res = await chai.request(server)
                .post(config.makeOrderPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    //do nothing
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });
    })
})

////////////////////////////////////////////////////////////////////////////////////////////////
/*
  Maintainer Deputize
*/
////////////////////////////////////////////////////////////////////////////////////////////////

describe('Maintainer', async () => {
    describe('Deputize', async () => {

        beforeEach(function (done) {
            this.timeout(4000)
            setTimeout(function () {
                done()
            }, 3000)
        })


        var data = {
            'clasificationID1': '',
            'hashID1': '',
            'clasificationID2': '',
            'hashID2': '',
            'aclasificationID': '',
            'ahashID': ''
        }

        it('Query Identity 1: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qIdentity)

            var list = res.body.result.value.identities.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "immutableMetaTraits1") {
                    data.clasificationID1 = value.value.id.value.classificationID.value.idString
                    data.hashID1 = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Query Identity 2: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qIdentity)

            var list = res.body.result.value.identities.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "immutableMetaTraits2") {
                    data.clasificationID2 = value.value.id.value.classificationID.value.idString
                    data.hashID2 = value.value.id.value.hashID.value.idString
                }
            });
        });

        it('Query Asset: ', async () => {

            var err, res = await chai.request(server)
                .get(config.qAsset)

            var list = res.body.result.value.assets.value.list
            list.forEach(function (value) {

                if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString == "ASSET_PD") {
                    data.aclasificationID = value.value.id.value.classificationID.value.idString
                    data.ahashID = value.value.id.value.hashID.value.idString
                }
            });
        });


        it('Maintainer Deputize: ', async () => {

            var num = helper.generateRandomInteger(0, 10000)

            let obj = {
                "type": config.deputizeType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": data.clasificationID1 + '|' + data.hashID1,
                    "toID": data.clasificationID2 + '|' + data.hashID2,
                    "classificationID": data.aclasificationID + '|' + data.ahashID,
                    "maintainedTraits": "maintainerTraits:S|maintainerTraits",
                    "addMaintainer": true,
                    "removeMaintainer": false,
                    "mutateMaintainer": false
                }
            }

            var err, res = await chai.request(server)
                .post(config.deputizePath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')

            txHash = res.body.txhash

        });

        it('Query Tx: ', async () => {
            var err, res = await chai.request(server)
                .get('/txs/' + txHash)
            var data1 = JSON.stringify(res.body)

            async.retry({ times: config.retry_count, interval: config.timeout }, function (callbackretry) {

                if ((data1.indexOf('error') != -1)) {
                    callbackretry('failed')
                } else {
                    //continue
                }
            }, function (err, response) {
                if (err) {
                    console.log("err: " + err)
                    console.log("response: " + response)
                    console.log("failed to send txHash query")
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                } else {
                    expect(res.body.raw_log).to.not.contain('failed')
                    expect(res.body).to.not.contain('error')
                }
            })
        });
    })
})