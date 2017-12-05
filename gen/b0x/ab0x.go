
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



// FileExecGoTpl is "exec.go.tpl"
var FileExecGoTpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\xdc\x58\x51\x6f\xdb\x38\x12\x7e\x96\x7e\xc5\xd4\x48\x5b\x29\x70\xa5\xb6\xe8\xbd\xa4\xf0\x43\x9b\x4b\x71\xc1\xdd\x25\xd9\xa6\xdd\x97\xa0\x28\x68\x6a\x24\x13\xa1\x48\x97\xa4\x1c\x1b\x86\xff\xfb\x62\x48\xc9\x56\x1c\xbb\xcd\xee\x66\xbb\xdd\xbe\xd9\x9c\xe1\xf0\x9b\x8f\xdf\xcc\x48\x9a\x32\x7e\xcd\x2a\x84\xe5\x32\xbb\x08\x3f\x57\xab\x38\x16\xf5\x54\x1b\x07\x49\x1c\x0d\xb8\x56\x0e\xe7\x6e\x10\x47\x83\x82\x39\x36\x66\x16\x73\xfb\x45\x6e\xff\xcf\x0b\x23\x66\x68\x68\xb9\xac\xbd\xb7\xc1\x52\x22\xf7\x3f\x1b\x65\x59\x89\xf4\xab\x12\x6e\xd2\x8c\x33\xae\xeb\x7c\xaa\x2d\x2a\x34\xb9\x36\xf5\x20\x8e\x01\x00\x06\xcb\x65\xf6\x61\x31\xc5\xec\xd4\x1f\x7f\xc1\xdc\x64\xb5\x1a\xc4\x69\x1c\xe7\x39\x9c\xcc\x91\x03\x37\xc8\x1c\x5a\x60\xe0\xd8\x58\x22\x94\xda\x80\x9b\x20\x54\x62\x86\x0a\xac\x33\x0d\x77\x71\xd9\x28\x0e\xc9\x18\x0e\x8f\xbd\xf7\xdb\x46\xc8\x02\x4d\xea\x23\x24\x29\x24\xf6\x8b\xcc\xde\xa3\x6d\xa4\x1b\x02\x1a\xa3\x4d\x0a\xcb\x38\xb2\xae\x76\x43\x60\xa6\xb2\x70\x34\x02\x82\xa8\x54\x56\x08\x46\x39\x64\x21\x52\xf2\x64\x9c\x4d\x99\x61\xb5\x4d\xe3\xa8\xf5\x90\xba\x4a\x06\xc1\x7c\x04\x4f\x1f\xcf\x9e\xc2\xe3\xd9\x60\x08\x9b\x68\x69\x1c\x19\x74\x8d\x51\xeb\x98\xe3\x8c\x90\x1c\x07\x5e\x93\x96\xdf\x73\xf3\x96\xf1\xeb\xca\xe8\x46\x15\x49\x77\x4c\x76\xec\xe6\x69\x3f\x58\x96\x65\x69\xbc\xf2\x7c\x7c\x69\xd0\x2c\x40\x58\x68\x2c\x16\x30\x5e\x78\x1e\x2e\xd1\xc3\xfd\xc5\xdb\x98\x2a\xba\x85\xff\x89\x5a\x38\x20\x62\x9c\xd0\xca\x6e\x28\x0a\xf6\x35\x45\x3e\x68\xc2\xdd\x1c\x5a\x58\x59\x0b\x33\x85\xe4\xd0\xf3\xa6\x6f\xec\xbd\x59\x0b\xc1\xf7\xb2\xe6\x51\xde\x9f\x34\xef\xbe\x66\xcd\xcd\xf7\xf0\xe2\x75\x22\x94\x45\xe3\xac\xe7\x84\x54\x0a\x4e\xf7\x74\xd2\xe9\x76\x43\xc3\xa9\xf7\xbf\xa7\x52\x44\x09\x12\xd5\xe6\x92\xde\x58\x2b\x2a\x55\xa3\x72\x36\x85\xd1\x08\x9e\x93\x53\x87\x5f\x09\x39\x84\xb2\x76\xd9\x09\xed\x2f\x93\x81\xd2\x6e\x22\x54\x45\x88\x02\xca\x41\x1a\x47\xab\x6f\x30\x19\xf0\xed\x65\x32\x98\xbf\xaf\xfe\xfe\x10\xcf\x1f\xa7\xc5\xfd\x2b\xf2\xc1\x78\x6e\xfc\xa9\xf7\xe1\x39\xe0\xdb\xcb\x73\x30\xff\x0d\x3c\x9b\x46\xb5\x24\xa3\x44\x87\x60\x1d\x73\x48\x4c\x80\x56\xc0\xb6\xe8\xce\x36\x7c\xff\xdb\xbb\x3f\x44\x07\x0c\x91\xf6\x32\x13\xcc\xdf\x97\x99\xd0\xe5\x3a\xe9\xdd\x16\xda\x56\x5f\xf3\xae\x94\xf9\xd5\xa7\x6e\xc6\x9c\xcc\xdd\x19\xab\x71\xb5\xea\x51\x40\x43\x88\x9a\xdf\xd1\x08\xbe\x09\xca\x3b\x9b\xae\x1d\x06\xc6\xd6\xed\x33\xf5\xea\xa5\xf5\x47\x23\x92\xe6\x1d\xa5\xa2\x31\x5e\x8d\x05\x96\x68\x7c\x98\xec\x58\x6a\x8b\x49\x1a\xc7\x51\x9e\x03\xce\x9d\x61\xdc\x79\x0b\x49\x38\x8c\xb6\xc6\xa0\x8d\xa3\x19\x33\xc0\xa4\x84\x1d\xb9\xc4\x11\x4d\x44\x1f\xee\x8c\xa8\xf5\xf7\x4a\x40\xf3\x1c\xf8\x04\xf9\x75\x97\x17\x70\xa6\x38\x4a\xc9\x68\x1e\x04\x97\x16\x30\xe5\xee\xe6\x54\x46\x49\xfa\xba\x9f\x42\x17\xca\xe7\xbd\x9d\x0a\x2d\xae\xe2\x28\x12\x0e\xeb\x1e\x21\xd6\x5f\x83\x36\xd9\x3b\x61\xac\x4b\xb6\x34\x45\x98\x93\x74\x08\x33\x26\x1b\xb4\xc9\x21\xe1\x4e\x03\xb1\x3d\x40\x3d\x06\xef\x52\x48\x1c\x46\x44\xc6\x08\xd8\x74\x8a\xaa\x48\x98\x94\x43\x38\x24\x1c\xbe\xde\x29\xd0\x72\x49\xb1\x68\x1c\x06\xba\xfe\xc3\xec\xb9\xc2\x0f\xfa\xff\x4c\x2d\xde\x63\xe0\xa0\x35\x5d\x18\x51\x33\xb3\xf8\x2f\x2e\xe0\xd9\x2a\xec\x0e\xd1\xc7\x99\xc1\xa2\xe1\x48\xf1\xd3\x2e\x2a\xaa\xc2\xbb\x75\xb8\xfc\xd9\x9e\x7d\x4f\x5f\xab\xd3\x63\xdd\x28\x07\xac\x28\x80\x01\xf7\xbf\xb9\x96\x4d\xad\xba\x66\xe9\x55\xb3\x57\xba\x7e\xf7\x6d\xe9\x86\xbb\xf6\x86\x3f\x23\xde\xde\xf5\x04\x58\x23\x70\xa6\xc1\x1f\x49\xd8\xbd\x4c\x7f\x6c\x69\x87\x5b\xfa\x09\xf4\x1d\x12\xf9\xfd\x22\xf7\x2c\xb4\xe4\x85\x51\x55\x86\x15\x7d\x03\x6e\xc2\x1c\xd4\xcc\xf1\x09\xda\x8d\xe2\x33\xda\x76\x5a\x82\xd2\xde\xe9\x8e\x7d\x08\x4c\xc1\x89\x31\x67\xda\xbd\x23\x05\xc3\x8d\x90\x12\xc6\xd8\x1e\x82\x85\x0f\xf0\x61\x22\x2c\x70\xca\x21\x5c\xbb\x05\xa6\x16\x30\x65\x95\x9f\xfd\x74\xf0\x0d\xb3\x60\xd1\xc1\x8d\x70\x13\x8a\x4e\xbb\x6e\xd5\x18\x4c\x0d\xce\x84\x6e\xac\x5c\x64\x7b\xcb\x30\x34\x30\x7a\x0c\x7e\xd0\x01\xb2\x5e\xb9\x60\x15\xb6\x8f\xe9\x23\x78\xb1\xc3\x76\x5e\x96\x94\xc5\x08\x9e\xff\x85\x05\x5a\x7a\xa2\x8f\x46\xfd\x3a\x8b\x5b\xad\x3e\x0a\xc6\x65\xbc\xab\x52\xb4\xa9\xb3\xde\x5d\xc5\xa1\x68\xfc\xce\x07\x9a\x08\xdb\x59\xed\xab\xd7\xb6\x5c\xdb\xe5\x70\xfa\x6d\xb5\xc6\xa1\x5e\x76\x16\xc4\xbe\xeb\x6f\x9b\x3f\x85\xb3\xbb\x06\x6f\xba\x6b\x91\xb8\xa7\x9e\x96\xc4\x51\x84\x73\x61\x9d\x85\x11\xd4\xec\x1a\x93\x9a\x4d\xaf\x3a\xf7\x0d\x82\xad\xed\x9f\xee\x4a\x2d\x0d\x97\x49\x39\xee\x9c\xfe\x69\xe8\x92\x9f\x87\x20\xfc\x2d\x32\x55\x21\x04\xd0\xa4\x03\xe2\x90\x70\x90\x2d\x00\xba\x12\xd9\x0e\x1c\xed\xf9\xaf\x5b\xef\x9e\x92\xda\x9e\x10\x02\x1f\x7c\x1e\xc2\x41\x49\xc1\x42\x84\xf7\x24\x29\x54\x1c\x6d\x68\x18\xd1\xa6\x39\x1d\x94\xc1\xe5\x52\x0a\x8e\xad\x35\x70\x92\x2d\x97\x07\x65\x7b\xe2\xa6\xcb\xdd\x31\x0d\x41\xf4\xff\xfa\xc7\xc0\x28\x8a\xfa\xfd\x69\xeb\xdf\x0a\x50\x5a\x5c\xf7\xd5\x4d\x6c\x83\x6e\x08\x22\x5d\x23\xf8\x3a\x09\x30\x82\x27\x06\xdd\x15\xbd\x8a\x18\x74\xe9\xb3\x17\x9f\x42\x5b\xde\x88\xcc\xa0\x23\x5d\x7d\x5d\x3b\xa1\xb1\x6e\x0b\xa8\x37\xe0\xd2\xdd\xcb\x0f\x26\xa2\xfe\x51\xbb\x64\x74\x6b\xd6\xfe\x43\x85\xd4\xcf\xe4\xa7\x11\x56\x7b\xf8\xaa\x93\xd8\xae\xc9\xb2\xf3\x93\xcd\xd6\x42\xfb\x4e\x4d\xae\xa3\xf5\x45\xf4\x5a\x68\xe7\xde\x0b\x9b\xf6\xd1\x70\x37\x6f\x87\x7d\x68\xcf\x20\x2c\x30\x98\x30\x7e\xdd\x3d\xc9\x76\x5f\x88\xba\x2f\x71\x79\x0e\x56\x28\x8e\xde\x68\x36\x06\x28\x34\x5a\x50\xda\x01\xce\xa7\xda\x92\xb2\x9e\x5a\x90\xcc\x3a\xae\xa5\x6d\xc3\xd3\x50\x01\x06\x37\x6c\x41\xe1\xe9\x15\x97\xe2\x31\xe0\x8d\x75\xba\x06\xcb\x99\x52\x68\xba\xa3\x2f\x39\x53\x50\xa3\x9b\xe8\xf0\x5c\x70\x89\x08\xc2\xda\x06\x61\xe2\xdc\xd4\x1e\xe5\x79\xef\xd3\x63\xa5\x25\x53\x55\x5e\xe9\xdc\xbb\xd8\xfc\xe5\xcb\x7f\xbd\x7a\x15\xd8\x6d\x47\x8f\x59\x27\x43\x75\x19\x3e\x6e\x66\xbf\x92\x8d\x58\xa3\xc4\x74\x8d\xd0\x54\x72\xb1\x66\x80\x71\x8e\xb6\x97\x46\x29\x50\x16\x71\x64\xfc\xeb\x74\xfb\x35\x34\x84\x38\x2f\x93\x27\x26\xcd\x4e\x24\xd6\x44\xb1\xf1\xd2\x37\x36\x7b\x47\x3b\xde\x2e\xfc\x14\x1c\x74\x71\x06\xed\x93\xb3\x9e\xa1\xe1\x74\xe8\x94\x29\xc1\x6f\x07\xcc\x4e\x95\x43\x53\x32\x8e\x47\xf4\x1c\x44\xcc\xb6\x97\xe6\xd3\x01\x3d\x76\x4c\x28\x2c\xa0\x34\xba\x86\x46\x11\xed\xc6\xd1\x7f\x3a\x91\x88\x0e\xd4\x79\x2c\x1b\xb0\x67\x78\xf3\xc6\x25\x26\x94\x1c\xcd\xe5\xf0\x21\x37\xbb\xd0\x82\xce\x23\xcb\x47\xbf\xf2\xa6\x28\x4c\x92\xa6\xbd\x94\x5a\x01\x97\x1b\x64\x49\x9a\x25\xb7\x99\xa4\x69\xfc\x5b\x00\x00\x00\xff\xff\xe2\x94\x10\x0c\x80\x16\x00\x00")

