#!/bin/sh


## The endpoint we wish to query
endpoint="ws://localhost:8000"


echo "[run-load-test.sh]: Beginning run-load-test.sh"
# Check for loadtest script
if [ -z "$1" ]; then
    echo "[run-load-test.sh]: Error no load test script provided script"
    echo "[run-load-test.sh]: Usage: ./run-load-test.sh <load-test-file> <ws>"
    echo "[run-load-test.sh]: Example: ./run-load-test.sh heathcheck/healthcheck-loadtest.js"
    echo "[run-load-test.sh]: Example: ./run-load-test.sh blockbynumber/multiple-request-bodies.js ws"
    echo "[run-load-test.sh]: Exiting...."
    exit
fi  

# Check for time duration 
if [ -z "$LOAD_TEST_ENDPOINT" ]; then
    echo "[run-load-test.sh]: LOAD_TEST_ENDPOINT Environment variable not set, using localhost:8000"
    # Spin up local webserver in background
    pushd $WORKSPACE > /dev/null 2>&1 
    echo "[run-load-test.sh]: Executing make binrun to run infura webserver"
    #make binrun > /dev/null 2>&1 &
    popd > /dev/null 2>&1 
    sleep 3
    else
      echo "[run-load-test.sh]: LOAD_TEST_ENDPOINT: $LOAD_TEST_ENDPOINT"
fi  

## Check for the ws extension, update route if nessecary
if [ -z "$2" ]; then 
    continue
elif [ "ws" = "$2" ]; then 
    endpoint+="/ws"
    continue
else
    echo "[run-load-test.sh]: Error invalid websocket parameter provided. Must be ws"
    echo "[run-load-test.sh]: Example: ./run-load-test.sh heathcheck/healthcheck-loadtest.js"
    echo "[run-load-test.sh]: Example: ./run-load-test.sh blockbynumber/multiple-request-bodies.js ws"
    echo "[run-load-test.sh]: Usage: ./run-load-test.sh <load-test-file> <ws>"
    echo "[run-load-test.sh]: Exiting...."
    exit
fi

if [ -z "$WORKSPACE" ]; then
    echo "[run-load-test.sh]: Error WORKSPACE environment variable not set"
    echo "[run-load-test.sh]: Please source .envrc in root of repository"
    echo "[run-load-test.sh]: Exiting...."
    exit
fi 

echo "[run-load-test.sh]: Using loadtest file: ${1}"
echo "[run-load-test.sh]: Number of Users: ${2}"
echo "[run-load-test.sh]: Time duration: ${3}"
echo "[run-load-test.sh]: Testing with ${endpoint}"

# Create of temporary file loadtest file to SED
temp_sed_filename="sed-loadtest.js"
mydir=`mktemp -d`
temp_loadtest_path="$mydir/$temp_sed_filename"
echo "[run-load-test.sh]: Creating temp loadtest file: ${temp_loadtest_path}"
touch ${temp_loadtest_path}

# Use arg 1 and replace the URL with localhost for binary testing
sed 's|SED-URL|'${endpoint}'|g' ${1} > ${temp_loadtest_path}

# Run the test and output to db
echo "[run-load-test.sh]: Running k6 tool"

# If we are processing data to influx run with appropriate args 
if [ -z "$INFLUX_DB_SETUP" ]; then
    k6 run ${temp_loadtest_path}
    else
    k6 run --out influxdb=http://localhost:8086/k6 ${temp_loadtest_path}
fi  

# Remove the file
echo "[run-load-test.sh]: Removing temp loadtest file: ${temp_loadtest_path}"
rm ${temp_loadtest_path}
echo "[run-load-test.sh]: Completed."
