#!/usr/bin/env bash
# Copyright 2018 The Anticycle Authors. All rights reserved.
# Use of this source code is governed by a GPL-style
# license that can be found in the LICENSE file.
#
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
