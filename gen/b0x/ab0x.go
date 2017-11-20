// Code generated by fileb0x at "2017-11-20 23:41:10.823433421 +0200 IST m=+0.011011642" from config file "b0x.yml" DO NOT EDIT.

package b0x


import (
  "bytes"
  "compress/gzip"
  "io"
  "log"
  "net/http"
  "os"
  "path"

  "golang.org/x/net/webdav"
  "golang.org/x/net/context"


)

var ( 
  // CTX is a context for webdav vfs
  CTX = context.Background()

  
  // FS is a virtual memory file system
  FS = webdav.NewMemFS()
  

  // Handler is used to server files through a http handler
  Handler *webdav.Handler

  // HTTP is the http file system
  HTTP http.FileSystem = new(HTTPFS)
)

// HTTPFS implements http.FileSystem
type HTTPFS struct {}



// FileAPIGoTpl is "api.go.tpl"
var FileAPIGoTpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\x8c\x90\xcd\x4e\xeb\x30\x10\x85\xd7\xd7\x4f\x31\xca\x2a\xa9\xae\xd2\x67\x40\xa5\x48\x91\x0a\xa4\x24\xac\x10\x0b\x13\x4e\xda\x40\xfe\x34\x99\x2e\xa2\xc8\xef\x8e\x3a\x4e\x2a\x44\x59\xb0\xfb\xfc\x79\x66\x3c\x3e\xbd\x2d\x3e\xed\x01\x34\x4d\x71\xea\xd1\x39\x63\xaa\xa6\xef\x58\x28\x34\xff\x82\xb2\x91\xc0\x18\x22\xa2\x60\x9a\xe2\x7c\xec\x11\x27\x7a\x9b\x5a\x39\x3a\x17\x98\xc8\x98\xf5\x9a\x6e\xd2\x84\xaa\x81\xe4\x08\xaa\x5a\x01\x97\xb6\x00\x75\xa5\x8a\xc7\xa7\x7b\xea\xde\x3e\x50\x88\x91\xb1\x87\xaf\xbd\x14\x4d\x3a\x7c\xc3\xb0\x82\x30\xa2\x55\xee\x51\x6d\x86\x1a\x85\xa8\xf5\xa8\x36\x69\x07\xb0\xb7\x1e\xd5\x3e\xf7\xef\xcb\x04\x8f\x6a\x6f\x51\x63\xb6\x1e\xbf\x4d\x58\xfe\xf3\x60\x1b\x38\x17\xae\x96\xf3\xdd\xa9\xae\xbd\xfb\xf5\x89\xbf\xb5\xcd\x3b\x38\x4d\x67\x7f\x02\x57\xe0\xeb\x84\xca\x8e\xc9\x52\xb6\xdd\x6d\x37\x39\x65\xfb\x1d\x0d\x62\x05\x0d\xda\x39\xaa\x4b\xe3\x8f\xb8\xca\x46\xe2\x4c\xb8\x6a\x0f\x60\x15\xe7\xc2\x31\x8c\x28\x7c\x79\xbd\xde\xe7\x3f\x81\xb9\xe3\xe8\xbc\xcd\x57\x00\x00\x00\xff\xff\x1c\x4b\xea\x72\xf1\x01\x00\x00")

// FileAPICommonGo is "api_common.go"
var FileAPICommonGo = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\x3c\xcb\x41\x0e\x82\x30\x10\x05\xd0\x75\xe7\x14\x93\xae\x20\x31\x70\x0a\x77\x6e\x94\x13\x8c\xf0\x8b\x8d\x85\x96\xe9\x90\x68\x8c\x77\x37\x6e\xd8\xbf\x57\x64\x7c\xca\x0c\xb6\x92\x2a\x51\x5c\x4a\x56\xe3\x86\x9c\x9f\xc4\xe4\x2e\x15\x7d\xdd\x92\x27\xe7\xc3\x62\x9e\x5a\xa2\xbe\xe7\xf3\x0b\x23\x94\x63\x65\x7b\x80\xe3\x6a\xd0\x20\x23\x38\x64\xe5\xe1\x7a\xe1\xbd\x4c\x62\xe0\x5c\xa0\x62\x31\xaf\x95\xec\x5d\x70\xb4\xc3\x7f\xc8\x85\xc5\xba\xc1\x34\xae\x33\x94\xdc\x5f\x34\x2d\x37\x75\x4b\xdd\x0d\x75\x4f\x76\x62\xa8\x66\x6d\xe9\x4b\xbf\x00\x00\x00\xff\xff\x6f\x29\xd7\x44\xaa\x00\x00\x00")

// FileCommonGo is "common.go"
var FileCommonGo = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\x34\x8d\xb1\x8e\x83\x30\x10\x44\x6b\xef\x57\x8c\x5c\x81\xc4\x1d\xa7\x2b\x91\xac\x7c\x41\x9a\xb4\x51\x0a\x0b\x2d\xc8\x02\x16\xc7\x5e\xd2\x44\xfc\x7b\x04\x84\x6e\x57\xf3\xde\x4c\xf4\xed\xe0\x7b\x86\xc6\x31\x13\x85\x29\xce\x49\x61\xb3\xa6\x20\x7d\xb6\x44\xdd\x22\x2d\x9e\x57\x9f\x86\x5c\x08\x82\x68\x89\x23\xc4\x9b\x4c\xe8\x20\x70\x0e\x7f\xdb\x63\x12\xeb\x92\x04\xd6\x92\x59\xc9\xec\x0e\x1a\xf7\xc5\xf3\xef\x8d\x23\x7b\x2d\xec\xa5\x82\xad\x20\xe5\x89\xb8\xa3\xfe\xde\x8c\x2c\xc5\x7e\x96\x3f\xff\x0f\xd4\x35\x12\x4f\xf3\x8b\x31\xfa\xac\x9b\x62\xe9\x9c\xd8\x29\x5a\xe9\x13\x00\x00\xff\xff\x56\x56\x3e\xf3\xbd\x00\x00\x00")

// FileCreateGoTpl is "create.go.tpl"
var FileCreateGoTpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\x24\x8d\x51\x0e\x82\x30\x10\x44\xbf\xed\x29\xe6\x4f\xf0\x03\x0e\xe1\x05\x4c\xe4\x00\xac\x75\xa1\x8d\xb0\x92\xed\xa2\x31\x4d\xef\x6e\xb4\x7f\x6f\x26\x93\x79\x1b\xf9\x07\xcd\x8c\x9c\xbb\x4b\xc5\x52\x9c\x9b\x76\xf1\x68\x3c\x4e\xc3\x59\x99\x8c\x5b\x5c\x4d\xa3\xcc\x4d\x8b\xf4\x07\x64\x07\x00\x7d\x8f\x3a\x40\x32\x32\x5e\x59\x0c\x81\x12\x08\x4b\x14\xc6\xf4\x54\x30\xf9\x80\x17\x69\xa4\xdb\xc2\x78\x47\x0b\x88\x76\x4c\x10\x5a\x19\x24\xf7\x9a\xec\xb3\x71\xe7\x0e\xca\xb6\xab\x60\xcc\xb9\x1b\x7e\x4d\x3d\xaf\xee\x52\x46\x57\xdc\x37\x00\x00\xff\xff\x74\x86\x73\x3b\xb0\x00\x00\x00")

// FileCreateCommonGo is "create_common.go"
var FileCreateCommonGo = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\x3c\x8d\xcd\x6a\xeb\x30\x10\x46\xd7\x9e\xa7\xf8\x30\x5c\x22\x85\xe0\xec\x03\x77\x51\x42\x96\xa5\xe0\xe6\x05\x64\x79\x62\x9b\xca\x96\x33\x1a\x87\x96\x92\x77\x2f\xfe\x21\xdb\x99\xef\x9c\x33\x3a\xff\xe5\x1a\x86\x8e\x21\x11\x75\xfd\x18\x45\x61\x28\xcb\x6b\xa7\xae\x72\x89\x8f\xe9\x1e\x72\xb2\x44\xc7\x23\xae\x67\x61\xa7\x8c\x2e\xc1\x21\xa9\x4c\x5e\xa1\xad\x53\xb4\x31\xd4\x09\x33\x82\x5b\x14\x68\xcb\x38\x97\x97\xb7\xeb\x05\x49\x9d\x72\xcf\x83\x92\xfe\x8c\xfc\x32\x6c\xf0\x2f\x65\x97\x6f\xf6\x2c\x94\x45\xe9\xb1\xff\x28\xdf\xe9\xb9\xa4\xe6\x33\xfc\x32\x9e\x63\xea\xaa\xc0\x2f\x77\xd3\x3d\x78\xd8\x1c\x74\x9b\x06\x0f\xe3\xb1\xdf\xdc\x76\x61\x8d\x85\x49\xf7\x50\x94\x9c\xa6\xa0\x07\xb0\x48\x14\x3b\x07\x93\xf6\x8a\xd3\x7f\xf8\xe2\x53\xa5\x1b\x1a\x63\x29\xf3\x45\x94\xbe\x08\xb1\x31\xf9\x2a\x39\x61\xf7\xef\xb1\xcb\x0f\x98\xd7\x96\x32\x61\x9d\x64\xc0\xba\xab\xab\x62\x49\xac\xbf\x27\xfd\x05\x00\x00\xff\xff\x4a\x1c\x8b\x76\x44\x01\x00\x00")

