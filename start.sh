#!/bin/bash
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#
# Exit on first error, print all commands.
set -ev


docker-compose -f docker-compose.yml down

export COMPOSE_PROJECT_NAME="bhoomi"

docker-compose -f docker-compose.yml up -d

# wait for Hyperledger Fabric to start
# incase of errors when running later commands, issue export 
# FABRIC_START_TIMEOUT=<larger number>

export FABRIC_CFG_PATH=$PWD

export FABRIC_START_TIMEOUT=15

#echo ${FABRIC_START_TIMEOUT}

sleep ${FABRIC_START_TIMEOUT}

# Create the channel
docker exec -e "CORE_PEER_LOCALMSPID=MojaniMSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@mojani.cts.com/msp" peer0.mojani.cts.com peer channel create -o orderer.cts.com:7050 -c mychannel -f /etc/hyperledger/configtx/channel.tx

# Join peer0.mojani.cts.com to the channel.
docker exec -e "CORE_PEER_LOCALMSPID=MojaniMSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@mojani.cts.com/msp" peer0.mojani.cts.com peer channel join -b mychannel.block

# fetch the genesis block for joining 2nd peer to channel

docker exec -e "CORE_PEER_LOCALMSPID=KaveriMSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@kaveri.cts.com/msp" peer0.kaveri.cts.com peer channel fetch 0 mychannel.block -o orderer.cts.com:7050 -c mychannel
# join second peer to channel
    docker exec -e "CORE_PEER_LOCALMSPID=KaveriMSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@kaveri.cts.com/msp" peer0.kaveri.cts.com peer channel join -b mychannel.block

#fetch the genesis block to join 3rd peer to channel

docker exec -e "CORE_PEER_LOCALMSPID=FinanceMSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@finance.cts.com/msp" peer0.finance.cts.com peer channel fetch 0 mychannel.block -o orderer.cts.com:7050 -c mychannel

#Join the 3rd peer to the channel
docker exec -e "CORE_PEER_LOCALMSPID=FinanceMSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@finance.cts.com/msp" peer0.finance.cts.com peer channel join -b mychannel.block

# Install chain code on the channel
docker exec -e "CORE_PEER_LOCALMSPID=MojaniMSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/mojani.cts.com/users/Admin@mojani.cts.com/msp" cli peer chaincode install -n bcc -v 1.0 -p github.com/bhoomi


# Instantiate chain code on the channel with a sample record for unit testing
docker exec -e "CORE_PEER_LOCALMSPID=MojaniMSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/mojani.cts.com/users/Admin@mojani.cts.com/msp" cli peer chaincode instantiate -o orderer.cts.com:7050 -C mychannel -n bcc -v 1.0 -c '{"Args":[""]}' -P "OR ('MojaniMSP.member','KaveriMSP.member','FinanceMSP.member')"

sleep 10

# Sample query to test instantiated record query from ledger
docker exec -e "CORE_PEER_LOCALMSPID=MojaniMSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/mojani.cts.com/users/Admin@mojani.cts.com/msp" cli peer chaincode invoke -o orderer.cts.com:7050 -C mychannel -n bcc -c '{"function":"queryLandRecord","Args":["999999999"]}'

