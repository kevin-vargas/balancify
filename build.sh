#!/bin/bash

DIRECTORY=$(pwd)
AUTHORIZATION_SERVER="${DIRECTORY}/authorization-server"
BFF_SERVER="${DIRECTORY}/bff-server"
BALANCIFY_SERVER="${DIRECTORY}/balancify"
WEB_SERVER="${DIRECTORY}/web-server"
SECRET_MANIFEST="${DIRECTORY}/deploy/secrets.yaml"
DEPLOY_MANIFEST="${DIRECTORY}/deploy/deployment.yaml"
buildAndPushBackend (){
   local REGISTRY=docker.fast.ar
   local IMAGE_TAG="${REGISTRY}/$1"
   echo $IMAGE_TAG
   echo $2
   docker build -t $IMAGE_TAG $2
   docker push $IMAGE_TAG
}

apply() {
    kubectl apply -f $1
}

log() {
    echo -e "${1}"
}

logt() {
    log "\t ♥♥♥ ${1} \t ♥♥♥"
}

{

    logt "building backend and pusing image"
    buildAndPushBackend authorize $AUTHORIZATION_SERVER
    buildAndPushBackend bff $BFF_SERVER
    buildAndPushBackend balancify $BALANCIFY_SERVER
    buildAndPushBackend web-server $WEB_SERVER

    logt "deploying kubernetes objects"
    apply $SECRET_MANIFEST
    apply $DEPLOY_MANIFEST
}