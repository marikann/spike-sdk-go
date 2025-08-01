#!/usr/bin/env bash

#    \\ SPIKE: Secure your secrets with SPIFFE.
#  \\\\\ Copyright 2024-present SPIKE contributors.
# \\\\\\\ SPDX-License-Identifier: Apache-2.0

VERSION="v0.5.19"

git tag -s "$VERSION" -m "$VERSION"
git push origin --tags