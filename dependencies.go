package main

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

func fullTypeName(fd *descriptor.FileDescriptorProto, typeName string) string {
	return fmt.Sprintf(".%s.%s", fd.GetPackage(), typeName)
}

type dependencyResolver struct {
	v map[string]*descriptor.FileDescriptorProto
}

func (d *dependencyResolver) AddFile(fd *descriptor.FileDescriptorProto) {
	for _, enum := range fd.GetEnumType() {
		d.set(fd, enum.GetName())
	}
	for _, service := range fd.GetService() {
		d.set(fd, service.GetName())
	}
	for _, message := range fd.GetMessageType() {
		d.addMessage(fd, message)
	}
}

func (d *dependencyResolver) addMessage(fd *descriptor.FileDescriptorProto, message *descriptor.DescriptorProto) {
	name := message.GetName()
	tsInterface := typeToInterface(name)
	jsonInterface := typeToJSONInterface(name)

	d.set(fd, name)
	d.set(fd, tsInterface)
	d.set(fd, jsonInterface)

	message.GetEnumType()

	for _, nm := range message.GetNestedType() {
		*nm.Name = fmt.Sprintf("%s.%s", message.GetName(), nm.GetName())
		d.addMessage(fd, nm)
	}

	for _, ne := range message.GetEnumType() {
		*ne.Name = fmt.Sprintf("%s.%s", message.GetName(), ne.GetName())
		d.set(fd, ne.GetName())
	}
}

func (d *dependencyResolver) set(fd *descriptor.FileDescriptorProto, messageName string) {
	if d.v == nil {
		d.v = make(map[string]*descriptor.FileDescriptorProto)
	}
	typeName := fullTypeName(fd, messageName)

	d.v[typeName] = fd
}

func (d *dependencyResolver) resolve(typeName string) (*descriptor.FileDescriptorProto, error) {
	fp := d.v[typeName]
	if fp == nil {
		panic(fmt.Errorf("missing type %q", typeName))
	}
	return fp, nil
}

type tsType struct {
	Import *importValues
	Name   string
}

func (d *dependencyResolver) TsType(fd *descriptor.FileDescriptorProto, typeName string) *tsType {
	switch typeName {
	case ".google.protobuf.Timestamp":
		// Google WKT Timestamp is a special case here:
		//
		// Currently the value will just be left as jsonpb RFC 3339 string.
		// JSON.stringify already handles serializing Date to its RFC 3339 format.
		return &tsType{Name: "Date"}
	case ".google.protobuf.Struct":
		// Structs map directly to JS objects
		return &tsType{Name: "object"}
	default:
		orig, err := d.resolve(typeName)
		var iv *importValues
		tsName := underscoreize(strings.Replace(typeName, "."+orig.GetPackage()+".", "", -1))
		if err == nil {
			if !samePackage(fd, orig) {
				iv = getImport(orig)
				tsName = fmt.Sprintf("%s.%s", iv.Name, tsName)
			}
		}
		return &tsType{Import: iv, Name: tsName}
	}
}
