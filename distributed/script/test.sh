#!/bin/bash

go test -v -run TestPut ./test/
go test -v -run TestGet ./test/
go test -v -run TestGetAllVersion ./test/
go test -v -run TestLocat ./test/
go test -v -run Test_putBigFile ./test/
go test -v -run Test_getBigFile ./test/
go test -v -run TestDel ./test/


