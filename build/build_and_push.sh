#!/bin/bash
set -e
set +x

: ${REPO:=$1/compose-unpacker}
: ${TAG:=$2}

BUILDS=('linux/amd64' 'linux/arm64' 'linux/arm' 'windows/amd64')


docker_image_build_and_push()
{
  os=${1?required}
  arch=${2?required}
  repo=${3?required}
  tag=${4?required}

  dockerfile="build/linux/Dockerfile"
  build_args=""
  if [[ ${os} == "windows" ]]; then
      dockerfile="build/windows/Dockerfile"
      build_args="--build-arg OSVERSION=1809"
  fi

  docker buildx build --push -f ${dockerfile} ${build_args} --platform ${os}/${arch} --tag ${repo}:${tag}-${os}-${arch} .
}

docker_manifest_create_and_push()
{
  repo=${1?required}
  tag=${2?required}

  for build in "${BUILDS[@]}"
  do
    IFS='/' read -ra build_parts <<< "$build"
    os=${build_parts[0]}
    arch=${build_parts[1]}

    image="${repo}:${tag}-${os}-${arch}"
    docker manifest create --amend ${repo}:${tag} $image
    docker manifest annotate ${repo}:${tag} ${img} --os ${os} --arch ${arch}
  done  
  
  docker manifest push ${repo}:${tag}
}


for build in "${BUILDS[@]}"
do
  echo "Creating build $build ..."
  IFS='/' read -ra build_parts <<< "$build"
  os=${build_parts[0]}
  arch=${build_parts[1]}

  make clean
  make PLATFORM=${os} ARCH=${arch}
  docker_image_build_and_push ${os} ${arch} ${REPO} ${TAG} 
done


docker_manifest_create_and_push ${REPO} ${TAG}
