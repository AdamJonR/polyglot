# Polyglot

Polyglot is a command line tool that automatically parses embedded DSLs that have been implemented with Dialect, a recursive-descent parser for Domain Specific Languages (DSLs) that is implemented using Go and facilitates parsing through use of Parsing Expression Grammars (PEGs).

## Motivation

Providing the ability to embed multiple DSLs in files containing other code allows you to use the best tool for the job. When you can benefit from the clarity and brevity of a DSL, using Polyglot you can embed the DSL(s) and use it alongside the more powerful and general programming facility.

## Supported DSLs

Currently, Polyglot only has one DSL implemented, but more are on the way.

- (QForm: a DSL for creating HTML5 forms)[https://github.com/AdamJonR/qform]
- H5: a DSL for creating HTML5 (coming soon)
- EZBC: a DSL for performing bcmath calculations in PHP (coming soon)

## Using Polyglot

You can automatically parse files that contain DSL source code by calling Polyglot and passing in the path to the config.json file. An example config file is below. In the file, the input directory (the directory that contains files to by parsed) and the output directory (the directory that contains files having been parsed.) The extensions object contains an array of parsing implementations by file extension.

In the examples below, we'll configure Polyglot to search for and parse QForm DSL blocks in files with the html extension. Once completed, polyglot-log.txt contains detailed parsing info for each file that contained a DSL block to parse.

### Example config.json

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

### Example HTML File in Input Directory

```
<!DOCTYPE html>
<html>
<head>
  <title>Example page with embedded DSL</title>
</head>
<body>
  <h1>Example page with embedded DSL</h1>
<!--qform:o-->
- method post

text
- name name
- maxlength 30
- required

email
- name email

textarea
- name my-message

submit
- value Send message
<!--qform:c-->
</body>
<html>
```

### Example HTML File in Output Directory After Parsing

```
<!DOCTYPE html>
<html>
<head>
  <title>Example page with embedded DSL</title>
</head>
<body>
  <h1>Example page with embedded DSL</h1>
<form method="post">
  <div class="form-group">
    <label for="name">Name</label>
    <input type="text" name="name" maxlength="30" required="required" id="name" />
  </div>
  <div class="form-group">
    <label for="email">Email</label>
    <input type="email" name="email" id="email" />
  </div>
  <div class="form-group">
    <label for="my-message">My-message</label>
    <textarea name="my-message" id="my-message"></textarea>
  </div>
  <div class="form-group">
    <input type="submit" name="field4" id="field4" value="Send message" />
  </div>
</form>
</body>
<html>
```

### Example polyglot-log.txt

```
File: /Users/adam/Desktop/testInput/form-example.html
form attribute*, form field*
| hyphen, name, value?, newline
| found
| hyphen, name, value?, newline
| missing hyphen on line 2
| newline?, field type, field attribute*
| | field name, newline
| | found
| | hyphen, name, value?, newline?
| | found
| | hyphen, name, value?, newline?
| | found
| | hyphen, name, value?, newline?
| | found
| | hyphen, name, value?, newline?
| | missing hyphen on line 7
| | hyphen, array, newline?
| | missing hyphen on line 7
| found
| newline?, field type, field attribute*
| | field name, newline
| | found
| | hyphen, name, value?, newline?
| | found
| | hyphen, name, value?, newline?
| | missing hyphen on line 10
| | hyphen, array, newline?
| | missing hyphen on line 10
| found
| newline?, field type, field attribute*
| | field name, newline
| | found
| | hyphen, name, value?, newline?
| | found
| | hyphen, name, value?, newline?
| | missing hyphen on line 13
| | hyphen, array, newline?
| | missing hyphen on line 13
| found
| newline?, field type, field attribute*
| | field name, newline
| | found
| | hyphen, name, value?, newline?
| | found
| found
found
```
