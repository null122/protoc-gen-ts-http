package plugin

import (
	"fmt"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
	"net/http"
	"path"
	"regexp"
	"strings"
	"unicode"
)

const (
	ProtoGenContent = `// Code generated by protoc-gen-ts-axios. DO NOT EDIT.
// source: %s

%s

%s`

	IMPORT = `import {%s} from "./%s";`

	ControllerClass = `
export class %s {
    send: <T = any, R = any>({ method, url, data }: { method: string, url: string, data: T }) => Promise<R>;
    fromRequest: <T = any>(data: T) => any;
    fromResponse: <T = any>(data: T) => any;
    constructor(
        send: <T = any, R = any>({ method, url, data }: { method: string, url: string, data: T }) => Promise<R>,
        fromRequest: <T = any>(data: T) => any,
        fromResponse: <T = any>(data: T) => any,
    ) {
        this.send = send;
        this.fromRequest = fromRequest;
        this.fromResponse = fromResponse;
    }
    
    %s
}`
	ControllerMethod = `
    public async %s(data: %s): Promise<%s> {
        const request_data = this.fromRequest(data)
        return new Promise<%s>((resolve, reject) => {
            %s
            this.send({
                method: "%s",
                url: ` + "`" + `%s` + "`" + `,
                data: request_data,
            }).then((data) => {
                resolve(this.fromResponse(data) as %s)
            }).catch((error) => {
                reject(error)
            })
        })
    }
`
	GetUrlParameter = `
            const %s = request_data.%s`
)

type httpRule struct {
	Method string
	Path   string
	Body   string
}

func GenerateFile(request *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error) {

	var res pluginpb.CodeGeneratorResponse

	for _, file := range request.ProtoFile {
		fileName := strings.Replace(file.GetName(), ".proto", "", 1)

		ts := fmt.Sprintf(ProtoGenContent, file.GetName(),
			generateImport(file, fileName),
			generateClass(file, fileName),
		)

		res.File = append(res.File, &pluginpb.CodeGeneratorResponse_File{
			Name:    proto.String(fileName + ".ts"),
			Content: proto.String(ts),
		})
	}

	res.SupportedFeatures = proto.Uint64(uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL))
	return &res, nil
}

// 生成import
func generateImport(file *descriptorpb.FileDescriptorProto, fileName string) string {
	var messageNames []string
	messageList := file.GetMessageType()
	for _, descriptorProto := range messageList {
		messageNames = append(messageNames, descriptorProto.GetName())
	}

	return fmt.Sprintf(IMPORT, strings.Join(messageNames, ","), fileName)
}

// 生成类
func generateClass(file *descriptorpb.FileDescriptorProto, fileName string) string {
	return fmt.Sprintf(ControllerClass, strings.ToUpper(fileName[:1])+fileName[1:]+"Api", generateMethod(file))
}

func generateMethod(file *descriptorpb.FileDescriptorProto) string {
	var methodList []string
	service := file.GetService()
	for _, s := range service {
		mList := s.GetMethod()
		for _, m := range mList {
			input := getTypeName(m.GetInputType())
			output := getTypeName(m.GetOutputType())

			if input == "Empty" {
				input = "undefined"
			}
			if output == "Empty" {
				output = "undefined"
			}

			rule, _ := proto.GetExtension(m.GetOptions(), annotations.E_Http).(*annotations.HttpRule)
			ruleInfo := getHttpRuleInfo(rule)

			parameter := generateUrlParameter(ruleInfo.Path)

			methodList = append(methodList, fmt.Sprintf(ControllerMethod, m.GetName(), input, output, output, parameter, ruleInfo.Method, ruleInfo.Path, output))
		}
	}

	return strings.Join(methodList, "\n")
}

// 生成url常量
func generateUrlParameter(url string) string {
	var parameter []string
	pattern := `\${([^}]+)}`
	regex := regexp.MustCompile(pattern)
	matches := regex.FindAllStringSubmatch(url, -1)
	for _, match := range matches {
		parameter = append(parameter, fmt.Sprintf(GetUrlParameter, match[1], toLowerCamelCase(match[1])))
		//parameter = append(parameter, match[1])
	}

	if len(parameter) == 0 {
		return ""
	}

	return strings.Join(parameter, "")
}

func getHttpRuleInfo(rule *annotations.HttpRule) httpRule {
	var (
		url    string
		method string
		body   string
	)

	switch pattern := rule.Pattern.(type) {
	case *annotations.HttpRule_Get:
		url = pattern.Get
		method = http.MethodGet
	case *annotations.HttpRule_Put:
		url = pattern.Put
		method = http.MethodPut
	case *annotations.HttpRule_Post:
		url = pattern.Post
		method = http.MethodPost
	case *annotations.HttpRule_Delete:
		url = pattern.Delete
		method = http.MethodDelete
	case *annotations.HttpRule_Patch:
		url = pattern.Patch
		method = http.MethodPatch
	case *annotations.HttpRule_Custom:
		url = pattern.Custom.Path
		method = pattern.Custom.Kind
	}
	body = rule.Body
	url = strings.ReplaceAll(url, "{", "${")

	return httpRule{
		Method: method,
		Path:   url,
		Body:   body,
	}
}

func getTypeName(typePath string) string {
	return path.Base(strings.ReplaceAll(typePath, ".", "/"))
}

func toLowerCamelCase(str string) string {
	words := strings.FieldsFunc(str, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

	for i := range words {
		if i == 0 {
			words[i] = strings.ToLower(words[i])
		} else {
			words[i] = strings.Title(words[i])
		}
	}

	return strings.Join(words, "")
}
