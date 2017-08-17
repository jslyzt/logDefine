@echo off
md out
::del /s /q out\*
logdefine.exe -idir ./ -odir out -model go

md out\java
del /s /q out\java\*
logdefine.exe -idir ./ -odir out/java -model java

md out\cpp
del /s /q out\cpp\*
logdefine.exe -idir ./ -odir out/cpp -model cpp
pause