@echo off
md out
del /s /q out\*

logdefine.exe -idir ./ -odir out/go -model go
logdefine.exe -idir ./ -odir out/java -model java
logdefine.exe -idir ./ -odir out/cpp -model cpp
logdefine.exe -idir ./ -odir out/js -model js

pause