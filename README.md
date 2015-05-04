# Polyglot

Polyglot is a command line tool that automatically parses DSLs that have been implemented with Dialect, a recursive-descent parser for Domain Specific Languages (DSLs) that is implemented using Go and facilitates parsing through use of Parsing Expression Grammars (PEGs).

## Using Polyglot

You can automatically parse files that contain DSL source code by calling Polyglot and passing in the path to the config.json file. An example config file is below. In the file, the input directory (the directory that contains files to by parsed) and the output directory (the directory that contains files having been parsed.) The extensions object contains an array of parsing implementations by file extension.

```
{
  "inputDir":"/Users/adam/Desktop/testInput/",
  "outputDir":"/Users/adam/Desktop/testOutput/",
  "extensions":{
    ".html":[
      {
        "dialect":"qform",
        "start":"<!--qform:o-->\n",
        "stop":"<!--qform:c-->\n"
      }
    ]
  }
}
```
