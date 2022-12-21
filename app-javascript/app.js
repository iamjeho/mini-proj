/*
 * Copyright IBM Corp. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

"use strict";

const FabricCAServices = require("fabric-ca-client");
const { Gateway, Wallets } = require("fabric-network");
const fs = require("fs");
const path = require("path");

async function main() {
    try {
        // load the network configuration
        const ccpPath = path.resolve(__dirname, "ccp", "connection-org1.json");//ca정보들이 json에 저장되어있음
        const ccp = JSON.parse(fs.readFileSync(ccpPath, "utf8"));// utf8로 코딩해서 ccp에 저장후 밑 함수에서 사용

        // Create a new CA client for interacting with the CA.
        const caInfo = ccp.certificateAuthorities["ca.org1.example.com"];
        const caTLSCACerts = caInfo.tlsCACerts.pem;
        const ca = new FabricCAServices(
            caInfo.url,
            { trustedRoots: caTLSCACerts, verify: false },
            caInfo.caName
        );//ca에 접속하기 위한 정보들

        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), "wallet");
        const wallet = await Wallets.newFileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the admin user.
        const identity = await wallet.get("admin");
        if (identity) {
            console.log(
                'An identity for the admin user "admin" already exists in the wallet'
            );
            return;
        }

        // Enroll the admin user, and import the new identity into the wallet.
        const enrollment = await ca.enroll({//enroll함수로 ca접속
            enrollmentID: "admin",
            enrollmentSecret: "adminpw",
        });
        const x509Identity = {
            credentials: {
                certificate: enrollment.certificate,
                privateKey: enrollment.key.toBytes(), //private key
            },
            mspId: "Org1MSP",
            type: "X.509", //표준
        };
        await wallet.put("admin", x509Identity);
        console.log(
            'Successfully enrolled admin user "admin" and imported it into the wallet'
        );  

        const userIdentity = await wallet.get("appuser13");
        if (userIdentity) {
            console.log(
                'An identity for the user "appuser13" already exists in the wallet'
            );
            return;
        }

        // Check to see if we've already enrolled the admin user.
        const adminIdentity = await wallet.get("admin");
        if (!adminIdentity) {
            console.log(
                'An identity for the admin user "admin" does not exist in the wallet'
            );
            console.log("Run the enrollAdmin.js application before retrying");
            return;
        }

        // build a user object for authenticating with the CA
        const provider = wallet
            .getProviderRegistry()
            .getProvider(adminIdentity.type);
        const adminUser = await provider.getUserContext(adminIdentity, "admin");

        // Register the user, enroll the user, and import the new identity into the wallet.
        const secret = await ca.register( //인증서 발급 요청 -> 이 함수는 발급받기위한 정보 등록
            {
                affiliation: "org1.department1", // 소속기관
                enrollmentID: "appuser13",
                role: "client",
            },
            adminUser //register를 adminUser로 발급 
        );

        const cenrollment = await ca.enroll({ // 발급 register를 했으므로 발급 가능
            enrollmentID: "appuser13",
            enrollmentSecret: secret,
        });
        const cx509Identity = {
            credentials: {
                certificate: cenrollment.certificate,
                privateKey: cenrollment.key.toBytes(),
            },
            mspId: "Org1MSP",
            type: "X.509",
        };
        await wallet.put("appuser13", cx509Identity);
        console.log(
            'Successfully registered and enrolled admin user "appuser13" and imported it into the wallet'
        );
                // Check to see if we've already enrolled the user.
        const identity_user = await wallet.get("appuser13");
        if (!identity_user) {
            console.log(
                'An identity for the user "appuser13" does not exist in the wallet'
            );
            console.log("Run the registerUser.js application before retrying");
            return;
        }
        else{
            console.log("success!!");
        }

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccp, {
            wallet,
            identity: "appuser13",
            discovery: { enabled: true, asLocalhost: true },
        });

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork("mychannel");

        // Get the contract from the network.
        const contract = network.getContract("miniCC");

        // Submit the specified transaction.
        // createCar transaction - requires 5 argument, ex: ('createCar', 'CAR12', 'Honda', 'Accord', 'Black', 'Tom')
        // changeCarOwner transaction - requires 2 args , ex: ('changeCarOwner', 'CAR12', 'Dave')
        await contract.submitTransaction(
            'createAsset',
            '12332323',
            2,
            '김뫄뫄',
            '성내3동 주민센터',
            '2022-12-19-13:44',
        );
        console.log("Transaction has been submitted");

        // Evaluate the specified transaction.
        // queryCar transaction - requires 1 argument, ex: ('queryCar', 'CAR4')
        // queryAllCars transaction - requires no arguments, ex: ('queryAllCars')
        const result = await contract.evaluateTransaction("getAllAssets");
        console.log(
            `Transaction has been evaluated, result is: ${result.toString()}`
        );


        // Disconnect from the gateway.
        await gateway.disconnect();
    } catch (error) {
        console.error(`Failed to enroll admin user "admin": ${error}`);
        process.exit(1);
    }
}

main();