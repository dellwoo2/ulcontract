:START
for %%I IN ("C:\VPMS\AMP\Demo Workshop\Process\*.xls") DO call bin\process.bat %%~nI %%~xI
for /d %%I IN ("C:\VPMS\AMP\Deploy\*") DO call bin\deploy1.bat "%%I"
for /d %%I IN ("C:\VPMS\AMP_UAT\Deploy\*") DO call bin\deploy2.bat "%%I"
timeout 1
GOTO START