// FileOrmGoTpl is "orm.go.tpl"
var FileOrmGoTpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\xd4\x57\x4f\x4f\xe4\xb8\x13\x3d\x77\x3e\x45\xfd\x10\x1a\x75\x46\x3d\xe1\xde\x12\x07\x06\x7a\x24\xf4\x63\x19\x86\x66\x4e\xa3\xd1\x8e\x93\x54\x1a\xef\x26\x76\xc6\x76\x60\x51\x94\xef\xbe\x2a\x97\xd3\x24\xfd\x0f\xd8\xd9\xcb\x72\xc1\xa9\xb2\xcb\xef\xbd\xb2\xab\xdc\xb5\xc8\xfe\x14\x2b\x84\xb6\x4d\x6e\x78\xd8\x75\x51\x24\xab\x5a\x1b\x07\xd3\x08\x00\xa0\x6d\xc1\x08\xb5\x42\x38\xfe\x7d\x06\xc7\xc1\x35\x3f\x85\xe4\xee\xa9\xc6\xe4\xd2\x7f\x5b\xf8\xd0\x75\x7e\xf6\x51\xdb\x86\x39\x5d\x77\xd4\xaf\x47\x95\x43\xef\xcf\x85\x13\xa9\xb0\x78\x62\x7f\x96\x47\xd1\xe4\x68\x25\xdd\x7d\x93\x26\x99\xae\x4e\x6a\x6d\x51\xa1\x39\xd1\xa6\xda\xef\x39\xc9\x74\x55\x69\x75\x60\x42\x2e\x45\x89\x99\x3b\x8a\xe2\x28\x3a\x39\x01\x27\xd2\x12\x41\x5a\x58\x7e\xb9\x0a\x1f\x4a\x54\x18\x65\x5a\x59\x17\x0c\xa7\x84\x9b\x09\xdd\x91\x81\xb0\xd3\xda\xcc\xa0\x70\x78\xae\xcb\xa6\x52\x76\xe9\x84\xc3\x0a\x95\xb3\x20\x0c\x42\xc6\x56\xc8\xb1\x90\x4a\x3a\xa9\x95\x05\xa9\x20\x97\x45\x81\x06\x95\x83\x80\xc3\x46\x0f\xc2\xec\x8d\x74\x0a\x95\xa8\xbf\x59\x67\xa4\x5a\x7d\xe7\x7f\xed\x0e\xd9\x73\xaf\xf8\x45\x88\x38\x52\x3b\x4f\xae\x45\x45\x88\xe7\xe1\x73\x73\x13\x38\xf6\xcc\xba\xee\x68\x36\x4c\x08\xc5\xe8\x3c\xcb\xb3\x9b\x4b\xd2\xc7\xdd\x23\x48\xe5\xd0\x14\x22\x43\xd0\x85\x37\x7c\xbe\xfd\x0d\x74\xfa\x07\x66\x2e\x72\x4f\x35\xf2\xdc\xf5\x24\x86\x7a\x5e\x6a\x8b\xd3\x18\xd0\x18\x6d\xd8\xe2\xd9\x4e\x63\x78\xcf\xa3\x8f\x8d\x2c\x73\x64\xdf\x12\x89\x03\xf9\x78\x34\xf4\x5d\x2a\x8b\xc6\xfb\x78\x34\xf4\x7d\xad\xf3\x10\x93\x47\x43\xdf\x05\x96\xc8\x3e\x1e\xf5\x3e\xef\xbc\xd2\xab\x15\x9a\xa9\x36\x55\xc2\xc3\x38\xf0\xfe\xd2\xa0\x91\x68\xb6\xb9\x17\xda\x80\x80\xe5\xe2\x6a\x71\x7e\xe7\x8f\x8d\xed\xb5\x64\x11\xd6\x0b\x37\x84\x20\xfb\xd3\x34\x86\xe9\xb7\xef\xfd\x71\x5a\xfc\xe5\x38\x3d\x33\x96\xa7\xdf\xfb\x5c\x37\xb4\xf8\x2d\x7b\x7b\x57\x46\xeb\xa4\x5a\x41\xdd\x18\x3a\xf4\x96\x11\xad\xc3\x6d\xa6\x86\xec\x63\x44\x0c\xc7\x3b\x36\x30\x7d\x92\xc6\xfe\x13\x4c\x2b\x74\x1e\x92\x56\xe5\x13\xad\xa4\x58\x05\xc5\x02\xe9\xb0\x4a\x40\x16\xa0\xb4\x1f\x43\x25\x5c\x76\x8f\x1c\xff\x27\xa9\x35\x03\xa1\xe0\x07\xe5\x66\x61\xcc\xb5\x76\x9f\x74\xa3\xf2\x1f\xf0\x28\xcb\x12\x52\x04\x83\xae\x31\x0a\xf3\x84\x69\xae\x11\x0e\x68\x4e\xbc\x91\x38\xbe\x7f\x49\xf4\xcf\x35\x2a\xd0\x35\xd2\xad\x0d\x65\x08\x32\xad\x14\x66\x74\x7d\xa3\xa2\x51\x99\x9f\x33\xcd\x8d\x7c\x40\x43\x31\x66\x7e\xe6\x52\x37\x26\x43\xfa\x06\xbe\xa3\x31\x4c\xcf\x6e\x2e\xfb\xf0\x84\x22\x4f\xfd\x17\xdd\x53\xfb\xb3\x4c\x5e\x08\x13\x47\x13\x59\xf8\xf9\xff\x3b\x05\x25\x4b\x8a\x30\x61\xb2\xf4\xe9\x43\x45\x93\x2e\x9a\xe4\xeb\xa8\xa1\x9c\x24\xd7\xf8\x38\x08\xfc\xea\x40\xc1\xf6\x8e\xf8\xb6\x21\xd6\x1c\xf2\x19\xe4\xe9\x1c\xf2\xb4\x9b\xd1\xf4\xa0\xd3\x35\x3e\x06\xe5\x2d\xa5\x87\x96\x84\x22\x00\x85\xd1\x15\x08\xc8\x53\x90\xca\x3a\xa1\x32\x64\xd9\xc6\xa8\x82\x4a\x14\x1c\x28\xb7\x17\x1f\x77\x08\xf6\x2f\x31\xa3\x63\xfe\x16\x72\x5c\x91\x9e\xf9\x41\xca\xc5\x82\x2a\x9e\x50\xfe\x78\x9f\xdf\x2e\xce\xee\x16\x83\x5b\xef\x29\x4e\x33\x78\x4f\x1b\xc4\xfb\xca\x1b\x21\xec\x91\x8c\x1c\x84\xbc\x16\x46\x54\x76\x0e\xdc\xbe\x12\xf6\xdf\x78\x23\xb9\x89\x85\x6f\x3d\x73\x6e\x49\xb3\x60\xdb\x2c\xe6\xf3\x7d\x9d\xe4\x5b\x96\xac\x75\x14\x15\x4e\xe3\xef\x5c\xed\xe9\xaf\x9b\x45\x3e\x18\xa1\x9f\x43\xc6\x8e\xbe\xfa\x73\x15\x3e\xa8\x47\xb8\xf9\x7b\xf5\xd8\x53\xd2\x49\x0f\x4b\x09\x7e\x37\x32\xef\x50\x83\xfd\x41\x8d\xa1\x0c\x84\x7c\xb2\x86\x1d\xb2\x6d\x13\x5e\xde\x77\x3a\x38\x85\x77\x36\xb1\x3e\x46\x68\x40\x21\x0d\x36\x70\xe4\x6e\x72\x90\xe3\xe5\xf5\x72\x71\x7b\x80\xe3\x9e\xd6\x34\xcc\xf9\xc8\xb1\x83\x25\xfb\x5f\xc3\x92\x51\x73\x9f\x3b\x88\xfa\xeb\xcd\xc5\xc1\x93\xba\xa7\x69\x0e\x51\x8f\x1c\x3b\x50\xb3\xff\x45\xd4\xc3\x23\xc5\x4d\xf8\x20\xf0\x8b\xc5\xd5\xe2\x10\xf0\x3d\x1d\x7d\x08\x7c\xe4\xd8\x01\x9c\xfd\x6f\x02\xce\x19\x1a\x77\xca\x61\x25\xdc\x75\x50\xd6\xe4\x0a\x59\x96\x98\xc3\xa3\x74\xf7\xf0\x20\xca\x06\xad\xe7\x0b\x2b\xf9\x80\x7d\xfd\x0c\x34\xd3\x8d\x53\x14\xef\xdc\x79\x5a\xc3\x76\x4b\xdb\x75\x00\xdb\xf6\xc3\xf0\xb1\x58\x3c\x3f\xcf\x3f\x49\x2c\x73\x4b\x8f\x6f\x3f\x49\x16\x70\x5c\x24\x97\x76\x89\x8e\x1f\xbe\xe4\x48\xfb\xeb\x74\x66\xad\x5c\x29\x5f\x4d\x92\xb3\x3c\x9f\xd2\x73\xb2\x08\x97\x8c\x5e\x8f\x50\x27\xde\xc2\x40\xda\xb6\x8f\x76\x8b\xfe\xc5\x9b\x51\x38\x9e\xe1\xf7\xbe\x31\xb2\x12\xe6\xe9\xff\xf8\xb4\x5e\x81\x2a\xef\xba\x98\xb1\xf4\x6f\xd0\xf5\x47\xf7\xdc\xa5\x52\xca\xc7\xf3\x15\xd8\x48\x48\xc3\xf7\x22\x68\xec\x1f\x27\x65\x49\x2d\xa7\xa1\x06\xe5\x19\x3f\xeb\x3c\x3a\xde\xf1\xce\x80\xfb\x74\xde\xba\x32\xaf\xd1\x99\x5f\xd8\xff\x31\xa9\x5f\x43\x2c\x6a\xdb\x6d\x52\x14\xd5\xf7\x10\x37\x80\x0b\x16\x9d\xe5\xf4\x84\x37\x2b\xd1\x82\x11\x45\xfa\xa9\x44\x8f\xc0\x3d\x65\x77\xfb\x82\x8c\xb7\x98\x72\x74\x6f\x59\xa2\x7b\xf9\x8a\xbc\x5e\x79\x1f\x39\x1e\xeb\xf3\x4b\x14\xf7\xd4\xe8\xed\xb3\xf9\x36\x8a\x5b\xa7\xf3\xd7\x28\x0e\x7e\x15\x0e\x86\x7f\x07\x00\x00\xff\xff\x35\x66\x07\x61\x20\x10\x00\x00")

