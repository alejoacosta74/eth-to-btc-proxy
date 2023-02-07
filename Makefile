QTUM_CONTAINER_NAME=regtest

QTUMD_FLAGS = \
	-regtest \
	-rpcbind=0.0.0.0:3889 \
	-rpcallowip=0.0.0.0/0 \
	-logevents \
	-addrindex \
	-reindex \
	-txindex \
	-rpcuser=qtum \
	-rpcpassword=qtum\
	-deprecatedrpc=accounts \
	-printtoconsole

QTUM_CONTAINER_FLAGS = \
	-d \
	--name ${QTUM_CONTAINER_NAME}  \
	-v ${shell pwd}/${QTUM_CONTAINER_NAME}:/root/.qtum \
	-p 3889:3889

.DEFAULT_GOAL := help
.PHONY: all build run test clean help regtest

### This help screen
help:
	@printf "Available targets:\n\n"
	@awk '/^[a-zA-Z\-\_0-9%:\\]+/ { \
	helpMessage = match(lastLine, /^## (.*)/); \
	if (helpMessage) { \
	helpCommand = $$1; \
	helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
	gsub("\\\\", "", helpCommand); \
	gsub(":+$$", "", helpCommand); \
	printf "  \x1b[32;01m%-35s\x1b[0m %s\n", helpCommand, helpMessage; \
	} \
        } \
        { lastLine = $$0 }' $(MAKEFILE_LIST) | sort -u
	@printf "\n"



# check docker is installed
docker_check: 
	@docker info > /dev/null 2>&1 || (echo "Docker is not installed" && exit 1)

# check golang is installed
golang_check:
	@go version info > /dev/null 2>&1 || (echo "Go not found" && exit 1)



# pull qtum image
pull-qtum:
	@ docker pull qtum/qtum:latest > /dev/null || (echo "error pulling qtum docker image" && exit 1) 

seed-qtum:
	@ printf "\n(2) Importing test accounts...\n\n"
	docker cp ${shell pwd}/test/fill_user_account.sh ${QTUM_CONTAINER_NAME}:.

	@ printf "\n(3) Filling test accounts wallets...\n\n"
	docker exec ${QTUM_CONTAINER_NAME} /bin/sh -c ./fill_user_account.sh
	@ printf "\n... Done\n\n"

## Start Qtum regtest docker container
start-regtest: pull-qtum 
	@ printf "\nRunning qtum on docker...\n\n"
	@ docker run ${QTUM_CONTAINER_FLAGS} qtum/qtum qtumd ${QTUMD_FLAGS} > /dev/null || (echo "error running qtumd on docker" && exit 1)


stop-regtest: #-- stop regtest qtum node
	@ printf "\nStopping and removing qtum container...\n\n"
	@ docker stop $(QTUM_CONTAINER_NAME)
	@ docker rm $(QTUM_CONTAINER_NAME)

# build qtum contract
build:
	qtum-cli -regtest -rpcuser=qtum -rpcpassword=testpasswd -rpcport=3889 createcontract 0x$(shell solc --bin --optimize --overwrite -o bin --combined-json bin,abi,metadata contracts/HelloWorld.sol | jq -r '.contracts["contracts/HelloWorld.sol:HelloWorld"].bin')

# run qtum contract
run:
	qtum-cli -regtest -rpcuser=qtum -rpcpassword=testpasswd -rpcport=3889 callcontract $(contract) '["Hello World"]'
#-- Testing
test: ## test qtum contract
	qtum-cli -regtest -rpcuser=qtum -rpcpassword=testpasswd -rpcport=3889 callcontract $(contract) '["Hello World"]' | jq -r '.result.executionResult.output'

## Run integration tests against a docker qtum regtest node
integration-test: 
	@ printf "\nChecking environment for integratioon tests...\n\n"
	make docker_check
	make golang_check
	@ docker ps -a | grep -q ${QTUM_CONTAINER_NAME} || make regtest
	make seed-qtum
	@ printf "\nRunning tests...\n\n"

## clean everything
clean: stop-regtest
	@ printf "\nRemoving files and folders...\n\n"
	@ rm -rf bin
	@ rm -rf regtest