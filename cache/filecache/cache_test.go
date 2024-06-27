package filecache

import (
	"os"
	"reflect"
	"sync"
	"testing"
)

func TestNewFileCache(t *testing.T) {
	type args struct {
		dir string
	}
	tests := []struct {
		name    string
		args    args
		want    *fileCache
		wantErr bool
	}{
		{
			name: "TestNewFileCache",
			args: args{
				dir: "/Users/yandongxiao/github/yandongxiao/github-image-proxy/cache",
			},
			want: &fileCache{
				cache: sync.Map{},
				dir:   "/Users/yandongxiao/github/yandongxiao/github-image-proxy/cache",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFileCache(tt.args.dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFileCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFileCache() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fileCache_Read(t *testing.T) {
	type fields struct {
		cache sync.Map
		dir   string
	}
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		fields  fields
		data    []byte
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "prepare cache, and try to read",
			fields: fields{
				cache: sync.Map{},
				dir:   "/tmp/cache",
			},
			args: args{
				filePath: "/xxx/test.txt",
			},
			data:    []byte("hw"),
			want:    []byte("hw"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache, err := NewFileCache(tt.fields.dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFileCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			err = cache.Write(tt.args.filePath, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			cacheName := cache.getName(tt.args.filePath)
			cacheFilePath := cache.getCacheFilePath(cacheName)
			_, err = os.Stat(cacheFilePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Stat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			got, err := cache.Read(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Read() got = %v, want %v", got, tt.want)
			}
		})
	}
}