// FileOrmCommonGo is "orm_common.go"
var FileOrmCommonGo = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\xac\x55\x4d\x6f\xe3\x36\x10\x3d\x9b\xbf\x62\xe0\x93\xbc\x28\xe8\x5f\xe0\xc3\x26\x56\x51\x03\x41\xb6\x4d\x5c\x6c\xaf\x14\x35\x96\xd9\x52\xa4\x4a\x8e\xd6\x0e\x16\xf9\xef\x05\x3f\x1c\x47\x4a\xa4\x14\x8b\xcd\x85\xf1\x70\xe6\xcd\xe3\x7b\xa4\xa6\x13\xf2\x1f\xd1\x20\x50\xa7\x3d\x63\xaa\xed\xac\x23\x28\xd8\x62\x29\xad\x21\x3c\xd3\x92\xb1\xc5\xb2\x51\x74\xec\x2b\x2e\x6d\xbb\xee\xac\x47\x83\x6e\x6d\x5d\xbb\x9c\xdc\x59\x4b\xdb\xb6\xd6\xcc\x24\xd4\x4a\x68\x94\xb4\x64\x2b\xc6\xd6\x6b\x90\xd6\x18\x70\xd8\x39\xf4\x68\xc8\x83\x80\xed\x4d\x0c\xa2\x24\x65\x0d\x1c\xac\x83\x56\x18\xd5\xf5\x5a\x90\x32\x0d\x08\x68\xd4\x37\x34\xe0\xc9\xf5\x92\x78\xc0\xf8\xac\x35\x1c\x7a\x13\x0b\x3c\x88\x6f\x42\x69\x51\x69\x04\xb2\xa0\x0c\xa1\x13\x92\xe0\xa4\xe8\x08\xc2\xc0\xe3\x1f\x77\x40\x69\xf7\x28\x08\x94\x07\x87\x5a\x10\xd6\x01\x88\x2c\xd0\x51\xf9\x8c\xfd\x0b\x08\x87\x50\x5b\x83\x50\x3d\x85\x5a\x65\x3c\x09\x23\x11\xec\x21\xe5\xd9\xea\x6f\xcc\x1c\xf6\x16\x1a\xa4\x71\x96\x75\x2d\xf4\x1e\xe1\x4b\x87\x06\xac\x83\x7b\x3c\x5d\x99\x72\x46\x4f\x1d\x26\x05\x52\x47\xf8\xce\x16\x59\x1f\xc8\x2b\xdf\xa6\x95\x2d\xea\x0a\xe2\x9f\x75\x2d\xdf\xde\xb0\x85\xb6\x4d\x83\x2e\xfd\xbe\x8b\xff\xb3\x67\xc6\x02\x3a\x14\x12\x3e\x05\xd8\x15\xdc\x6a\xeb\xb1\x58\x01\x3a\x67\x5d\x80\x77\x48\xbd\x33\x20\x79\x5d\xf1\xbc\x19\xca\xd6\x6b\x48\x18\xe0\x31\xba\x90\xd1\xa3\x22\x99\x63\xbe\x30\xe3\x16\xa9\xae\xc8\x05\x57\x36\xab\xd0\x4e\xf2\x1c\xdf\x64\xc4\xdc\xec\xd6\xa1\x20\xbc\xe9\x95\xae\xd1\x41\x15\x56\x7f\xb1\xe7\xf6\xa1\xfc\xbc\x2f\xc1\x93\x20\x6c\xd1\x10\x74\xc2\x89\x16\x09\x9d\x4f\x8a\x0d\x8b\xaf\xd2\xc5\x3c\x0f\xe9\x0a\xf2\x94\xf5\x7b\x8c\xb1\x45\x3c\x01\x24\xce\x99\xc3\xee\x70\x6f\xa9\x3c\x2b\x4f\x3e\x9d\x7a\xf7\x2b\xdc\x7f\xd9\x43\xf9\xd7\xee\x71\xff\x18\x6f\x5e\x38\x7b\xa6\x13\x98\xbd\x50\xca\x1a\x54\xf0\x69\xc0\x65\xf5\x1a\xb3\x58\x8d\x76\x03\xc5\x8a\x27\x92\xfc\x75\xf3\x0d\x90\xeb\xf1\xc5\x9b\xea\xa2\x51\x7a\x87\x89\x5b\x76\x21\x06\x2e\xcc\x02\xa5\x7f\x7b\x74\x4f\x93\x74\x32\x44\x21\xe9\x7c\xa9\xe6\x39\x36\xcb\xee\x96\xce\xb0\x01\x49\xe7\x37\xa4\x76\xc6\xa3\xa3\xb7\xc6\xed\xee\x1f\xcb\x87\xfd\x8c\x69\xc3\xc2\x29\xd3\x52\xd6\x8c\x69\x3f\x22\xca\xa0\xf5\x47\xa2\x0c\x79\xfe\x3f\x51\xfe\xec\xea\xb7\xb7\x39\x10\xf9\x50\x95\x61\xe5\x94\x2a\x29\x6b\x46\x95\xaf\x47\x74\x78\xd5\xe4\xeb\x6f\xe5\xc3\xeb\x07\x94\x5f\xf1\x3b\xd2\x0c\xfa\xaf\x12\x4e\x71\x8a\x68\xb9\x77\x0c\xad\x46\x99\x03\x5d\x52\xf3\x0d\xc4\xb2\x9f\x72\x8b\x47\xac\xe6\x0d\x9b\x26\x36\x6d\xd8\x16\x35\xbe\x6b\xd8\xb6\xbc\x2b\x67\xbf\x3d\xc3\xca\x29\xc3\x52\xd6\x87\x86\x89\xae\xd3\x0a\x3d\x5c\x04\x37\xb5\x4a\x33\xcc\x9a\x29\x69\x06\xfd\x5f\x0c\x1b\x9b\x35\x64\xf9\xae\x59\x3f\xc5\xa8\x11\x9b\x79\xa3\xa6\x49\x4d\x1b\xa5\x6d\x03\xea\x70\x9d\x47\x27\x11\xbf\xd5\xe3\x19\xa4\x6d\x53\xc4\xb1\xad\x4c\x13\xc6\x76\xe3\x81\x73\x1e\x87\xff\x41\x48\xfc\xfe\x1c\x87\x91\x3a\xc0\x75\x1e\x6d\xc0\x28\x1d\xa2\xb9\x25\x5b\x3c\x5f\xc7\x55\xe1\x13\x0a\xe7\x3c\x8c\xc7\xff\x02\x00\x00\xff\xff\x4d\x16\xa8\xf0\x2b\x09\x00\x00")

