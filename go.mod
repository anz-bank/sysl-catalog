module github.com/anz-bank/sysl-catalog

go 1.14

require (
	github.com/anz-bank/protoc-gen-sysl v0.0.0-20200421063430-ac292bed0e56
	github.com/anz-bank/sysl v0.62.0
	github.com/gohugoio/hugo v0.69.2
	github.com/radovskyb/watcher v1.0.7
	github.com/sirupsen/logrus v1.5.0
	github.com/spf13/afero v1.2.2
	golang.org/x/net v0.0.0-20200319234117-63522dbf7eec // indirect
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
)

replace github.com/anz-bank/sysl => ../sysl