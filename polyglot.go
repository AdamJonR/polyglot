package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/adamjonr/dialects"
	"github.com/adamjonr/qform"
)

type Config struct {
	ConfigPath   string
	InputDir     string
	OutputDir    string
	InputDirAbs  string
	OutputDirAbs string
	Extensions   map[string][]Lexicon `json:"extensions"`
}

type Lexicon struct {
	Dialect string `json:"dialect"`
	Start   string `json:"start"`
	Stop    string `json:"stop"`
}

type File struct {
	Path    string
	Name    string
	Dir     string
	Ext     string
	PathRel string
}

func main() {
	// store args
	args := os.Args
	// ensure the config argument is present
	if len(args) < 2 {
		fmt.Println("missing json config file argument")
		os.Exit(1)
	}
	// create config object
	config, err := NewConfig(args[1])
	// exit with error
	if err != nil {
		os.Exit(1)
	}
	// declare log summary
	logSummary := ""
	// walk inputDir
	err = filepath.Walk(config.InputDirAbs, func(path string, f os.FileInfo, err error) error {
		// skip directories
		if f.IsDir() {
			return nil
		}
		// prep File
		dir, name := filepath.Split(path)
		ext := filepath.Ext(path)
		pathRel, err := filepath.Rel(config.InputDirAbs, path)
		if err != nil {
			fmt.Printf("file %q not parsed due to relative path error\n", path)
			return nil
		}
		// create file pointer
		file := &File{Path: path, Name: name, Dir: dir, PathRel: pathRel, Ext: ext}
		// skip files without extensions
		if ext == "" {
			CopyFile(config, file)
			return nil
		}
		// skip files with extensions that are not configured with dialects
		_, ok := config.Extensions[ext]
		if !ok {
			CopyFile(config, file)
			return nil
		}
		// skip hidden files
		if name[0:1] == "." {
			CopyFile(config, file)
			return nil
		}
		// call parser
		logSummary = logSummary + ParseFile(config, file)
		// standard return
		return nil
	})
	// handle severe error
	if err != nil {
		os.Exit(1)
	}
	// write log file
	_ = ioutil.WriteFile("polyglot-log.txt", []byte(logSummary), 0777)
}

func NewConfig(path string) (*Config, error) {
	// create pointer to config
	config := &Config{ConfigPath: path}
	// read in config file
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		// provide feedback and return early
		fmt.Printf("json config file %q could not be read\n", path)
		return config, err
	}
	// declare top-level JSON object into map
	var objmap map[string]*json.RawMessage
	// marshall into map
	err = json.Unmarshal(bytes, &objmap)
	// handle error
	if err != nil {
		fmt.Printf("json config file %q could not be parsed\n", path)
		return config, err
	}
	// check for required keys and return early with feedback if absent
	if _, ok := objmap["inputDir"]; !ok {
		fmt.Printf("json config file %q missing inputDir key\n", path)
		return config, err
	}
	if _, ok := objmap["outputDir"]; !ok {
		fmt.Printf("json config file %q missing outputDir key\n", path)
		return config, err
	}
	if _, ok := objmap["extensions"]; !ok {
		fmt.Printf("json config file %q missing extensions key\n", path)
		return config, err
	}
	// marshall values into config
	err = json.Unmarshal(*objmap["inputDir"], &config.InputDir)
	if err != nil {
		fmt.Printf("json config file %q contained invalid inputDir string\n", path)
		return config, err
	}
	err = json.Unmarshal(*objmap["outputDir"], &config.OutputDir)
	if err != nil {
		fmt.Printf("json config file %q contained invalid outputDir string\n", path)
		return config, err
	}
	err = json.Unmarshal(*objmap["extensions"], &config.Extensions)
	if err != nil {
		fmt.Printf("json config file %q contained invalid extensions object\n", path)
		return config, err
	}
	// clean directories
	config.InputDir = filepath.Clean(config.InputDir)
	config.OutputDir = filepath.Clean(config.OutputDir)
	// add absolut directories
	config.InputDirAbs, err = filepath.Abs(config.InputDir)
	config.OutputDirAbs, err = filepath.Abs(config.OutputDir)
	// ensure input and output directories exist
	if _, err := os.Stat(config.InputDirAbs); err != nil {
		fmt.Printf("input directory %q does not exist\n", config.InputDirAbs)
		return config, err
	}
	if _, err := os.Stat(config.OutputDirAbs); err != nil {
		fmt.Printf("output directory %q does not exist\n", config.OutputDirAbs)
		return config, err
	}
	// otherwise, all is well
	return config, nil
}

