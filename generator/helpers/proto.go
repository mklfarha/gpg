package helpers

import (
	"fmt"

	"github.com/iancoleman/strcase"
	"github.com/maykel/gpg/entity"
)

func ProtoEnumOptions(protoType string, opts []entity.OptionValue) []string {
	res := []string{strcase.ToScreamingSnake(fmt.Sprintf("%s_%s", protoType, "invalid"))}
	for _, opt := range opts {
		res = append(res, strcase.ToScreamingSnake(fmt.Sprintf("%s_%s", protoType, opt.Identifier)))
	}
	return res
}
