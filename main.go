package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"sort"
)

func main() {
	var (
		inputFileName, outputFileName string
	)

	flag.StringVar(&inputFileName, "f", "input filename", "")
	flag.StringVar(&inputFileName, "o", "output filename", "")
	flag.Parse()

	keys := flag.Args()
	keysMap := make(map[string]int, len(keys))
	for i, k := range keys {
		keysMap[k] = len(keys) - i
	}

	var (
		in io.Reader = os.Stdin
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

	var document yaml.MapSlice

	if err := yaml.NewDecoder(in).Decode(&document); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var recursiveSort func(document yaml.MapSlice)
	recursiveSort = func(document yaml.MapSlice) {
		sort.Sort(sorter{
			len: len(document),
			swap: func(i, j int) {document[i], document[j] = document[j], document[i]},

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
			fmt.Printf("%T\n", document[i])
			d, ok := document[i].Value.(yaml.MapSlice)
			if !ok {
				continue
			}
			recursiveSort(d)
			document[i].Value = d
		}
	}

	recursiveSort(document)

	if err := yaml.NewEncoder(out).Encode(document); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type sorter struct{
	len int
	less func(i, j int) bool
	swap func(i, j int)
}

func (s sorter) Len() int { return s.len }
func (s sorter) Swap(i, j int) { s.swap(i, j) }
func (s sorter) Less(i, j int) bool { return s.less(i, j) }

