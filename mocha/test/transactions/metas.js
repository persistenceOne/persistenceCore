//During the test the env variable is set to test
process.env.NODE_ENV = 'test';
let chai = require('chai');
let chaiHttp = require('chai-http');
let async = require('async');
let should = chai.should();
let expect = chai.expect
const crypto = require("crypto");
const {request} = require('chai')
var assert = chai.assert

var config = require('../config.json');
var temp = require('../helpers/helpers');
const {type} = require('os');
var server = config.ip + config.port
var identity = require('../helpers/identities')
var cls = require('../helpers/classifications')

chai.use(chaiHttp);


describe('Metas', async () => {

    describe('Reveal a meta', async () => {

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
                    "mutableTraits": "ASSET_PA:S|AAA, burn:H|4",
                    "immutableTraits": "ASSET_PB:D|0.344,ASSET_PE:I|ID,ASSET_PF:S|A",
                    "mutableMetaTraits": "ASSET_PC:S|ABBCBBC",
                    "immutableMetaTraits": "ASSET_PD:S|QQQQQ"
                }
            }

            var err, res = await chai.request(server)
                .post(config.defineAssetPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });

        it('Mint Asset: ', async () => {

            let identityID = await identity.queryIdentity("immutableMetaTraits1")
            let clsID = await cls.queryClassification("ASSET_PD")

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
                    "mutableProperties": "ASSET_PA:S|AAA, burn:H|4",
                    "immutableProperties": "ASSET_PB:D|0.344,ASSET_PE:I|ID,ASSET_PF:S|A",
                    "mutableMetaProperties": "ASSET_PC:S|ABBCBBC",
                    "immutableMetaProperties": "ASSET_PD:S|QQQQQ"
                }
            }

            var err, res = await chai.request(server)
                .post(config.mintAssetPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
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
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });
    })

    describe('Reveal a meta of id type', async () => {

        beforeEach(function (done) {
            this.timeout(5000)
            setTimeout(function () {
                done()
            }, 4000)
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
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });
    })

    describe('Reveal a meta of string type', async () => {

        beforeEach(function (done) {
            this.timeout(5000)
            setTimeout(function () {
                done()
            }, 4000)
        })

        it('Meta Reveal: ', async () => {

            let obj = {
                "type": config.metaRevealType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "metaFact": "S|A"
                }
            }

            var err, res = await chai.request(server)
                .post(config.metaRevealPath)
                .send(obj)

            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });
    })

    describe('Reveal a meta of dec type', async () => {

        beforeEach(function (done) {
            this.timeout(5000)
            setTimeout(function () {
                done()
            }, 4000)
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
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });
    })

    describe('Reveal a meta of height type', async () => {

        beforeEach(function (done) {
            this.timeout(5000)
            setTimeout(function () {
                done()
            }, 4000)
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
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });
    })

    describe('Reveal an already revealed meta', async () => {

        beforeEach(function (done) {
            this.timeout(5000)
            setTimeout(function () {
                done()
            }, 4000)
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
            expect(res.body.raw_log).to.contain('failed')
            expect(res.body.raw_log).to.contain('EntityAlreadyExists')
        });
    })
})
