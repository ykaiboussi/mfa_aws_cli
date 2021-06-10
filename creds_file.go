package main

import (
	"io/ioutil"
	"os"
	"strings"
)

type Section struct {
	IDEnv     string
	AWSRegion string
	Output    string
	KeyEnv    string
	Token     string
}

func EditCredFile(filePath, profile, mfaPrf, IDEnv, KeyEnv, token string) (map[string]Section, error) {
	st := make(map[string]Section)
	body, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	b := string(body)
	c := strings.TrimSpace(b)
	sections := strings.Split(c, "\n\n") // A new line between aws profiles is required to parse the file correctly
	for _, s := range sections {
		var v Section
		sl := strings.Split(s, "\n")
		for _, sv := range sl {
			if strings.Contains(sv, "aws_access_key_id") {
				v.IDEnv = sv
			}
			if strings.Contains(sv, "aws_secret_access_key") {
				v.KeyEnv = sv
			}
			if strings.Contains(sv, "aws_session_token") {
				v.Token = sv
			}
			if strings.Contains(sv, "region") {
				v.AWSRegion = sv
			}
			if strings.Contains(sv, "output") {
				v.Output = sv
			}
		}
		st[sl[0]] = v
	}

	if err = os.Remove(filePath); err != nil {
		return nil, err
	}

	f, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	var e = make(map[string]Section)
	for key, value := range st {
		e[key] = value
		if strings.Contains(key, "["+profile+"]") {
			value.IDEnv = "aws_access_key_id=" + IDEnv
			value.KeyEnv = "aws_secret_access_key=" + KeyEnv
			value.Token = "aws_session_token=" + token
			e["["+mfaPrf+"]"] = value
		}
	}

	for key, value := range e {
		if strings.Contains(key, "["+mfaPrf+"]") {
			f.WriteString("\n")
			f.WriteString(key + "\n")
			f.WriteString(value.IDEnv + "\n")
			f.WriteString(value.KeyEnv + "\n")
			f.WriteString(value.Token + "\n")
			if value.AWSRegion != "" {
				f.WriteString(value.AWSRegion + "\n")
			}
			if value.Output != "" {
				f.WriteString(value.Output + "\n")
			}

		} else {
			f.WriteString("\n")
			f.WriteString(key + "\n")
			f.WriteString(value.IDEnv + "\n")
			f.WriteString(value.KeyEnv + "\n")
			if value.AWSRegion != "" {
				f.WriteString(value.AWSRegion + "\n")
			}
			if value.Output != "" {
				f.WriteString(value.Output + "\n")
			}
		}
	}
	return e, nil
}
