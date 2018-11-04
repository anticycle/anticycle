#!/usr/bin/env bash
dist=$1
artifacts=($(ls ${dist}))

mkdir -p "$dist/release"
cd ${dist}

for ar in "${artifacts[@]}"
do
    if [[ ${ar} == *"windows"* ]]; then
        cp ${ar} anticycle.exe
        zip -r "release/${ar%.exe}.zip" anticycle.exe
        rm anticycle.exe
    else
        cp ${ar} anticycle
        tar -zcvf "release/$ar.tar.gz" anticycle
        rm anticycle
    fi
done
