# yURL: Universal Links / AASA File Validator
[![CircleCI Build Status](https://circleci.com/gh/chayev/yurl.svg?style=shield)](https://circleci.com/gh/chayev/yurl) [![GitHub License](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/chayev/yurl/master/LICENSE)

yURL is a CLI (Command-Line Interface) and [webapp](https://yurl.chayev.com/) that allows you to validate whether a URL is properly configured for Universal Links. This allows you to check if the apple-app-site-association (AASA) file exists and is in the proper format as [defined by Apple](https://developer.apple.com/documentation/safariservices/supporting_associated_domains).

## macOS Install Instructions

### Install with Brew (recommended)

Install yURL with [Brew](https://brew.sh/):

```
brew install chayev/tap/yurl
```

### Install using cURL 

Run the below command:

```
curl -sSL "https://github.com/chayev/yurl/releases/download/v0.3.2/yurl-v0.3.2-macOS-amd64.tar.gz" | sudo tar -xz -C /usr/local/bin yurl
```

Note: You will be prompted to enter your password because of the `sudo` command. `0.3.2` may need to be replaced with your desired version. 

## Linux Install Instructions 

### Install with Snap (recommended)

Install yURL with [Snap](https://snapcraft.io/):

```
sudo snap install yurl
```

### Install using cURL 

Run the below command:

```
curl -sSL "https://github.com/chayev/yurl/releases/download/v0.3.2/yurl-v0.3.2-linux-amd64.tar.gz" | sudo tar -xz -C /usr/local/bin yurl
```

Note: You will be prompted to enter your password because of the `sudo` command. `0.3.2` may need to be replaced with your desired version.

## Windows Install Instructions

You could download the executable from the [releases](https://github.com/chayev/yurl/releases) page. More instructions coming soon! 

We are planning on supporting [chocolatey](chocolatey.org) package manager as well.

## Usage and Example

Run `yurl help` for information on how to use yURL.

Example:

```bash
yurl aasa validate "https://www.google.com/search?q=gothamhq"
```

## Contributing

Contributions to yURL of any kind are welcome! Feel free to open [PRs](https://github.com/chayev/yurl/pulls) or an [issue](https://github.com/chayev/yurl/pulls). 

### Asking Support Questions

Feel free to open an issue if you have a question. 

### Reporting Issues

If you believe you have found a defect in yURL or its documentation, create an issue to report the problem.
When reporting the issue, please provide the version of yURL in use (`yurl version`).

## License

This repository is licensed under the MIT license.
The license can be found [here](./LICENSE).
