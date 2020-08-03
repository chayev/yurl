# yURL: Universal Links / AASA File Validator

yURL is a terminal utility that allows you to validate whether a URL is properly formatted for Universal Links. This allows you to check if the apple-app-site-association (AASA) file exists and is in the proper configuration as [defined by Apple](https://developer.apple.com/documentation/safariservices/supporting_associated_domains).

## macOS Install Instructions

### Install with Brew

Install yURL with [Brew](https://brew.sh/):

```
brew install chayev/tap/yurl
```

### Install using cURL 

Run the below command. Note that you will be prompted to enter your password.

```
curl -sSL "https://github.com/chayev/yurl/releases/download/v0.1.0/yurl-v0.1.0-macOS-amd64.tar.gz" | sudo tar -xz -C /usr/local/bin yurl
```

Note: that you will be prompted to enter your password.

## Usage

Run `yurl help` for information on how to use yURL.

## License

This repository is licensed under the MIT license.
The license can be found [here](./LICENSE).
