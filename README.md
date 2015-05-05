# Polyglot

Polyglot is a command line tool that automatically parses embedded DSLs that have been implemented with Dialect, a recursive-descent parser for Domain Specific Languages (DSLs) that is implemented using Go and facilitates parsing through use of Parsing Expression Grammars (PEGs).

## Motivation

Providing the ability to embed multiple DSLs in files containing other code allows you to use the best tool for the job. When you can benefit from the clarity and brevity of a DSL, using Polyglot you can embed the DSL(s) and use it alongside the more powerful and general programming facility.

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

Once parsed, Polyglot outputs polyglot-log.txt with information on each parsed file extension.
