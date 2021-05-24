#!/usr/bin/env bash

rm -rf ~/Downloads/AStar/*

bee pack -be GOOS=linux

mv AStar.tar.gz ~/Downloads/AStar

tar xf ~/Downloads/AStar/AStar.tar.gz -C ~/Downloads/AStar

rm -f AStar
mv ~/Downloads/AStar/AStar .