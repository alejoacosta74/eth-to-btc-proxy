#!/bin/bash

set -eEu pipefail
# set -x
trap onError ERR

onError(){
	echo "error: Script failed: see failed command above"
	cleanup
    	exit 1
}

cleanup() {
	if [ "${CLEANUP:-false}" = true ]; then
		# remove qtumd docker container
		if docker ps | grep -q $QTUM_CONTAINER_NAME; then
			echo "Removing qtumd docker container..."
			docker rm -f $QTUM_CONTAINER_NAME
		fi
		# remove ethcli repo
		rm -rf $ETHCLI_DIR
		# remove qtum regtest data
		sudo rm -rf "$QTUM_CONTAINER_NAME"

		QPROXY_PID=$(ps aux | grep qproxy | grep -v grep | awk '{print $2}')
		if [ -n "$QPROXY_PID" ]; then
			echo "Killing qproxy..."
			kill -9 $QPROXY_PID
		fi
	else
		echo "... skipping clean up. Set env CLEANUP=true if you want cleanup to run after test finishes"
		return
	fi
}


spinner() {
	bg_process_name=$1
	bg_process_pid=$!
	spinner_chars="/-\|"

	while ps -p $bg_process_pid > /dev/null; do
	spinner_char=${spinner_chars:0:1}
	spinner_chars=${spinner_chars:1}${spinner_chars:0:1}
	printf "\r%s Waiting for "$bg_process_name" process to complete..." "$spinner_char"
	sleep 0.3
	done
}

# vars for qproxy
QPROXY_ADDRESS=127.0.0.1:8085

# vars for qtumd docker image
QTUM_IMAGE="qtum/qtum"
QTUM_CONTAINER_NAME="regtest"
USER="qtum"
PASSWORD="testpasswd"
QTUMD_FLAGS="qtumd -regtest \
	-rpcbind=0.0.0.0:3889 \
	-rpcallowip=0.0.0.0/0 \
	-logevents \
	-addrindex \
	-reindex \
	-txindex \
	-rpcuser=$USER \
	-rpcpassword=$PASSWORD\
	-deprecatedrpc=accounts \
	-printtoconsole "

QTUM_CONTAINER_FLAGS="-d --name $QTUM_CONTAINER_NAME \
-v $(pwd)/$QTUM_CONTAINER_NAME:/root/.qtum \
-p 3889:3889"

# vars for qproxy
QTUM_CONTAINER_URL="http://localhost:3889"

# vars for ethcli
ETHCLI_REPO="https://github.com/alejoacosta74/ethcli.git"
ETHCLI_BRANCH="master"
ETHCLI_DIR="ethcli"
ETHCLI="ethcli/bin/ethcli"
ETH_NODE_URL="$QPROXY_ADDRESS/rpc"
KEY_B58="cMbgxCJrTYUqgcmiC1berh5DFrtY1KeU4PXZ6NZxgenniF1mXCRk"
KEY="00821d8c8a3627adc68aa4034fea953b2f5da553fab312db3fa274240bd49f35"
QTUM_ADDRESS="qUbxboqjBRp96j3La8D1RYkyqx5uQbJPoW" # qtum address hex: 7926223070547d2d15b2ef5e7383e541c338ffe9

clone_eth_client() {
	if [ ! -d $ETHCLI_DIR ]; then
		echo "Cloning ethcli repo..."
		git clone $ETHCLI_REPO -b $ETHCLI_BRANCH
	fi
	if [ ! -f $ETHCLI ]; then
		cd $ETHCLI_DIR
		make build
		cd ..
	fi
}

run_qtumd() {
	# check if docker is installed
	if ! [ -x "$(command -v docker)" ]; then
		echo 'Error: docker is not installed.' >&2
		exit 1
	fi
	## if qtumd container is running, do nothing
	if [  -n "$(docker ps -qa -f name="$QTUM_CONTAINER_NAME")" ]; then
		echo "(qtumd docker container is already running)"
		return
	else
		# check if qtumd docker image is installed. If not, pull it
		if [ -z $(docker images -q "$QTUM_IMAGE" 2> /dev/null) ]; then
			echo "Pulling qtumd docker image..."
			docker pull $QTUM_IMAGE
			if [ $? -ne 0 ]; then
				echo "Error: could not pull qtumd docker image"
				exit 1
			fi
		fi
		# start qtumd docker container
		echo "Running qtumd docker container..."
		docker run $QTUM_CONTAINER_FLAGS $QTUM_IMAGE $QTUMD_FLAGS > /dev/null
		if [ $? -ne 0 ]; then
			echo "Error: could not run qtumd docker container"
			exit 1
		fi
		# seed regtest addresses with QTUM
		sleep 5
		echo "Seeding regtest addresses with QTUM..."
		docker cp fill_user_account.sh "$QTUM_CONTAINER_NAME":.
		(docker exec "$QTUM_CONTAINER_NAME" /bin/sh -c ./fill_user_account.sh ) > /dev/null &
		spinner "seeding"
		echo "...done seeding regtest addresses with QTUM"
	fi
}

