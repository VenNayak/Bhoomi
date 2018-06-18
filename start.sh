#!/bin/bash
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#
# Exit on first error, print all commands.
set -ev


docker-compose -f docker-compose.yml down

docker-compose -f docker-compose.yml up -d ca.mojani.cts.com ca.kaveri.cts.com ca.finance.cts.com orderer.cts.com peer0.mojani.cts.com peer0.kaveri.cts.com  peer1.kaveri.cts.com  peer0.finance.cts.com couchdb cli

# wait for Hyperledger Fabric to start
# incase of errors when running later commands, issue export 
# FABRIC_START_TIMEOUT=<larger number>

export FABRIC_CFG_PATH=$PWD

export FABRIC_START_TIMEOUT=20

#echo ${FABRIC_START_TIMEOUT}

sleep ${FABRIC_START_TIMEOUT}

# Create the channel
docker exec -e "CORE_PEER_LOCALMSPID=MojaniMSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@mojani.cts.com/msp" peer0.mojani.cts.com peer channel create -o orderer.cts.com:7050 -c mychannel -f /etc/hyperledger/configtx/channel.tx

# Join peer0.mojani.cts.com to the channel.
docker exec -e "CORE_PEER_LOCALMSPID=MojaniMSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@mojani.cts.com/msp" peer0.mojani.cts.com peer channel join -b mychannel.block

# Install chain code on the channel
docker exec -e "CORE_PEER_LOCALMSPID=MojaniMSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/mojani.cts.com/users/Admin@mojani.cts.com/msp" cli peer chaincode install -n bhoomi -v 1.0 -p github.com/bhoomi


# Instantiate chain code on the channel with a sample record for unit testing
docker exec -e "CORE_PEER_LOCALMSPID=MojaniMSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/mojani.cts.com/users/Admin@mojani.cts.com/msp" cli peer chaincode instantiate -o orderer.cts.com:7050 -C mychannel -n bhoomi -v 1.0 -c '{"Args":[""]}' -P "OR ('MojaniMSP.member','KaveriMSP.member')"

sleep 10

# Sample query to test instantiated record query from ledger S
docker exec -e "CORE_PEER_LOCALMSPID=MojaniMSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/mojani.cts.com/users/Admin@mojani.cts.com/msp" cli peer chaincode invoke -o orderer.cts.com:7050 -C mychannel -n bhoomi -c '{"function":"queryLandRecord","Args":["999999999"]}'

