#!/bin/bash
echo "Building baddayduck/services/usersvc:run" 
docker build -t baddayduck/services/usersvc:run . -f Dockerfile.usersvc

echo "Building baddayduck/services/authsvc:run" 
docker build -t baddayduck/services/authsvc:run . -f Dockerfile.authsvc
