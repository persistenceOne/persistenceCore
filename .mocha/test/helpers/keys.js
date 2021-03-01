const bip39 = require("bip39");
const bip32 = require("bip32");
const fs = require("fs");
const tmSig = require("@tendermint/sig");
const config = require("../config.json");
const crypto = require("crypto");

const passwordHashAlgorithm = "sha512";

function getWallet(mnemonic, bip39Passphrase = "") {
    const seed = bip39.mnemonicToSeedSync(mnemonic, bip39Passphrase);
    const masterKey = bip32.fromSeed(seed);
    const walletPath = getWalletPath();
    return tmSig.createWalletFromMasterKey(masterKey, config.prefix, walletPath);
}

function createRandomWallet(bip39Passphrase = "") {
    const mnemonic = bip39.generateMnemonic(256);
    const walletInfo = getWallet(mnemonic, bip39Passphrase);
    return {
        address: walletInfo.address,
        mnemonic: mnemonic
    };
}

function createWallet(mnemonic, bip39Passphrase = "") {
    let validateMnemonic = bip39.validateMnemonic(mnemonic);
    if (validateMnemonic) {
        const walletInfo = getWallet(mnemonic, bip39Passphrase);
        return {
            address: walletInfo.address,
            mnemonic: mnemonic
        };
    } else {
        return {error: "Invalid mnemonic."}
    }

}

function getWalletPath() {
    return "m/44'/118'/0'/0/0"
}

function createStore(mnemonic, password, filename = "key") {
    try {
        const key = crypto.randomBytes(32);
        const iv = crypto.randomBytes(16);
        let cipher = crypto.createCipheriv("aes-256-cbc", Buffer.from(key), iv);
        let encrypted = cipher.update(mnemonic);
        encrypted = Buffer.concat([encrypted, cipher.final()]);

        fs.writeFileSync(
            filename + ".json",
            JSON.stringify({
                hashpwd: crypto.createHash(passwordHashAlgorithm).update(password).digest("hex"),
                iv: iv.toString("hex"),
                salt: key.toString("hex"),
                crypted: encrypted.toString("hex"),
            })
        );
        return {
            success: true
        };
    } catch (exception) {
        return {
            success: false,
            error: exception.message
        };
    }
}

function decryptStore(file, password) {
    const data = fs.readFileSync(file, { encoding: "utf8", flag: "r" });
    if (
        JSON.parse(data).hashpwd === crypto.createHash(passwordHashAlgorithm).update(password).digest("hex")
    ) {
        let iv = Buffer.from(JSON.parse(data).iv, "hex");
        let encryptedText = Buffer.from(JSON.parse(data).crypted, "hex");

        let decipher = crypto.createDecipheriv(
            "aes-256-cbc",
            Buffer.from(JSON.parse(data).salt, "hex"),
            iv
        );
        let decrypted = decipher.update(encryptedText);
        decrypted = Buffer.concat([decrypted, decipher.final()]);
        return {
            mnemonic: decrypted.toString(),
        };
    } else {
        return {
            error: "Incorrect password."
        };
    }
}

module.exports = {
    createRandomWallet,
    createStore,
};