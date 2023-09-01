package cmd

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

const Layout = "2006-01-02T15:04:05"

var UseHTTPS bool
var trace bool
var debug bool
var log bool

var Loc *time.Location
var wg sync.WaitGroup

var rootCmd = &cobra.Command{
	Use:   "webhookstestclient",
	Short: `webhookstestclient is a CLI to run and manage the webhooks test client`,
	Long:  "webhookstestclient is a CLI to run and manage the webhooks test client",
}
var runCmd = &cobra.Command{
	Use:   "run {server:port}",
	Short: `Start the Webhooks test client`,
	Long:  `Start the Webhooks test client`,
	Run: func(cmds *cobra.Command, args []string) {
		if len(args) >= 2 {
			if args[1] == "trace" {
				trace = true
			}
			if args[1] == "debug" {
				debug = true
			}
			if args[1] == "log" {
				log = true
			}
		}
		StartGinServer(args)
	},
}

func InitCobraTestClient() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	runCmd.PersistentFlags().BoolVarP(&UseHTTPS, "https", "s", false, "Use HTTPS")
	rootCmd.AddCommand(runCmd)

	Loc, _ = time.LoadLocation("Local")

}
func ExecuteCobraTestClient() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func StartGinServer(args []string) {

	mode := gin.DebugMode
	gin.SetMode(mode)

	router := gin.New()

	router.POST("/test", testQuery)
	router.POST("/subscriptionEndPoint", subPush)
	router.POST("/changeEndpoint", changePush)
	router.GET("/file", fileTest)

	wg.Add(1)
	if !UseHTTPS {
		fmt.Printf("Started Webhooks Test Client on %s. Using HTTPS = %t\n", args[0], UseHTTPS)
		err := router.Run(args[0])
		if err != nil {
			os.Exit(2)
		}
	} else {
		cert := &x509.Certificate{
			SerialNumber: big.NewInt(1658),
			Subject: pkix.Name{
				Organization:  []string{"ORGANIZATION_NAME"},
				Country:       []string{"COUNTRY_CODE"},
				Province:      []string{"PROVINCE"},
				Locality:      []string{"CITY"},
				StreetAddress: []string{"ADDRESS"},
				PostalCode:    []string{"POSTAL_CODE"},
			},
			NotBefore:    time.Now(),
			NotAfter:     time.Now().AddDate(10, 0, 0),
			SubjectKeyId: []byte{1, 2, 3, 4, 6},
			ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
			KeyUsage:     x509.KeyUsageDigitalSignature,
		}
		priv, _ := rsa.GenerateKey(rand.Reader, 2048)
		pub := &priv.PublicKey

		// Sign the certificate
		certificate, _ := x509.CreateCertificate(rand.Reader, cert, cert, pub, priv)

		certBytes := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certificate})
		keyBytes := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

		// Generate a key pair from your pem-encoded cert and key ([]byte).
		x509Cert, _ := tls.X509KeyPair(certBytes, keyBytes)

		tlsConfig := &tls.Config{Certificates: []tls.Certificate{x509Cert}}
		server := http.Server{Addr: args[0], Handler: router, TLSConfig: tlsConfig}

		fmt.Printf("Started Webhooks Test Client on %s. Using HTTPS = %t\n", args[0], UseHTTPS)
		err := server.ListenAndServeTLS("", "")
		if err != nil {
			fmt.Println("Unable to start HTTPS server with local certificates and key")
			os.Exit(2)
		}

	}

	wg.Wait()

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func fileTest(c *gin.Context) {
	d1 := []byte("hello\ngo\n")
	err := os.WriteFile("dat1", d1, 0777)
	check(err)

	f, _ := os.OpenFile("dat1", os.O_RDONLY, 0755)
	fi, err := f.Stat()
	if err != nil {
		// Could not obtain stat, handle error
	}
	c.DataFromReader(200, fi.Size(), "application/text", f, nil)

	defer func() {
		err = f.Close()
		if err != nil {
			fmt.Println(err.Error())
		} else {
			err := os.Remove("dat1")
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}()

}
func testQuery(c *gin.Context) {

	jsonData, _ := io.ReadAll(c.Request.Body)
	fmt.Println(string(jsonData[:]))

}
func changePush(c *gin.Context) {

	path := "./WebhookLogs"
	os.MkdirAll(path, os.ModePerm)

	jsonData, _ := io.ReadAll(c.Request.Body)
	if debug {
		fmt.Println(string(jsonData[:500]))
	} else if trace {
		fmt.Println(string(jsonData))
	}
	if log {
		file, errs := os.CreateTemp("./WebhookLogs", "changelog-*.json")
		if errs != nil {
			fmt.Println(errs)
			return
		}
		file.WriteString(string(jsonData))
		file.Close()
	}
	fmt.Printf("Change message received at %s\n", time.Now().Format(Layout))

}
func subPush(c *gin.Context) {

	path := "./WebhookLogs"
	os.MkdirAll(path, os.ModePerm)

	jsonData, _ := io.ReadAll(c.Request.Body)
	if debug {
		fmt.Println(string(jsonData[:500]))
	} else if trace {
		fmt.Println(string(jsonData))
	}
	if log {
		file, errs := os.CreateTemp("./WebhookLogs", "pushlog-*.json")
		if errs != nil {
			fmt.Println(errs)
			return
		}
		file.WriteString(string(jsonData))
		file.Close()
	}
	fmt.Printf("Subcription message received at %s\n", time.Now().Format(Layout))

}