// FileParseGo is "parse.go"
var FileParseGo = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\xac\x92\x4d\x8b\xea\x30\x14\x86\xd7\xcd\xaf\x38\x37\x20\xb6\x58\x4b\x14\x75\x21\x78\x17\x77\x21\xdc\x9d\x0c\xe3\x4a\x5c\xd4\x4e\x5a\x0e\x93\x26\x25\x39\x1d\x18\x06\xff\xfb\x90\x58\x6a\x5d\xcc\x87\xc3\xec\x52\x78\x3f\x9e\xbc\x69\x93\x17\xcf\x79\x25\x81\x1a\xe5\x18\xc3\xba\x31\x96\x20\x66\x11\x57\xa6\xe2\x2c\xe2\x8e\x6c\x61\xf4\x4b\x77\x44\x5d\x39\x7f\x24\xac\x25\x67\x09\x63\x65\xab\x0b\x68\x72\xeb\xe4\x7f\x4d\xb1\x83\xc3\xf1\xf4\x4a\x32\x01\xd4\xb4\x5a\xc0\x1b\x8b\x30\x05\x69\x2d\xac\x37\xd0\x25\x65\xbb\x5e\x1d\xf2\x62\x97\xa4\x30\x13\x29\xac\x16\x09\x8b\xb0\x0c\xf2\x3f\x1b\xd0\xa8\xbc\x3f\x52\xa6\xca\x76\x16\x35\x95\x31\xdf\xe6\xa8\xe4\x53\xe8\x43\x5d\xc1\xc8\x01\x19\x5f\xc5\x53\xe8\xc3\x12\x16\x9d\x59\x64\x25\xb5\x56\x03\xb2\xf3\x90\x71\x7f\x0b\xd9\x7e\x45\xb9\xc7\xdf\xc3\x6c\xef\xe0\xdc\x2a\x93\x0f\x41\x4b\xff\xfd\x29\x69\xe7\xb8\xa2\xfe\x14\x33\x54\x7d\x97\xf3\x11\x6b\xd9\x63\xa6\xd0\x58\x59\xa0\x43\xa3\xfd\x9b\x24\xe0\x7f\x92\xcc\x4b\x7c\x73\x69\x6c\x9d\x93\xe7\xe6\x73\x21\x56\x53\x31\x9b\x8a\x39\xcc\x96\x6b\xb1\x58\x8b\x25\x0f\xac\x57\xff\x5f\x10\x01\xb7\x73\x4d\x36\xc0\x33\x0e\x93\x8e\xca\x65\x0f\xb2\x91\x39\xc5\x5c\xf0\x41\xeb\x85\x94\xfa\x81\x42\x7f\x58\x27\xbe\xe4\xdc\xde\xea\xae\x71\xc6\x23\x37\xf6\xf3\xf4\x77\xfa\x70\x22\xba\x9d\xe8\x9f\x31\x6a\xf0\x92\x27\x63\x42\x57\x27\x76\x07\x71\xf4\x04\x82\x9d\xd9\x7b\x00\x00\x00\xff\xff\xeb\xe7\x7c\x37\x8a\x03\x00\x00")

