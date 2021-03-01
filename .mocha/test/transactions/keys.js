//During the test the env variable is set to test
process.env.NODE_ENV = 'test';
let chai = require('chai');
let chaiHttp = require('chai-http');
let should = chai.should();
let expect = chai.expect
const crypto = require("crypto");
const {request} = require('chai')
let assert = chai.assert

let config = require('../config.json');
let server = config.ip + config.port

chai.use(chaiHttp);

const name1 = crypto.randomBytes(16).toString("hex");
const name2 = crypto.randomBytes(16).toString("hex");

let names1 = {
    address: '',
    typekey: '',
    value: '',
    signature: ''
}

let names2 = {
    address: '',
    typekey: '',
    value: '',
    signature: ''
}

////////////////////////////////////////////////////////////////////////////////////////////////
/*
Account Creation for User1
*/
////////////////////////////////////////////////////////////////////////////////////////////////
describe('Keys', async () => {
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

            let err, res = await chai.request(server)
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

            let err, res = await chai.request(server)
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


            let err, res = await chai.request(server)
                .post(config.broadcastTx)
                .send(obj)


            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
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

            let err, res = await chai.request('http://localhost:1317')
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

            let err, res = await chai.request(server)
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


            let err, res = await chai.request(server)
                .post(config.broadcastTx)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
        });
    });
})
