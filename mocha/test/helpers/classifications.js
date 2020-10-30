"use strict";
const config = require('../config.json');
const helper = require('../helpers/helpers');
const request = require('request');
const Promise = require('promise');
//During the test the env variable is set to test
process.env.NODE_ENV = 'test';

function queryClassification(id) {

    let options = {
        'method': 'GET',
        'url': config.ip + config.port + config.qClassification,
        'headers': {
        }
    };
    return new Promise(function(resolve, reject) {
        request(options, async function (error, res) {
            if (error) throw new Error(error);
            let result = JSON.parse(res.body)
            let list = result.result.value.classifications.value.list
            let find = await helper.FindInResponse("classifications", list, id)
            let clsID = find.chainID + '.' + find.hashID
            resolve(clsID)
        });
    });
}

module.exports = {
    queryClassification
};