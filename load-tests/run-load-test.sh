#!/bin/bash

## Basic default ENDPOINT
ENDPOINT="http://localhost:8000"

## Path to user config file
LOCAL_LOAD_CONFIG="$WORKSPACE/load-tests/load-test-configuration/local-load-config.json"
REMOTE_LOAD_CONFIG="$WORKSPACE/load-tests/load-test-configuration/remote-load-config.json"


echo "[run-load-test.sh]: Beginning run-load-test.sh"
if [ -z "$WORKSPACE" ]; then
    echo "[run-load-test.sh]: Error WORKSPACE environment variable not set"
    echo "[run-load-test.sh]: Please source .envrc in root of repository"
    echo "[run-load-test.sh]: Exiting...."
    exit
fi 

# Check for loadtest script
if [ -z "$1" ]; then
    echo "[run-load-test.sh]: Error no load test script provided script"
    echo "[run-load-test.sh]: Usage: ./run-load-test.sh <load-test-file> <ws>"
    echo "[run-load-test.sh]: Example: ./run-load-test.sh heathcheck/healthcheck-loadtest.js"
    echo "[run-load-test.sh]: Example: ./run-load-test.sh blockbynumber/multiple-request-bodies.js ws"
    echo "[run-load-test.sh]: Exiting...."
    exit
fi  

# EXECUTED_CONFIG will be set by the script, leave it empty
EXECUTED_CONFIG=""

# Check for load test endpoint 
if [ -z "$LOAD_TEST_ENDPOINT" ]; then
    echo "[run-load-test.sh]: LOAD_TEST_ENDPOINT Environment variable not set, using localhost:8000"

    EXECUTED_CONFIG=$LOCAL_LOAD_CONFIG
    ## If running locally and socket2socket replace with http://localhost:8000 with ws://localhost:8000
    if [[ "$1" =~ "socket2socket" ]]; then
        echo "[run-load-test.sh]: Using ws://localhost:8000 for socket2socket locally"
        ENDPOINT="ws://localhost:8000"
    fi

    # Spin up local webserver in background
    pushd $WORKSPACE > /dev/null 2>&1 
    echo "[run-load-test.sh]: Executing make binrun to run infura webserver"
    make binrun > /dev/null 2>&1 &
    popd > /dev/null 2>&1 
    sleep 3

else
    EXECUTED_CONFIG=$REMOTE_LOAD_CONFIG
    ENDPOINT=$LOAD_TEST_ENDPOINT
fi



## Check for the ws extension, update route if nessecary
if [ -z "$2" ]; then 
    : ## Do nothing with colon
elif [ "ws" = "$2" ]; then 
    ENDPOINT+="/ws"
else
    echo "[run-load-test.sh]: Error invalid websocket parameter provided. Must be ws"
    echo "[run-load-test.sh]: Example: ./run-load-test.sh heathcheck/healthcheck-loadtest.js"
    echo "[run-load-test.sh]: Example: ./run-load-test.sh blockbynumber/multiple-request-bodies.js ws"
    echo "[run-load-test.sh]: Usage: ./run-load-test.sh <load-test-file> <ws>"
    echo "[run-load-test.sh]: Exiting...."
    exit
fi


echo "[run-load-test.sh]: Using loadtest file: ${1}"
echo "[run-load-test.sh]: Using config from EXECUTED_CONFIG: $EXECUTED_CONFIG"
echo "[run-load-test.sh]: Testing with ${ENDPOINT}"

# Create of temporary file loadtest file to SED
temp_sed_filename="sed-loadtest.js"
mydir=`mktemp -d`
temp_loadtest_path="$mydir/$temp_sed_filename"
echo "[run-load-test.sh]: Creating temp loadtest file: ${temp_loadtest_path}"
touch ${temp_loadtest_path}

# Use arg 1 and replace the URL with localhost for binary testing
sed 's|SED-URL|'${ENDPOINT}'|g' ${1} > ${temp_loadtest_path}

# Run the test and output to db
echo "[run-load-test.sh]: Running k6 tool"

# If we are processing data to influx run with appropriate args 
influx_cmd=""
if [ "$INFLUX_DB_SETUP" == "true" ]; then
    influx_cmd="--out influxdb=http://localhost:8086/k6"
fi


echo "[run-load-test.sh]: Running Command: k6 run --config ${EXECUTED_CONFIG} ${influx_cmd} ${temp_loadtest_path}"
k6 run --config ${EXECUTED_CONFIG} ${influx_cmd} ${temp_loadtest_path}

# Remove the file
echo "[run-load-test.sh]: Removing temp loadtest file: ${temp_loadtest_path}"
rm ${temp_loadtest_path}
echo "[run-load-test.sh]: Completed."
