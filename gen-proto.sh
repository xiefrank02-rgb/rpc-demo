@echo off
SET PROTO_DIR=api

:: 查找所有 .proto 文件
for /r %PROTO_DIR% %%f in (*.proto) do (
    echo Generating Go code for %%f
    protoc --go_out=. --go-grpc_out=. --proto_path=api %%f
)

echo Go code generation completed for all .proto files.
pause