// FileSelectGoTpl is "select.go.tpl"
var FileSelectGoTpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\xb4\x56\x4b\x8f\xdb\x36\x10\x3e\x9b\xbf\x62\x20\xf8\x20\x2d\x5c\xea\x52\xf4\x50\xc0\x87\xc6\xdd\x14\x29\x16\xc9\xb6\x0e\xda\x43\x10\x14\x34\x39\xf2\xb2\x2b\x91\x0a\x49\x25\x36\x04\xfd\xf7\x82\x0f\x19\x96\x5f\x8b\x2d\x90\xcb\x2e\x35\x9a\xf9\xe6\xfb\x38\x0f\xab\x65\xfc\x99\x6d\x11\xfa\x9e\x3e\xc6\xe3\x30\x10\x22\x9b\x56\x1b\x07\x39\x01\x00\xc8\xb8\x56\x0e\x77\x2e\x23\xb3\x4c\x30\xc7\x36\xcc\x62\x69\xbf\xd4\xa5\x30\xf2\x2b\x9a\x8c\x90\x59\xb6\x95\xee\xa9\xdb\x50\xae\x9b\xb2\xd5\x16\x15\x9a\x52\x9b\xa6\xe4\xba\x69\xb4\xca\x22\x4c\xdf\xd3\x8f\xfb\x16\xe9\xbb\x00\xfe\xc8\xdc\xd3\x30\x64\xa4\x20\xc4\xed\x5b\x84\x35\x67\x4a\xa1\x01\xa9\x1c\x9a\x8a\x71\x84\x3e\x84\xad\x74\xdd\x35\xca\xe6\x05\x7c\xfa\x6c\x9d\x91\x6a\x1b\xcc\x6f\xa5\xb1\x2e\x17\x92\xd5\xc8\x1d\xc4\x17\x0b\xf8\xca\xea\x0e\x2d\x7c\xfa\x1c\xa9\xd1\xbf\xfc\x73\x01\xf9\x5d\xdf\xcf\x63\xf2\xfb\x9d\x7b\xcf\x9e\x51\x0c\xc3\x02\xd0\x18\x6d\x0a\x32\x10\xd2\xf7\x60\x98\xda\x22\xcc\xff\x59\xc0\xbc\x82\x9f\x97\x90\xfc\xff\xc4\x0a\x0d\x2a\x8e\x16\x86\x21\x32\xed\xfb\x79\x45\xdf\xb3\x06\x87\xe1\xbb\x93\xae\x6e\xb0\xee\x7b\x40\x25\x3c\x2d\x52\x96\x30\xde\x6e\x24\xb6\xd2\x9d\x72\x20\x2d\x30\x9f\xa6\xe3\x0e\x2a\x6d\x80\x7b\xab\x54\x5b\x30\xfa\x9b\x05\x5d\x41\xd2\x73\x1c\x49\x2e\xd8\x22\x5a\x02\x8a\x0a\x47\x87\xc0\x2c\xc4\x45\xdd\x21\xad\x72\x3f\xfd\x48\x22\xad\x35\x7a\xad\x6f\x3a\x59\x0b\x34\xb0\xf1\xff\x2d\x30\x05\xeb\x3f\x1e\x60\x7d\xff\x70\xbf\xfa\x08\xd6\x31\x87\x0d\x2a\x07\x2d\x33\xac\x41\x87\xc6\xa6\x9e\x98\x04\x1f\xf2\xcf\x82\x9f\x85\xd8\x5d\x34\x7a\x3d\x06\x1b\x99\x71\xad\x14\xdc\xf9\xbf\x64\x66\xc3\x1b\x6d\x60\x3c\x78\x52\x55\xa7\x38\xe4\x1b\xb8\x9b\xa0\x17\x63\x03\xe6\x87\x53\x52\x6a\xd0\x75\x46\xc1\x86\xc6\xac\x34\xd5\x96\xe6\x77\x23\x6a\x91\xb4\xfe\xfd\x84\x06\x81\xb5\x6d\x2d\xd1\xc2\xb7\xf0\xc4\xb5\x12\xd2\x49\xad\x2c\x68\x05\xee\x09\xe1\x4b\x87\x66\x7f\x95\x45\xc0\xc8\xc7\xd8\xa0\x2f\x98\x8a\x13\x4f\x7f\x0d\x07\x4a\x31\xf1\x32\xa6\x24\xb3\x91\x71\xa2\xf5\x20\x1b\xe9\x0e\xb4\x42\xe9\xeb\x60\x3a\x26\x04\x06\x6d\xab\x95\xc5\xab\xcc\x02\x4c\x1e\x23\x43\x85\x6f\x52\x7a\x64\x5b\xa4\x31\xf3\x32\xa6\x3b\xe3\xe5\x5d\xa6\xb4\x74\x55\x59\x74\xc0\x94\xf8\x5f\x0c\x3d\x60\x1e\x31\x16\xf0\x2a\xa2\x1f\x62\xe2\x65\x62\xf0\x0a\x1d\x17\x36\x47\x1c\x8c\xb7\x12\x7d\xab\xff\x30\x84\x41\x95\x15\x28\xed\x60\x5e\xd1\x77\xf6\xb0\x52\xfc\xe8\x1e\x46\xe4\x68\xab\x00\x13\xc2\x1e\xaf\x19\x70\x3a\xdc\x43\x6c\x38\x14\xc0\x43\x0f\xfa\x09\x66\x2f\xf4\xd3\x19\x78\x7e\xe1\x36\x7c\x9b\x6f\xe8\xd8\xce\xf4\x9c\xd0\x12\x9c\xe9\x70\x32\x0e\x69\x01\xd5\x76\x94\xf1\xbb\x96\xea\x44\x04\x30\xf8\x57\x4b\x95\xea\xe7\x37\xd0\x91\xc3\x55\xca\x27\x40\xb9\x1d\xc7\xf1\x6c\xef\xbe\x2c\xe5\x94\xd4\x12\x12\xda\x25\x2d\x69\x99\x4e\xca\x15\x6a\xb9\xae\x25\xc7\x50\xca\xb2\x84\x0f\x46\xa0\x79\xb3\x3f\x46\xf5\xbd\xa3\xbd\x79\x2c\xd4\xa1\x61\xbb\xda\x59\x60\x9c\x6b\x23\xfc\xde\x75\x7a\x2c\x5d\x08\x8f\xab\xe4\xc6\x55\x9c\xe7\xca\x85\x34\xe3\x62\x08\x6f\x7f\x95\xd7\xaf\x21\xf5\x70\xf0\xb3\xf4\x17\x21\xf2\x6c\x92\x37\x5b\x80\x90\xa6\x38\xbd\x0b\xaf\xf2\x37\xa3\xbb\x76\xaa\xb2\x61\xcf\x78\xa4\x6e\xeb\x3d\x60\xb3\x7f\x9d\xa2\x73\xdc\xfc\x25\xfa\x21\xe4\x22\xfd\xe2\x4a\x15\xd3\xd0\x8d\x47\xaf\x67\x15\xbf\x61\x7c\xa9\x6c\x10\x91\x3e\x6a\x42\x57\xfa\x67\xff\x83\x74\x7b\x94\x12\x44\xce\xdd\x6e\x8c\xa6\xc9\x76\x73\xc3\xac\xdc\x0e\x96\xc0\xdd\x6e\xb2\x39\xfe\x0b\x00\x00\xff\xff\x8d\x92\xf5\x18\x79\x09\x00\x00")

