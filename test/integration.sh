#!/bin/bash

set -euox pipefail
trap onError ERR

onError(){
	echo "error: Script failed: see failed command above"
	cleanup
    	exit 1
}

ETH_CLIENT_REPO="https://github.com/alejoacosta74/ethcli.git"
ETH_CLIENT_BRANCH="master"
ETH_CLIENT_DIR="ethcli"
ETH_CLIENT_BIN="ethcli"

# 1. parse arguments
args () {
    while getopts "hvdnstcp" opt; do
	case $opt in
	    h)
		echo "Usage: integration.sh [testname]"
		echo "  testname: name of the test to run, if not specified, run all tests"
		echo "  -h: print this help message"
		echo "  -v: verbose mode"
		echo "  -d: debug mode"
		echo "  -n: dry run mode"
		echo "  -s: skip cleanup"
		echo "  -c: cleanup only"
		echo "  -t: test only"
		echo "  -p: print test name only"
		exit 0
		;;
	    v)
		set -x
		;;
	    d)
		DEBUG=1
		;;
	    n)
		DRYRUN=1
		;;
	    s)
		SKIP_CLEANUP=1
		;;
	    c)
		CLEANUP_ONLY=1
		;;
	    t)
		TEST_ONLY=1
		;;
	    p)
		PRINT_TEST_NAME_ONLY=1
		;;
	    \?)
		echo "Invalid option: -$OPTARG" >&2
		exit 1
		;;
	esac
    done
    shift "$((OPTIND-1))"
    TEST_NAME=${1:-}
}



# 2. git clone ethereum client repo and build it
clone_eth_client () {
    if [ ! -d $ETH_CLIENT_DIR ]; then
	git clone $ETH_CLIENT_REPO -b $ETH_CLIENT_BRANCH
    fi
    cd $ETH_CLIENT_DIR
    make
    cd ..
}


# 3. run test
run_test () {
    if [ -n "${TEST_NAME:-}" ]; then
	if [ -n "${PRINT_TEST_NAME_ONLY:-}" ]; then
	    echo $TEST_NAME
	else
	    $TEST_NAME
	fi
    else
	for test in $(declare -F | grep -o 'test_[^ ]*'); do
	    if [ -n "${PRINT_TEST_NAME_ONLY:-}" ]; then
		echo $test
	    else
		$test
	    fi
	done
    fi
}

# 4. cleanup
cleanup () {
    if [ -n "${SKIP_CLEANUP:-}" ]; then
	return
    fi
    rm -rf $ETH_CLIENT_DIR
}

# 5. main
main () {

    args "$@"
    if [ -z "${TEST_ONLY:-}" ]; then
	clone_eth_client
    fi
    run_test
    if [ -z "${CLEANUP_ONLY:-}" ]; then
	cleanup
    fi
}

main "$@"

# Path: test/test_01.sh

# 1. test_01: test ethcli help
test_01 () {
    $ETH_CLIENT_BIN help
}



# read cli arguments from json file
read_cli_args () {
    local json_file=$1
    local cli_args=$(cat $json_file | jq -r '.cli_args | join(" ")')
    echo $cli_args
}

# run cli with args
run_cli () {
    local cli_args=$1
    $ETH_CLIENT_BIN $cli_args
}
# Path: test/test_02.sh