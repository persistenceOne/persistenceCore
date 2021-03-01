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

chai.use(chaiHttp);


describe('Splits', async () => {

    describe('Send split of an asset', async () => {

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
                    "mutableTraits": "AssetDef1:S|Hello" + ",burn:H|10",
                    "immutableTraits": "AssetBDef2:S|",
                    "mutableMetaTraits": "AssetCDef3:S|",
                    "immutableMetaTraits": "AssetDDef4:S|"
                }
            }

            let err, res = await chai.request(server)
                .post(config.defineAssetPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });

        it('Mint asset: ', async () => {

            let identityID = await identity.queryIdentity("immutableMetaTraits1")
            let clsID = await cls.queryClassification("AssetDDef4")

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
                    "mutableProperties": "AssetDef1:S|Hello" + ",burn:H|10",
                    "immutableProperties": "AssetBDef2:S|",
                    "mutableMetaProperties": "AssetCDef3:S|",
                    "immutableMetaProperties": "AssetDDef4:S|"

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

        it('Send split of an asset: ', async () => {

            let identityID = await identity.queryIdentity("immutableMetaTraits1")
            let identityID1 = await identity.queryIdentity("immutableMetaTraits3")
            let assetID = await assets.queryAsset("AssetDDef4")

            let obj = {
                "type": config.sendSplitType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": identityID,
                    "toID": identityID1,
                    "ownableID": assetID,
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
    })

    describe('Wrap a coin', async () => {

        beforeEach(function (done) {
            this.timeout(5000)
            setTimeout(function () {
                done()
            }, 4000)
        })

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
    })

    describe('Unwrap a coin', async () => {

        beforeEach(function (done) {
            this.timeout(5000)
            setTimeout(function () {
                done()
            }, 4000)
        })

        it('Unwrap a coin: ', async () => {

            let identityID = await identity.queryIdentity("immutableMetaTraits1")

            let obj = {
                "type": config.unwrapCoinType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": identityID,
                    "ownableID": "stake",
                    "split": "100"
                }
            }

            let err, res = await chai.request(server)
                .post(config.unwrapCoinPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });
    })

    describe('Send split of an coin', async () => {

        beforeEach(function (done) {
            this.timeout(5000)
            setTimeout(function () {
                done()
            }, 4000)
        })

        it('Send split of an coin: ', async () => {

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
    })
})

