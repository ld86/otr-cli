OTR with CLI
============

It's simple cli application for communication over OTR protocol.

Conversation procedure
===========================

Alice (as server)
	
	$ ./otr-cli -l -s SecretForOTR :1234
Bob (as client)

	$ ./otr-cli -s SecretForOTR 127.0.0.1:1234
	[!] Their fingerprint: A5FFC7DB70FE61C6CC3E2DB3BE99586DC4E49F64
	[!] Answer a question 'Do you know secret?'
	[!] Answer is correct
Alice 

	[!] Their fingerprint: 8F48A7F30ED884678D7EE93BE0CB941E1EC4C77B
	[!] Asking a question with answer: SecretForOTR
	[!] Answer is correct

Now you can chat with each other.

Building
=============

	sudo apt-get install git golang
	git clone https://github.com/ld86/otr-cli
	cd otr-cli
	export GOPATH=`pwd`
	go get
	go build
