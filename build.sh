
rm -r -f $GOPATH/bin/linux_amd64
rm -r -f $GOPATH/bin/sms*
rm -r -f $GOPATH/bin/windows_amd64

GOOS=windows GOARCH=amd64 go install
GOOS=linux GOARCH=amd64 go install
GOOS=darwin GOARCH=amd64 go install
