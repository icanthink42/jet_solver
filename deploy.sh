#!/bin/bash

if [ $# -ne 1 ]; then
    echo "Usage: $0 <server_ip>"
    exit 1
fi

SERVER_IP=$1
BINARY_NAME="jet_solver"
REMOTE_PATH="/usr/local/bin"

echo "Building production binary..."
go build -tags prod -o $BINARY_NAME

if [ $? -ne 0 ]; then
    echo "Build failed"
    exit 1
fi

echo "Copying binary to $SERVER_IP..."
scp $BINARY_NAME "ec2-user@$SERVER_IP:/home/ec2-user/"
ssh ec2-user@$SERVER_IP "sudo mv /home/ec2-user/$BINARY_NAME $REMOTE_PATH/ && sudo chown root:root $REMOTE_PATH/$BINARY_NAME"

if [ $? -ne 0 ]; then
    echo "Failed to copy binary"
    rm $BINARY_NAME
    exit 1
fi

echo "Cleaning up..."
rm $BINARY_NAME

echo "Done! Binary deployed to $SERVER_IP:$REMOTE_PATH/$BINARY_NAME"
