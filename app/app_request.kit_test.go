package apputil

import (
	"testing"

	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/stretchr/testify/require"

	headerpkg "github.com/eden/go-kratos-pkg/header"
)

// go test -v -count=1 ./app -test.run=TestRequest_Codec
func TestRequest_Codec(t *testing.T) {
	accept := headerpkg.ContentTypeMultipartForm
	subtype := ContentSubtype(accept)
	t.Log("subtype: ", subtype)
	codec := encoding.GetCodec(subtype)
	require.NotNil(t, codec)
	t.Log("codec: ", codec.Name())
	require.Equal(t, codecNameMultipartForm, codec.Name())

}
