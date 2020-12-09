"use strict";
const helper = require('../helpers/helpers');
const config = require('../config.json');
const request = require('request');
const Promise = require('promise');
//During the test the env variable is set to test
process.env.NODE_ENV = 'test';

function queryIdentity(id) {

    let options = {
        'method': 'GET',
        'url': config.ip + config.port + config.qIdentity,
        'headers': {
        }
    };
    return new Promise(function(resolve, reject) {
        request(options, async function (error, res) {
            if (error) throw new Error(error);
            let result = JSON.parse(res.body)
            let list = result.result.list
            let find = await helper.FindInResponse("identities", list, id)
            let identityID = find.classificationID + "|" + find.hashID
            resolve(identityID)
        });
    });
}


module.exports = {
    queryIdentity
};