"use strict";

//During the test the env variable is set to test
process.env.NODE_ENV = 'test';

class helpers {
    async generateRandomInteger(min, max) {
        return Math.floor(min + Math.random() * (max + 1 - min))
    }

    async delay(interval) {
        return it('should delay', done => {
            setTimeout(() => done(), interval)

        }).timeout(interval + 100) // The extra 100ms should guarantee the test will not fail due to exceeded timeout
    }

    async sleep(ms) {
        return new Promise((resolve) => {
            setTimeout(resolve, ms);
        });
    }
}

module.exports = helpers