// FileSelectorGoTpl is "selector.go.tpl"
var FileSelectorGoTpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\xbc\x58\x4f\x6f\xdc\xb6\x12\x3f\x5b\x9f\x62\x22\xac\xf3\xa4\x44\xd1\x5e\x1e\xde\x61\x1f\x7c\x08\xfc\x12\x20\xaf\x6d\x9a\xc6\x6e\x2f\x86\x51\xd0\x12\xe5\xa5\x23\x91\x0a\x49\x79\xbb\x50\xf4\xdd\x8b\xe1\x1f\x89\xbb\xab\x5d\xbb\x48\x51\x5f\x2c\x92\xc3\xe1\xcc\x6f\x7e\x9c\x19\x6e\x4b\x8a\x2f\xe4\x9e\x42\xdf\xe7\x9f\xec\xe7\x30\x44\x11\x6b\x5a\x21\x35\x24\xd1\x59\x5c\x12\x4d\xee\x88\xa2\x4b\xf5\xb5\x5e\x96\x92\x3d\x52\x19\x47\x67\x71\xd5\xe8\x38\x02\x00\xe8\x7b\x90\x84\xdf\x53\x58\xfc\x9e\xc1\xc2\x6d\x5c\x5d\x40\x7e\xbd\x6d\x69\xfe\xc1\x8c\x15\xbc\x19\x06\x23\x1d\xf7\xbd\x93\x19\x86\x71\x3f\xe5\x25\xf8\xf5\x7b\xa6\xd7\xdd\x5d\x5e\x88\x66\xd9\x0a\x45\x39\x95\x4b\x21\x9b\x65\x21\x9a\x46\xf0\x38\x4a\xa3\xa8\x10\x5c\x69\xa0\x52\xfe\xa4\xee\xe1\x02\xe2\x42\xf0\x47\x2a\x35\xe3\xf7\x70\xae\x56\x50\x88\xba\x6b\x38\x9c\x97\xb0\x61\x7a\x0d\x8f\xa4\xee\x28\x9c\x3f\x42\xa2\xb7\x2d\x85\xf3\xeb\x14\xb4\x80\x73\x15\x47\xd1\x72\x09\x8a\xd6\xb4\xd0\x42\xba\x0f\xe5\x76\x2b\xa8\x84\x84\xab\x5f\x7e\x84\xaf\x1d\x95\x8c\x2a\x20\xbc\x34\x73\x2d\x91\x0a\x4f\xc2\x35\x29\x36\x2a\x32\x5a\x27\x35\x5a\x76\x85\x86\x7e\x0f\x19\x96\xc1\xa2\x9a\x40\x79\xcf\x68\x5d\x4e\x98\xf4\x3d\xb0\x0a\xb8\xd0\xb0\xa8\xf2\x0f\xea\x33\xad\xa8\xa4\xbc\xa0\xa3\xc0\x95\x51\xdf\xf7\x8b\x2a\xff\x48\x1a\x3a\x0c\x70\x27\x44\x3d\xa2\x57\xab\x49\xf4\xff\x82\xf1\x50\x30\xf8\xbe\x2a\x08\xe7\x54\x86\xa0\x07\x16\x84\xc3\x42\x74\x5c\x9b\x33\x60\xb9\x84\x4e\x51\xeb\xbb\xfa\x5a\xc3\xe5\xcf\xbf\x7e\xbc\x4e\x5e\xa5\x0e\xa8\x68\x30\x30\x5e\x3a\xd4\x88\xa4\xa0\xd7\x14\x38\x69\xa8\x02\x51\x39\x60\x68\xe9\x71\x8d\xaa\x8e\x17\x90\x28\x78\xe5\x21\x4b\xfd\xe6\x24\x85\x9b\x5b\xa5\x25\xa2\xdb\x47\x67\x8f\x44\xe2\x26\x35\x4e\xce\x90\xed\x3b\x20\x65\x15\xa8\xfc\x10\x57\x1b\x38\x8b\x41\xad\xe0\x02\x48\xdb\x52\x5e\x26\x38\xca\x0c\x79\xab\xdc\xda\x3b\x0c\x71\x6a\x84\x87\x27\x00\x3d\x93\x54\x77\x92\x1b\x85\x0e\x2e\x0c\x92\x05\xeb\x41\x30\x0e\xa2\xd5\x4c\x70\x83\x17\x82\x87\x94\xdb\xce\x01\x65\xb6\x19\x98\xec\x6d\xc8\x71\xe2\x13\x91\xa4\x51\x1e\xb1\x07\xa3\x79\x46\xe0\x34\x7a\x23\x3e\xfb\x08\x2e\x2a\x2b\x70\x55\xb3\x3d\xf0\x3c\xe3\x57\x17\xa0\xf2\x3d\xd6\xfd\x77\x5a\x7e\x71\x01\x9c\xd5\x01\xae\xcb\x25\x90\xb2\xc4\x20\x5b\xdf\x2b\x10\x9c\xe2\x95\x6c\x08\xdf\x82\xa4\x35\x41\x30\x32\x7f\x8d\x51\x44\xaf\xa9\x04\x73\xcf\x5a\xc1\xb8\x56\x28\xad\xd7\x4c\x99\xb9\x51\xb1\x75\x7d\x8c\x98\x19\xa2\x9a\x3d\x20\x26\x4b\xf0\xef\xbd\x90\x94\xdd\xf3\x1f\xe8\x76\xe5\x45\xa7\xa9\x5d\x51\x67\xfc\x64\xd7\x68\x02\xe8\x35\xd1\x26\x74\x81\xa9\x4c\x59\x6b\xd1\x53\xc1\x0f\x34\x59\x12\xad\x1c\xa5\xa6\x33\x27\x76\x65\x73\xa7\x07\x07\x68\x72\x57\xd3\x03\x99\xcf\xb4\xba\xc6\x85\x19\xcd\x7e\xe9\x19\xba\x9d\x93\xc6\xaf\x53\x98\x07\xa7\x1e\x75\x68\x5c\x3b\x38\x77\xd8\x1d\xda\xdb\xe8\xb2\xc1\x6a\xa4\x50\x3e\xe6\x87\x49\x7c\xd8\xbf\x7a\x61\x0a\xfc\x4e\x76\x1a\x5a\xba\x88\xfa\x00\x4f\x18\x10\x3e\xc1\xf4\x2f\x05\xad\x64\x0d\x91\x5b\xf8\x42\xb7\xa1\x8e\x71\xa7\x82\x1b\x6b\xfe\x6d\xa0\xe3\xc6\x47\xe2\x36\xbf\x19\xd1\xb9\xfd\x47\x89\xbc\x1b\xad\x13\x9c\xdb\xe7\x93\xc9\x06\xc7\x69\x74\x40\x04\x23\xff\xc9\xa2\x74\x82\xde\x7f\x23\x15\xa6\x5e\x62\x1a\xf9\x1c\x6c\xb0\x1c\x6b\x16\x16\x39\x8c\x93\xec\x28\x6c\xd6\x94\x03\xd9\x2f\x6f\xa0\xd6\xa2\xab\x4b\xb8\xa3\x98\xb4\x68\x69\x6f\xc1\x89\x14\x6d\x94\x26\xa9\xad\x9d\x16\x76\x77\xb4\xca\x4d\x55\x75\x87\xbf\x67\x52\x69\x6b\x81\x2a\x08\x16\x03\xee\x5b\x0a\xc3\x32\x6c\xc7\x0c\x74\x96\xb2\x4e\xd2\x34\x17\x73\xc7\x4e\xea\x92\x92\x11\x9c\x05\x5b\x34\x33\xec\x7f\xb0\x20\xd8\xc6\x2d\xff\x0d\xbb\xa1\x14\x92\x57\x33\xfa\x33\xec\xa9\x50\x9b\xb5\x5b\x6d\x98\x2e\xd6\xe0\xf5\xf9\x96\xe6\x4d\x58\x42\xfc\xe2\xea\x02\x16\xf9\xff\xec\x40\x79\xf8\x0b\xa2\xa8\x61\x81\x93\x72\x67\xc5\xab\x31\x76\x23\x34\x88\xc1\x81\x60\x82\xa6\xa7\x73\xd5\xb5\xa4\x15\xe9\x6a\x7d\xa0\x88\xb3\x3a\x83\xaa\xd1\xf9\x3b\x74\xa4\x4a\xe2\x8e\xab\xae\xc5\x66\x93\x96\xa3\x23\xe7\x2a\xce\xfc\xc0\x73\x67\x18\x63\xf2\xac\x70\x3c\x15\x89\xbf\x1a\x84\x77\x7f\x68\xab\x78\x2f\x04\x4c\xd3\xc6\x4c\xd9\x3c\x76\x18\x65\xab\x39\xf5\x69\x0f\x25\x0f\x52\x5a\x88\x0d\x95\x32\xb8\x2d\x6e\xe5\x25\x1e\x93\xef\x7a\x98\xa1\x3c\x72\x75\xb7\x65\x38\x1a\x6f\x6c\xa5\xe7\x62\xe8\xe0\x3c\x9c\x0f\x90\x5d\x3c\x13\xda\xf9\x03\x92\xa3\xd8\x2e\x9e\x62\x38\x36\x4c\xc9\x84\x93\xd8\xc0\xdc\xa6\x51\x80\xd4\x35\x60\x18\xac\x45\x6f\xeb\x3a\x49\xc7\x35\x06\xcc\x49\xa6\xc7\x3a\xad\xc5\x4e\xa3\xfa\x54\x9f\x3a\xd5\x32\x3c\xf6\xdb\xb7\x27\xba\x55\x56\x19\x26\xdc\xb0\xdb\x89\x00\x01\x56\x97\xf6\x89\x64\xc0\xb9\x14\x25\x85\x45\xe5\x0e\x98\xd8\x60\xd4\xbc\x7e\x3d\x53\x58\x27\x0b\x83\x5e\x70\xd7\xc0\x1e\xab\x9e\xfa\xc2\x5a\x7c\x28\x60\xd5\xc1\x8a\xe8\xdf\x08\xc7\xb5\xcf\xe5\xea\xa9\x3d\xb7\xcf\x90\xc9\x49\x97\x91\x1e\x49\x8d\x78\x3a\x7f\x73\xf3\xae\x4b\xc3\xce\x1d\xf3\x0e\xe3\xfa\x3f\xff\x5e\xed\x14\x15\x29\x36\xb9\x4d\xa4\x66\xf3\xae\xfc\xcd\xed\xdd\x56\xd3\xe3\x1b\xf0\xdd\x47\x3f\x70\x8d\x74\x9b\xc2\x7e\x90\x88\x4e\x24\x23\xfb\x60\xcd\x20\xf6\x35\x26\xce\x80\x65\xde\x8d\xe0\x23\x36\xb6\x67\xce\xa4\x38\x3d\x19\xa7\x27\xd8\x16\x34\xf6\x53\xc0\x1e\x8e\xb4\x45\x0f\x87\xc9\x43\x37\xed\x98\x80\x1e\x6c\x02\x4a\x66\x52\xba\xb7\x7e\x75\x9b\x86\x94\x9c\xcd\x47\xc7\x72\xd2\xae\x87\xa7\x9e\x1f\x3e\x36\xe1\x55\x18\xdb\xa5\xbd\x85\x0c\x5d\x48\x43\xb5\x3b\xbd\xe2\xbc\x2a\xdd\xb4\x3b\x3b\x82\xd2\x33\xcb\x56\x9f\x48\xa5\xd8\xf8\xc4\x19\x48\x8c\x3f\x33\xbc\xad\x6b\x27\xea\x9a\x0e\x73\xb1\x7c\xab\xb1\x21\x0a\x54\x4b\x0b\x56\xb1\x82\xd4\xf5\x76\x7c\x3b\xcf\xe6\xc2\x29\x07\xcd\xf4\x1a\xcf\xf8\xe1\xa1\xef\x8f\xe5\x1d\x78\x31\x97\x6a\x5e\xbe\x84\xbe\xa7\xbc\xc4\x9d\xe6\x1f\x8a\xf9\x8e\xe6\xcf\x00\x00\x00\xff\xff\xe1\xb3\x63\x2f\x44\x12\x00\x00")