// FileDeleteGoTpl is "delete.go.tpl"
var FileDeleteGoTpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\x34\xcd\x4d\x0e\x82\x30\x10\xc5\xf1\xb5\x73\x8a\x17\x56\xa0\x49\xb9\x04\x75\x85\xf1\x83\x5e\xa0\x92\x11\x8c\xd8\x90\x76\x88\x31\x4d\xef\x6e\x94\xba\x7b\x8b\x7f\x7e\x6f\xb6\xfd\xc3\x0e\x8c\x18\xd5\x69\x9d\x29\x11\xd5\x35\x3a\xf1\x77\x37\xc0\xb3\x2c\xde\x05\xc8\xc8\xe8\xce\x2d\x1a\xdd\x6a\xa3\x11\xc4\x0a\x3f\xd9\x09\xc2\xaf\xa3\xdb\xe2\x7a\x94\x01\x5b\xd3\xf0\xc4\xc2\x55\x06\xca\x2a\x17\x88\xb4\x59\x31\x14\x19\xd9\x5f\x8e\x87\xef\xb1\x79\xcf\xac\x8c\xbd\x4e\x9c\x12\x0a\xec\x10\xd4\x6b\x64\xcf\xea\x4f\x50\x22\xfa\x04\x00\x00\xff\xff\x3d\x00\x23\x24\xa9\x00\x00\x00")

// FileDeleteCommonGo is "delete_common.go"
var FileDeleteCommonGo = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\x6c\x90\x41\x6b\xdb\x40\x10\x85\xcf\x9a\x5f\xf1\x30\x84\x68\x8d\x59\xdd\x03\x39\x94\xd6\xb7\x96\x42\x1c\xe8\x79\xa3\x9d\xca\xa2\xab\x5d\x79\x76\x64\x37\x94\xfc\xf7\xb2\x5a\xc7\xb9\xe4\x24\x31\xfb\xde\xfb\xde\xcc\xec\xfa\x3f\x6e\x60\xe8\x1c\x32\xd1\x38\xcd\x49\x14\x2d\x35\x1b\xef\xd4\xbd\xb8\xcc\x5d\x3e\x85\x0d\x19\xa2\xae\xc3\x81\x03\xf7\x8a\x31\x43\x8f\x8c\xac\xb2\xf4\x0a\x3d\x3a\xc5\x31\x05\x5f\xa7\x87\xfd\xf7\xfd\xd7\x67\x14\x3b\xe9\xeb\xcc\x78\xfe\xc6\x81\xf5\x26\xff\x47\xcd\xfe\x2f\xf7\x2c\xd4\x24\x99\x00\x6c\x7f\x3e\xfd\xa0\xe6\x72\x64\x61\x6c\x7f\x95\x0f\xbd\xad\xb4\xf5\x1f\x6e\x9e\xc3\xc8\x19\x55\xd0\xa7\xe8\x47\x1d\x53\xcc\x48\x71\xe5\x9d\x16\x96\x57\xfa\xbd\xc4\x1e\xad\xc7\xf6\x4a\x33\xd5\xdd\x5e\xae\x91\xe6\xf6\x52\x0a\x78\x5b\xd3\x1e\x71\xa1\x46\x58\x17\x89\xf0\x57\x6a\x29\x07\x59\x62\xdd\xc6\xbf\x77\x77\xca\x13\x47\x2d\x54\x87\x61\x3c\x73\xc4\xfb\x85\xec\x27\xf4\x92\xd2\x1a\xb4\xf9\x14\xec\x13\xe7\x25\xe8\x0e\x2c\x92\xc4\x14\x7e\xd7\xa1\x17\x76\x25\xb8\x5e\xf4\x96\x4f\x4d\xd6\x49\xf1\xf0\x08\x6f\x0f\x2a\x63\x1c\x5a\x43\x8d\x93\x21\xd7\xd9\xda\xdb\x7e\x91\x21\x97\xb9\xb7\x49\x26\x1b\xd2\xd0\x6e\x2a\xf8\x01\xf7\x77\xe7\x7b\xdc\x9d\x37\x3b\x94\xa0\x1d\x8a\xd5\x7c\x2c\xb9\x1a\xfc\x8b\x5d\xfb\x7d\x28\xac\xb5\x86\xde\xe8\x7f\x00\x00\x00\xff\xff\x78\xc4\xfc\x61\x0d\x02\x00\x00")

// FileInsertGoTpl is "insert.go.tpl"
var FileInsertGoTpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\x84\x92\x41\x6f\xd4\x30\x10\x85\xcf\xf1\xaf\x78\x8a\x8a\x64\xaf\x82\xf7\x8e\xc4\x01\xa1\x22\x2d\xa2\xdb\xaa\x09\x5c\xa9\x9b\x4c\xb6\xa6\x89\x1d\x6c\xa7\x12\xb2\xf2\xdf\x91\xe3\x2c\x82\x52\xc4\x6d\x32\x13\x7f\xef\xcd\xb3\x27\xd5\x3e\xaa\x13\x21\x46\x79\x93\xcb\x65\x61\x4c\x8f\x93\x75\x01\x9c\x01\x40\xd9\x8f\xa1\xcc\x95\x0f\x4e\x9b\x93\x2f\x59\xfe\x8c\x51\x36\x3f\x26\x92\x87\xf5\xf7\x1b\x15\x1e\x96\xa5\x64\x82\xb1\x7e\x36\x2d\xb8\xc6\xae\x39\x18\x4f\x2e\x08\xd4\xeb\x49\x2e\x90\x11\x88\xac\x70\x14\x66\x67\xd0\x8f\x41\xd6\x93\xd3\x26\xf4\xfc\xee\x70\xac\x2f\x6f\x1b\x1c\x8e\xcd\x35\xce\xf4\x46\xdd\x0f\xb4\x2c\xe0\xaf\xbc\xc0\x97\x77\x9f\x3e\x5f\xd6\x6b\x7d\x57\xb1\xa2\xd8\x1c\xc9\x8f\x56\x1b\xae\x65\x6b\x07\x5f\xa1\xac\x50\x8a\x34\xfd\x7e\xa5\xdc\xa3\xe7\x03\xa5\xd9\x93\x1a\x66\xf2\x22\x0d\x04\x5b\x18\xdb\xef\x91\xdd\x9d\x85\x8e\x6a\x4c\x3a\xad\x23\x15\xc8\x43\x19\x6c\x76\x7c\x50\x81\x46\x32\x01\xaa\x6d\xad\xeb\xd2\x02\xc1\x22\x3c\x10\x4e\xfa\x89\x0c\xec\xfd\x37\x6a\xc3\xb6\xb6\xc5\xee\xfa\xf6\x4a\xbc\x08\xe7\x13\x76\xe7\xce\x87\x79\x18\x72\x57\xfc\x0a\x2a\xe5\xa2\xf1\xe6\x2d\xac\xcc\x0d\x2e\x58\x11\xe3\x6b\x38\x65\x4e\x84\x8b\xaf\x15\x2e\xfa\x34\xdf\x10\x9a\x86\xce\x2f\x0b\x2b\xb4\x54\x5d\xc7\xcb\x18\x2f\x7a\xf9\xde\x0e\xf3\x68\x32\xba\xac\x30\xc9\xb5\xbb\x49\x65\x1c\x99\x2e\x9d\xda\xae\x40\xa7\x38\x62\xfc\x8f\xc6\x7e\x8f\x9a\xc2\x6f\x2c\x78\x0a\x1e\x6b\xac\xe8\xad\x43\xbb\xea\xe2\x2f\x0f\xd0\x66\xcd\xea\x79\x9a\x2f\x3d\x93\x3f\x04\x78\x66\xaf\x9d\xe4\xe5\x79\x52\x67\xfb\xff\xde\x7d\x05\xa4\xdb\x8e\x31\xaf\xfc\x33\x00\x00\xff\xff\x5c\x4d\x98\xcc\xf0\x02\x00\x00")

