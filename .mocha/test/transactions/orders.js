//During the test the env variable is set to test
process.env.NODE_ENV = 'test';
let chai = require('chai');
let chaiHttp = require('chai-http');
let should = chai.should();
let expect = chai.expect
const {request} = require('chai')
let assert = chai.assert

let config = require('../config.json');
let server = config.ip + config.port
let identity = require('../helpers/identities')
let cls = require('../helpers/classifications')
let assets = require('../helpers/assets')
let orders = require('../helpers/orders')

chai.use(chaiHttp);

describe('Orders', async () => {

    describe('Create an asset make order', async () => {

        beforeEach(function (done) {
            this.timeout(5000)
            setTimeout(function () {
                done()
            }, 4000)
        })

        it('Define Asset: ', async () => {

            let identityID = await identity.queryIdentity("immutableMetaTraits1")

            let obj = {
                "type": config.defineAssetType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": identityID,
                    "mutableTraits": "A_P1:S|" + ",burn:H|1",
                    "immutableTraits": "A_P2:S|",
                    "mutableMetaTraits": "A_P3:S|",
                    "immutableMetaTraits": "A_P4:S|"
                }
            }

            let err, res = await chai.request(server)
                .post(config.defineAssetPath)
                .send(obj);

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });

        it('Mint Asset: ', async () => {

            let identityID = await identity.queryIdentity("immutableMetaTraits1")
            let clsID = await cls.queryClassification("A_P4")

            let obj = {
                "type": config.mintAssetType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "toID": identityID,
                    "fromID": identityID,
                    "classificationID": clsID,
                    "mutableProperties": "A_P1:S|,burn:H|1",
                    "immutableProperties": "A_P2:S|",
                    "mutableMetaProperties": "A_P3:S|",
                    "immutableMetaProperties": "A_P4:S|"
                }
            }

            let err, res = await chai.request(server)
                .post(config.mintAssetPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });

        it('Define Order: ', async () => {

            let identityID = await identity.queryIdentity("immutableMetaTraits1")

            let obj = {
                "type": config.defineOrderType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": identityID,
                    "immutableMetaTraits": "Name:S|,Gifts:S|Exchange,OrderID:S|",
                    "immutableTraits": "Which Gifts:S|,What Gifts:S|",
                    "mutableMetaTraits": "exchangeRate:D|1,makerOwnableSplit:D|0.000000000000000001,expiry:H|1000000,takerID:I|ID,makerSplit:D|0.000000000000000001",
                    "mutableTraits": "descriptions:S|"
                }
            }

            let err, res = await chai.request(server)
                .post(config.defineOrderPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });

        it('Asset Make Order: ', async () => {

            let identityID = await identity.queryIdentity("immutableMetaTraits1")
            let clsID = await cls.queryClassification("Name")
            let assetID = await assets.queryAsset("A_P4")

            let obj = {
                "type": config.makeOrderType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": identityID,
                    "classificationID": clsID,
                    "makerOwnableID": assetID,
                    "takerOwnableID":"stake",
                    "expiresIn":"100000",
                    "makerOwnableSplit":"0.000000000000000001",
                    "immutableMetaProperties": "Name:S|Board,Gifts:S|Exchange,OrderID:S|12345",
                    "immutableProperties": "Which Gifts:S|Christmas Gift,What Gifts:S|kitty",
                    "mutableMetaProperties": "exchangeRate:D|1,makerSplit:D|0.000000000000000001",
                    "mutableProperties": "descriptions:S|awesomekitty"
                }
            }

            let err, res = await chai.request(server)
                .post(config.makeOrderPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });
    })

    describe('Cancel an asset order', async () => {

        beforeEach(function (done) {
            this.timeout(5000)
            setTimeout(function () {
                done()
            }, 4000)
        })

        it('Cancel Order: ', async () => {

            let identityID = await identity.queryIdentity("immutableMetaTraits1")
            let orderID = await orders.queryOrder("Name")

            let obj = {
                "type": config.cancelOrderType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": identityID,
                    "orderID": orderID
                }
            }

            let err, res = await chai.request(server)
                .post(config.cancelOrderPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });
    })

    describe('Take an asset take order', async () => {

        beforeEach(function (done) {
            this.timeout(5000)
            setTimeout(function () {
                done()
            }, 4000)
        })

        it('Make Order: ', async () => {

            let identityID = await identity.queryIdentity("immutableMetaTraits1")
            let clsID = await cls.queryClassification("Name")
            let assetID = await assets.queryAsset("A_P4")

            let obj = {
                "type": config.makeOrderType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": identityID,
                    "classificationID": clsID,
                    "makerOwnableID": assetID,
                    "takerOwnableID":"stake",
                    "expiresIn":"100000",
                    "makerOwnableSplit":"0.000000000000000001",
                    "immutableMetaProperties": "Name:S|Board,Gifts:S|Exchange,OrderID:S|12345",
                    "immutableProperties": "Which Gifts:S|Christmas Gift,What Gifts:S|kitty",
                    "mutableMetaProperties": "exchangeRate:D|1,makerSplit:D|0.000000000000000001",
                    "mutableProperties": "descriptions:S|awesomekitty"
                }
            }

            let err, res = await chai.request(server)
                .post(config.makeOrderPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });

        it('Splits send: ', async () => {

            let identityID = await identity.queryIdentity("immutableMetaTraits1")
            let identityID1 = await identity.queryIdentity("immutableMetaTraits2")

            let obj = {
                "type": config.sendSplitType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": identityID,
                    "toID": identityID1,
                    "ownableID": "stake",
                    "split": config.splitVal
                }
            }

            let err, res = await chai.request(server)
                .post(config.sendSplitPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });

        it('Take Order: ', async () => {

            let identityID = await identity.queryIdentity("immutableMetaTraits2")
            let orderID = await orders.queryOrder("Name")

            let obj = {
                "type": config.takeOrderType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": identityID,
                    "takerOwnableSplit": config.makerOwnableSplit,
                    "orderID": orderID
                }
            }

            let err, res = await chai.request(server)
                .post(config.takeOrderPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });
    })

    describe('Create a coin make order', async () => {

        beforeEach(function (done) {
            this.timeout(5000)
            setTimeout(function () {
                done()
            }, 4000)
        })

        it('Define Asset: ', async () => {

            let identityID = await identity.queryIdentity("immutableMetaTraits1")

            let obj = {
                "type": config.defineAssetType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": identityID,
                    "mutableTraits": "A_PP1:S|" + ",burn:H|1",
                    "immutableTraits": "A_PP2:S|",
                    "mutableMetaTraits": "A_PP3:S|",
                    "immutableMetaTraits": "A_PP4:S|"
                }
            }

            let err, res = await chai.request(server)
                .post(config.defineAssetPath)
                .send(obj);

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });

        it('Mint Asset: ', async () => {

            let identityID = await identity.queryIdentity("immutableMetaTraits1")
            let clsID = await cls.queryClassification("A_PP4")

            let obj = {
                "type": config.mintAssetType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "toID": identityID,
                    "fromID": identityID,
                    "classificationID": clsID,
                    "mutableProperties": "A_PP1:S|,burn:H|1",
                    "immutableProperties": "A_PP2:S|",
                    "mutableMetaProperties": "A_PP3:S|",
                    "immutableMetaProperties": "A_PP4:S|"
                }
            }

            let err, res = await chai.request(server)
                .post(config.mintAssetPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });

        it('Define Order: ', async () => {

            let identityID = await identity.queryIdentity("immutableMetaTraits1")

            let obj = {
                "type": config.defineOrderType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": identityID,
                    "immutableMetaTraits": "Name1:S|,Gifts1:S|Exchange,OrderID1:S|",
                    "immutableTraits": "Which Gifts1:S|,What Gifts1:S|",
                    "mutableMetaTraits": "exchangeRate:D|1,makerOwnableSplit:D|0.000000000000000001,expiry:H|1000000,takerID:I|ID,makerSplit:D|0.000000000000000001",
                    "mutableTraits": "descriptions1:S|"
                }
            }

            let err, res = await chai.request(server)
                .post(config.defineOrderPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });

        it('Wrap a coin: ', async () => {

            let identityID = await identity.queryIdentity("immutableMetaTraits1")

            let obj = {
                "type": config.wrapCoinType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": identityID,
                    "coins": config.coins
                }
            }


            let err, res = await chai.request(server)
                .post(config.wrapCoinPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });

        it('Coin Make Order: ', async () => {

            let identityID = await identity.queryIdentity("immutableMetaTraits1")
            let clsID = await cls.queryClassification("Name1")

            let obj = {
                "type": config.makeOrderType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": identityID,
                    "classificationID": clsID,
                    "makerOwnableID": "stake",
                    "takerOwnableID":"stake",
                    "expiresIn":"100000",
                    "makerOwnableSplit":"0.000000000000000001",
                    "immutableMetaProperties": "Name1:S|Board,Gifts1:S|Exchange,OrderID1:S|12345",
                    "immutableProperties": "Which Gifts1:S|Christmas Gift,What Gifts1:S|kitty",
                    "mutableMetaProperties": "exchangeRate:D|1,makerSplit:D|0.000000000000000001",
                    "mutableProperties": "descriptions1:S|awesomekitty"
                }
            }

            let err, res = await chai.request(server)
                .post(config.makeOrderPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });
    })

    describe('Cancel a coin order', async () => {

        beforeEach(function (done) {
            this.timeout(5000)
            setTimeout(function () {
                done()
            }, 4000)
        })

        it('Cancel Order: ', async () => {

            let identityID = await identity.queryIdentity("immutableMetaTraits1")
            let orderID = await orders.queryOrder("Name1")

            let obj = {
                "type": config.cancelOrderType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": identityID,
                    "orderID": orderID
                }
            }

            let err, res = await chai.request(server)
                .post(config.cancelOrderPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });
    })

    describe('Take a coin take order', async () => {

        beforeEach(function (done) {
            this.timeout(5000)
            setTimeout(function () {
                done()
            }, 4000)
        })

        it('Coin Make Order: ', async () => {

            let identityID = await identity.queryIdentity("immutableMetaTraits1")
            let clsID = await cls.queryClassification("Name1")

            let obj = {
                "type": config.makeOrderType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": identityID,
                    "classificationID": clsID,
                    "makerOwnableID": "stake",
                    "takerOwnableID":"stake",
                    "expiresIn":"100000",
                    "makerOwnableSplit":"0.000000000000000001",
                    "immutableMetaProperties": "Name1:S|Board,Gifts1:S|Exchange,OrderID1:S|12345",
                    "immutableProperties": "Which Gifts1:S|Christmas Gift,What Gifts1:S|kitty",
                    "mutableMetaProperties": "exchangeRate:D|1,makerSplit:D|0.000000000000000001",
                    "mutableProperties": "descriptions1:S|awesomekitty"
                }
            }

            let err, res = await chai.request(server)
                .post(config.makeOrderPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });


        it('Coin Take Order: ', async () => {

            let identityID = await identity.queryIdentity("immutableMetaTraits1")
            let orderID = await orders.queryOrder("Name1")

            let obj = {
                "type": config.takeOrderType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": identityID,
                    "takerOwnableSplit": config.makerOwnableSplit,
                    "orderID": orderID
                }
            }


            let err, res = await chai.request(server)
                .post(config.takeOrderPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });
    })

    describe('Create an order with correct takerID', async () => {

        beforeEach(function (done) {
            this.timeout(5000)
            setTimeout(function () {
                done()
            }, 4000)
        })

        it('Define Asset: ', async () => {

            let identityID = await identity.queryIdentity("immutableMetaTraits1")

            let obj = {
                "type": config.defineAssetType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": identityID,
                    "mutableTraits": "A_PPP1:S|" + ",burn:H|1",
                    "immutableTraits": "A_PPP2:S|",
                    "mutableMetaTraits": "A_PPP3:S|",
                    "immutableMetaTraits": "A_PPP4:S|"
                }
            }

            let err, res = await chai.request(server)
                .post(config.defineAssetPath)
                .send(obj);

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });

        it('Mint Asset: ', async () => {

            let identityID = await identity.queryIdentity("immutableMetaTraits1")
            let clsID = await cls.queryClassification("A_PPP4")

            let obj = {
                "type": config.mintAssetType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "toID": identityID,
                    "fromID": identityID,
                    "classificationID": clsID,
                    "mutableProperties": "A_PPP1:S|,burn:H|1",
                    "immutableProperties": "A_PPP2:S|",
                    "mutableMetaProperties": "A_PPP3:S|",
                    "immutableMetaProperties": "A_PPP4:S|"
                }
            }

            let err, res = await chai.request(server)
                .post(config.mintAssetPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });

        it('Define Order: ', async () => {

            let identityID = await identity.queryIdentity("immutableMetaTraits1")

            let obj = {
                "type": config.defineOrderType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": identityID,
                    "immutableMetaTraits": "Name2:S|,Gifts2:S|Exchange,OrderID2:S|",
                    "immutableTraits": "Which Gifts2:S|,What Gifts2:S|",
                    "mutableMetaTraits": "exchangeRate:D|1,makerOwnableSplit:D|0.000000000000000001,expiry:H|1000000,takerID:I|ID,makerSplit:D|0.000000000000000001",
                    "mutableTraits": "descriptions2:S|"
                }
            }

            let err, res = await chai.request(server)
                .post(config.defineOrderPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });

        it('Asset Make Order: ', async () => {

            let identityID = await identity.queryIdentity("immutableMetaTraits1")
            let identityID1 = await identity.queryIdentity("immutableMetaTraits2")
            let clsID = await cls.queryClassification("Name2")
            let assetID = await assets.queryAsset("A_PPP4")

            let obj = {
                "type": config.makeOrderType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": identityID,
                    "classificationID": clsID,
                    "makerOwnableID": assetID,
                    "takerOwnableID":identityID1,
                    "expiresIn":"100000",
                    "makerOwnableSplit":"0.000000000000000001",
                    "immutableMetaProperties": "Name2:S|Board,Gifts2:S|Exchange,OrderID2:S|12345",
                    "immutableProperties": "Which Gifts2:S|Christmas Gift,What Gifts2:S|kitty",
                    "mutableMetaProperties": "exchangeRate:D|1,makerSplit:D|0.000000000000000001",
                    "mutableProperties": "descriptions2:S|awesomekitty"
                }
            }

            let err, res = await chai.request(server)
                .post(config.makeOrderPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });

        it('Take Order with Incorrect ID: ', async () => {

            let identityID1 = await identity.queryIdentity("immutableMetaTraits3")
            let orderID = await orders.queryOrder("Name2")

            let obj = {
                "type": config.takeOrderType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": identityID1,
                    "takerOwnableSplit": config.makerOwnableSplit,
                    "orderID": orderID
                }
            }

            let err, res = await chai.request(server)
                .post(config.takeOrderPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.contain('failed')
        });

        it('Take Order with correct ID: ', async () => {

            let identityID1 = await identity.queryIdentity("immutableMetaTraits2")
            let orderID = await orders.queryOrder("Name2")

            let obj = {
                "type": config.takeOrderType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": identityID1,
                    "takerOwnableSplit": config.makerOwnableSplit,
                    "orderID": orderID
                }
            }

            let err, res = await chai.request(server)
                .post(config.takeOrderPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.contain('failed')
        });
    })

    describe('Sell assets with splits, where taker gives more splits than he is supposed to', async () => {

        beforeEach(function (done) {
            this.timeout(5000)
            setTimeout(function () {
                done()
            }, 4000)
        })


        it('Define Asset: ', async () => {

            let identityID = await identity.queryIdentity("immutableMetaTraits1")

            let obj = {
                "type": config.defineAssetType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": identityID,
                    "mutableTraits": "ASSETS10:S|num1" + ",burn:H|1",
                    "immutableTraits": "ASSETS11:S|",
                    "mutableMetaTraits": "ASSETS12:S|num3",
                    "immutableMetaTraits": "ASSETS13:S|num4"
                }
            }


            let err, res = await chai.request(server)
                .post(config.defineAssetPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });

        it('Mint Asset: ', async () => {

            let identityID = await identity.queryIdentity("immutableMetaTraits1")
            let clsID = await cls.queryClassification("ASSETS13")

            let obj = {
                "type": config.mintAssetType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "toID": identityID,
                    "fromID": identityID,
                    "classificationID": clsID,
                    "mutableProperties": "ASSETS10:S|num1" + ",burn:H|1",
                    "immutableProperties": "ASSETS11:S|abc",
                    "mutableMetaProperties": "ASSETS12:S|num3",
                    "immutableMetaProperties": "ASSETS13:S|num4"
                }
            }

            let err, res = await chai.request(server)
                .post(config.mintAssetPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });

        it('Define Order: ', async () => {

            let identityID = await identity.queryIdentity("immutableMetaTraits1")

            let obj = {
                "type": config.defineOrderType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": identityID,
                    "immutableMetaTraits": "Name3:S|,Gifts3:S|Exchange,OrderID3:S|",
                    "immutableTraits": "Which Gifts3:S|,What Gifts3:S|",
                    "mutableMetaTraits": "exchangeRate:D|1,makerOwnableSplit:D|0.000000000000000001,expiry:H|1000000,takerID:I|ID,makerSplit:D|0.000000000000000001",
                    "mutableTraits": "descriptions3:S|"
                }
            }

            let err, res = await chai.request(server)
                .post(config.defineOrderPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });

        it('Make Order: ', async () => {

            let identityID = await identity.queryIdentity("immutableMetaTraits1")
            let clsID = await cls.queryClassification("Name3")

            let obj = {
                "type": config.makeOrderType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": identityID,
                    "classificationID": clsID,
                    "makerOwnableID": "stake",
                    "takerOwnableID":"stake",
                    "expiresIn":"100000",
                    "makerOwnableSplit":"100",
                    "immutableMetaProperties": "Name3:S|Board,Gifts3:S|Exchange,OrderID3:S|12345",
                    "immutableProperties": "Which Gifts3:S|Christmas Gift,What Gifts3:S|kitty",
                    "mutableMetaProperties": "exchangeRate:D|1,makerSplit:D|0.000000000000000001",
                    "mutableProperties": "descriptions3:S|awesomekitty"
                }
            }

            let err, res = await chai.request(server)
                .post(config.makeOrderPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });

        it('Take Order: ', async () => {

            let identityID = await identity.queryIdentity("immutableMetaTraits1")
            let orderID = await orders.queryOrder("Name3")

            let obj = {
                "type": config.takeOrderType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": identityID,
                    "takerOwnableSplit": "200",
                    "orderID": orderID
                }
            }

            let err, res = await chai.request(server)
                .post(config.takeOrderPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });
    })
})
