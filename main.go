package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"reflect"
	"sort"
)

func main() {
	var (
		inputFileName, outputFileName string
		showHelp                      bool
	)

	config := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	config.BoolVar(&showHelp, "h", false, "help text")
	config.StringVar(&inputFileName, "f", "", "input filename (default stdin)")
	config.StringVar(&inputFileName, "o", "", "output filename (default stdout)")
	config.ErrorHandling()
	if err := config.Parse(os.Args[1:]); err != nil || showHelp {
		fmt.Printf("\nExample:\n\n  %s -f input.yml -o output.yml first_key second_key third_key\n\n", os.Args[0])
		os.Exit(1)
	}

	keys := config.Args()
	keysMap := make(map[string]int, len(keys))
	for i, k := range keys {
		keysMap[k] = len(keys) - i
	}

	var (
		in  io.Reader = os.Stdin
		out io.Writer = os.Stdout
	)

	if inputFileName != "" {
		f, err := os.Open(inputFileName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer func() {
			_ = f.Close()
		}()
		in = f
	}

	if outputFileName != "" {
		f, err := os.Create(outputFileName)
		if err != nil && !os.IsExist(err) {
			fmt.Println(err)
			os.Exit(1)
		}
		defer func() {
			_ = f.Close()
		}()
		out = f
	}

	var recursiveSort func(interface{})
	recursiveSort = func(doc interface{}) {
		document, isDocument := doc.(yaml.MapSlice)

		if array := reflect.ValueOf(doc); !isDocument && array.Kind() == reflect.Slice {
			length := array.Len()

			for i := 0; i < length; i++ {
				elem := array.Index(i).Interface()
				recursiveSort(elem)
				array.Index(i).Set(reflect.ValueOf(elem))
			}

			return
		}

		if !isDocument {
			return
		}

		sort.Sort(sorter{
			len:  len(document),
			swap: func(i, j int) { document[i], document[j] = document[j], document[i] },

			less: func(i, j int) bool {
				iKey, iIsString := document[i].Key.(string)
				jKey, jIsString := document[j].Key.(string)
				if !jIsString || !iIsString {
					return false
				}
				return keysMap[iKey] > keysMap[jKey]
			},
		})

		for i := range document {
			_, isMap := document[i].Value.(yaml.MapSlice)
			if !isMap && reflect.ValueOf(document[i].Value).Kind() != reflect.Slice {
				continue
			}
			recursiveSort(document[i].Value)
		}
	}

	var document yaml.MapSlice

	if err := yaml.NewDecoder(in).Decode(&document); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	recursiveSort(document)

	if err := yaml.NewEncoder(out).Encode(document); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type sorter struct {
	len  int
	less func(i, j int) bool
	swap func(i, j int)
}

func (s sorter) Len() int           { return s.len }
func (s sorter) Swap(i, j int)      { s.swap(i, j) }
func (s sorter) Less(i, j int) bool { return s.less(i, j) }