// FileInsertCommonGo is "insert_common.go"
var FileInsertCommonGo = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\x64\x91\x41\x6b\xe3\x30\x10\x85\xcf\x9a\x5f\xf1\x30\x94\xda\xc5\x28\x7b\x2e\xe4\x98\x43\x0f\xbb\x0b\x69\x6f\xa5\x07\xd5\x1e\x27\x62\x65\xc9\x95\xc6\x61\x97\x34\xff\x7d\x91\xec\x86\x2c\x7b\x13\xc3\xbc\x37\xdf\x7b\x9a\x4c\xf7\xcb\x1c\x18\x32\xb9\x44\x64\xc7\x29\x44\x41\x4d\xaa\xea\x8d\x98\x77\x93\x78\x93\x3e\x5c\x45\xaa\x1a\x46\xa9\xa8\x21\xda\x6c\xf0\xf2\xe4\x13\x47\x81\x4d\x30\x48\x12\xe7\x4e\x20\x01\xc7\xe0\x7a\x58\x3f\x84\x38\x1a\xb1\xc1\x63\x08\x11\xc6\xe3\xe9\xc7\xf3\x6e\xff\x82\x24\x46\x78\x64\x2f\x24\x7f\x26\xbe\x9a\xac\xfa\x33\xa9\xdd\x6f\xee\x38\x92\x0a\x71\x04\x80\x87\x9f\xfb\xef\xa4\xba\xe0\x12\x80\xd7\xb7\x24\xd1\xfa\x03\xa9\x93\x71\x33\x27\xbc\xbe\x59\x2f\x1c\x07\xd3\xf1\xf9\x42\x97\xc2\x95\x0d\x60\x8b\x6d\x82\x1c\x19\x39\x43\x26\xcb\xef\x83\x3d\xb1\xc7\x57\x2a\x1a\x66\xdf\xa1\xb6\x78\x58\x39\x9a\xa2\xae\x1b\xd4\xe9\xc3\xe9\x3d\xa7\xd9\x49\x0b\x8e\x31\xc4\x26\xc3\xd9\x01\x8e\x7d\x6d\x75\x06\x6a\xb0\xdd\xe2\x1b\x3e\x3f\xd7\xd9\xc2\xb4\x4e\xcf\xa4\x54\x64\x99\xa3\x87\xb7\xae\xc5\x30\x8a\xde\x65\x9f\xa1\xae\x7c\x90\xa3\xf5\x87\xcc\xb4\x70\x56\x0d\xa9\x0b\x91\x4a\x32\x0a\x1e\xb7\xb0\xfa\xb9\xe4\xac\x1b\x52\x56\x87\x38\x6a\x17\x0e\x75\xb5\x30\x3e\xe2\xfe\xee\x74\x8f\xbb\x53\xd5\x22\x0b\x5a\x5c\x4f\xd3\xd7\xc9\x45\xd4\xbf\xeb\x12\xe7\xdf\x2d\xad\x75\x93\xab\xfa\x3f\xbb\xe9\xfb\xda\x9b\x91\xb1\xb4\xdc\xa2\x08\x70\x53\x71\x73\xdd\x2e\x6d\x94\x1a\xb0\x85\x99\x26\xf6\xfd\x5a\x4b\x8b\xec\x51\xc8\xd7\x5f\xba\x59\x58\x26\xab\xf3\x0d\x2f\x5d\xe8\x6f\x00\x00\x00\xff\xff\x6e\x72\xc1\xa0\x84\x02\x00\x00")

// FileLogGo is "log.go"
var FileLogGo = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\x5c\x8c\x41\x0a\xc3\x20\x10\x45\xd7\xce\x29\xfe\xd2\x94\x32\x37\xc8\x0d\x5a\x0a\xbd\x81\x84\xc9\x20\x15\x95\xd1\x2c\x4a\xf0\xee\x25\x84\x6c\xba\xfd\xef\xbf\x57\xc3\xf2\x09\x2a\xe8\x35\x35\xa2\xfe\xad\x82\x47\x51\x15\xc3\xba\xe5\xc5\xb7\x6e\x31\xeb\x1d\xcc\x1c\x73\x17\x5b\xc3\x22\xfb\x98\x88\x0e\x0a\x5f\x70\x7b\xbd\x9f\x13\x52\x51\xdf\x70\x9d\x83\x69\xfb\x37\xb0\x93\x8b\x2b\x0a\xa7\xb3\x3e\xcf\xc8\x31\x1d\xab\x33\xe9\x9b\x65\x72\x83\xdc\x85\x7d\x3b\x2b\xcc\x3c\xd1\xa0\x5f\x00\x00\x00\xff\xff\x26\x19\xa5\x52\xa5\x00\x00\x00")

// FileOpGo is "op.go"
var FileOpGo = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\x2c\xca\x31\x0a\x02\x31\x10\x46\xe1\xda\x39\xc5\xcf\x56\x5a\xed\x05\xb2\x76\x41\xc4\xe0\x22\x9e\x20\x84\xb0\x84\xd5\x64\x4c\xa6\xf1\xf6\x12\x33\xe5\xc7\x7b\xec\xc3\xee\xb7\x08\xe1\x57\x23\x9a\x67\xac\x8c\xd4\xe0\x33\x9e\x0f\x87\x50\xde\xec\x6b\x6a\x25\xa3\x70\xac\x5e\x52\xc9\x24\x5f\x8e\x7d\x6b\x52\x53\xde\x88\x42\xc9\x4d\x70\xa4\xc3\xca\xf6\x03\x60\xc1\xb4\x4c\x5d\xf7\x38\x64\xce\x7f\x5e\x64\x50\x65\x55\xe3\x75\x1a\xcd\x90\x46\xa3\x31\xed\xb1\xd3\x5d\x6f\x76\xa2\x13\xfd\x02\x00\x00\xff\xff\xc5\x67\xbb\x45\xb5\x00\x00\x00")

// FileOrmGo is "orm.go"
var FileOrmGo = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\xac\x93\x41\x8f\xda\x30\x14\x84\xcf\xf8\x57\x8c\xf6\x50\x25\x68\x15\xee\x48\x3d\x14\x92\xc3\x4a\xec\xae\x1a\xe8\x0f\x70\x92\x47\x4a\x1b\xec\x60\xbf\x88\xae\xa2\xfc\xf7\xca\xb1\xa1\xd9\x48\x5b\xf5\x50\x2e\xa0\xd1\x9b\xf9\xc6\xf6\xa3\x95\xe5\x4f\x59\x13\xb8\x6d\xac\x10\xa7\x73\xab\x0d\xe3\xa1\x92\x2c\x0b\x69\x69\x65\x2f\xcd\x83\x10\xfc\xd6\x12\xd2\x0d\x4e\x8a\xc9\x1c\x65\x49\xe8\xc5\x22\xfb\x45\x65\x74\xe9\xc8\xbc\xc1\xb2\x39\xa9\xfa\x11\xd2\xd4\x16\x49\x92\xdc\xe7\xfa\x21\x46\x64\x2f\x4d\x92\x93\xed\x1a\x7e\x04\x19\xa3\x4d\x2c\x16\x5f\x9d\xef\x9f\xdc\xcb\xd1\xae\xaf\xf6\x6e\x1e\x84\x38\x76\xaa\xc4\x0b\x5d\xa3\xaa\x40\xba\x89\xb1\x7c\xcd\x9f\x5d\x27\x43\xdc\x19\x85\x4f\xaf\xf9\x73\x5f\x15\x6b\x54\xc5\xe0\xc6\xc7\xfe\x6e\xc4\xb2\xe9\x4a\x76\x93\x55\x01\xf7\x49\x37\x62\xd1\xe8\xba\x26\x83\xdd\xf8\xe5\xc6\x57\x2b\x6c\x0d\x49\x26\xf8\x3c\x0b\x79\x73\x1e\xb5\x81\xc4\x36\xcf\xbe\x1c\x32\x58\x96\x4c\x67\x52\xec\xfb\x44\x7a\xec\x11\x07\x73\x14\x63\x79\x08\x39\x93\x6a\x41\xea\xb5\x39\xaf\xa1\x87\xc0\xdb\x53\x43\x25\xff\xe1\x29\xe8\xe2\x87\x53\x58\xa3\xf4\x19\x12\xfb\x6c\x97\x6d\x0f\x1f\x62\x7d\xc6\x88\x0d\x71\x53\xac\x97\x66\xd8\x27\x65\xc9\x4c\xb0\x50\x74\xc5\xd3\xcb\x3e\xcb\x3f\xc6\x78\xcf\x88\x09\xf6\x29\xc6\x4b\xff\x01\xf3\xad\xad\x6e\x97\xe8\x7f\xbe\xc3\x78\xa9\x0f\xb8\x35\x66\xdc\x1b\x38\xa5\x86\xa6\xcf\x78\xbf\x56\xff\x8e\x69\xb6\xcb\xfe\xf2\x8e\xde\x3d\x56\x08\x41\xd3\x0a\x5e\x9a\x9d\xd4\x2f\x11\x2c\xb1\x3b\x66\xd8\x2c\xd6\xe0\xef\x7e\x01\xc3\xdf\x6d\x06\xf2\xae\xe8\xdd\x22\xc6\x0e\xa6\x93\xa0\x7d\x0e\x59\x62\x10\xbf\x03\x00\x00\xff\xff\x71\x90\x5a\xa6\xb2\x03\x00\x00")

