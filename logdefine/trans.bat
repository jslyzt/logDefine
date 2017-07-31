@echo off
md out
del /s /q out\*
logdefine.exe -idir ./ -odir out -model go
pause