// FileWhereGoTpl is "where.go.tpl"
var FileWhereGoTpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\xac\x52\x4d\x8b\xdb\x30\x10\x3d\x5b\xbf\x62\x30\x39\xd8\xe0\xca\xf7\xc2\x5e\xb6\xb8\xd4\x97\xb4\x94\x40\x0e\x21\x14\xad\x3d\xb6\xc5\x5a\x23\x23\x2b\x49\x8b\x98\xff\x5e\x22\xc7\xed\x42\xbb\x10\xd8\xbd\xd9\x33\xef\x43\xef\x49\x93\x6a\x9e\x55\x8f\x10\x82\xfc\xb6\x7c\x32\x0b\xa1\xcd\x64\x9d\x87\x4c\x00\x00\x84\x00\x4e\x51\x8f\xb0\xf9\x51\xc0\xe6\xb6\xfa\xf8\x00\x72\xf7\x6b\x42\x59\xc7\xff\x19\x3e\x30\x47\x74\x1a\xc2\x0d\xc3\x9c\xae\x7c\xa4\x16\x98\x45\x92\xf6\xda\x0f\xa7\x27\xd9\x58\x53\x4e\x76\x46\x42\x57\x5a\x67\xca\xc6\x1a\x63\x29\x15\xb9\x10\x21\xbc\x30\xeb\xfe\xfa\x7c\xd6\x38\xb6\x33\xb3\x28\x4b\xd8\x0f\xe8\x30\x84\x4d\x27\xb7\xca\x20\x33\xa8\xb6\x9d\x41\x41\x63\xa9\xd5\x5e\x5b\x02\x4b\xf0\x72\xef\x2d\xf8\x01\x61\xff\xa5\xfa\x5e\xc1\xec\x95\x47\x83\xe4\x45\x77\xa2\xe6\x1f\xb1\xcc\x4e\xb0\x9c\x47\x7e\x9d\x0a\x38\xab\x71\x91\x8a\xa7\xa8\x7e\xfa\x05\x95\xaf\x98\x48\x87\x20\x12\x87\xfe\xe4\x68\x1d\x6f\xf1\x12\x37\x99\x9d\x8a\xd8\xc9\xc2\xdf\xa9\xa7\x11\x99\xd3\x65\xd6\xc9\x4f\x76\x3c\x19\x8a\x83\xb3\x1a\x73\xc1\xe2\x7f\xf9\x6a\xba\x25\x24\xa8\xb7\xef\x93\xb2\xa6\xec\xac\xc6\x19\xa4\x94\x77\xa5\x53\xae\x9f\xaf\x97\x61\xd4\x33\x66\x87\xa3\x26\x8f\xae\x53\x0d\x06\x2e\x60\xc4\x45\x2c\xcf\x45\xd2\x59\x07\xfa\x0a\x5c\x6e\x31\x7a\x04\x91\x44\xfe\x41\x1f\xe1\x21\x8e\x0e\xfa\x28\x12\x7e\xad\xb3\x9a\xb2\xbb\x1a\xbb\x6a\x4a\x29\x5f\x6b\xed\x11\xfd\x05\x71\xad\x0e\x1e\xab\xdd\xbe\xaa\xde\xa9\xbe\x9b\x76\x36\xda\x4b\x01\x83\xee\x87\xb7\xbc\x91\x55\xec\xae\xd0\x7f\x1c\xaf\xb1\x43\x40\x6a\x99\xc5\xef\x00\x00\x00\xff\xff\x18\x22\x51\xf0\xc6\x03\x00\x00")



