#!/usr/bin/env bash

DOC_DIR=/usr/share/hub_project/doc

rm -rf ${DOC_DIR}/*

cp index.html /usr/share/hub_project/doc/index.html

cp -R handlers models ${DOC_DIR}/

xdg-open 'http://localhost:9080/hub_project/doc/'
