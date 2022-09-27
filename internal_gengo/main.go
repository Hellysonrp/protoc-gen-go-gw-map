// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package internal_gengo is internal to the protobuf module.
package internal_gengo

import (
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"

	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

// SupportedFeatures reports the set of supported protobuf language features.
var SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

// GenerateVersionMarkers specifies whether to generate version markers.
var GenerateVersionMarkers = true

var mapFileForMapCreated = make(map[protogen.GoImportPath]struct{})

// Standard library dependencies.
const (
	databaseSqlDriverPackage = protogen.GoImportPath("database/sql/driver")
)

// GenerateFile generates the contents of a .pb.go file.
func GenerateFile(gen *protogen.Plugin, file *protogen.File) {
	f := newFileInfo(file)
	if len(f.Services) == 0 {
		return
	}

	_, ok := mapFileForMapCreated[file.GoImportPath]
	if !ok {
		filename := file.GeneratedFilenamePrefix + ".pb.gw.map.init.go"
		mf := gen.NewGeneratedFile(filename, file.GoImportPath)

		mf.P("package ", f.GoPackageName)
		mf.P()
		mf.P("type GWRoute struct {")
		mf.P("Method string")
		mf.P("Path string")
		mf.P("}")
		mf.P()
		mf.P("var MapGWRoutes = make(map[GWRoute]string)")
		mf.P()

		mapFileForMapCreated[file.GoImportPath] = struct{}{}
	}

	filename := file.GeneratedFilenamePrefix + ".pb.gw.map.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)

	g.P("package ", f.GoPackageName)
	g.P()
	g.P("func init() {")

	for _, service := range f.Proto.Service {
		genService(g, f, service)
	}

	g.P("}")
}

func genService(g *protogen.GeneratedFile, f *fileInfo, s *descriptorpb.ServiceDescriptorProto) {
	for _, m := range s.GetMethod() {
		methodFullName := "/" + f.Proto.GetPackage() + "." + s.GetName() + "/" + m.GetName()
		options, err := extractAPIOptions(m)
		if err != nil {
			panic(err)
		}

		if options == nil {
			continue
		}

		var method, path string
		switch c := options.Pattern.(type) {
		case *annotations.HttpRule_Get:
			method = "GET"
			path = c.Get

		case *annotations.HttpRule_Post:
			method = "POST"
			path = c.Post

		case *annotations.HttpRule_Put:
			method = "PUT"
			path = c.Put

		case *annotations.HttpRule_Patch:
			method = "PATCH"
			path = c.Patch

		case *annotations.HttpRule_Delete:
			method = "DELETE"
			path = c.Delete

		case *annotations.HttpRule_Custom:
			method = c.Custom.Kind
			path = c.Custom.Path

		}
		g.P("MapGWRoutes[GWRoute{")
		g.P("Method: `", method, "`,")
		g.P("Path: `", path, "`,")
		g.P("}] = `", methodFullName, "`")
	}
}

func extractAPIOptions(meth *descriptorpb.MethodDescriptorProto) (*annotations.HttpRule, error) {
	if meth.Options == nil {
		return nil, nil
	}
	if !proto.HasExtension(meth.Options, annotations.E_Http) {
		return nil, nil
	}
	ext := proto.GetExtension(meth.Options, annotations.E_Http)
	opts, ok := ext.(*annotations.HttpRule)
	if !ok {
		return nil, fmt.Errorf("extension is %T; want an HttpRule", ext)
	}
	return opts, nil
}
