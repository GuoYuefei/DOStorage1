package test

import (
	"fmt"
	"io"
	"os"
	"testing"
)

const somefile = "somefile.txt"
const fortest = "fortest.txt"
const miya = "IMG_3874-1.jpg"

func TestDel(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "del 1", args: args{name: somefile}, wantErr: false},
		{name: "del 3", args: args{miya}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Del(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("Del() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGet(t *testing.T) {
	type args struct {
		name string
		out	 io.Writer
	}
	out, _ := os.Create("./data-IMG.jpg")
	fmt.Println(out.Stat())
	defer out.Close()
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "get 1", args: args{somefile, nil}, wantErr: false},
		{name: "get by version 1", args: args{somefile +"?version=1", nil}, wantErr: false},
		{name: "get 2", args: args{fortest, nil}, wantErr: false},
		{name: "get 3", args: args{miya, out}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Get(tt.args.name, tt.args.out); (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPut(t *testing.T) {
	type args struct {
		file string
		ok   bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "the correct try1", args: args{somefile,true}, wantErr: false},
		{name: "the incorrect try1", args: args{somefile,false}, wantErr: true},
		{name: "the correct try2", args: args{fortest,true}, wantErr: false},
		{name: "the incorrect try2", args: args{fortest,false}, wantErr: true},
		{name: "the correct try3", args: args{miya, true}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Put(tt.args.file, tt.args.ok); (err != nil) != tt.wantErr {
				t.Errorf("Put() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetAllVersion(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "get all version 0", args: args{somefile}, wantErr: false},
		{name: "get all version 1", args: args{""}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GetAllVersion(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("GetAllVersion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLocat(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "locate 1", args: args{name: fortest}, wantErr: false},
		{name: "locate 3", args: args{miya}, wantErr: false},
	}
		for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Locat(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("Locat() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_putBigFile(t *testing.T) {
	type args struct {
		file string
		size int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"big put test", args{
			file: miya,
			size: 1<<18,
		}, false},
		//{"big put test1", args{
		//	file: somefile,
		//	size: 4,
		//}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := putBigFile(tt.args.file, tt.args.size); (err != nil) != tt.wantErr {
				t.Errorf("putBigFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