// FilePageGo is "page.go"
var FilePageGo = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\x7c\x8f\xc1\x4a\x03\x31\x10\x86\xcf\x99\xa7\xf8\x09\x58\xba\x52\xb6\x1e\xc4\x5b\x8f\x16\x0a\x15\x95\xf5\x05\x96\x75\x52\x83\x4d\x36\x26\xb3\x87\xb2\xec\xbb\x4b\xb2\xab\xe2\x41\x8f\x93\xf9\xff\xef\xcb\x84\xb6\x7b\x6f\x4f\x0c\x09\xe7\x44\x64\x5d\xe8\xa3\x40\x1b\x27\x9a\x68\xbb\xc5\x53\xde\x45\x0e\x91\x13\x7b\x49\x68\x3d\x9a\xe7\x23\x8e\x87\x87\xc3\x0b\x92\xb4\xc2\x8e\xbd\x90\x5c\x02\xcf\xd9\x24\x71\xe8\x04\x23\xa9\xb3\x75\x56\x00\xeb\xe5\xee\x96\x54\x6f\x4c\x62\x59\xa6\xa9\xb0\x1b\x89\xd6\x9f\x10\x59\x86\xe8\x13\xe4\x8d\x0b\xfb\x63\xe0\x78\xf9\x71\xb6\x62\x7b\x8f\xde\x94\x7d\x56\x90\x19\x7c\x87\x75\xc0\x75\x9e\xaa\x05\xb3\xae\xb2\x3a\xf3\x46\x52\xd6\x20\xd4\xb3\x7f\xb7\xc3\x0d\x56\x2b\x84\x7a\xf9\x41\x79\x18\x49\xa9\xd9\x0b\xad\x49\x4d\x4b\xe5\x8f\x88\x71\x52\x37\x21\x5a\x2f\x66\xad\xe7\xd3\xaf\x5e\xf5\xe6\xcb\x51\x15\xc0\x7f\x59\x3c\xee\xf7\xcd\xfd\xef\xd6\xe6\xdb\x57\xd1\x44\x9f\x01\x00\x00\xff\xff\x24\xfe\x03\xc8\x87\x01\x00\x00")

// FileSelectGoTpl is "select.go.tpl"
var FileSelectGoTpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\x8c\x94\x41\x6f\xdb\x38\x10\x85\xcf\xe2\xaf\x78\x2b\xa4\xad\x58\xa8\x74\xf7\x9a\xc2\x87\x22\x48\x81\x5d\xb4\xdd\x76\x9d\x9b\x61\x2c\x18\x69\xa4\x10\xa5\x28\x95\xa4\xd3\x14\x82\xfe\xfb\x82\xa4\x2c\x3b\xae\x0f\xcd\x29\x22\x67\x1e\xdf\xc7\xe1\xf3\x20\xab\x6f\xb2\x25\x8c\xa3\xf8\x92\xfe\x9d\x26\xc6\x54\x37\xf4\xd6\xa3\x60\x59\xee\xbc\x55\xa6\x75\x39\x63\x00\x90\x8f\xa3\xb8\xfb\x39\x90\xf8\x2b\x56\x7c\x91\xfe\x61\x9a\x72\xc6\x19\x5b\xad\xb0\x89\xa5\xb0\xe4\xf7\xd6\x38\xf8\x07\xc2\xe6\xeb\x47\x7c\xdf\x93\xfd\x89\xa4\xc3\x9a\xbd\xa9\x50\x38\xbc\xbe\xdb\x90\xa6\xca\xf3\xb9\xab\xe0\x73\x05\xc6\x78\x50\x12\x99\xd7\x9c\xf8\xbb\x57\xa6\xd8\xee\xd2\x67\xaa\x88\x76\x36\xb7\x1f\x6f\x6f\xee\xf2\x12\x4e\xb8\xa8\x77\x50\x2b\x91\x7f\xf8\xf7\x9f\x4f\x38\xf8\xbd\x93\xf7\x9a\xa6\x29\x2f\x97\x5e\x27\x7e\x3c\x90\x25\xb1\x74\x9c\xec\x0c\xb2\x3d\xdf\x98\x4a\xe4\xc8\x39\x63\x53\x64\xbd\x7d\xa2\x0a\x76\x3f\x63\x7e\x8d\x88\xbd\x81\x44\xab\x1e\xc9\xa0\x96\x5e\xde\x4b\x47\xe2\x02\x70\x2c\x2e\x38\x8a\xed\xee\xe0\xee\xc3\x5e\xeb\xcf\xb2\xa3\x69\x2a\x41\xd6\xf6\x96\x63\x64\xd9\x6a\x85\xca\x92\xf4\x84\xc4\x06\xe7\xa5\xa7\x8e\x8c\x67\x99\xf3\x9d\xc7\xf5\x1a\x6e\xb1\xc9\x32\x69\x5b\x97\xd6\x12\xd9\x7b\xdb\xba\xb0\xee\x44\x6f\x3b\xa1\xfb\xb6\xc8\xe3\xd9\xd7\x78\xf5\xe2\xf1\x15\x5e\x3c\x86\x7b\xf3\x9d\x2f\x11\x3a\x39\xcb\x6c\xff\xc3\x45\x03\x49\x25\x74\xd5\xf7\x22\xf9\x3d\x16\x0a\x21\x38\xcb\x54\x13\x0b\xff\x58\xc3\x28\x1d\xdc\x66\xf3\xc8\x8c\xd2\x51\x83\x65\x13\xcb\x6a\x6a\xc8\x22\xe8\x8a\x1b\xdd\x3b\x2a\x38\x8b\x5c\xf4\xe4\xad\xac\x7c\xdc\x81\xef\xc3\x9c\xf7\x95\xdf\x5b\x72\x2c\x7b\x94\x16\x52\x6b\x5c\xba\x1e\x96\x35\xfd\xac\xf7\x99\x9e\x7c\x11\xef\x29\x76\x28\x5c\xaa\x3e\xd8\xbc\x5e\xa7\xa6\x4d\x25\x4d\xe1\x84\xab\xa4\x89\xb7\xf3\x52\xf1\x80\xf3\xee\x9c\xe5\x57\x98\x40\x93\x05\x5b\x6b\xc8\x61\x20\x53\x17\x52\xeb\x12\x8a\x47\xce\xb9\x3c\x2e\xc5\x83\x6e\xad\x2d\x78\x78\x2b\xe3\x08\x2b\x4d\x4b\xb8\xfa\xaf\xc4\x55\x13\x9c\xcc\x36\x15\xe9\xda\xe1\xcd\x34\xc5\xec\xc4\x09\x8f\xe3\x55\x23\x92\x77\xbc\xaf\x6b\x9c\x7e\xfb\x3e\x3e\xb5\xf4\x14\xa8\x46\xd5\xeb\x7d\x67\xd0\x37\x90\x29\x62\x97\xb2\x75\xae\x5a\xf0\x65\x77\x0e\x9a\x13\x49\xc8\x1d\xc1\x96\xa5\x32\xe4\xfd\xaa\x11\x37\xf1\x33\x29\xe4\xfc\x59\x3e\xd9\x14\x08\xc9\xd4\x11\x24\x90\x1c\xee\x16\xd2\x12\xb4\x72\x3e\x58\x6c\x12\xac\xef\x71\x4f\x73\x44\x0e\x3c\xdf\x35\xc2\x54\x50\xf5\x5d\x27\x4d\x7d\x81\x62\x99\xd6\x80\xd7\xbf\x0e\x99\x63\xbb\x53\xc6\x93\x6d\x64\x45\xe3\x14\xc6\xa7\x1a\x68\x32\x47\x0e\x8e\xf5\x1a\x6f\x71\xfc\xdd\x58\xad\x20\xeb\x3a\x58\x88\x99\x09\x63\x0d\x5e\x66\x97\x7d\x83\x61\x29\x9d\x41\x9f\x9d\x71\x14\x0a\x7f\xbf\x33\xe1\xd3\xfa\x97\x83\x38\x99\x48\x79\xae\x15\xee\xf2\xa4\x63\x8a\xef\xab\x4b\x89\x9c\x79\x3e\xc9\xe1\x34\xef\x9d\xfc\x46\xc5\x33\x83\xe5\x19\x3f\x67\xd9\xef\xb8\x0c\x17\xa7\xa2\xe2\xf6\xc2\xe0\x77\xef\xa0\x42\x46\xde\xc6\x84\x84\xc3\xb7\xea\xcd\x9f\x3b\xac\xcf\x88\xa2\xe1\x93\x47\xb1\x84\xc3\xb6\xe1\xb9\xfc\x1f\x00\x00\xff\xff\x7c\xf2\xca\x76\x72\x06\x00\x00")

