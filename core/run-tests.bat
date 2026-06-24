@echo off
pushd "%~dp0"
set GINKGO=%GOPATH%\bin\ginkgo.exe

if "%1"=="" goto usage
if "%1"=="help" goto usage

:: List scenarios
if "%1"=="list" (
    set LIST_MODE=true& set DRY_RUN=true
    if "%2"=="" (
        %GINKGO% --no-color --succinct -r ./tests/scenarios 2>nul | findstr /v /c:"Running Suite" /c:"Random Seed" /c:"Will run" /c:"PASS" /c:"Ginkgo ran" /c:"Test Suite" /c:"SUCCESS" /c:"Passed" /c:"seconds" /c:"-----" /c:"BeforeSuite" /c:"AfterSuite" /c:"Starting" /c:"Loaded" /c:"Finished" /c:"No test" /c:"specs" /c:"===" /c:"SS" /c:"+"
    ) else (
        %GINKGO% --no-color --succinct --label-filter="%~2" -r ./tests/scenarios 2>nul | findstr /v /c:"Running Suite" /c:"Random Seed" /c:"Will run" /c:"PASS" /c:"Ginkgo ran" /c:"Test Suite" /c:"SUCCESS" /c:"Passed" /c:"seconds" /c:"-----" /c:"BeforeSuite" /c:"AfterSuite" /c:"Starting" /c:"Loaded" /c:"Finished" /c:"No test" /c:"specs" /c:"===" /c:"SS" /c:"+"
    )
    set LIST_MODE=& set DRY_RUN=
    goto end
)

:: Run tests
if "%1"=="run" (
    set DATA_MODE=%2& set DATA_NAME=%3& set DATA_ENV=%4
    if not exist reports mkdir reports
    goto dorun
)
goto skiprun
:dorun
set RUN_ID=%date:~-4%%date:~3,2%%date:~0,2%_%time:~0,2%%time:~3,2%%time:~6,2%
set RUN_ID=%RUN_ID: =0%
%GINKGO% -v -r --output-dir="reports" --json-report="%DATA_MODE%_%DATA_NAME%_%DATA_ENV%_%RUN_ID%.json" --label-filter="%DATA_NAME% && %DATA_ENV%" ./tests/scenarios
goto end
:skiprun

:: Fast run (cached build)
if "%1"=="fast" (
    set DATA_MODE=%2& set DATA_NAME=%3& set DATA_ENV=%4
    if not exist reports mkdir reports
    goto dofast
)
goto skipfast
:dofast
set RUN_ID=%date:~-4%%date:~3,2%%date:~0,2%_%time:~0,2%%time:~3,2%%time:~6,2%
set RUN_ID=%RUN_ID: =0%
go test -c -o tests\scenarios\test.exe ./tests/scenarios 2>nul
tests\scenarios\test.exe --ginkgo.v --ginkgo.json-report="reports\%DATA_MODE%_%DATA_NAME%_%DATA_ENV%_%RUN_ID%.json" --ginkgo.label-filter="%DATA_NAME% && %DATA_ENV%"
goto end
:skipfast

:: Build
if "%1"=="build" (
    echo Building test binary...
    go test -c -o tests\scenarios\test.exe ./tests/scenarios
    echo Done.
    goto end
)

:usage
echo.
echo Usage:
echo   run-tests.bat list                              List all scenarios
echo   run-tests.bat list "products && sit"            List with filter
echo   run-tests.bat run [mode] [name] [env]           Run tests
echo   run-tests.bat fast [mode] [name] [env]          Fast run (cached)
echo   run-tests.bat build                             Pre-build binary
echo.
echo Examples:
echo   run-tests.bat run service products sit
echo   run-tests.bat run project PROJ-123 sit
echo   run-tests.bat fast service auth sit
echo   run-tests.bat list "products && PROJ-123 && sit"
echo.

:end
popd
