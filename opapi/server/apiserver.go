package server

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"

	"github.com/daveontour/opapi/opapi/globals"
	"github.com/daveontour/opapi/opapi/models"
	"github.com/daveontour/opapi/opapi/repo"

	"fmt"
	"io"

	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func StartGinServer(demoMode bool) {

	mode := gin.ReleaseMode
	if globals.ConfigViper.GetBool("DebugService") {
		mode = gin.DebugMode
	}
	gin.SetMode(mode)

	router := gin.New()

	// Configure all the endpoints for the HTTP Server

	// Test purposes only to just printout whatever was received by the server
	if globals.ConfigViper.GetBool("TestHTTPServer") {
		router.POST("/test", func(c *gin.Context) {
			if globals.Logger.Level == logrus.TraceLevel {
				globals.Logger.Info("Received message on test HTTP Server")
				jsonData, _ := io.ReadAll(c.Request.Body)
				fmt.Println(string(jsonData[:]))
			} else {
				globals.Logger.Info("Received message on test HTTP Server")
			}
		})
	}
	router.GET("/getFlights/:apt", repo.GetRequestedFlightsAPI)
	router.GET("/getAllocations/:apt", repo.GetResourceAPI)
	router.GET("/getConfiguredResources/:apt/:resourceType", repo.GetConfiguredResources)
	router.GET("/getConfiguredResources/:apt", repo.GetConfiguredResources)

	router.GET("/admin/repoMetricsReport/:apt", metricsReport)
	router.GET("/admin/repoMetricsReportNow/:apt", metricsReportNow)
	router.GET("/admin/enableMetrics", func(c *gin.Context) {
		if hasAdminToken(c) {
			globals.MetricsLogger.SetLevel(logrus.InfoLevel)
			globals.MetricsLogger.Info("Performance Metrics Reporting Enabled")
			c.JSON(http.StatusOK, gin.H{"PerformanceMetricsReporting": "Enabled"})
		} else {
			c.JSON(http.StatusForbidden, gin.H{"Error": "Not Authorized"})
		}
	})
	router.GET("/admin/disableMetrics", func(c *gin.Context) {
		if hasAdminToken(c) {
			globals.MetricsLogger.Info("Performance Metrics Reporting Disabled")
			globals.MetricsLogger.SetLevel(logrus.ErrorLevel)
			c.JSON(http.StatusOK, gin.H{"PerformanceMetricsReporting": "Disabledd"})
		} else {
			globals.MetricsLogger.Info("Performance Metrics Enabled")
		}
	})
	router.GET("/help", func(c *gin.Context) {
		data, err := os.ReadFile("./help.html")
		if err != nil {
			return
		}
		c.Header("Content-Type", "text/html")
		_, _ = c.Writer.Write(data)
	})
	router.GET("/adminhelp", func(c *gin.Context) {
		data, err := os.ReadFile("./adminhelp.htm")
		if err != nil {
			return
		}
		c.Header("Content-Type", "text/html")
		_, _ = c.Writer.Write(data)
	})

	if demoMode {
		//These endpoints are only started in Demo Mode and allow the service to be populated with data from the opapiseeder program
		router.POST("/demoMessageAppend", func(c *gin.Context) {
			xmlData, err := io.ReadAll(c.Request.Body)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"PostDemoMessageError": fmt.Errorf(err.Error())})
			}
			fmt.Println("Demo flight append message received")
			repo.UpdateFlightEntry(string(xmlData), true, true)
		})
		router.POST("/demoMessageUpdate", func(c *gin.Context) {
			xmlData, err := io.ReadAll(c.Request.Body)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"PostDemoMessageError": fmt.Errorf(err.Error())})
			}
			fmt.Println("Demo flight uppdate message received")
			repo.UpdateFlightEntry(string(xmlData), false, true)
		})

		router.GET("/demoPrepare", func(c *gin.Context) {
			repo.PerfTestInit()
		})

	}

	// Start it up with the configured security mode
	if !globals.ConfigViper.GetBool("UseHTTPS") && !globals.ConfigViper.GetBool("UseHTTPSUntrusted") {

		// Plain old HTTP
		err := router.Run(globals.ConfigViper.GetString("ServiceIPPort"))
		if err != nil {
			globals.Logger.Fatal("Unable to start HTTP server.")
			globals.Wg.Done()
			os.Exit(2)
		}

	} else if globals.ConfigViper.GetBool("UseHTTPS") && globals.ConfigViper.GetString("KeyFile") != "" && globals.ConfigViper.GetString("CertFile") != "" {

		// HTTPS with a supplied Certificate file and Key File
		server := http.Server{Addr: globals.ConfigViper.GetString("ServiceIPPort"), Handler: router}
		err := server.ListenAndServeTLS(globals.ConfigViper.GetString("CertFile"), globals.ConfigViper.GetString("KeyFile"))
		if err != nil {
			globals.Logger.Fatal("Unable to start HTTPS server. Likely cause is that the keyFile or certFile were not found")
			globals.Wg.Done()
			os.Exit(2)
		}

	} else if globals.ConfigViper.GetBool("UseHTTPS") && (globals.ConfigViper.GetString("KeyFile") == "" && globals.ConfigViper.GetString("CertFile") == "") {

		// HTTPS was configured, but the certificate and key file were not configured
		globals.Logger.Fatal("Unable to start HTTPS server. Trusted HTTPS was configured but The keyFile or certFile were not configured")
		globals.Wg.Done()
		os.Exit(2)

	} else if globals.ConfigViper.GetBool("UseHTTPSUntruste") {

		//Use HTTPS with a dodgy local certificate
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
		server := http.Server{Addr: globals.ConfigViper.GetString("ServiceIPPort"), Handler: router, TLSConfig: tlsConfig}

		err := server.ListenAndServeTLS("", "")
		if err != nil {
			globals.Logger.Fatal("Unable to start HTTPS server with local certificates and key")
			globals.Wg.Done()
			os.Exit(2)
		}
	}

}

