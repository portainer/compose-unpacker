#!/bin/bash
set -e
set +x

: ${REPO:=$1/compose-unpacker}
: ${TAG:=$2}

docker_image_build_and_push()
{
  arch=${1?required}
  os=${2?required}
  repo=${3?required}
  tag=${4?required}
  build_args=""

  dockerfile="build/linux/Dockerfile"
  if [[ ${os} == "windows" ]]; then
      dockerfile="build/windows/Dockerfile"
      build_args="--build-arg OSVERSION=1809"
  fi

  echo docker buildx build -o type=docker -f ${dockerfile} ${build_args} --platform ${os}/${arch} -t ${repo}:${tag}-${os}-${arch} .
  docker buildx build -o type=docker -f ${dockerfile} ${build_args} --platform ${os}/${arch} -t ${repo}:${tag}-${os}-${arch} .

  echo docker image push ${repo}:${tag}-${os}-${arch}
  docker image push ${repo}:${tag}-${os}-${arch}
}

docker_manifest_create_and_push()
{
  repo=${1?required}
  tag=${2?required}

  images=$(docker image ls "${repo}:${tag}*" --format '{{.Repository}}:{{.Tag}}')

  docker manifest create --amend ${repo}:${tag} $images
  for img in $images; do    
    if [[ "$img" == *"win"* ]]; then
      os="windows"
    else
      os="linux"
    fi

    case ${img} in
      *"amd64"*)
        arch="amd64";;
      *"arm64"*)
        arch="arm64";;
      *"arm32"*)
        arch="arm32";;
      *)
        continue;;
    esac

    docker manifest annotate ${repo}:${tag} ${img} --os ${os} --arch ${arch}
  done
  
  docker manifest push ${repo}:${tag}
}

make clean
make PLATFORM=linux ARCH=amd64
docker_image_build_and_push amd64 linux ${REPO} ${TAG} 

make clean
make PLATFORM=linux ARCH=arm64
docker_image_build_and_push arm64 linux ${REPO} ${TAG} 

make clean
make PLATFORM=linux ARCH=arm
docker_image_build_and_push arm linux ${REPO} ${TAG} 

make clean
make PLATFORM=windows ARCH=amd64
docker_image_build_and_push amd64 windows ${REPO} ${TAG} 

docker_manifest_create_and_push ${REPO} ${TAG}
