// Copyright(C) 2023 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2023/9/2

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/vmihailenco/msgpack/v5"
	"gopkg.in/yaml.v3"
)

var input = flag.String("i", "", "Input file path or url")
var fromEncoding = flag.String("f", "", "Encoding from.\nallow: "+allEncodings)
var toEncoding = flag.String("t", "", "Encoding to.\nallow: "+allEncodings)
var timeout = flag.String("timeout", "30s", "Timeout for HTTP Requests")

func main() {
	flag.Parse()

	if _, ok := decoders[*toEncoding]; !ok {
		log.Fatalf("invalid -t=%q, support types %q\n", *toEncoding, ecTypes)
	}
	from, err := getFromEncoding()
	if err != nil {
		log.Fatalln(err)
	}
	if *toEncoding == from {
		log.Fatalln("encoding from and to are same")
	}
	content, err := fetchContent()
	if err != nil {
		log.Fatalln(err)
	}

	var data any
	if err = decoders[from](content, &data); err != nil {
		log.Fatalf("decode input content as %q failed: %v\n", from, err)
	}

	bf, err := encoders[*toEncoding](data)
	if err != nil {
		log.Fatalf("encoding content to %q failed: %v\n", *toEncoding, err)
	}

	fmt.Println(string(bf))
}

func getFromEncoding() (string, error) {
	if *fromEncoding != "" {
		if _, ok := decoders[*fromEncoding]; !ok {
			log.Fatalf("invalid -f=%q, support types %q\n", *fromEncoding, ecTypes)
		}
		return *fromEncoding, nil
	}
	fileName := *input
	if fileName != "" {
		index := strings.LastIndex(*input, ".")
		if index > 0 {
			ext := fileName[index+1:]
			if _, ok := decoders[ext]; ok {
				return ext, nil
			}
		}
	}
	return "", errors.New("encoding from ( -f ) is required")
}

func fetchContent() ([]byte, error) {
	if *input == "" {
		return io.ReadAll(os.Stdin)
	}

	if strings.HasPrefix(*input, "http://") || strings.HasPrefix(*input, "https://") {
		tm, err := time.ParseDuration(*timeout)
		if err != nil {
			return nil, fmt.Errorf("invalid timeout %q: %w", *timeout, err)
		}
		c := &http.Client{
			Timeout: tm,
		}
		resp, err := c.Get(*input)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		return io.ReadAll(resp.Body)
	}
	return os.ReadFile(*input)
}

var ecTypes []string

func init() {
	for name := range encoders {
		ecTypes = append(ecTypes, name)
	}
}

const allEncodings = "json, toml, yml, msgpack"

var encoders = map[string]func(v any) ([]byte, error){
	"json": func(v any) ([]byte, error) {
		return json.MarshalIndent(v, "", "    ")
	},
	// "xml": xml.Marshal,
	"toml": func(v any) ([]byte, error) {
		bf := &bytes.Buffer{}
		err := toml.NewEncoder(bf).Encode(v)
		return bf.Bytes(), err
	},
	"yml":     yaml.Marshal,
	"msgpack": msgpack.Marshal,
}

var decoders = map[string]func(data []byte, v any) error{
	"json": json.Unmarshal,
	// "xml":     xml.Unmarshal,
	"toml":    toml.Unmarshal,
	"yml":     yaml.Unmarshal,
	"msgpack": msgpack.Unmarshal,
}
