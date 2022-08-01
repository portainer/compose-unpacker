#!/bin/bash
set -e

: ${REPO:=$1/compose-unpacker}
: ${TAG:=$2}

docker_image_build_and_push()
{
    arch=${1?required}
    platform=${2?required}
    repo=${3?required}
    tag=${4?required}
    
    dockerfile="build/linux/Dockerfile"
    if [[ ${platform} == "windows" ]]; then
        dockerfile="build/windows/Dockerfile"
    fi

    docker buildx build -o type=docker -f ${dockerfile} --platform ${platform}/${arch} -t ${repo}:${tag}-${arch} .
    docker image push ${repo}:${tag}-${arch}
}

docker_manifest_create_and_push()
{
  images=$(docker image ls $1 --format '{{.Repository}}:{{.Tag}}')
  echo docker manifest create --amend ${2?required} $images
   for img in $images; do
     docker manifest annotate $2 $1-${img##*-} --os linux --arch ${img##*-}
   done
   docker manifest push $2
}

# echo docker_image_build_and_push amd64 amd64:latest ${TAG} ${REPO} $(dirname $0)/.
docker_image_build_and_push amd64 linux ${REPO} ${TAG} 
docker_image_build_and_push arm64 linux ${REPO} ${TAG} 
#docker_image_build_and_push arm64 windows ${REPO} ${TAG} 
#docker_image_build_and_push amd64 windows ${REPO} ${TAG} 

# docker_image_build_and_push arm64  arm64v8/alpine:latest ${TAG} ${KUBERNETES_RELEASE} ${REPO} $(dirname $0)/.
# docker_image_build_and_push arm    arm32v7/alpine:latest ${TAG} ${KUBERNETES_RELEASE} ${REPO} $(dirname $0)/.

docker_manifest_create_and_push ${REPO} ${REPO}:${TAG}
