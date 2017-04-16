/*
 * Copyright (c) 2017 Scott Brumley
 *
 *
 *
 * Contributors:
 *    Scott Brumley
 */


package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"time"
	"flag"
	"github.com/go-ini/ini"
	"strings"
	"os"
)

func NewTLSConfig(cfg *ini.File) *tls.Config {
	certCA := cfg.Section("Certs").Key("CertFile").String()
	clientCRT := cfg.Section("Certs").Key("CertFile").String()
	clientKey := cfg.Section("Certs").Key("PrivateKey").String()

	// Import trusted certificates from CAfile.pem.
	// Alternatively, manually add CA certificates to
	// default openssl CA bundle.
	certpool := x509.NewCertPool()
	pemCerts, err := ioutil.ReadFile(certCA)
	if err == nil {
		certpool.AppendCertsFromPEM(pemCerts)
	}

	// Import client certificate/key pair
	cert, err := tls.LoadX509KeyPair(clientCRT, clientKey)
	if err != nil {
		panic(err)
	}

	// Just to print out the client certificate..
	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		panic(err)
	}
	//fmt.Println(cert.Leaf)

	// Create tls.Config with desired tls properties
	return &tls.Config{
		// RootCAs = certs used to verify server cert.
		RootCAs: certpool,
		// ClientAuth = whether to request cert from server.
		// Since the server is set up for SSL, this happens
		// anyways.
		ClientAuth: tls.NoClientCert,
		// ClientCAs = certs used to validate client cert.
		ClientCAs: nil,
		// InsecureSkipVerify = verify that cert contents
		// match server. IP matches what is in cert etc.
		InsecureSkipVerify: true,
		// Certificates = list of certs client sends to server.
		Certificates: []tls.Certificate{cert},
	}
}

var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func connect(tlsconfig *tls.Config,cfg *ini.File){
	opts := MQTT.NewClientOptions()
	// Get Broker List
	dxlCfg := cfg.Section("Brokers").KeyStrings()

	for _,element := range dxlCfg {
		// index is the index where we are
		// element is the element from someSlice for where we are
		brokerValue := element + " = " + cfg.Section("Brokers").Key(element).String()

		// Split Broker Values
		brokerProperties := strings.Split(brokerValue,",")

		// Add all brokers to List of potential brokers
		if (len(brokerProperties) > 0){
			brokerStr := "ssl://" + brokerProperties[2] + ":" + brokerProperties[1]
			opts.AddBroker(brokerStr)
		} else {
			println("You need a broker list.")
		}

	}

	opts.SetClientID("dxlclient").SetTLSConfig(tlsconfig)
	opts.SetDefaultPublishHandler(f)

	// Start the connection
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if c.IsConnected() {
		println("Connected!")
	}

	if token := c.Subscribe("/mcafee/service/tie/file/reputation", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	i := 0
	for _ = range time.Tick(time.Duration(1) * time.Second) {
		if i == 5 {
			break
		}
		fmt.Printf("this is msg #%d!\n", i)
		//text := fmt.Sprintf("this is msg #%d!", i)
		//c.Publish("/go-mqtt/sample", 0, false, text)
		i++
	}

	if token := c.Unsubscribe("/mcafee/service/tie/file/reputation"); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	c.Disconnect(250)

	if ! c.IsConnected() {
		println("Disconnected!")
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func dxlConfigParser(configFile string)(cfg *ini.File){
	dxlTmpFile := "/tmp/dxlclient.config.tmp"

	// Read in dxlclient.config file
	contents,err := ioutil.ReadFile(configFile)
	check(err)
	// Replace all semi-colons with commas because semi-colons are comments for go-ini
	configStr := strings.Replace(string(contents), ";", ",", -1)

	// Write to temporary file
	f, err := os.Create(dxlTmpFile)
	check(err)
	defer f.Close()

	_, err = f.WriteString(configStr)
	check(err)

	// Parse Temporary file
	cfg, err = ini.Load(dxlTmpFile)
	check(err)

	return cfg
}

func main() {
	var configFlag = flag.String("dxlconfig", "dxlclient.config", "DXL Client Configuration File i.e. dxlclient.config")
	flag.Parse()
	configStr := *configFlag

	// Parse DXL Configuration File default: dxlclient.config
	dxlCfg := dxlConfigParser(configStr)

	// Create TLS Config for MQTT
	tlsconfig := NewTLSConfig(dxlCfg)

	// Connect to Broker
	connect(tlsconfig, dxlCfg)

}