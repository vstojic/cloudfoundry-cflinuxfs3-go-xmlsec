package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/crewjam/go-xmlsec"
)

// cat crs_payload.xml | ./xmldsig -k dip.key -s | ./xmldsig -c dip.crt -v
func main() {

	doVerify := flag.Bool("v", false, "verify the document")
	doSign := flag.Bool("s", false, "sign the document")
	keyPath := flag.String("k", "", "the path to the key")
	certPath := flag.String("c", "", "the path to the certificate")
	flag.Parse()

	if !*doVerify && !*doSign {
		fmt.Println("Please, specify -v to verify or -s to sign")
		os.Exit(1)
	}

	buf, _ := ioutil.ReadAll(os.Stdin)

	if *doSign {

		if *keyPath == "" {
			fmt.Println("Please, specify a private key file (PEM format)")
			os.Exit(1)
		}

		key, err := ioutil.ReadFile(*keyPath)
		if err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}

		signedBuf, err := xmlsec.Sign(key, buf, xmlsec.SignatureOptions{})
		if err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}
		os.Stdout.Write(signedBuf)
	}

	if *doVerify {

		if *certPath == "" {
			fmt.Println("Please, specify a certification file (PEM format)")
			os.Exit(1)
		}

		cert, err := ioutil.ReadFile(*certPath)
		if err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}

		err = xmlsec.Verify(cert, buf, xmlsec.SignatureOptions{})
		if err == xmlsec.ErrVerificationFailed {
			fmt.Println("The signature is NOT correct")
			os.Exit(1)
		}
		if err != nil {
			fmt.Printf("error: %s\n", err)
			os.Exit(1)
		}
		fmt.Println("The signature is correct")
	}
}
