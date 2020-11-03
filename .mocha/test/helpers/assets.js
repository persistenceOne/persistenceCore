"use strict";
const helper = require('../helpers/helpers');
const config = require('../config.json');
const request = require('request');
const Promise = require('promise');
//During the test the env variable is set to test
process.env.NODE_ENV = 'test';

async function queryAsset(id) {

    let options = {
        'method': 'GET',
        'url': config.ip + config.port + config.qAsset,
        'headers': {
        }
    };
    return new Promise(function(resolve, reject) {
        request(options, async function (error, res) {
            if (error) throw new Error(error);
            let result = JSON.parse(res.body)
            let list = result.result.value.assets.value.list
            let find = await helper.FindInResponse("assets", list, id)
            let assetID = find.classificationID + "|" + find.hashID
            resolve(assetID)
        });
    });
}


module.exports = {
    queryAsset
};