#!/bin/bash
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#
export PATH=$GOPATH/src/github.com/hyperledger/fabric/build/bin:${PWD}/../bin:${PWD}:$PATH
export FABRIC_CFG_PATH=${PWD}
CHANNEL_NAME=mychannel

# remove previous crypto material and config transactions
rm -fr config/*
rm -fr crypto-config/*

# generate crypto material
cryptogen generate --config=./crypto-config.yaml
if [ "$?" -ne 0 ]; then
  echo "Failed to generate crypto material..."
  exit 1
fi
mkdir config
# generate genesis block for orderer
configtxgen -profile ThreeOrgOrdererGenesis -outputBlock ./config/genesis.block
if [ "$?" -ne 0 ]; then
  echo "Failed to generate orderer genesis block..."
  exit 1
fi

# generate channel configuration transaction
configtxgen -profile ThreeOrgChannel -outputCreateChannelTx ./config/channel.tx -channelID $CHANNEL_NAME
if [ "$?" -ne 0 ]; then
  echo "Failed to generate channel configuration transaction..."
  exit 1
fi

# generate anchor peer transaction
configtxgen -profile ThreeOrgChannel -outputAnchorPeersUpdate ./config/MojaniMSPanchors.tx -channelID $CHANNEL_NAME -asOrg MojaniMSP
if [ "$?" -ne 0 ]; then
  echo "Failed to generate anchor peer update for MojaniMSP..."
  exit 1
fi

# generate anchor peer transaction
configtxgen -profile ThreeOrgChannel -outputAnchorPeersUpdate ./config/KaveriMSPanchors.tx -channelID $CHANNEL_NAME -asOrg KaveriMSP
if [ "$?" -ne 0 ]; then
  echo "Failed to generate anchor peer update for KaveriMSP..."
  exit 1
fi

# generate anchor peer transaction
configtxgen -profile ThreeOrgChannel -outputAnchorPeersUpdate ./config/FinanceMSPanchors.tx -channelID $CHANNEL_NAME -asOrg FinanceMSP
if [ "$?" -ne 0 ]; then
  echo "Failed to generate anchor peer update for FinanceMSP..."
  exit 1
fi
