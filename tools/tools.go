//go:build tools
// +build tools

package tools

import (
	_ "github.com/bufbuild/buf/cmd/buf"
	_ "github.com/deepmap/oapi-codegen/cmd/oapi-codegen"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
)