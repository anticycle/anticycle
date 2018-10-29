#!/usr/bin/env bash
# Build binaries for each OS and architecture
declare -a OSARCHS=("linux/amd64" "linux/arm" "darwin/amd64" "windows/amd64")
version=($("$(dirname "$0")/version.sh"))
ld=(
    "-X main.version=${version[0]}"
    "-X main.build=${version[1]}"
)
out=$1
in=$2

for osarch in "${OSARCHS[@]}"
do
  echo "Build artifacts: ${osarch}"

  oa=(${osarch//// })  # replace slash to space and split to array
  os_name=${oa[0]}
  os_arch=${oa[1]}
  mkdir -p ./dist/${os_name}

  filename="anticycle_${os_arch}"
  if [[ ${os_name} == "windows" ]]; then
    filename="${filename}.exe"
  fi

  env GOOS=${os_name} GOARCH=${os_arch} go build -ldflags="${ld[*]}" \
                                                 -o ${out}/${os_name}/${filename} ${in}
done