// FileSelectCommonGo is "select_common.go"
var FileSelectCommonGo = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\x94\x93\x4f\x6b\xdc\x30\x10\xc5\xcf\xd6\xa7\x78\xec\xa5\x76\x6a\x76\x5b\x28\x3d\x14\x7c\x2a\xb9\x94\x84\xfe\x49\xa0\x87\xb0\x07\xd5\x1e\xef\x8a\xd8\x92\x2b\x8d\x49\x42\xc9\x77\x2f\x23\x69\x13\x27\xd9\xd2\x76\x61\xc1\x68\x34\x6f\x7e\x6f\x46\x33\xe9\xf6\x5a\xef\x08\x3c\x0d\x41\x29\x33\x4e\xce\x33\x56\x81\xbd\xb1\xbb\xb0\x52\x6a\xb3\xc1\xe5\x05\x0d\xd4\x32\x4c\x00\xef\x09\x81\xfd\xdc\x32\x78\xaf\x19\x7b\x37\x74\xe9\xf4\xe2\xf4\xec\xf4\xe3\x25\x3a\xcd\x5a\xf1\xdd\x44\x0f\x59\xf9\xfa\x2f\x55\x6c\x36\xf8\x3a\x93\x37\xe4\xf1\x54\x55\xe3\x67\x3a\xaf\xf1\x63\x16\x65\x13\x60\x2c\x93\xef\x75\x4b\xf1\xc2\xcc\x6e\x47\x96\xbc\x66\xea\x54\xe1\xfc\x08\xf9\x9d\x7c\xfe\x76\xae\x8a\xd6\x0d\xf3\x68\x03\xae\xb6\x89\x5a\x15\x37\x7b\xf2\x24\xf1\xef\xf2\xa1\x8a\x49\x0c\x02\xf8\xa2\x77\xa4\xee\xa3\xa7\x18\x81\x9e\xa6\xc1\x50\x40\x4a\x68\x9d\xed\x0c\x1b\x67\x03\x9c\x8d\xa6\x04\xeb\x4e\xf5\xb3\x6d\x51\x06\x9c\x64\xe4\x2a\x65\x97\x29\x2b\x15\xa9\x1e\xa2\xe2\x34\xac\x53\xac\x49\xca\xaa\xf0\xc4\xb3\xb7\x08\xb9\xfa\x99\x19\x0d\x3f\x54\xf7\xee\x26\x60\x88\x47\xcb\xba\xf0\x14\x26\x67\x03\x1d\x01\x88\x02\x65\xca\x31\x96\xdf\xbf\x7b\x5e\x5f\x2c\xaf\x53\xbc\x49\xda\x2f\x20\xa4\x1b\x4f\x19\x5c\xdf\x07\x62\x68\xdb\xfd\x27\x8e\x48\x95\x29\xbb\xc6\xdf\xa8\x72\x95\x26\x97\xfb\x47\xda\x10\x65\x2e\xe2\x84\x91\x62\xe9\xe1\x1d\xc6\xcf\x2e\xdf\x41\xef\xfc\xf2\x49\x06\xd6\x4c\x23\x59\x3e\x02\xbe\x54\x2d\x2b\xa4\x07\x24\xac\xa6\xc7\x40\xb6\x0c\xeb\x2c\x5f\xa1\x69\xf0\x46\x22\x07\xb0\xd5\xc9\x4a\x15\xf7\x8f\x9c\x69\x63\xd6\x9f\x9c\x59\xa4\xd5\x58\xd5\x58\x55\xd9\x83\xee\x3a\xf9\xcb\x83\x4f\x71\x81\x8e\x2b\x75\x58\x95\x3f\x93\xea\xae\x2b\x73\x52\x2a\xf5\xbc\xb7\x87\x3e\x34\x32\x53\xb2\xdd\x12\x22\x7d\x54\x2f\x9a\x9a\x2f\x9c\xeb\x29\xad\xe1\xa8\xa7\xc8\x33\xea\x29\x1c\x10\xad\x1e\x49\x38\x0d\xbf\x0a\x28\x07\x13\x64\xb6\x1d\xdd\xe2\x35\xde\x56\x22\x62\xfa\x7c\xf5\x5c\x4f\x57\xe9\x6b\x1b\xbb\x55\x2f\xe6\x83\xce\x51\x80\x75\x0c\xba\x35\x81\x65\xbf\x97\xce\x33\xc8\x11\xdf\x8f\x88\x65\x25\x5c\x57\xc9\xfd\xd6\xd8\xe8\x7b\xc4\x87\x06\xa3\xbe\xa6\xf2\x69\xac\x7e\x36\xbe\x4a\x15\xf2\x2c\x4c\xec\x85\xe4\x78\x6d\x77\x84\xc7\xb6\xc9\x64\x47\xc1\xdf\xa2\x81\x11\x73\xcb\xe9\x8e\xea\x5e\xfd\x0e\x00\x00\xff\xff\xc3\x86\xdf\x63\x2a\x05\x00\x00")

// FileUpdateGoTpl is "update.go.tpl"
var FileUpdateGoTpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\x84\x92\x5f\x6b\x9c\x40\x14\xc5\x9f\x9d\x4f\x71\x90\x0d\xe8\xb2\x9d\x7d\x2f\xf4\xa1\xb4\x29\x04\xfa\x67\xe9\x9a\xe7\x66\xa2\x57\xd7\x56\x47\x99\xb9\x93\x52\x86\xf9\xee\x65\x1c\x0d\x4d\xda\x92\xb7\xcb\xb9\xfa\x3b\xf7\x1c\x9d\x55\xfd\x43\x75\x04\xef\xe5\x29\x8d\x21\x08\xd1\x8f\xf3\x64\x18\x85\x00\x80\xbc\x1d\x39\x17\x69\xf4\x5e\x56\xbf\x66\x92\x37\xcb\x03\x27\xc5\x97\x10\x72\x51\x0a\xd1\x3a\x5d\xa3\x70\xd8\x57\xb7\x73\xa3\x98\x4a\x9c\xd9\xf4\xba\x2b\x4a\xd8\x65\x80\x17\x99\x21\x76\x46\xa3\x1d\x59\x9e\x67\xd3\x6b\x6e\x8b\xbb\xdb\xd3\xfb\xb7\xd5\x35\x36\x70\xa5\xee\x07\x0a\x01\xe7\xeb\x0a\x57\x16\x57\xf6\xee\x20\xb2\x68\xed\xa4\xb2\xb6\xef\xf4\x48\x9a\x3f\xf6\x96\x8b\xf2\x71\xf1\xf3\x42\x86\xe4\xe6\x77\x10\x59\x29\x82\x10\xc7\x23\x6e\xb4\x25\xc3\x1b\xfa\xb3\x1a\x23\xb9\x36\xa4\x98\x2c\x94\xc6\xea\x6d\x59\x31\x45\x2e\x54\x5d\x4f\xa6\x89\xd7\xf2\x04\xbe\x10\xba\xfe\x81\x34\xa6\xfb\xef\x54\xf3\x9a\x71\xc2\xfe\xcb\xd7\x4f\x25\x52\xce\xa7\xf0\x62\xc6\x7e\x53\x3e\xb8\x61\x48\x6a\xf9\xd8\x4a\x2c\xc1\xe1\xf5\x1b\x4c\x32\x09\x45\x29\x32\xef\x5f\xc1\x28\xdd\x11\x76\xdf\x0e\xd8\xb5\x71\xbf\x22\x7a\x1a\x1a\x1b\x82\xc8\x9c\x54\x4d\x53\xe4\xde\xef\x5a\xf9\x6e\x1a\xdc\xa8\x13\x3a\x3f\x60\x96\x8b\xba\x5a\x25\x1c\xe9\x26\xbe\xb5\xf6\xed\x62\x1d\xde\xbf\xe0\x71\x3c\xe2\x4c\xfc\x07\x0b\x96\xd8\xe2\x41\x0d\x8e\xd0\x4e\x06\xf5\xe2\x8b\xbf\x6e\x40\xaf\x97\xae\x9e\xb7\xf9\xaf\x7f\xe2\x89\x41\x91\xd8\x8b\x12\x6f\x79\xde\xd4\x76\xfe\xff\xb3\x2f\x80\xf8\xb5\xbd\x4f\x91\x7f\x07\x00\x00\xff\xff\xbd\xe2\xcb\x8e\xcf\x02\x00\x00")