func hasAdminToken(c *gin.Context) bool {
	// See if there is a valid admin token supplied in the header of the request
	keys := c.Request.Header["Token"]
	if keys == nil {
		return false
	}
	if keys[0] == globals.ConfigViper.GetString("AdminToken") {
		return true
	} else {
		return false
	}
}

// func reinit(c *gin.Context) {

// 	if !hasAdminToken(c) {
// 		c.JSON(http.StatusForbidden, gin.H{"Error": fmt.Sprintf("Not Authorized")})
// 		return
// 	} else {
// 		globals.RequestLogger.Info(fmt.Sprintf("User: %s IP: %s Request:%s", "admin", c.RemoteIP(), c.Request.RequestURI))
// 	}

// 	apt := c.Param("apt")
// 	repo.ReInitAirport(apt)
// }

func metricsReport(c *gin.Context) {
	// Get the profile of the user making the request

	// if !hasAdminToken(c) {
	// 	c.JSON(http.StatusForbidden, gin.H{"Error": fmt.Sprintf("Not Authorized")})
	// 	return
	// } else {
	globals.RequestLogger.Info(fmt.Sprintf("User: %s IP: %s Request:%s", "admin", c.RemoteIP(), c.Request.RequestURI))
	// }

	apt := c.Param("apt")

	metrics := models.MetricsReport{}
	metrics.Airport = apt

	repo := repo.GetRepo(apt)

	if repo == nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Airport not found"})
		return
	}

	metrics.NumberOfFlights = repo.FlightLinkedList.Len()
	metrics.NumberOfCheckins = repo.CheckInList.Len()

	metrics.NumberOfGates = repo.GateList.Len()
	metrics.NumberOfStands = repo.StandList.Len()
	metrics.NumberOfCarousels = repo.CarouselList.Len()
	metrics.NumberOfChutes = repo.ChuteList.Len()

	metrics.TotalNumberOfCheckinAllocations = repo.CheckInList.NumberOfFlightAllocations()
	metrics.TotalNumberOfStandAllocations = repo.StandList.NumberOfFlightAllocations()
	metrics.TotalNumberOfGateAllocations = repo.GateList.NumberOfFlightAllocations()
	metrics.TotalNumberOfCarouselAllocations = repo.CarouselList.NumberOfFlightAllocations()
	metrics.TotalNumberOfChuteAllocations = repo.ChuteList.NumberOfFlightAllocations()

	metrics.CheckInAllocationMetrics = repo.CheckInList.AllocationsMetrics()
	metrics.GateAllocationMetrics = repo.GateList.AllocationsMetrics()
	metrics.StandAllocationMetrics = repo.StandList.AllocationsMetrics()
	metrics.CarouselAllocationMetrics = repo.CarouselList.AllocationsMetrics()
	metrics.ChuteAllocationMetrics = repo.ChuteList.AllocationsMetrics()

	// var m runtime.MemStats
	// runtime.ReadMemStats(&m)
	// metrics.MemAllocMB = int(m.Alloc / 1024 / 1024)
	// metrics.MemSysMB = int(m.Sys / 1024 / 1024)
	// metrics.MemTotaAllocMB = int(m.TotalAlloc / 1024 / 1024)
	// metrics.MemHeapAllocMB = int(m.HeapAlloc / 1024 / 1024)
	// metrics.MemNumGC = int(m.NumGC)

	c.JSON(http.StatusOK, gin.H{"RepositoryMetrics": metrics})

}

