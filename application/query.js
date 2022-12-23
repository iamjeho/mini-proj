const express = require('express');
const app = express();

// 패브릭연결설정
const fs = require('fs');
const path = require('path');
const FabricCAServices = require("fabric-ca-client");
const { Gateway, Wallets } = require("fabric-network");

// 서버설정
const PORT = 5050;
const HOST = '0.0.0.0';

// use static file
app.use(express.static(path.join(__dirname, 'views2')));

// configure app to use body-parser
app.use(express.json());
app.use(express.urlencoded({ extended: false }));

// / GET index.html 페이지 라우팅
app.get('/', (req, res)=>{
    res.sendFile(__dirname + 'views2/index.html');
})


async function main() {
    try {
        // load the network configuration
        const ccpPath = path.resolve(__dirname, "ccp", "connection-org1.json");
        const ccp = JSON.parse(fs.readFileSync(ccpPath, "utf8"));

        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), "wallet");
        const wallet = await Wallets.newFileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        const identity = await wallet.get("admin");
        if (!identity) {
            console.log(
                'An identity for the user "appUser" does not exist in the wallet'
            );
            console.log("Run the registerUser.js application before retrying");
            return;
        }

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccp, {
            wallet,
            identity: "admin",
            discovery: { enabled: true, asLocalhost: true },
        });

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork("mychannel");

        // Get the contract from the network.
        const contract = network.getContract("basic");

        // Evaluate the specified transaction.
        // queryCar transaction - requires 1 argument, ex: ('queryCar', 'CAR4')
        // queryAllCars transaction - requires no arguments, ex: ('queryAllCars')
/*         const result = await contract.evaluateTransaction('GetAllVoters');
        console.log(
            `Transaction has been evaluated, result is: ${result.toString()}`
        )
        const result2 = await contract.evaluateTransaction('GetAllAssets');
        console.log(
            `Transaction has been evaluated, result is: ${result2.toString()}`
        ) */
        const result1 = await contract.evaluateTransaction('QueryCandidate', '1');

        console.log(
            `Transaction has been evaluated, result is: ${result1}`
        )

        const result2 = await contract.evaluateTransaction('QueryCandidate', '2');

        console.log(
            `Transaction has been evaluated, result is: ${result2}`
        )
        const result3 = await contract.evaluateTransaction('QueryCandidate', '3');
        
        var jsonObj1 = JSON.parse(result1);
        var jsonObj2 = JSON.parse(result2);
        var jsonObj3 = JSON.parse(result3);
        console.log(jsonObj1.VoteCount)
        console.log(jsonObj2.VoteCount)
        console.log(jsonObj3.VoteCount)
        var VoteCount1 = jsonObj1.VoteCount;
        var VoteCount2 = jsonObj2.VoteCount;
        var VoteCount3 = jsonObj3.VoteCount;
        var VoteCount_arr = [VoteCount1, VoteCount2,  VoteCount3];
   
          
  
    } catch (error) {
        console.error(`Failed to evaluate transaction: ${error}`);
        process.exit(1);
    }
}


main();

// 서버시작
app.listen(PORT, HOST);
console.log(`Running on http://${HOST}:${PORT}`);