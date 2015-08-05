#!/bin/bash

set -e

godep get
godep go build
