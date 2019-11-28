package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/SimonBaeumer/cmd"
	"github.com/stretchr/testify/require"
)

var (
	// The path to the statically compiled binary (CLI command) based on github.com/crewjam/go-xmlsec
	// See https://github.com/vstojic/cloudfoundry-cflinuxfs3-go-xmlsec
	xmldsigCmdPath, _ = filepath.Abs("./xmldsig")
)

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

	//cmd := cmd.NewCommand("cat " + signatureTemplateFile.Name() + " | " + XmldsigCmdPath + " -k " + privateKeyFile.Name() + " -s > " + signedFile.Name())
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
		return err
	}

	return nil
}

// Use the following command to generate the self-signed CSR:
// openssl req \
//   -newkey rsa:2048 -sha256 -nodes -keyout dip.key \
//   -x509 -days 365 -out dip.crt \
//   -subj '/CN=and/C=DE'
var (
	privateKeyPath, _        = filepath.Abs("dip.key")
	certPath, _              = filepath.Abs("dip.crt")
	signatureTemplatePath, _ = filepath.Abs("crs_payload.xml")
)

func Test(t *testing.T) {

	privateKey, err := ioutil.ReadFile(privateKeyPath)
	require.NoError(t, err)
	require.NotEmpty(t, privateKey)

	cert, err := ioutil.ReadFile(certPath)
	require.NoError(t, err)
	require.NotEmpty(t, cert)

	signatureTemplate, err := ioutil.ReadFile(signatureTemplatePath)
	require.NoError(t, err)
	require.NotEmpty(t, signatureTemplate)

	signed, err := SignXML(signatureTemplate, string(privateKey))
	require.NoError(t, err)
	require.NotEmpty(t, signed)

	err = ValidateXMLSignature(signed, cert)
	require.NoError(t, err)
}

/*
func main() {

	privateKey, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	cert, err := ioutil.ReadFile(certPath)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	signatureTemplate, err := ioutil.ReadFile(signatureTemplatePath)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	signed, err := SignXML(signatureTemplate, string(privateKey))
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	//println("%s\n", string(signed))

	err = ValidateXMLSignature(signed, cert)
	if err != nil {
		println("%s\n", err)
		os.Exit(1)
	}

	//println("The signature is correct")
}
*/
