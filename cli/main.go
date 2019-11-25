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
	enveloped := flag.Bool("e", false, "the document is enveloped")
	keyPath := flag.String("k", "", "the path to the key")
	certPath := flag.String("c", "", "the path to the certificate")

	// SignatureOptions represents additional, less commonly used, options for Sign and Verify.
	//
	// Specify the name of ID attributes for specific elements. This
	// may be required if the signed document contains Reference elements
	// that define which parts of the document are to be signed.
	//
	// https://www.aleksey.com/xmlsec/faq.html#section_3_2
	// http://www.w3.org/TR/xml-id/
	// http://xmlsoft.org/html/libxml-valid.html#xmlAddID
	//
	// XMLIDOption represents the definition of an XML reference element
	// (See http://www.w3.org/TR/xml-id/)
	SignatureOptions_XMLIDOption_ElementName := flag.String("element-name", "AuthnRequest", "the Reference element containing ID attribute")
	SignatureOptions_XMLIDOption_ElementNamespace := flag.String("element-ns", "urn:oasis:names:tc:SAML:2.0:protocol", "the namespace of the Reference element")
	SignatureOptions_XMLIDOption_AttributeName := flag.String("attribute-name", "ID", "the name of the ID attribute")

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

		opts := xmlsec.SignatureOptions{}
		if *enveloped {

			var params []xmlsec.XMLIDOption

			params = append(params, xmlsec.XMLIDOption{
				ElementName:      *SignatureOptions_XMLIDOption_ElementName,
				ElementNamespace: *SignatureOptions_XMLIDOption_ElementNamespace,
				AttributeName:    *SignatureOptions_XMLIDOption_AttributeName,
			})
			opts = xmlsec.SignatureOptions{XMLID: params}
		}

		signedBuf, err := xmlsec.Sign(key, buf, opts)
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

		opts := xmlsec.SignatureOptions{}
		if *enveloped {

			var params []xmlsec.XMLIDOption

			params = append(params, xmlsec.XMLIDOption{
				ElementName:      *SignatureOptions_XMLIDOption_ElementName,
				ElementNamespace: *SignatureOptions_XMLIDOption_ElementNamespace,
				AttributeName:    *SignatureOptions_XMLIDOption_AttributeName,
			})
			opts = xmlsec.SignatureOptions{XMLID: params}
		}

		err = xmlsec.Verify(cert, buf, opts)
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
