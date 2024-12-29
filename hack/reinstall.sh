#!/bin/bash

gh extension remove aggregate-issue
go build .
gh extension install .
