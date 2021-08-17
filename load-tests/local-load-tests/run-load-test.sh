#!/bin/sh
echo "[run-load-test.sh]: Beginning run-load-test.sh"
# Check for loadtest script
if [ -z "$1" ]; then
    echo "[run-load-test.sh]: Error no load-test script provided"
    echo "[run-load-test.sh]: Usage: ./run-load-test.sh <load-test-file> <# users> <time> <ws>"
    echo "[run-load-test.sh]: Example: ./run-load-test.sh heathcheck/healthcheck-loadtest.js 1 5m"
    echo "[run-load-test.sh]: Example: ./run-load-test.sh blockbynumber/multiple-request-bodies.js 100 10s ws"
    echo "[run-load-test.sh]: Exiting...."
    exit
fi  

# Check for virtual users
if [ -z "$2" ]; then
    echo "[run-load-test.sh]: Error no virtual users provided"
    echo "[run-load-test.sh]: Usage: ./run-load-test.sh <load-test-file> <# users> <time> <ws>"
    echo "[run-load-test.sh]: Example: ./run-load-test.sh heathcheck/healthcheck-loadtest.js 1 5m"
    echo "[run-load-test.sh]: Example: ./run-load-test.sh blockbynumber/multiple-request-bodies.js 100 10s ws"
    echo "[run-load-test.sh]: Exiting...."
    exit
fi  

# Check for time duration 
if [ -z "$3" ]; then
    echo "[run-load-test.sh]: Error no time parameter provided"
    echo "[run-load-test.sh]: Example: ./run-load-test.sh heathcheck/healthcheck-loadtest.js 1 5m"
    echo "[run-load-test.sh]: Example: ./run-load-test.sh blockbynumber/multiple-request-bodies.js 100 10s ws"
    echo "[run-load-test.sh]: Usage: ./run-load-test.sh <load-test-file> <# users> <time> <ws>"
    echo "[run-load-test.sh]: Exiting...."
    exit
fi  

## Check for the ws extension, update route if nessecary
endpoint="http://localhost:8000"
if [ -z "$4" ]; then 
    continue
elif [ "ws" = "$4" ]; then 
    endpoint="http://localhost:8000/ws"
    continue
else
    echo "[run-load-test.sh]: Error invalid websocket parameter provided. Must be ws"
    echo "[run-load-test.sh]: Example: ./run-load-test.sh heathcheck/healthcheck-loadtest.js 1 5m"
    echo "[run-load-test.sh]: Example: ./run-load-test.sh blockbynumber/multiple-request-bodies.js 100 10s ws"
    echo "[run-load-test.sh]: Usage: ./run-load-test.sh <load-test-file> <# users> <time> <ws>"
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

pushd $WORKSPACE > /dev/null 2>&1 
echo "[run-load-test.sh]: Executing make binrun to run infura webserver"
make binrun > /dev/null 2>&1 &
popd > /dev/null 2>&1 
# Let the server spin up
sleep 5

# Use arg 1 and replace the URL with localhost for binary testing
sed 's|SED-URL|'${endpoint}'|g' ${1} > ${temp_loadtest_path}

# Run the test and output to db
echo "[run-load-test.sh]: Running k6 tool"
k6 run --vus ${2} --duration ${3} --out influxdb=http://localhost:8086/k6 ${temp_loadtest_path}

# Remove the file
echo "[run-load-test.sh]: Removing temp loadtest file: ${temp_loadtest_path}"
rm ${temp_loadtest_path}
echo "[run-load-test.sh]: Completed."
