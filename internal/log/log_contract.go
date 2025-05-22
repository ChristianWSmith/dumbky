package log

//lint:ignore U1000 Ignore used to maintain parity between debug and release
var _initFunc func() = Init

//lint:ignore U1000 Ignore used to maintain parity between debug and release
var _debugFunc func(string) = Debug

//lint:ignore U1000 Ignore used to maintain parity between debug and release
var _infoFunc func(string) = Info

//lint:ignore U1000 Ignore used to maintain parity between debug and release
var _warnFunc func(error) = Warn

//lint:ignore U1000 Ignore used to maintain parity between debug and release
var _errorFunc func(error) = Error
