# sigbypass4xx

[![release](https://img.shields.io/github/release/signedsecurity/sigbypass4xx?style=flat&color=0040ff)](https://github.com/signedsecurity/sigbypass4xx/releases) ![maintenance](https://img.shields.io/badge/maintained%3F-yes-0040ff.svg) [![open issues](https://img.shields.io/github/issues-raw/signedsecurity/sigbypass4xx.svg?style=flat&color=0040ff)](https://github.com/signedsecurity/sigbypass4xx/issues?q=is:issue+is:open) [![closed issues](https://img.shields.io/github/issues-closed-raw/signedsecurity/sigbypass4xx.svg?style=flat&color=0040ff)](https://github.com/signedsecurity/sigbypass4xx/issues?q=is:issue+is:closed) [![license](https://img.shields.io/badge/license-MIT-gray.svg?colorB=0040FF)](https://github.com/signedsecurity/sigbypass4xx/blob/master/LICENSE) [![twitter](https://img.shields.io/badge/twitter-@signedsecurity-0040ff.svg)](https://twitter.com/signedsecurity)

sigbypass4xx is a utility to automate well-know techniques used to bypass access control restrictions.

**This project has been merged into [sigurlscann3r](https://github.com/signedsecurity/sigurlscann3r), therefore, this repository won't be maintained**

## Resources

* [Usage](#usage)
* [Installation](#installation)
    * [From Binary](#from-binary)
    * [From source](#from-source)
    * [From github](#from-github)
* [Contribution](#contribution)

## Usage

To display help message for sigbypass4xx use the `-h` flag:

```bash
sigbypass4xx -h
```

```
     _       _                               _  _             
 ___(_) __ _| |__  _   _ _ __   __ _ ___ ___| || |__  ____  __
/ __| |/ _` | '_ \| | | | '_ \ / _` / __/ __| || |\ \/ /\ \/ /
\__ \ | (_| | |_) | |_| | |_) | (_| \__ \__ \__   _>  <  >  < 
|___/_|\__, |_.__/ \__, | .__/ \__,_|___/___/  |_|/_/\_\/_/\_\ v1.0.0
       |___/       |___/|_| 

USAGE:
  bypass403 [OPTIONS]

OPTIONS:
  -c         concurrency level (default: 20)
  -delay     delay between requests (default: 100ms)
  -iL        urls with 403 to bypass (use `iL -` to read from stdin)
  -nC        no color mode


```

## Installation

#### From Binary

You can download the pre-built binary for your platform from this repository's [releases](https://github.com/signedsecurity/sigbypass4xx/releases/) page, extract, then move it to your `$PATH`and you're ready to go.

#### From Source

sigbypass4xx requires **go1.14+** to install successfully. Run the following command to get the repo

```bash
GO111MODULE=on go get -u -v github.com/signedsecurity/sigbypass4xx/cmd/sigbypass4xx
```

#### From Github

```bash
git clone https://github.com/signedsecurity/sigbypass4xx.git; cd sigbypass4xx/cmd/sigbypass4xx/; go build; mv sigbypass4xx /usr/local/bin/; sigbypass4xx -h
```

## Contribution

[Issues](https://github.com/signedsecurity/sigbypass4xx/issues) and [Pull Requests](https://github.com/signedsecurity/sigbypass4xx/pulls) are welcome!
