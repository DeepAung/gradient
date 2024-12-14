package types

import "github.com/DeepAung/gradient/grader-server/proto"

func ProtoResultToChar(res proto.ResultType) (string, bool) {
	switch res {
	case proto.ResultType_COMPILATION_ERROR:
		return "C", true
	case proto.ResultType_PASS:
		return "P", true
	case proto.ResultType_INCORRECT:
		return "-", true
	case proto.ResultType_RUNTIME_ERROR:
		return "X", true
	case proto.ResultType_TIME_LIMIT_EXCEEDED:
		return "T", true
	case proto.ResultType_MEMORY_LIMIT_EXCEEDED:
		return "M", true
	default:
		return "", false
	}
}

func ProtoLanguageToString(language proto.LanguageType) (string, bool) {
	switch language {
	case proto.LanguageType_CPP:
		return "cpp", true
	case proto.LanguageType_C:
		return "c", true
	case proto.LanguageType_GO:
		return "go", true
	case proto.LanguageType_PYTHON:
		return "python", true
	default:
		return "", false
	}
}

func ProtoLanguageToExtension(language proto.LanguageType) (string, bool) {
	switch language {
	case proto.LanguageType_CPP:
		return ".cpp", true
	case proto.LanguageType_C:
		return ".c", true
	case proto.LanguageType_GO:
		return ".go", true
	case proto.LanguageType_PYTHON:
		return ".py", true
	default:
		return "", false
	}
}

func StringToProtoLanguage(str string) (proto.LanguageType, bool) {
	switch str {
	case "cpp":
		return proto.LanguageType_CPP, true
	case "c":
		return proto.LanguageType_C, true
	case "go":
		return proto.LanguageType_GO, true
	case "python":
		return proto.LanguageType_PYTHON, true
	default:
		return proto.LanguageType(0), false
	}
}
