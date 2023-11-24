package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyTokenString(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDA4NjgzMjQsIm9wYXF1ZV91c2VyX2lkIjoiVTgwNzc4NDMyMCIsInVzZXJfaWQiOiI4MDc3ODQzMjAiLCJjaGFubmVsX2lkIjoiODA3Nzg0MzIwIiwicm9sZSI6ImJyb2FkY2FzdGVyIiwiaXNfdW5saW5rZWQiOmZhbHNlLCJwdWJzdWJfcGVybXMiOnsibGlzdGVuIjpbImJyb2FkY2FzdCIsIndoaXNwZXItVTgwNzc4NDMyMCIsImdsb2JhbCJdLCJzZW5kIjpbImJyb2FkY2FzdCIsIndoaXNwZXItKiJdfX0.HyFRNRcs79Z4Es8DfJHymcw25MQ5FIcbyoswXljI8yQ"
	claims, err := verifyTokenString(token)
	assert.NoError(t, err)
	t.Log(claims)
}
