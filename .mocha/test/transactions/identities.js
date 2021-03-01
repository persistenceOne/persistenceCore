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
let keys = require('../helpers/keys')
let identity = require('../helpers/identities')
let cls = require('../helpers/classifications')

chai.use(chaiHttp);

describe('Identity', async () => {
    let randomWallet = keys.createRandomWallet();

    describe('Nub Tx', async () => {

        beforeEach(function (done) {
            this.timeout(6000)
            setTimeout(function () {
                done()
            }, 5000)
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

            let err, res = await chai.request(server)
                .post(config.nubPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('EntityAlreadyExists')
        });
    });

    describe('Issue identity 1', async () => {

        beforeEach(function (done) {
            this.timeout(5000)
            setTimeout(function () {
                done()
            }, 4000)
        })

        it('Define identity: ', async () => {

            let identityID = await identity.queryIdentity(config.nubID)

            let obj = {
                "type": config.defineIdentityType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": identityID,
                    "mutableTraits": "mutableTraits1:S|num1",
                    "immutableTraits": "immutableTraits1:S|",
                    "mutableMetaTraits": "mutableMetaTraits1:S|num3",
                    "immutableMetaTraits": "immutableMetaTraits1:S|num4"
                }
            }

            let err, res = await chai.request(server)
                .post(config.defineIdentityPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });


        it('Issue identity: ', async () => {

            let identityID = await identity.queryIdentity(config.nubID)
            let clsID = await cls.queryClassification("immutableMetaTraits1")

            let obj = {
                "type": config.issueIdentityType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "to": config.testAccountAddress,
                    "fromID": identityID,
                    "classificationID": clsID,
                    "mutableProperties": "mutableTraits1:S|num1",
                    "immutableProperties": "immutableTraits1:S|",
                    "mutableMetaProperties": "mutableMetaTraits1:S|num3",
                    "immutableMetaProperties": "immutableMetaTraits1:S|num4"
                }
            }

            let err, res = await chai.request(server)
                .post(config.issueIdentityPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });
    });

    describe('Provision key', async () => {

        beforeEach(function (done) {
            this.timeout(5000)
            setTimeout(function () {
                done()
            }, 4000)
        })

        it('Provision Key: ', async () => {

            let identityID = await identity.queryIdentity(config.nubID)

            let obj = {
                "type": config.provisionKeyType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "to": randomWallet.address,
                    "identityID": identityID
                }
            }

            let err, res = await chai.request(server)
                .post(config.provisionKeyPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')

        });
    });

    describe('Unprovision key', async () => {

        beforeEach(function (done) {
            this.timeout(5000)
            setTimeout(function () {
                done()
            }, 4000)
        })

        it('Unprovision key: ', async () => {

            let identityID = await identity.queryIdentity(config.nubID)

            let obj = {
                "type": config.unprovisionKeyType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "to": randomWallet.address,
                    "identityID": identityID
                }
            }

            let err, res = await chai.request(server)
                .post(config.unprovisionKeyPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });
    });

    describe('Provision an unprovision Key', async () => {

        beforeEach(function (done) {
            this.timeout(5000)
            setTimeout(function () {
                done()
            }, 4000)
        })

        it('Provision an unprovision Key: ', async () => {

            let identityID = await identity.queryIdentity(config.nubID)

            let obj = {
                "type": config.provisionKeyType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "to": randomWallet.address,
                    "identityID": identityID
                }
            }

            let err, res = await chai.request(server)
                .post(config.provisionKeyPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.contain('DeletionNotAllowed')
        });
    });

    describe('Issue identity 2', async () => {

        beforeEach(function (done) {
            this.timeout(5000)
            setTimeout(function () {
                done()
            }, 4000)
        })

        it('Define identity: ', async () => {

            let identityID = await identity.queryIdentity(config.nubID)

            let obj = {
                "type": config.defineIdentityType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": identityID,
                    "mutableTraits": "mutableTraits2:S|num1",
                    "immutableTraits": "immutableTraits2:S|",
                    "mutableMetaTraits": "mutableMetaTraits2:S|num3",
                    "immutableMetaTraits": "immutableMetaTraits2:S|num4"
                }
            }
            let err, res = await chai.request(server)
                .post(config.defineIdentityPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });

        it('Issue identity 2: ', async () => {

            let identityID = await identity.queryIdentity(config.nubID)
            let clsID = await cls.queryClassification("immutableMetaTraits2")

            let obj = {
                "type": config.issueIdentityType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "to": config.testAccountAddress,
                    "fromID": identityID,
                    "classificationID": clsID,
                    "mutableProperties": "mutableTraits2:S|num1",
                    "immutableProperties": "immutableTraits2:S|",
                    "mutableMetaProperties": "mutableMetaTraits2:S|num3",
                    "immutableMetaProperties": "immutableMetaTraits2:S|num4"
                }
            }

            let err, res = await chai.request(server)
                .post(config.issueIdentityPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });
    });

    describe('Issue identity 3', async () => {

        beforeEach(function (done) {
            this.timeout(5000)
            setTimeout(function () {
                done()
            }, 4000)
        })

        it('Define identity: ', async () => {

            let identityID = await identity.queryIdentity(config.nubID)

            let obj = {
                "type": config.defineIdentityType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": identityID,
                    "mutableTraits": "mutableTraits3:S|num1",
                    "immutableTraits": "immutableTraits3:S|",
                    "mutableMetaTraits": "mutableMetaTraits3:S|num3",
                    "immutableMetaTraits": "immutableMetaTraits3:S|num4"
                }
            }

            let err, res = await chai.request(server)
                .post(config.defineIdentityPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });

        it('Issue identity 3: ', async () => {

            let identityID = await identity.queryIdentity(config.nubID)
            let clsID = await cls.queryClassification("immutableMetaTraits3")

            let obj = {
                "type": config.issueIdentityType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "to": config.testAccountAddress,
                    "fromID": identityID,
                    "classificationID": clsID,
                    "mutableProperties": "mutableTraits3:S|num1",
                    "immutableProperties": "immutableTraits3:S|",
                    "mutableMetaProperties": "mutableMetaTraits3:S|num3",
                    "immutableMetaProperties": "immutableMetaTraits3:S|num4"
                }
            }
            let err, res = await chai.request(server)
                .post(config.issueIdentityPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });
    });
})
