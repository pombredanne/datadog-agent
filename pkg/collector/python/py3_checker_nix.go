// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2018 Datadog, Inc.

// +build python,!windows

package python

import (
	"path/filepath"
)

const (
	pythonBin = "python"
)

var (
	pythonPath = filepath.Join("..", "..", "embedded", "bin", pythonBin)
)
