# pup

`pup` is a command line tool for processing HTML. It reads from stdin,
prints to stdout, and allows the user to filter parts ot the page using
[CCS selectors](http://www.w3schools.com/cssref/css_selectors.asp).

Inspired by [`jq`](http://stedolan.github.io/jq/), `pup` aims to be a
fast and flexible way of exploring HTML from the terminal.

## Install

	go get github.com/ericchiang/pup

## Basic Usage

```bash
$ cat index.html | pup [selectors and flags]
```

or

```bash
$ pup < index.html [selectors and flags]
```

## Examples

Download a webpage with `wget`.

```bash
$ wget http://en.wikipedia.org/wiki/Robots_exclusion_standard -O robots.html
```

####Clean and indent

By default `pup` will fill in missing tags and properly indent the page.

```bash
$ cat robots.html
# nasty looking HTML
$ cat robots.html | pup --color
# cleaned, indented, and colorful HTML
```

####Filter by tag
```bash
$ pup < robots.html title
<title>
 Robots exclusion standard - Wikipedia, the free encyclopedia
</title>
```

####Filter by id
```bash
$ pup < robots.html span#See_also
<span class="mw-headline" id="See_also">
 See also
</span>
```

####Chain selectors together

The following two commands are equivalent. (NOTE: pipes do not work with the
`--color` flag)

```bash
$ pup < robots.html table.navbox ul a | tail
```

```bash
$ pup < robots.html table.navbox | pup ul | pup a | tail
```

Both produce the ouput:

```bash
</a>
<a href="/wiki/Stop_words" title="Stop words">
 Stop words
</a>
<a href="/wiki/Poison_words" title="Poison words">
 Poison words
</a>
<a href="/wiki/Content_farm" title="Content farm">
 Content farm
</a>
```

####How many nodes are selected by a filter?

```bash
$ pup < robots.html a -n
283
```

####Limit print level

```bash
$ pup < robots.html table -l 2
<table class="metadata plainlinks ambox ambox-content" role="presentation">
 <tbody>
  ...
 </tbody>
</table>
<table style="background:#f9f9f9;font-size:85%;line-height:110%;max-width:175px;">
 <tbody>
  ...
 </tbody>
</table>
<table cellspacing="0" class="navbox" style="border-spacing:0;">
 <tbody>
  ...
 </tbody>
</table>
```

## Implemented Selectors

For further examples of these selectors head over to [w3schools](
http://www.w3schools.com/cssref/css_selectors.asp).

```bash
cat index.html | pup .class
# '#' indicates comments at the command line so you have to escape it
cat index.html | pup \#id
cat index.html | pup element
cat index.html | pup [attribute]
cat index.html | pup [attribute=value]
```

You can mix and match selectors as you wish.

```bash
cat index.html | pup element#id[attribute=value]
```

## Functions

Non-HTML selectors which effect the output type are implemented as functions
which can be provided as a final argument.

As of now, `text{}` is the only implemented function.

#### `text{}`

Print all text from selected nodes and children in depth first order.

```bash
$ cat robots.html | pup .mw-headline text{}
History
About the standard
Disadvantages
Alternatives
Examples
Nonstandard extensions
Crawl-delay directive
Allow directive
Sitemap
Host
Universal "*" match
Meta tags and headers
See also
References
External links
```

## Flags

```bash
-c --color         print result with color
-f --file          file to read from
-h --help          display this help
-i --indent        number of spaces to use for indent or character
-n --number        print number of elements selected
-l --limit         restrict number of levels printed
--version          display version
```

## TODO:

* Print attribute function `attr{attr1, attr2}`
* Print as json function `json{}`
* Switch `-n` from a flag to a function
