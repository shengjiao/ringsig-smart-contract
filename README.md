# ringsig-smart-contract

This smart contract is the chaincode to be deployed in hyperledger fabric network. It uses the unique ring signature (https://github.com/abovemealsky/urs forked from https://github.com/monero-project/urs to resolve some dependency changes).


It offers the following invoke functions:

## configureTopic

## setStage
Payload: {"topic":"Election","stage":"prepare"}

This is used to set different stage of the topic, including prepare, start, finish. In prepare stage, participants can init their public key (urs) on the ledger. In start stage, participants can submit to the smart contract some message signed using urs. In finish stage, no more invoke action is allowed.

## initPublicKey
Payload: {"topic":"Election","uid": "u01", "x":"", "y":""}

This is used to initialize participants urs public key and store them on the ledger in the prepare stage. x and y are big int which compose the URS public key.



To bring up a example fabric network, please refer to https://github.com/hyperledger/fabric-samples
