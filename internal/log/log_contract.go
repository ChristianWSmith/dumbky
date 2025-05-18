package log

var _initFunc func() = Init
var _debugFunc func(string, ...any) = Debug
var _infoFunc func(string, ...any) = Info
var _warnFunc func(string, ...any) = Warn
var _errorFunc func(string, ...any) = Error
