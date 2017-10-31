#!/bin/bash
echo "Building baddayduck/services/usersvc:run" 

docker build -t baddayduck/services/usersvc:run . -f Dockerfile.usersvc