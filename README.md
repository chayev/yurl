# yURL: Universal Links / AASA File Validator

yURL is a terminal utility that allows you to validate whether a URL is properly formatted for Universal Links. This allows you to check if the apple-app-site-association (AASA) file exists and is in the proper configuration as [defined by Apple](https://developer.apple.com/documentation/safariservices/supporting_associated_domains).

## macOS Install Instructions

### Install with Brew

Install yURL with [Brew](https://brew.sh/):

```
brew install chayev/tap/yurl
```

### Install using cURL 

Run the below command:

```
curl -sSL "https://github.com/chayev/yurl/releases/download/v0.1.0/yurl-v0.1.0-macOS-amd64.tar.gz" | sudo tar -xz -C /usr/local/bin yurl
```

Note: you will be prompted to enter your password because of the `sudo` command.

## Usage

Run `yurl help` for information on how to use yURL.

## Contributing

Contributions to yURL of any kind are welcome! Feel free to open PRs, log issues, open feature requests, etc. 

### Asking Support Questions

Feel free to open an issue if you have a question. 

### Reporting Issues

If you believe you have found a defect in yURL or its documentation, create an issue to report the problem.
When reporting the issue, please provide the version of yURL in use (`yurl --version`).

## License

This repository is licensed under the MIT license.
The license can be found [here](./LICENSE).
