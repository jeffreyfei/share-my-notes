if [ $1 == "build" ]; then
    mkdir -p bin
    pushd bin
        go build ../server/cmd/server
    popd
elif [ $1 == "up" ]; then
    pushd server
        glide up
    popd
elif [ $1 == "run-server" ]; then
    export PGHOST=/var/run/postgresql
    export PGENV=development
    export PGUSER=postgres
    export SERVER_PORT=8080
    export BASE_URL=http://localhost:8080
    export SESSION_KEY=dev_session_key
    ./bin/share-my-notes
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
