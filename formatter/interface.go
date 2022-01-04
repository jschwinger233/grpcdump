package formatter

import "github.com/jhump/protoreflect/dynamic"

type Formatter interface {
	Format(dynamic.Message)
}