func metricsReportNow(c *gin.Context) {
	// Get the profile of the user making the request

	// if !hasAdminToken(c) {
	// 	c.JSON(http.StatusForbidden, gin.H{"Error": fmt.Sprintf("Not Authorized")})
	// 	return
	// } else {
	globals.RequestLogger.Info(fmt.Sprintf("User: %s IP: %s Request:%s", "admin", c.RemoteIP(), c.Request.RequestURI))
	// }

	apt := c.Param("apt")

	metrics := models.MetricsReportNow{}
	metrics.Airport = apt

	repo := repo.GetRepo(apt)

	if repo == nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Airport not found"})
		return
	}

	metrics.NumberOfFlights = repo.FlightLinkedList.Len()
	metrics.NumberOfCheckins = repo.CheckInList.Len()

	metrics.NumberOfGates = repo.GateList.Len()
	metrics.NumberOfStands = repo.StandList.Len()
	metrics.NumberOfCarousels = repo.CarouselList.Len()
	metrics.NumberOfChutes = repo.ChuteList.Len()

	metrics.TotalNumberOfCheckinAllocationsNow = repo.CheckInList.NumberOfFlightAllocationsNow()
	metrics.TotalNumberOfStandAllocationsNow = repo.StandList.NumberOfFlightAllocationsNow()
	metrics.TotalNumberOfGateAllocationsNow = repo.GateList.NumberOfFlightAllocationsNow()
	metrics.TotalNumberOfCarouselAllocationsNow = repo.CarouselList.NumberOfFlightAllocationsNow()
	metrics.TotalNumberOfChuteAllocationsNow = repo.ChuteList.NumberOfFlightAllocationsNow()

	metrics.CheckInAllocationMetricsNow = repo.CheckInList.AllocationsMetricsNow()
	metrics.GateAllocationMetricsNow = repo.GateList.AllocationsMetricsNow()
	metrics.StandAllocationMetricsNow = repo.StandList.AllocationsMetricsNow()
	metrics.CarouselAllocationMetricsNow = repo.CarouselList.AllocationsMetricsNow()
	metrics.ChuteAllocationMetricsNow = repo.ChuteList.AllocationsMetricsNow()

	c.JSON(http.StatusOK, gin.H{"RepositoryMetrics": metrics})

}
