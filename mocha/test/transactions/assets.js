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
const {type} = require('os');
var server = config.ip + config.port

var assets = require('../helpers/assets')
var cls = require('../helpers/classifications')
const identity = require("../helpers/identities");

chai.use(chaiHttp);

describe('Assets', async () => {

    describe('Mint Asset', async () => {

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
                    "mutableTraits": "ASSET1:S|" + ",burn:H|1",
                    "immutableTraits": "ASSET2:S|",
                    "mutableMetaTraits": "ASSET3:S|",
                    "immutableMetaTraits": "ASSET4:S|"
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
            let clsID = await cls.queryClassification("ASSET4")

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
                    "mutableProperties": "ASSET1:S|num1" + ",burn:H|1",
                    "immutableProperties": "ASSET2:S|num2",
                    "mutableMetaProperties": "ASSET3:S|num3",
                    "immutableMetaProperties": "ASSET4:S|num4"
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
    });

    describe('Mutate Asset', async () => {

        beforeEach(function (done) {
            this.timeout(5000)
            setTimeout(function () {
                done()
            }, 4000)
        })

        it('Mutate Asset: ', async () => {

            let identityID = await identity.queryIdentity("immutableMetaTraits1")
            let assetID = await assets.queryAsset("ASSET4")

            let obj = {
                "type": config.mutateAssetType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": identityID,
                    "assetID": assetID,
                    "mutableProperties": "ASSET1:S|",
                    "mutableMetaProperties": "ASSET3:S|num3"
                }
            }


            let err, res = await chai.request(server)
                .post(config.mutateAssetPath)
                .send(obj);

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });
    });

    describe('Mint Asset with meta properties', async () => {

        beforeEach(function (done) {
            this.timeout(5000)
            setTimeout(function () {
                done()
            }, 4000)
        })

        it('Mint Asset with meta properties: ', async () => {

            let identityID = await identity.queryIdentity("immutableMetaTraits1")
            let clsID = await cls.queryClassification("ASSET4")

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
                    "mutableProperties": "ASSET1:S|num5" + ",burn:H|1",
                    "immutableProperties": "ASSET2:S|num6",
                    "mutableMetaProperties": "ASSET3:S|num7",
                    "immutableMetaProperties": "ASSET4:S|num8"
                }
            }

            let err, res = await chai.request(server)
                .post(config.mintAssetPath)
                .send(obj);

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });
    });

    describe('Mutate asset non meta properties to meta properties', async () => {

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
                    "mutableTraits": "ASSETS1:S|" + ",burn:H|1",
                    "immutableTraits": "ASSETS2:S|",
                    "mutableMetaTraits": "ASSETS3:S|",
                    "immutableMetaTraits": "ASSETS4:S|"
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
            let clsID = await cls.queryClassification("ASSETS4")

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
                    "mutableProperties": "ASSETS1:S|num5" + ",burn:H|1",
                    "immutableProperties": "ASSETS2:S|num6",
                    "mutableMetaProperties": "ASSETS3:S|num7",
                    "immutableMetaProperties": "ASSETS4:S|num8"
                }
            }

            let err, res = await chai.request(server)
                .post(config.mintAssetPath)
                .send(obj);

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
                    "metaFact": "S|num5"
                }
            }

            let err, res = await chai.request(server)
                .post(config.metaRevealPath)
                .send(obj);

            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });

        it('Mutate Asset non meta properties to meta properties: ', async () => {

            let identityID = await identity.queryIdentity("immutableMetaTraits1")
            let assetID = await assets.queryAsset("ASSETS4")

            let obj = {
                "type": config.mutateAssetType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": identityID,
                    "assetID": assetID,
                    "mutableProperties": "ASSETS1:S|",
                    "mutableMetaProperties": "ASSETS3:S|num5"
                }
            }


            let err, res = await chai.request(server)
                .post(config.mutateAssetPath)
                .send(obj);

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });
    });

    describe('Mint Asset with 22 properties', async () => {

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
                    "mutableTraits": "ASSETP1:S|A,ASSETP11:S|B,ASSETP12:S|C,ASSETP13:S|D,ASSETP14:S|E,burn:H|2",
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
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });

        it('Mint Asset: ', async () => {

            let identityID = await identity.queryIdentity("immutableMetaTraits1")
            let clsID = await cls.queryClassification("ASSETP4")

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
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });
    });

    describe('Add asset properties on mutation', async () => {

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
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });

        it('Mint Asset: ', async () => {

            let identityID = await identity.queryIdentity("immutableMetaTraits1")
            let clsID = await cls.queryClassification("ASSET_A8")

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
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });

        it('Add asset properties on mutation: ', async () => {

            let identityID = await identity.queryIdentity("immutableMetaTraits1")
            let assetID = await assets.queryAsset("ASSET_A8")

            let obj = {
                "type": config.mutateAssetType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": identityID,
                    "assetID": assetID,
                    "mutableProperties": "ASSET_A5:S|AA",
                    "mutableMetaProperties": "ASSET_A7:S|CC"
                }
            }


            var err, res = await chai.request(server)
                .post(config.mutateAssetPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });
    });

    describe('Mint Asset with more than 22 properties', async () => {

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
            expect(res.body.raw_log).to.contain('InvalidRequest')
            expect(res.body.raw_log).to.contain('failed')
        });
    });

    // describe('Mutate Asset to add more that 22 properties', async () => {
    //
    //     beforeEach(function (done) {
    //         this.timeout(6000)
    //         setTimeout(function () {
    //             done()
    //         }, 5000)
    //     })
    //
    //     it('Mutate Asset to add more that 22 properties: ', async () => {
    //
    //         let identityID = await identity.queryIdentity("immutableMetaTraits1")
    //         let assetID = await assets.queryAsset("P4")
    //
    //         let obj = {
    //             "type": config.mutateAssetType,
    //             "value": {
    //                 "baseReq": {
    //                     "from": config.testAccountAddress,
    //                     "chain_id": config.chain_id
    //                 },
    //                 "fromID": identityID,
    //                 "assetID": assetID,
    //                 "mutableProperties": "ASSET_A5:S|A,ASSET_A5:S|B,ASSET_A5:S|C,ASSET_A5:S|D,ASSET_A5:S|E,ASSET_A5:S|F,ASSET_A5:S|G,ASSET_A5:S|H,ASSET_A5:S|I,ASSET_A5:S|J,ASSET_A5:S|K,ASSET_A5:S|L,ASSET_A5:S|M,ASSET_A5:S|N",
    //                 "mutableMetaProperties": "ASSET_A7:S|O,ASSET_A7:S|P,ASSET_A7:S|Q,ASSET_A7:S|R,ASSET_A7:S|S,ASSET_A7:S|T,ASSET_A7:S|U,ASSET_A7:S|V,ASSET_A7:S|W"
    //             }
    //         }
    //         console.log(obj)
    //         var err, res = await chai.request(server)
    //             .post(config.mutateAssetPath)
    //             .send(obj)
    //
    //         res.should.have.status(200);
    //         res.body.should.be.a('object');
    //         expect(res.body.txhash).to.not.equal(null)
    //         expect(res.body.txhash).to.not.equal('')
    //         expect(res.body.raw_log).to.contain('failed')
    //     });
    // });

    describe('Mint Asset with burn greater than forseeable block height', async () => {

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
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });

        it('Mint Asset: ', async () => {

            let identityID = await identity.queryIdentity("immutableMetaTraits1")
            let clsID = await cls.queryClassification("AssetD")

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
                    "metaFact": "H|100000000"
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

        it('Burn Asset', async () => {

            let identityID = await identity.queryIdentity("immutableMetaTraits1")
            let assetID = await assets.queryAsset("AssetD")

            let obj = {
                "type": config.burnAssetType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": identityID,
                    "assetID": assetID
                }
            }

            var err, res = await chai.request(server)
                .post(config.burnAssetPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.contain('failed')
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
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });

        it('Mutate Asset: ', async () => {

            let identityID = await identity.queryIdentity("immutableMetaTraits1")
            let assetID = await assets.queryAsset("AssetD")

            let obj = {
                "type": config.mutateAssetType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": identityID,
                    "assetID": assetID,
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
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')

        });

        it('Burn Asset', async () => {

            let identityID = await identity.queryIdentity("immutableMetaTraits1")
            let assetID = await assets.queryAsset("AssetD")

            let obj = {
                "type": config.burnAssetType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": identityID,
                    "assetID": assetID
                }
            }

            var err, res = await chai.request(server)
                .post(config.burnAssetPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.contain('failed')
        });
    });

    describe('Send splits of an asset and then Mutate ', async () => {
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
                    "mutableTraits": "One:S|" + ",burn:H|1",
                    "immutableTraits": "Two:S|",
                    "mutableMetaTraits": "Three:S|",
                    "immutableMetaTraits": "Four:S|"
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
            let clsID = await cls.queryClassification("Four")

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
                    "mutableProperties": "One:S|One" + ",burn:H|1",
                    "immutableProperties": "Two:S|Two",
                    "mutableMetaProperties": "Three:S|Three",
                    "immutableMetaProperties": "Four:S|Four"
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

        it('Send Asset: ', async () => {

            let identityID = await identity.queryIdentity("immutableMetaTraits1")
            let identityID1 = await identity.queryIdentity("immutableMetaTraits2")
            let assetID = await assets.queryAsset("Four")

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
                    "split":"0.000000000000000001"
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

        it('Mutate Asset: ', async () => {

            let identityID = await identity.queryIdentity("immutableMetaTraits2")
            let assetID = await assets.queryAsset("Four")

            let obj = {
                "type": config.mutateAssetType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": identityID,
                    "assetID": assetID,
                    "mutableProperties": "One:S|One" + ",burn:H|1",
                    "mutableMetaProperties": "Three:S|Three",
                }
            }


            let err, res = await chai.request(server)
                .post(config.mutateAssetPath)
                .send(obj);

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.contain('failed')
        });

        it('Make toID as maintainer: ', async () => {

            let identityID = await identity.queryIdentity("immutableMetaTraits1")
            let identityID1 = await identity.queryIdentity("immutableMetaTraits2")
            let clsID = await cls.queryClassification("Four")

            let obj = {
                "type": config.deputizeType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "toID": identityID1,
                    "classificationID": clsID,
                    "fromID": identityID,
                    "maintainedTraits": "One:S|One,Three:S|Three,burn:H|1",
                    "addMaintainer": true,
                    "removeMaintainer": true,
                    "mutateMaintainer": true
                }
            }


            let err, res = await chai.request(server)
                .post(config.deputizePath)
                .send(obj);

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });

        it('Mutate Asset: ', async () => {

            let identityID = await identity.queryIdentity("immutableMetaTraits2")
            let assetID = await assets.queryAsset("Four")

            let obj = {
                "type": config.mutateAssetType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": identityID,
                    "assetID": assetID,
                    "mutableProperties": "One:S|One" + ",burn:H|1",
                    "mutableMetaProperties": "Three:S|Three",
                }
            }


            let err, res = await chai.request(server)
                .post(config.mutateAssetPath)
                .send(obj);

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.not.contain('failed')
        });
    })

    describe('Mint Asset with extra properties when mutable trait is not defined', async () => {

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
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });

        it('Mint Asset: ', async () => {

            let identityID = await identity.queryIdentity("immutableMetaTraits1")
            let clsID = await cls.queryClassification("ASSETA4")

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
                    "mutableProperties": "ASSETA1:S|A,burn:H|1,ASSETA11:S|B,ASSETA111:S|C",
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
            expect(res.body.raw_log).to.contain('failed')
            expect(res.body.raw_log).to.contain('NotAuthorized')
        });
    });

    describe('Burn Asset', async () => {

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
                    "mutableTraits": "Five:S|" + ",burn:H|1",
                    "immutableTraits": "Six:S|",
                    "mutableMetaTraits": "Seven:S|",
                    "immutableMetaTraits": "Eight:S|"
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
            let clsID = await cls.queryClassification("Eight")

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
                    "mutableProperties": "Five:S|One" + ",burn:H|1",
                    "immutableProperties": "Six:S|Two",
                    "mutableMetaProperties": "Seven:S|Three",
                    "immutableMetaProperties": "Eight:S|Four"
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
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });

        it('Burn Asset', async () => {

            let identityID = await identity.queryIdentity("immutableMetaTraits1")
            let assetID = await assets.queryAsset("Eight")

            let obj = {
                "type": config.burnAssetType,
                "value": {
                    "baseReq": {
                        "from": config.testAccountAddress,
                        "chain_id": config.chain_id
                    },
                    "fromID": identityID,
                    "assetID": assetID
                }
            }
            var err, res = await chai.request(server)
                .post(config.burnAssetPath)
                .send(obj)

            res.should.have.status(200);
            res.body.should.be.a('object');
            expect(res.body.txhash).to.not.equal(null)
            expect(res.body.txhash).to.not.equal('')
            expect(res.body.raw_log).to.not.contain('failed')
            expect(res.body.raw_log).to.not.contain('error')
        });
    });
});