// FileUpdateCommonGo is "update_common.go"
var FileUpdateCommonGo = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\x6c\x53\x51\x6b\xdb\x30\x10\x7e\xb6\x7e\xc5\x21\x08\x95\x82\x51\xfa\x5c\x08\x63\xb0\x3c\x74\x8c\x31\x9a\x94\x3d\x8c\x31\x54\x5b\x76\x44\x6d\xc9\x95\x4e\xd9\x46\x9a\xff\x3e\x4e\xb6\xd3\x74\xed\x83\x12\x45\x77\xdf\x77\x77\xdf\x77\x19\x74\xf5\xa8\x5b\x03\x38\x74\x91\x31\xdb\x0f\x3e\x20\x08\x56\xf0\x5a\xa3\x7e\xd0\xd1\xac\xe2\x53\xc7\x59\xc1\x9b\x1e\xe9\x2b\x62\xb0\xae\x8d\x9c\x49\xc6\x56\x2b\xd8\xdd\x0f\xb5\x46\x03\x36\x82\x86\x88\x21\x55\x08\xe8\x61\xef\xbb\x1a\xac\x6b\x7c\xe8\x35\x5a\xef\xa0\xf1\x01\xb4\x83\xdb\xaf\xdb\xcd\xdd\x0e\x22\x6a\x34\xbd\x71\xc8\xf0\xef\x60\xce\x24\x13\xfe\xc8\x8a\xdd\xad\x8b\x26\x20\x2b\x7e\xef\x4d\x30\xf0\x9d\x3e\xd9\x89\xb1\x26\xb9\x0a\x44\x82\xe5\x04\x91\x63\x48\x8c\x69\xcb\xfc\x43\x9e\xa3\xc4\x94\xd4\x18\x5b\xc3\x32\x5f\x58\x11\x0c\xa6\xe0\x20\x11\xdf\x6a\x05\x9b\x3f\xa6\x02\x9b\xcb\x45\xc0\xbd\x01\x1a\x9c\x66\xa0\x7b\x6b\x0f\xc6\xc1\x2c\xc5\x3b\xe5\x09\x2d\x24\x88\xf8\xd4\xa9\x3b\x13\x53\x87\x25\x98\x10\x7c\x90\x54\xdc\x36\xd0\x19\x27\x92\xaa\x7c\x17\x25\xac\xd7\x70\x0d\xcf\xcf\xd3\xdb\x41\x77\xc9\xcc\xaf\x47\x56\xcc\x8d\x39\xdb\x95\xd0\xf4\xa8\x36\xc4\xd3\x08\xee\x3c\xee\xad\x6b\xa9\xa7\x94\xcb\x72\xc9\x8a\x13\x63\x45\xc4\x1e\xe1\x66\x0d\x49\x6d\xb3\x2b\x42\xd2\xbc\x3e\xf4\xaa\xf3\xad\xe0\x63\x8f\x37\x70\xb5\x38\x5c\xc1\xe2\xc0\x4b\x20\x40\x09\xe7\xd2\x2f\x5a\x64\x50\xfd\xa0\xf2\x38\x63\x96\x1e\x06\xe3\xea\x73\x9f\x04\xcb\x02\x2a\x1d\xda\xa8\x94\x92\x74\x26\x0d\x75\x5d\xd3\xa1\x1d\xa8\x7c\x97\x7a\x07\xda\xd5\x90\x81\xb3\x92\xf7\xdf\x3e\x7d\xdc\x6d\x2e\x9c\x7f\xab\xa5\xae\x6b\xe1\x74\x9f\xf7\xc0\xba\xb6\x9c\x08\xac\x43\x13\x1a\x5d\x99\xe3\xe9\x7f\x6b\xa7\x35\x51\x33\x72\x82\xc8\xd7\x26\xbf\x53\x29\x46\xdb\x3a\x6a\xe3\x8b\x8d\x28\xe4\x54\x91\x48\x5f\x42\x91\xa4\xed\xf5\xa3\x11\x3f\x7e\xce\x1d\x5d\x97\x97\x86\x4a\x56\xd0\x5e\xff\x2a\x69\x6a\xca\x0e\xda\xb5\x06\xc6\x68\xb6\xf4\x92\x6d\x3d\x4b\x7a\xf1\x38\x1a\xbd\x1d\x82\x75\xd8\x08\xbe\xa0\xac\x0f\x3c\xf3\xc9\x6c\xf2\x3c\xc8\xf4\xb7\x53\x9f\xbd\x75\xaf\x09\x78\x09\x9c\x7c\xf8\x17\x00\x00\xff\xff\x95\x19\x92\x95\xc9\x03\x00\x00")

// FileWhereGoTpl is "where.go.tpl"
var FileWhereGoTpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\xac\x91\x41\x6b\xab\x40\x14\x85\xd7\xce\xaf\x38\x84\x2c\xf4\x21\x93\xfd\x83\x6c\xf2\xf0\x51\x37\x69\x29\x81\x2c\x42\x28\x83\x5e\x75\x88\xde\x11\x9d\x24\x94\x61\xfe\x7b\x51\x93\x52\x48\x02\x5d\x64\x27\xd7\x73\xbe\x39\xf7\xdc\x56\x65\x07\x55\x12\x9c\x93\x6f\xd3\xa7\xf7\x42\x38\xd7\x29\x2e\x09\xf3\x8f\x18\xf3\x02\x7f\x97\x90\x9b\xcf\x96\xe4\x7f\x4d\x75\xde\x7b\x2f\x16\x0b\x6c\x2b\xea\xc8\xb9\x79\x21\xd7\xaa\x21\xef\xa1\xf2\xbc\x87\x42\x66\x38\xd7\x56\x1b\x86\x61\xfc\xfc\x6f\x0d\x6c\x45\xd8\xbe\x24\xef\x09\x7a\xab\x2c\x35\xc4\x56\x14\x47\xce\x6e\x60\xa1\x69\xf1\xda\xc6\x38\xa9\x7a\x62\x0c\xcf\x7b\x1f\xe1\xcf\xa8\x84\x13\x41\x47\xf6\xd8\x31\x98\xce\xe3\x28\x34\x6d\x8c\xd9\xa8\xfd\x67\xea\x63\xc3\x13\x68\x36\x32\x22\xe1\xc5\xbd\xcc\x29\x5f\x52\x33\xd2\xf5\x73\x92\xa7\x1c\x9e\x54\xdd\x43\x4a\xf9\x20\xb8\xea\xca\x7e\xa8\xb4\x51\x07\x0a\x77\x7b\xcd\x96\xba\x42\x65\xe4\x7c\x8c\x9a\x26\x7b\x14\x89\xa0\x30\x1d\xf4\x20\x9c\x6e\x31\x52\x9d\x08\x46\xff\x4e\xef\xb1\x1c\x47\x3b\xbd\x17\x81\xbf\xa9\x23\xe5\xf0\x6e\x19\x83\x59\x4a\xf9\xa8\x90\x15\xd9\x33\xd1\xb5\x15\xac\x92\xcd\x36\x49\x9e\xd4\xcc\x85\x1d\xd6\xe6\x1c\xa3\xd2\x65\xf5\xcb\xcb\x5e\x7d\x77\xf7\xf9\x86\x0d\x1b\x39\x47\x9c\x7b\x2f\xbe\x02\x00\x00\xff\xff\xd4\x7d\xfc\x53\xd5\x02\x00\x00")

