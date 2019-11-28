package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/SimonBaeumer/cmd"
)

// Use the following command to generate the self-signed CSR:
// openssl req \
//   -newkey rsa:2048 -sha256 -nodes -keyout dip.key \
//   -x509 -days 365 -out dip.crt \
//   -subj '/CN=and/C=DE'
var (
	// The path to the statically compiled binary (CLI command) based on github.com/crewjam/go-xmlsec
	// See https://github.com/vstojic/cloudfoundry-cflinuxfs3-go-xmlsec
	xmldsigCmdPath, _        = filepath.Abs("./xmldsig")
	privateKeyPath, _        = filepath.Abs("./cli/dip.key")
	certPath, _              = filepath.Abs("./cli/dip.crt")
	signatureTemplatePath, _ = filepath.Abs("./cli/crs_payload.xml")
	authnRequestPath, _      = filepath.Abs("./cli/AuthnRequest.xml")
)

func main() {
	http.HandleFunc("/", handler)
	//http.HandleFunc("/enveloped", handlerEnv)

	var port string
	if port = os.Getenv("PORT"); len(port) == 0 {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {

	privateKey, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		fmt.Fprintf(w, "%s\n", err)
		return
	}

	cert, err := ioutil.ReadFile(certPath)
	if err != nil {
		fmt.Fprintf(w, "%s\n", err)
		return
	}

	signatureTemplate, err := ioutil.ReadFile(signatureTemplatePath)
	if err != nil {
		fmt.Fprintf(w, "%s\n", err)
		return
	}

	signed, err := SignXML(signatureTemplate, string(privateKey))
	if err != nil {
		fmt.Fprintf(w, "%s\n", err)
		return
	}

	//println("%s\n", string(signed))

	err = ValidateXMLSignature(signed, cert)
	if err != nil {
		fmt.Fprintf(w, "%s\n", err)
		return
	}

	//println("The signature is correct")
	fmt.Fprintf(w, "The signature is correct\n")
}

func SignXML(signatureTemplate []byte, privateKey string) ([]byte, error) {
	signatureTemplateFile, err := ioutil.TempFile(os.TempDir(), "")
	if err != nil {
		return nil, err
	}
	defer os.Remove(signatureTemplateFile.Name())
	_, err = signatureTemplateFile.Write(signatureTemplate)
	if err != nil {
		return nil, err
	}

	privateKeyFile, err := ioutil.TempFile(os.TempDir(), "")
	if err != nil {
		return nil, err
	}
	defer os.Remove(privateKeyFile.Name())
	_, err = privateKeyFile.Write([]byte(privateKey))
	if err != nil {
		return nil, err
	}

	// Pass the adequate switches to the command (-k ... -s)
	cmd := cmd.NewCommand("cat " + signatureTemplateFile.Name() + " | " + xmldsigCmdPath + " -k " + privateKeyFile.Name() + " -s")
	err = cmd.Execute()
	if err != nil {
		return nil, err
	}

	return []byte(cmd.Stdout()), nil
}

func ValidateXMLSignature(message []byte, certificate []byte) error {

	messageFile, err := ioutil.TempFile(os.TempDir(), "")
	if err != nil {
		return err
	}
	_, err = messageFile.Write(message)
	if err != nil {
		return err
	}
	defer os.Remove(messageFile.Name())

	certificateFile, err := ioutil.TempFile(os.TempDir(), "")
	if err != nil {
		return err
	}
	_, err = certificateFile.Write([]byte(certificate))
	if err != nil {
		return err
	}
	defer os.Remove(certificateFile.Name())

	// Pass the adequate switches to the command (-c ... -v)
	cmd := cmd.NewCommand("cat " + messageFile.Name() + " | " + xmldsigCmdPath + " -c " + certificateFile.Name() + " -v")
	err = cmd.Execute()
	if err != nil {
		return nil
	}

	return nil
}
