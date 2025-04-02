@echo off

rem 打开第一个终端并执行命令
start cmd /k "cd D:\GolandProgram\Vnollx && docker build -f dockerfiles/api/Dockerfile -t vnollx-api ."

rem 打开第一个终端并执行命令
start cmd /k "cd D:\GolandProgram\Vnollx && docker build -f dockerfiles/rpc/user/Dockerfile -t vnollx-user ."

rem 打开第一个终端并执行命令
start cmd /k "cd D:\GolandProgram\Vnollx && docker build -f dockerfiles/rpc/file/Dockerfile -t vnollx-file ."

rem 打开第一个终端并执行命令
start cmd /k "cd D:\GolandProgram\Vnollx && docker build -f dockerfiles/rpc/folder/Dockerfile -t vnollx-folder ."

rem 打开第一个终端并执行命令
start cmd /k "cd D:\GolandProgram\Vnollx && docker build -f dockerfiles/rpc/notification/Dockerfile -t vnollx-notification ."