run_qproxy() {
	#     if qproxy is already running, kill it
	set +eE
	if [ -n "$(ps aux | grep qproxy | grep -v grep)" ]; then
		kill -9 "$(ps aux | grep qproxy | grep -v grep | awk '{print $2}')"
	fi
	set -eE
	cd ..
	make build > /dev/null
	echo "Starting qproxy on $QPROXY_ADDRESS ..."
	./bin/qproxy --qtumrpc "$QTUM_CONTAINER_URL" --user $USER --pass $PASSWORD --address "$QPROXY_ADDRESS" &
	sleep 5
	if [ -z "$(ps aux | grep qproxy | grep -v grep)" ]; then
		echo "Error: qproxy failed to start"
		return 1
	fi
	cd test
}

# hex2dec() {
# 	local dec=$(echo "ibase=16; $(echo "$1" | tr '[:lower:]' '[:upper:]')" | bc)
# 	echo "$dec"
# }

compare_big_ints() {
	int1="$1"
	int2="$2"

	epsilon_ratio="0.001"
	diff=$(echo "scale=10; $int1-$int2" | bc)
	if [ $diff -lt 0 ]; then
	diff=$(echo "scale=10; $diff*-1" | bc)
	fi

	if (( $(echo "$int1 == 0" | bc -l) )); then
	ratio=$(echo "scale=10; $diff/$int2" | bc)
	elif (( $(echo "$int2 == 0" | bc -l) )); then
	ratio=$(echo "scale=10; $diff/$int1" | bc)
	else
	if (( $(echo "$int1 > $int2" | bc -l) )); then
	ratio=$(echo "scale=10; $diff/$int1" | bc)
	else
	ratio=$(echo "scale=10; $diff/$int2" | bc)
	fi
	fi

	if (( $(echo "$ratio < $epsilon_ratio" | bc -l) )); then
	echo "true"
	else
	echo "false"
	fi
}

test_send() {
	# run integration tests for eth_sendRawTransaction
	echo ""
	echo "Running integration test for eth_sendRawTransaction..."
	local AMOUNT=25000000000000000 # 0.025 QTUM = 0.025 ETH in wei
	local ADDRESS="0x7cB57B5A97eAbe94205C07890BE4c1aD31E486A8" # qtum address qUvnMQ4KUEM3ZqAT4AwKUB513pSgYyeoF7
	local balance="$($ETHCLI getbalance $ADDRESS -n "$ETH_NODE_URL" | grep -o '[0-9]\+')" 

	echo "...initial balance for 0x$ADDRESS (receiver): $balance wei"

	$ETHCLI importkey $KEY -n "$ETH_NODE_URL" > /dev/null
	echo "... sending $AMOUNT wei to $ADDRESS"
	$ETHCLI send $AMOUNT $ADDRESS -k $KEY -n "$ETH_NODE_URL" > /dev/null
	local new_balance="$($ETHCLI getbalance $ADDRESS -n "$ETH_NODE_URL" | grep -o '[0-9]\+')"


	echo "...new balance for 0x$ADDRESS (receiver): $new_balance wei"

	want="$( bc <<<"$balance + $AMOUNT" )"
	echo "...expected new balance: $want wei" ; echo ""

	if [ "$(compare_big_ints $want $new_balance)" = "true" ]; then
		echo "eth_sendRawTransaction integration test passed"
	else
		echo "Error: eth_sendRawTransaction integration test failed"
		SUCCESS=false
	fi

}

# test_getbalance() {
# 	local balance_hex="$($ETHCLI getbalance $ETH_ADDRESS -n "$ETH_NODE_URL" | grep -o '[0-9]\+')"
# 	local balance=$(hex2dec $balance_hex)
# 	# check if balance is greater than 0
# 	if [ $balance -ne 0 ]; then
# 		echo "Error: eth_getBalance integration test failed"
# 		SUCCESS=false
# 	else
# 		echo "eth_getBalance integration test passed"
# 	fi
# }

# 5. main
main () {
	SUCCESS=true
	clone_eth_client
	run_qtumd
	run_qproxy
	test_send

	if [ "$SUCCESS" = true ]; then
		echo "All integration tests passed"
		cleanup
		return 0
	else
		echo "Some integration tests failed"
		return 1
	fi
	
}

main
