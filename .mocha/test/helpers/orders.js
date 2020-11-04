"use strict";
const helper = require('../helpers/helpers');
const config = require('../config.json');
const request = require('request');
const Promise = require('promise');
//During the test the env variable is set to test
process.env.NODE_ENV = 'test';

async function queryOrder(id) {

    let options = {
        'method': 'GET',
        'url': config.ip + config.port + config.qOrder,
        'headers': {
        }
    };
    return new Promise(function(resolve, reject) {
        request(options, async function (error, res) {
            if (error) throw new Error(error);
            let result = JSON.parse(res.body)
            let list = result.result.value.orders.value.list
            let find = await helper.FindInResponse("orders", list, id)
            let orderID = find.classificationID + "*" + find.makerOwnableID + "*" + find.takerOwnableID + "*" + find.makerID + "*" + find.hashID
            resolve(orderID)
        });
    });
}

module.exports = {
    queryOrder
};