// FileWhereCommonGo is "where_common.go"
var FileWhereCommonGo = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\xac\x94\xdd\x4e\xdb\x30\x14\xc7\xaf\xe3\xa7\x38\xb2\xc4\xb0\xb7\x28\xec\x05\x10\x2a\xa3\xd2\x98\x58\xab\x51\x24\x2e\x10\x17\xa6\x75\x52\x0b\xd7\x0e\xb6\x83\x87\x4a\xdf\x7d\x72\xec\xa6\x69\xcb\xc7\x36\x71\xd5\xc6\x3e\x1f\xff\xdf\xf1\x39\xa7\x66\xd3\x7b\x56\x71\x70\xb5\xb4\x08\x89\x45\xad\x8d\x03\x82\x32\x5c\x2e\x1c\x46\x19\xb6\xce\x08\x55\x59\x8c\x28\x42\x47\x47\x70\x3d\xe7\x86\x03\x33\x1c\x74\xed\x84\x56\x16\x4a\x6d\x60\xf2\xeb\x02\xae\xbf\x0f\x2f\x87\x60\x1d\x73\x7c\xc1\x95\x43\xee\xa9\xe6\xc9\xdc\x3a\xd3\x4c\x1d\x2c\x51\x66\xdd\xc2\xc1\xcd\x6d\x0c\x8a\x32\x66\x2a\x0b\x37\xb7\x42\x39\x6e\x4a\x36\xe5\xcb\x15\x5a\xb5\x69\x14\xf7\xd1\xd5\x70\xd7\x18\x65\x81\x85\xa3\xbd\x1c\x65\xa3\xa6\x9d\x2d\xd1\x35\x8c\xeb\x1c\x1e\x99\x11\xec\x4e\xb6\x69\x85\xaa\xc2\x81\x6c\x38\xf4\x92\x50\xf8\x1c\xa3\x2f\x51\xf6\xc8\x0c\xf8\xa8\x13\x65\xbe\x68\x05\x1e\x03\xab\x6b\xae\x66\x24\x7e\xe7\x50\x2e\x5c\x31\xa9\x8d\x50\xae\x24\xf8\xc0\xc2\x81\x85\x13\xbc\xc9\x94\x83\xae\x29\x0d\xee\x2d\x50\xcf\x3d\x7c\x27\x01\x14\x65\x11\x06\x3e\xf9\x1d\xca\x73\xb5\xc3\x79\x18\x41\x3b\x92\xf3\x11\x90\xa2\x28\xe8\xe1\x6b\xec\xe7\x8a\xbc\x8c\x6d\xa1\x28\x8a\x8f\x44\x0f\x52\x0e\x2c\xdd\xa2\x7f\xf8\xc9\xcc\xbd\x25\x92\x07\x15\x21\x29\xa5\xef\x96\xc3\x06\x9e\xd7\x4b\x72\xca\x9d\xe7\xfc\xbd\xba\x9c\x0e\xaf\xae\x87\xc3\x11\x48\xed\x61\x30\x3a\x83\xb9\xa8\xe6\xaf\xd6\x28\x85\xdc\x2f\x94\xd4\x3e\x6f\x5d\x3f\xb4\x47\xd6\xda\x4e\x5a\x65\xfd\x76\x79\xab\x36\x9d\x96\x17\x6a\x33\x69\xf5\x76\x25\x71\x73\x9e\x06\x22\x8c\xdf\x0e\x34\xf1\x49\x3f\x4d\x6e\x84\x26\xde\x00\x24\x4a\xf0\x70\x7c\x0c\x4a\x48\x78\x7e\x86\xf0\x70\x11\x83\x86\xd3\xaf\xc1\x64\x9d\x1d\x63\x94\xad\x3a\x2d\x38\x26\xc4\xf0\x25\x45\xb3\xc5\x0f\x2d\x54\x57\x04\x0c\x98\x26\xb5\x83\xc0\x17\x16\x45\xd0\xc9\x4c\xd5\x04\x71\x71\x5d\xf0\xdf\x7c\xda\xb8\xa0\x85\xc1\x64\x78\x31\xfc\x76\x05\x0f\x0d\x37\x4f\xe0\x85\x9b\x03\x4b\x54\x53\xad\x66\x22\xec\x98\x3d\xa2\x10\x9a\xd0\xed\xc5\xb1\x8b\xd5\x43\x50\x42\xf6\x19\x62\xa9\x93\xca\xb1\x09\x2f\x20\x05\xb7\xc0\x14\x68\xb3\xc9\x0a\x77\xa9\x05\x9d\xd7\xe0\xdb\x56\xe8\xee\xec\x9e\xa4\xb1\x21\x46\x54\x73\xd7\x1d\x6c\xba\x27\xa5\xbd\x13\x8a\x99\x27\xe2\x73\x68\x0d\x73\xc0\xe3\xcb\x4d\xb1\xd4\xac\xaf\x83\xa9\xd9\xff\x0a\x19\xa8\xd9\xbf\x2b\x19\x8c\xce\xa2\x94\x36\x5a\xba\x97\xc9\x33\x07\xd3\xfd\xd3\x75\x7a\xf7\x7e\x58\xb9\x33\x11\xeb\xfd\xbe\xc4\x04\xaf\x72\x88\xd7\x71\xdc\x77\x4d\xe5\xba\x6f\xc2\x42\xd1\x75\x0e\x98\xe0\x37\xcc\xcc\xdf\x85\x6a\x6f\xb7\xe7\x4b\xa6\xf9\x32\xed\xef\xd6\xee\x91\x68\x85\xfe\x04\x00\x00\xff\xff\xa7\xe4\x4f\x65\x06\x07\x00\x00")



