# ringsig-smart-contract

This smart contract is the chaincode to be deployed in hyperledger fabric network. It uses the unique ring signature (https://github.com/abovemealsky/urs forked from https://github.com/monero-project/urs to resolve some dependency changes).


It offers the following invoke functions:

## configureTopic
TBW

## setStage
Payload: {"topic":"Election","stage":"prepare"}

This is used to set different stage of the topic, including prepare, start, finish. In prepare stage, participants can init their public key (urs) on the ledger. In start stage, participants can submit to the smart contract some message signed using urs. In finish stage, no more invoke action is allowed.

## initPublicKey
Payload: {"topic":"Election","uid": "u01", "x":"1", "y":"2"}

This is used to initialize participants urs public key and store them on the ledger in the prepare stage. x and y are big int which compose the URS public key.

## submit
Payload: {"topic":"Election","msg":"content","sig":{"hsx":"1","hsy":"2","c":["3","4"],"t":["5","6"]},"keyIndex":[{"uid":"u01"},{"uid":"u02"}]}

The urs signature is composed of hsx(big Int), hsy(big Int), c(array of big Int), t(array of big Int). The length of c and t is equal to the ring size. The keyIndex lists the user ids used to construct the urs key ring.

According to properties of urs, we can only verify the signature is someone from the keyIndex list, without knowing who actually signed it. And if a user submits the same result multiple times, it can be linked to some submission in the past and will be rejected. A valid submission will be stored in the ledger.


To bring up a example fabric network, please refer to https://github.com/hyperledger/fabric-samples
