#!/bin/bash

for i in */; do (cd $i/inp; sh submit); done
