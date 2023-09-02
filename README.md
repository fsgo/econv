# econv
Encoding convert


## Install
```bash
go install github.com/fsgo/econv@latest
```

### Usage
```bash
# econv -help
Usage of econv:
  -i string
    	Input file path
  -f string
    	Encoding from.
    	allow: json, toml, yml, msgpack
  -t string
    	Encoding to.
    	allow: json, toml, yml, msgpack
  -timeout string
    	Timeout for HTTP Requests (default "10s")
```

```bash
econv -i passport.toml -t json                 # toml -> json
cat passport.toml | econv -f toml -t json      # toml -> json
```