func ParseFile(config *Config, file *File) string {
	var logSummary string = "File: " + file.Path + "\n"
	var log string
	// store file contents
	bytes, err := ioutil.ReadFile(file.Path)
	// check for error
	if err != nil {
		// exit early if read error
		return log
	}
	// convert to string
	source := string(bytes)
	// create label so other lexicons can be handled
lexiconFor:
	// cycle through lexicons for this file extension
	for _, lexicon := range config.Extensions[file.Ext] {
		// store string array of sections
		sections := strings.Split(source, lexicon.Start)
		// continue to next lexicon if there is no start delimiter (i.e., only one section)
		if len(sections) < 2 {
			continue
		}
		var dsl dialects.Dialectable
		// create the appropriate dsl
		switch lexicon.Dialect {
		case "qform":
			dsl = new(qform.DSL)
		}
		// cycle through sections
		for i, section := range sections {
			// skip first section
			if i == 0 {
				continue
			}
			// get closing section
			selections := strings.Split(section, lexicon.Stop)
			// if only one section, there's no closing tag
			if len(selections) != 2 {
				logSummary = logSummary + "incomplete parsing using " + lexicon.Dialect + " because the start delimiter '" + lexicon.Start + "' has no stop delimiter or an extra stop delimiter\n\n"
				continue lexiconFor
			}
			// parse contents of first selection
			selections[0], err, log = dialects.Parse(dsl, selections[0])
			// check for parsing error
			if err != nil {
				logSummary = logSummary + "incomplete parsing using " + lexicon.Dialect + ": " + err.Error() + "\n\n"
				continue lexiconFor
			}
			// append log to logSummary
			logSummary = logSummary + log
			// put put parsed result back in selections
			sections[i] = strings.Join(selections, "")
		}
		// put parsed result back in source
		source = strings.Join(sections, "")
	}
	// convert source to bytes[]
	byteSource := []byte(source)
	// try to save source to outputDir
	err = ioutil.WriteFile(filepath.Clean(config.OutputDirAbs+"/"+file.PathRel), byteSource, 0777)
	// done if no error
	if err == nil {
		return logSummary
	}
	// otherwise, check if directories need to be created first
	dirRel := filepath.Dir(file.PathRel)
	os.MkdirAll(config.OutputDir+"/"+dirRel, 0777)
	// try to save source to outputDir
	err = ioutil.WriteFile(filepath.Clean(config.OutputDirAbs+"/"+file.PathRel), byteSource, 0777)
	// done if no error
	if err == nil {
		return logSummary
	}
	// log error and move on
	logSummary = logSummary + "output file " + file.Path + "could not be written to: " + err.Error() + "\n\n"
	// return logSummary
	return logSummary
}

func CopyFile(config *Config, file *File) {
	source, err := os.Open(file.Path)
	if err != nil {
		// note error and move on
		fmt.Printf("the input file %q could not be opened.\n", file.Path)
		return
	}
	defer source.Close()
	dest, err := os.Create(filepath.Clean(config.OutputDirAbs + "/" + file.PathRel))
	if err != nil {
		// note error and move on
		fmt.Printf("the output file %q could not be created.\n", file.Path)
		return
	}
	if _, err := io.Copy(dest, source); err != nil {
		dest.Close()
		// note error and move on
		fmt.Printf("the output file %q could not be created.\n", file.Path)
	}
	dest.Close()
	return
}