func init() {
  if CTX.Err() != nil {
		log.Fatal(CTX.Err())
	}



var err error



  




  
  var f webdav.File
  

  
  
  var rb *bytes.Reader
  var r *gzip.Reader
  
  

  
  
  
  rb = bytes.NewReader(FileAPIGoTpl)
  r, err = gzip.NewReader(rb)
  if err != nil {
    log.Fatal(err)
  }

  err = r.Close()
  if err != nil {
    log.Fatal(err)
  }
  
  

  f, err = FS.OpenFile(CTX, "api.go.tpl", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
  if err != nil {
    log.Fatal(err)
  }

  
  
  _, err = io.Copy(f, r)
  if err != nil {
    log.Fatal(err)
  }
  
  

  err = f.Close()
  if err != nil {
    log.Fatal(err)
  }
  
  
  
  rb = bytes.NewReader(FileAPICommonGo)
  r, err = gzip.NewReader(rb)
  if err != nil {
    log.Fatal(err)
  }

  err = r.Close()
  if err != nil {
    log.Fatal(err)
  }
  
  

  f, err = FS.OpenFile(CTX, "api_common.go", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
  if err != nil {
    log.Fatal(err)
  }

  
  
  _, err = io.Copy(f, r)
  if err != nil {
    log.Fatal(err)
  }
  
  

  err = f.Close()
  if err != nil {
    log.Fatal(err)
  }
  
  
  
  rb = bytes.NewReader(FileCommonGo)
  r, err = gzip.NewReader(rb)
  if err != nil {
    log.Fatal(err)
  }

  err = r.Close()
  if err != nil {
    log.Fatal(err)
  }
  
  

  f, err = FS.OpenFile(CTX, "common.go", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
  if err != nil {
    log.Fatal(err)
  }

  
  
  _, err = io.Copy(f, r)
  if err != nil {
    log.Fatal(err)
  }
  
  

  err = f.Close()
  if err != nil {
    log.Fatal(err)
  }
  
  
  
  rb = bytes.NewReader(FileCreateGoTpl)
  r, err = gzip.NewReader(rb)
  if err != nil {
    log.Fatal(err)
  }

  err = r.Close()
  if err != nil {
    log.Fatal(err)
  }
  
  

  f, err = FS.OpenFile(CTX, "create.go.tpl", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
  if err != nil {
    log.Fatal(err)
  }

  
  
  _, err = io.Copy(f, r)
  if err != nil {
    log.Fatal(err)
  }
  
  

  err = f.Close()
  if err != nil {
    log.Fatal(err)
  }
  
  
  
  rb = bytes.NewReader(FileCreateCommonGo)
  r, err = gzip.NewReader(rb)
  if err != nil {
    log.Fatal(err)
  }

  err = r.Close()
  if err != nil {
    log.Fatal(err)
  }
  
  

  f, err = FS.OpenFile(CTX, "create_common.go", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
  if err != nil {
    log.Fatal(err)
  }

  
  
  _, err = io.Copy(f, r)
  if err != nil {
    log.Fatal(err)
  }
  
  

  err = f.Close()
  if err != nil {
    log.Fatal(err)
  }
  
  
  
  rb = bytes.NewReader(FileDeleteGoTpl)
  r, err = gzip.NewReader(rb)
  if err != nil {
    log.Fatal(err)
  }

  err = r.Close()
  if err != nil {
    log.Fatal(err)
  }
  
  

  f, err = FS.OpenFile(CTX, "delete.go.tpl", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
  if err != nil {
    log.Fatal(err)
  }

  
  
  _, err = io.Copy(f, r)
  if err != nil {
    log.Fatal(err)
  }
  
  

  err = f.Close()
  if err != nil {
    log.Fatal(err)
  }
  
  
  
  rb = bytes.NewReader(FileDeleteCommonGo)
  r, err = gzip.NewReader(rb)
  if err != nil {
    log.Fatal(err)
  }

  err = r.Close()
  if err != nil {
    log.Fatal(err)
  }
  
  

  f, err = FS.OpenFile(CTX, "delete_common.go", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
  if err != nil {
    log.Fatal(err)
  }

  
  
  _, err = io.Copy(f, r)
  if err != nil {
    log.Fatal(err)
  }
  
  

  err = f.Close()
  if err != nil {
    log.Fatal(err)
  }
  
  
  
  rb = bytes.NewReader(FileInsertGoTpl)
  r, err = gzip.NewReader(rb)
  if err != nil {
    log.Fatal(err)
  }

  err = r.Close()
  if err != nil {
    log.Fatal(err)
  }
  
  

  f, err = FS.OpenFile(CTX, "insert.go.tpl", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
  if err != nil {
    log.Fatal(err)
  }

  
  
  _, err = io.Copy(f, r)
  if err != nil {
    log.Fatal(err)
  }
  
  

  err = f.Close()
  if err != nil {
    log.Fatal(err)
  }
  
  
  
  rb = bytes.NewReader(FileInsertCommonGo)
  r, err = gzip.NewReader(rb)
  if err != nil {
    log.Fatal(err)
  }

  err = r.Close()
  if err != nil {
    log.Fatal(err)
  }
  
  

  f, err = FS.OpenFile(CTX, "insert_common.go", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
  if err != nil {
    log.Fatal(err)
  }

  
  
  _, err = io.Copy(f, r)
  if err != nil {
    log.Fatal(err)
  }
  
  

  err = f.Close()
  if err != nil {
    log.Fatal(err)
  }
  
  
  
  rb = bytes.NewReader(FileLogGo)
  r, err = gzip.NewReader(rb)
  if err != nil {
    log.Fatal(err)
  }

  err = r.Close()
  if err != nil {
    log.Fatal(err)
  }
  
  

  f, err = FS.OpenFile(CTX, "log.go", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
  if err != nil {
    log.Fatal(err)
  }

  
  
  _, err = io.Copy(f, r)
  if err != nil {
    log.Fatal(err)
  }
  
  

  err = f.Close()
  if err != nil {
    log.Fatal(err)
  }
  
  
  
  rb = bytes.NewReader(FileOpGo)
  r, err = gzip.NewReader(rb)
  if err != nil {
    log.Fatal(err)
  }

  err = r.Close()
  if err != nil {
    log.Fatal(err)
  }
  
  

  f, err = FS.OpenFile(CTX, "op.go", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
  if err != nil {
    log.Fatal(err)
  }

  
  
  _, err = io.Copy(f, r)
  if err != nil {
    log.Fatal(err)
  }
  
  

  err = f.Close()
  if err != nil {
    log.Fatal(err)
  }
  
  
  
  rb = bytes.NewReader(FileOrmGo)
  r, err = gzip.NewReader(rb)
  if err != nil {
    log.Fatal(err)
  }

  err = r.Close()
  if err != nil {
    log.Fatal(err)
  }
  
  

  f, err = FS.OpenFile(CTX, "orm.go", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
  if err != nil {
    log.Fatal(err)
  }

  
  
  _, err = io.Copy(f, r)
  if err != nil {
    log.Fatal(err)
  }
  
  

  err = f.Close()
  if err != nil {
    log.Fatal(err)
  }
  
  
  
  rb = bytes.NewReader(FilePageGo)
  r, err = gzip.NewReader(rb)
  if err != nil {
    log.Fatal(err)
  }

  err = r.Close()
  if err != nil {
    log.Fatal(err)
  }
  
  

  f, err = FS.OpenFile(CTX, "page.go", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
  if err != nil {
    log.Fatal(err)
  }

  
  
  _, err = io.Copy(f, r)
  if err != nil {
    log.Fatal(err)
  }
  
  

  err = f.Close()
  if err != nil {
    log.Fatal(err)
  }
  
  
  
  rb = bytes.NewReader(FileSelectGoTpl)
  r, err = gzip.NewReader(rb)
  if err != nil {
    log.Fatal(err)
  }

  err = r.Close()
  if err != nil {
    log.Fatal(err)
  }
  
  

  f, err = FS.OpenFile(CTX, "select.go.tpl", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
  if err != nil {
    log.Fatal(err)
  }

  
  
  _, err = io.Copy(f, r)
  if err != nil {
    log.Fatal(err)
  }
  
  

  err = f.Close()
  if err != nil {
    log.Fatal(err)
  }
  
  
  
  rb = bytes.NewReader(FileSelectCommonGo)
  r, err = gzip.NewReader(rb)
  if err != nil {
    log.Fatal(err)
  }

  err = r.Close()
  if err != nil {
    log.Fatal(err)
  }
  
  

  f, err = FS.OpenFile(CTX, "select_common.go", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
  if err != nil {
    log.Fatal(err)
  }

  
  
  _, err = io.Copy(f, r)
  if err != nil {
    log.Fatal(err)
  }
  
  

  err = f.Close()
  if err != nil {
    log.Fatal(err)
  }
  
  
  
  rb = bytes.NewReader(FileUpdateGoTpl)
  r, err = gzip.NewReader(rb)
  if err != nil {
    log.Fatal(err)
  }

  err = r.Close()
  if err != nil {
    log.Fatal(err)
  }
  
  

  f, err = FS.OpenFile(CTX, "update.go.tpl", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
  if err != nil {
    log.Fatal(err)
  }

  
  
  _, err = io.Copy(f, r)
  if err != nil {
    log.Fatal(err)
  }
  
  

  err = f.Close()
  if err != nil {
    log.Fatal(err)
  }
  
  
  
  rb = bytes.NewReader(FileUpdateCommonGo)
  r, err = gzip.NewReader(rb)
  if err != nil {
    log.Fatal(err)
  }

  err = r.Close()
  if err != nil {
    log.Fatal(err)
  }
  
  

  f, err = FS.OpenFile(CTX, "update_common.go", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
  if err != nil {
    log.Fatal(err)
  }

  
  
  _, err = io.Copy(f, r)
  if err != nil {
    log.Fatal(err)
  }
  
  

  err = f.Close()
  if err != nil {
    log.Fatal(err)
  }
  
  
  
  rb = bytes.NewReader(FileWhereGoTpl)
  r, err = gzip.NewReader(rb)
  if err != nil {
    log.Fatal(err)
  }

  err = r.Close()
  if err != nil {
    log.Fatal(err)
  }
  
  

  f, err = FS.OpenFile(CTX, "where.go.tpl", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
  if err != nil {
    log.Fatal(err)
  }

  
  
  _, err = io.Copy(f, r)
  if err != nil {
    log.Fatal(err)
  }
  
  

  err = f.Close()
  if err != nil {
    log.Fatal(err)
  }
  
  
  
  rb = bytes.NewReader(FileWhereCommonGo)
  r, err = gzip.NewReader(rb)
  if err != nil {
    log.Fatal(err)
  }

  err = r.Close()
  if err != nil {
    log.Fatal(err)
  }
  
  

  f, err = FS.OpenFile(CTX, "where_common.go", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
  if err != nil {
    log.Fatal(err)
  }

  
  
  _, err = io.Copy(f, r)
  if err != nil {
    log.Fatal(err)
  }
  
  

  err = f.Close()
  if err != nil {
    log.Fatal(err)
  }
  


  Handler = &webdav.Handler{
    FileSystem: FS,
    LockSystem: webdav.NewMemLS(),
  }


}

// Open a file
func (hfs *HTTPFS) Open(path string) (http.File, error) {
  f, err := FS.OpenFile(CTX, path, os.O_RDONLY, 0644)
  if err != nil {
    return nil, err
  }

  return f, nil
}

// ReadFile is adapTed from ioutil
func ReadFile(path string) ([]byte, error) {
  f, err := FS.OpenFile(CTX, path, os.O_RDONLY, 0644)
  if err != nil {
    return nil, err
  }

  buf := bytes.NewBuffer(make([]byte, 0, bytes.MinRead))

  // If the buffer overflows, we will get bytes.ErrTooLarge.
  // Return that as an error. Any other panic remains.
  defer func() {
    e := recover()
    if e == nil {
      return
    }
    if panicErr, ok := e.(error); ok && panicErr == bytes.ErrTooLarge {
      err = panicErr
    } else {
      panic(e)
    }
  }()
  _, err = buf.ReadFrom(f)
  return buf.Bytes(), err
}

// WriteFile is adapTed from ioutil
func WriteFile(filename string, data []byte, perm os.FileMode) error {
  f, err := FS.OpenFile(CTX, filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
  if err != nil {
    return err
  }
  n, err := f.Write(data)
  if err == nil && n < len(data) {
    err = io.ErrShortWrite
  }
  if err1 := f.Close(); err == nil {
    err = err1
  }
  return err
}

// WalkDirs looks for files in the given dir and returns a list of files in it
// usage for all files in the b0x: WalkDirs("", false)
func WalkDirs(name string, includeDirsInList bool, files ...string) ([]string, error) {
	f, err := FS.OpenFile(CTX, name, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}

	fileInfos, err := f.Readdir(0)
	f.Close()
	if err != nil {
		return nil, err
	}

	for _, info := range fileInfos {
		filename := path.Join(name, info.Name())

		if includeDirsInList || !info.IsDir() {
			files = append(files, filename)
		}

		if info.IsDir() {
			files, err = WalkDirs(filename, includeDirsInList, files...)
			if err != nil {
				return nil, err
			}
		}
	}

	return files, nil
}


