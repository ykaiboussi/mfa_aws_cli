# mfa_aws_cli_setup

mfa_aws_cli_setup is a lightweight utility tool that generates [AWS MFA profile](https://aws.amazon.com/premiumsupport/knowledge-center/authenticate-mfa-cli/) using `aws-sdk-go-v2` and no external library when parsing secrets.

Features:
1. Creates a new MFA profile of the primary.
2. Refreshes session token from an existing MFA profile.

## Install 
```
go get github.com/ykaiboussi/mfa_aws_cli_setup
```

## Usage
```
➜  mfa_aws_cli_setup git:(main) mfa_aws_cli_setup -h
Usage of mfa_aws_cli_setup:
  -d string
    	primary profile
  -m string
    	MFA-enabled profile

➜  mfa_aws_cli_setup git:(main) ✗  mfa_aws_cli_setup -d default -m default-mfa
```

## Note
In order to parse `$HOME/.aws/credentials` corretly. A new line is required betweeen multiple profiles. As an example `test_data/credentials`


## Credits 
[go-aws-mfa](https://github.com/jdevelop/go-aws-mfa)
