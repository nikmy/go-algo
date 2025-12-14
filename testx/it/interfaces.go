package it

import "github.com/stretchr/testify/require"

type tHelper interface {
	Helper()
}

type testingT = require.TestingT
