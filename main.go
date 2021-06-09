package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func main() {
	defaultF := flag.String("d", "", "primary profile")
	mfaPrl := flag.String("m", "", "MFA-enabled profile")
	flag.Parse()

	if *defaultF == "" || *mfaPrl == "" {
		flag.Usage()
		os.Exit(1)
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(aws.ToString(defaultF)))
	if err != nil {
		log.Fatal("failed to load SDK configuration: ", err)
	}

	iamClient := iam.NewFromConfig(cfg)

	devices, err := iamClient.ListMFADevices(context.TODO(), &iam.ListMFADevicesInput{})
	if err != nil {
		log.Fatal(err)
	}

	if len(devices.MFADevices) == 0 {
		log.Fatal("No MFA devices configured")
	}

	sn := devices.MFADevices[0].SerialNumber

	fmt.Printf("Using device %1s\n", *sn)

	stsClient := sts.NewFromConfig(cfg)

	fmt.Printf("Enter MFA code: ")
	r := bufio.NewReader(os.Stdin)
	code, _, err := r.ReadLine()
	if err != nil {
		log.Fatal(err)
	}

	codeStr := string(code)

	res, err := stsClient.GetSessionToken(context.TODO(), &sts.GetSessionTokenInput{
		// Set the expiration
		//DurationSeconds: aws.Int32(900),
		TokenCode:    &codeStr,
		SerialNumber: sn,
	})
	if err != nil {
		log.Fatal(err)
	}

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	filePath := usr.HomeDir + "/.aws/credentials"
	if _, err = EditCredFile(filePath, aws.ToString(defaultF), aws.ToString(mfaPrl), aws.ToString(res.Credentials.AccessKeyId), aws.ToString(res.Credentials.SecretAccessKey), aws.ToString(res.Credentials.SessionToken)); err != nil {
		log.Fatal(err)
	}
}
