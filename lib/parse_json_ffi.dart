import 'dart:ffi' as ffi; // For FFI
import 'package:ffi/ffi.dart';

typedef JSONParseFunctionTemplate = ffi.Pointer<Utf8> Function(
    ffi.Pointer<Utf8> json,ffi.Pointer<Utf8> className);

typedef JsonParse = JSONParseFunctionTemplate;

final dylib = ffi.DynamicLibrary.open('./vendor/out/main.dll');

final JsonParse _parseJson = dylib
    .lookup<ffi.NativeFunction<JSONParseFunctionTemplate>>('ParseJson')
    .asFunction();

String jsonParse(String json,String className) {
 return  _parseJson(json.toNativeUtf8(),className.toNativeUtf8()).toDartString();
}