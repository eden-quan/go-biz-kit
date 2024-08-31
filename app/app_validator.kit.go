package apputil

import (
	errorv1 "github.com/eden-quan/go-biz-kit/common/def"
)

// validator ...
type validator interface {
	Validate() error
}

// Validate ...
func Validate(req validator) error {
	if err := req.Validate(); err != nil {
		msg := "Invalid Parameters"
		e := errorv1.ErrorInvalidParameter(msg)
		e.Metadata = map[string]string{"error": err.Error()}
		return e
	}
	return nil
}
