// Copyright 2022 Ilya Dmitrichenko
// SPDX-License-Identifier: Apache-2.0

package nested1

import "github.com/errordeveloper/cue-utils/template/testtypes"

defaults: {}
resource: testtypes.#Cluster
template: testtypes.#Cluster
