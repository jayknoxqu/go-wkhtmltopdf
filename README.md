# wkhtmltopdf as a web service

A docker web service written in [Go](https://golang.org/) that uses [wkhtmltopdf](http://wkhtmltopdf.org/) to convert HTML into documents (images or pdf files).

## Docker

### build

```
docker build -t jayknoxqu/go-wkhtmltopdf:1.0v .
```

### run

```
docker run -d --name go-wkhtmltopdf -p 9900:80 jayknoxqu/go-wkhtmltopdf:1.0v
```



## Usage

The service listens on port 80 for GET requests on the root path (`/`). Any other method returns a `405 not allowed` status. Any other path returns a `404 not found` status.

The body should contain a JSON-encoded object containing the following parameters:

- **url**: The URL of the page to convert.
- **cookies**: A list of key-value arguments that are passed on to the appropriate `wkhtmltopdf` binary as separate `cookie` arguments.
- **output**: The type of document to generate, can be either `jpg`, `png` or `pdf`. Defauts to `pdf` if not specified. Depending on the output type the appropriate binary is called.
- **options**: A list of key-value arguments that are passed on to the appropriate `wkhtmltopdf` binary. Boolean values are interpreted as flag arguments (e.g.: `--greyscale`).
- **fileName**: The name of the download file.

**Example:** get the following JSON as request param `param`:

```
{
  "url": "http://www.google.com",
  "options": {
    "margin-bottom": "1cm",
    "orientation": "Landscape",
    "disable-smart-shrinking":true
  },
  "cookies": {
    "foo": "bar",
    "baz": "foo"
  },
  "output":"pdf",
  "fileName":"test"
}
```
**the json param need url encode .**
```
http://localhost:9900?param=%7b%22url%22%3a%22http%3a%2f%2fwww.google.com%22%2c%22options%22%3a%7b%22margin-bottom%22%3a%221cm%22%2c%22orientation%22%3a%22Landscape%22%2c%22disable-smart-shrinking%22%3atrue%7d%2c%22cookies%22%3a%7b%22foo%22%3a%22bar%22%2c%22baz%22%3a%22foo%22%7d%2c%22output%22%3a%22pdf%22%2c%22fileName%22%3a%22test%22%7d
```


will have the effect of the following command-line being executed on the server:

```
/bin/wkhtmltopdf --margin-bottom 1cm --orientation Landscape --cookie foo bar --cookie baz foo http://www.google.com -
```

The `-` at the end of the command-line is so that the document contents are redirected to stdout so we can in turn redirect it to the web server's response stream.

When using `jpg` or `png` output, the set of options you can pass are actually more limited. Please check [wkhtmltopdf usage docs](http://wkhtmltopdf.org/docs.html) or rather just use `wkhtmltopdf --help` or `wkhtmltoimage --help` to get help on the available command-line arguments.



### multiple urls

- support for multiple urls combined in one PDF

**Example:** 

```
{
  "urls": {
    "http://www.google.com",
    "http://www.reddit.com",
  },
  "options": {
    "margin-bottom": "1cm",
    "orientation": "Landscape",
    "grayscale": true
  },
  "cookies": {
    "foo": "bar",
    "baz": "foo"
  }
}
```



## reference

https://github.com/Surnet/docker-wkhtmltopdf

https://github.com/mickaelperrin/docker-wkhtmltopdf-service
