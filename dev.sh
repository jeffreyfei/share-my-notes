if [ $1 == "build" ]; then
    mkdir -p bin
    pushd bin
        go build ../server/cmd/share-my-notes
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
    ./bin/share-my-notes
else
    echo "Unrecognized command"
fi
