if [ $1 == "build" ]; then
    mkdir -p bin
    pushd bin
        go build ../server/cmd/share-my-notes
    popd
elif [ $1 == "install" ]; then
    pushd server
        glide install
    popd
elif [ $1 == "run-server" ]; then
    ./bin/share-my-notes
else
    echo "Unrecognized command"
fi
