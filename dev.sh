if [ $1 == "build" ]; then
    mkdir -p bin
    pushd bin
        go build ../server/cmd/server
        go build ../server/cmd/load_balancer
    popd
elif [ $1 == "up" ]; then
    pushd server
        glide up
    popd
elif [ $1 == "run-server" ]; then
    export LB_PRI_URL=http://localhost:3001
    export LB_PUB_URL=http://localhost:3000
    export PGHOST=/var/run/postgresql
    export PGENV=development
    export PGUSER=postgres
    export SERVER_PORT=8080
    export BASE_URL=http://localhost:8080
    export SESSION_KEY=dev_session_key
    ./bin/server
elif [ $1 == "run-load-balancer" ]; then
    export CLIENT_PORT=3000
    export PROVIDER_PORT=3001
    ./bin/load_balancer
elif [ $1 == "test" ]; then
    export PGHOST=/var/run/postgresql
    export PGENV=testing
    export PGUSER=postgres
    pushd server
        go test -p 1 ./...
    popd
else
    echo "Unrecognized command"
fi
