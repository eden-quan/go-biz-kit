package errorutil

var DBTxError *ErrorsCode = NewErrorCode(500, 10000001, "create transaction failed")
var DBGetError *ErrorsCode = NewErrorCode(500, 10000002, "create db connection failed")
