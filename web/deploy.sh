#!/bin/sh

rm -rf dist/*
yarn build || exit

rsync -rcv dist/* one@channel_father:channel/html/
