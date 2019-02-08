package main

import (
	"fmt"
	"log"
	"path"
	"path/filepath"
	"strings"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

type packageFile struct {
	name string
	pf   []*protoFile
}

func (f *packageFile) addProto(pf *protoFile) {
	f.pf = append(f.pf, pf)
}

func (f *packageFile) protoFile() *protoFile {
	pf := &protoFile{
		Imports:  map[string]*importValues{},
		Messages: []*messageValues{},
		Services: []*serviceValues{},
		Enums:    []*enumValues{},
	}
	for i := range f.pf {
		for j := range f.pf[i].Imports {
			pf.Imports[j] = f.pf[i].Imports[j]
		}
		pf.Messages = append(pf.Messages, f.pf[i].Messages...)
		pf.Services = append(pf.Services, f.pf[i].Services...)
		pf.Enums = append(pf.Enums, f.pf[i].Enums...)
	}
	return pf
}

var (
	packageFiles = map[string]*packageFile{}
)

func addProtoToPackage(fileName string, pf *protoFile) {
	if _, ok := packageFiles[fileName]; !ok {
		packageFiles[fileName] = &packageFile{name: fileName}
	}
	packageFiles[fileName].addProto(pf)
}

func samePackage(a *descriptor.FileDescriptorProto, b *descriptor.FileDescriptorProto) bool {
	if a.GetPackage() != b.GetPackage() {
		return false
	}
	return true
}

func generate(req *plugin.CodeGeneratorRequest) (*plugin.CodeGeneratorResponse, error) {
	resolver := dependencyResolver{}

	res := &plugin.CodeGeneratorResponse{
		File: []*plugin.CodeGeneratorResponse_File{
			{
				Name:    &twirpFileName,
				Content: &twirpSource,
			},
		},
	}

	protoFiles := req.GetProtoFile()
	for _, pf := range protoFiles {
		if pf.GetPackage() == "" {
			return nil, fmt.Errorf("all files must have a package")
		}
		resolver.AddFile(pf)
	}

	for i := range protoFiles {
		file := protoFiles[i]

		pfile := &protoFile{
			Imports:  map[string]*importValues{},
			Messages: []*messageValues{},
			Services: []*serviceValues{},
			Enums:    []*enumValues{},
		}

		// Add enum
		for _, enum := range file.GetEnumType() {
			v := &enumValues{
				Name:   underscoreize(enum.GetName()),
				Values: []*enumKeyVal{},
			}

			for _, value := range enum.GetValue() {
				v.Values = append(v.Values, &enumKeyVal{
					Name:  value.GetName(),
					Value: value.GetNumber(),
				})
			}

			pfile.Enums = append(pfile.Enums, v)
		}

		// Add messages
		for _, message := range file.GetMessageType() {
			pfile.Messages = append(pfile.Messages, processMessage(&resolver, file, pfile, message))
		}

		// Add services
		for _, service := range file.GetService() {
			v := &serviceValues{
				Package:   file.GetPackage(),
				Name:      service.GetName(),
				Interface: typeToInterface(service.GetName()),
				Methods:   []*serviceMethodValues{},
			}

			for _, method := range service.GetMethod() {
				input := fieldTypeName(&resolver, file, descriptor.FieldDescriptorProto_TYPE_MESSAGE, method.GetInputType())
				if input.Import != nil {
					pfile.Imports[input.Import.Name] = input.Import
				}
				output := fieldTypeName(&resolver, file, descriptor.FieldDescriptorProto_TYPE_MESSAGE, method.GetOutputType())
				if output.Import != nil {
					pfile.Imports[output.Import.Name] = output.Import
				}

				v.Methods = append(v.Methods, &serviceMethodValues{
					Name:       method.GetName(),
					InputType:  input.Name,
					OutputType: output.Name,
				})
			}

			pfile.Services = append(pfile.Services, v)
		}

		// Add to appropriate file
		addProtoToPackage(tsFileName(file), pfile)
	}

	for packageName := range packageFiles {
		pf := packageFiles[packageName]

		// Compile to typescript
		content, err := pf.protoFile().Compile()
		if err != nil {
			log.Fatal("could not compile template: ", err)
		}

		// Add to file list
		res.File = append(res.File, &plugin.CodeGeneratorResponse_File{
			Name:    &pf.name,
			Content: &content,
		})
	}

	return res, nil
}

func processMessage(resolver *dependencyResolver, file *descriptor.FileDescriptorProto, pfile *protoFile, message *descriptor.DescriptorProto) *messageValues {
	name := underscoreize(message.GetName())
	tsInterface := typeToInterface(name)
	jsonInterface := typeToJSONInterface(name)

	v := &messageValues{
		Name:          name,
		Interface:     tsInterface,
		JSONInterface: jsonInterface,

		Fields:      []*fieldValues{},
		NestedTypes: []*messageValues{},
		NestedEnums: []*enumValues{},
	}

	if len(message.GetNestedType()) > 0 {
		for _, nt := range message.GetNestedType() {
			pfile.Messages = append(pfile.Messages, processMessage(resolver, file, pfile, nt))
		}
	}

	// Add nested enums
	for _, enum := range message.GetEnumType() {
		e := &enumValues{
			Name:   underscoreize(enum.GetName()),
			Values: []*enumKeyVal{},
		}

		for _, value := range enum.GetValue() {
			e.Values = append(e.Values, &enumKeyVal{
				Name:  value.GetName(),
				Value: value.GetNumber(),
			})
		}

		v.NestedEnums = append(v.NestedEnums, e)
	}

	// Add message fields
	for _, field := range message.GetField() {
		fieldType := fieldTypeName(resolver, file, field.GetType(), field.GetTypeName())
		if fieldType.Import != nil {
			pfile.Imports[fieldType.Import.Name] = fieldType.Import
		}

		v.Fields = append(v.Fields, &fieldValues{
			Name:  field.GetName(),
			Field: field.GetJsonName(),

			Type:       fieldType.Name,
			IsEnum:     field.GetType() == descriptor.FieldDescriptorProto_TYPE_ENUM,
			IsRepeated: isRepeated(field),
		})
	}

	return v
}

func isRepeated(field *descriptor.FieldDescriptorProto) bool {
	return field.Label != nil && *field.Label == descriptor.FieldDescriptorProto_LABEL_REPEATED
}

func upperCaseFirst(s string) string {
	return strings.ToUpper(s[0:1]) + s[1:]
}

func getImport(fd *descriptor.FileDescriptorProto) *importValues {
	return &importValues{
		Name: tsImportName(fd),
		Path: tsImportPath(fd),
	}
}

func underscoreize(s string) string {
	return strings.Replace(s, ".", "_", -1)
}

func tsImportPath(fd *descriptor.FileDescriptorProto) string {
	fileName := tsFileName(fd)
	return fileName[0 : len(fileName)-len(path.Ext(fileName))]
}

func tsImportName(fd *descriptor.FileDescriptorProto) string {
	return underscoreize(fd.GetPackage())
}

func tsFileName(fd *descriptor.FileDescriptorProto) string {
	return filepath.Join(filepath.Dir(fd.GetName()), fd.GetPackage()+".ts")
}

func fieldTypeName(resolver *dependencyResolver, fd *descriptor.FileDescriptorProto, typ descriptor.FieldDescriptorProto_Type, typeName string) *tsType {
	switch typ {
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE,
		descriptor.FieldDescriptorProto_TYPE_FIXED32,
		descriptor.FieldDescriptorProto_TYPE_FIXED64,
		descriptor.FieldDescriptorProto_TYPE_INT32,
		descriptor.FieldDescriptorProto_TYPE_INT64,
		descriptor.FieldDescriptorProto_TYPE_UINT32,
		descriptor.FieldDescriptorProto_TYPE_UINT64:
		return &tsType{Name: "number"}
	case descriptor.FieldDescriptorProto_TYPE_STRING,
		descriptor.FieldDescriptorProto_TYPE_BYTES:
		return &tsType{Name: "string"}
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		return &tsType{Name: "boolean"}
	case descriptor.FieldDescriptorProto_TYPE_ENUM,
		descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		return resolver.TsType(fd, typeName)
	default:
		log.Printf("unknown type %q", typ)
		return &tsType{Name: "string"}
	}

}

func fieldType(f *fieldValues) string {
	t := f.Type
	if t == "Date" {
		t = "string"
	}
	if f.IsRepeated {
		return t + "[]"
	}
	return t
}
