#!/bin/bash

if [ "${BASH_SOURCE-$0}" == "build.sh" ]; then
  echo "cd .."
  cd ..
fi

APP_NAME="db_diff"
ICON_PATH="./coco.jpg"
OUTPUT_DIR="./build"

echo "start package"

mkdir -p ${OUTPUT_DIR}

fyne package -name ${OUTPUT_DIR}/${APP_NAME} \
             -icon ${ICON_PATH} \
             -os windows

if [ $? -eq 0 ]; then
  echo "package success in ${OUTPUT_DIR}"
else
  echo "package failed"
fi