func init() {
  if CTX.Err() != nil {
		log.Fatal(CTX.Err())
	}



var err error



  




  
  var f webdav.File
  

  
  
  var rb *bytes.Reader
  var r *gzip.Reader
  
  

  
  
  
  rb = bytes.NewReader(FileExecGoTpl)
  r, err = gzip.NewReader(rb)
  if err != nil {
    log.Fatal(err)
  }

  err = r.Close()
  if err != nil {
    log.Fatal(err)
  }
  
  

  f, err = FS.OpenFile(CTX, "exec.go.tpl", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
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
  
  
  
  rb = bytes.NewReader(FileOrmGoTpl)
  r, err = gzip.NewReader(rb)
  if err != nil {
    log.Fatal(err)
  }

  err = r.Close()
  if err != nil {
    log.Fatal(err)
  }
  
  

  f, err = FS.OpenFile(CTX, "orm.go.tpl", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
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
  
  
  
  rb = bytes.NewReader(FileOrmCommonGo)
  r, err = gzip.NewReader(rb)
  if err != nil {
    log.Fatal(err)
  }

  err = r.Close()
  if err != nil {
    log.Fatal(err)
  }
  
  

  f, err = FS.OpenFile(CTX, "orm_common.go", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
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
  
  
  
  rb = bytes.NewReader(FileParseGo)
  r, err = gzip.NewReader(rb)
  if err != nil {
    log.Fatal(err)
  }

  err = r.Close()
  if err != nil {
    log.Fatal(err)
  }
  
  

  f, err = FS.OpenFile(CTX, "parse.go", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
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
  
  
  
  rb = bytes.NewReader(FileSelectorGoTpl)
  r, err = gzip.NewReader(rb)
  if err != nil {
    log.Fatal(err)
  }

  err = r.Close()
  if err != nil {
    log.Fatal(err)
  }
  
  

  f, err = FS.OpenFile(CTX, "selector.go.tpl", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
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


