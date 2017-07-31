@echo off
md out
del /s /q out\*
build.exe -idir ./ -odir out -model cpp;go
pause