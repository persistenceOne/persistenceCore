"use strict";

//During the test the env variable is set to test
process.env.NODE_ENV = 'test';


function generateRandomInteger(min, max) {
    return Math.floor(min + Math.random() * (max + 1 - min))
}

function delay(interval) {
    return it('should delay', done => {
        setTimeout(() => done(), interval)

    }).timeout(interval + 100) // The extra 100ms should guarantee the test will not fail due to exceeded timeout
}

function FindInResponse(type, list, id) {
    let data = {
        'classificationID': '',
        'hashID': ''
    }

    let ordersData = {
        'classificationID': '',
        'makerOwnableID':'',
        'takerOwnableID':'',
        'makerID':'',
        'hashID': ''
    }

    return new Promise(function(resolve, reject) {
        switch (type) {
            case 'assets':
                list.forEach(function (value) {
                    if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString === id) {
                        data.classificationID = value.value.id.value.classificationID.value.idString
                        data.hashID = value.value.id.value.hashID.value.idString
                        resolve(data);
                    }
                });
                break;
            case 'identities':
                list.forEach(function (value) {
                    if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString === id) {
                        data.classificationID = value.value.id.value.classificationID.value.idString
                        data.hashID = value.value.id.value.hashID.value.idString
                        resolve(data);
                    }
                });
                break;
            case 'classifications':
                list.forEach(function (value) {
                    if (value.value.immutableTraits.value.properties.value.propertyList[0].value.id.value.idString === id) {
                        data.chainID = value.value.id.value.chainID.value.idString
                        data.hashID = value.value.id.value.hashID.value.idString
                        resolve(data);
                    }
                });
                break;
            case 'orders':
                list.forEach(function (value) {
                    if (value.value.immutables.value.properties.value.propertyList[0].value.id.value.idString === id) {
                        ordersData.classificationID = value.value.key.value.classificationID.value.idString
                        ordersData.makerOwnableID = value.value.key.value.makerOwnableID.value.idString
                        ordersData.takerOwnableID = value.value.key.value.takerOwnableID.value.idString
                        ordersData.makerID = value.value.key.value.makerID.value.idString
                        ordersData.hashID = value.value.key.value.hashID.value.idString
                        resolve(ordersData);
                    }
                });
                break;
        }
    })
}


module.exports = {
    FindInResponse,
    generateRandomInteger,
    delay
};
