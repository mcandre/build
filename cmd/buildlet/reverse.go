// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"crypto/hmac"
	"crypto/md5"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"golang.org/x/build"
	"golang.org/x/build/revdial"
)

func dialCoordinator() error {
	devMode := !strings.HasPrefix(*coordinator, "farmer.golang.org")

	if *hostname == "" {
		*hostname, _ = os.Hostname()
	}

	modes := strings.Split(*reverse, ",")
	var keys []string
	for _, m := range modes {
		if devMode {
			keys = append(keys, string(devBuilderKey(m)))
			continue
		}
		keyPath := filepath.Join(homedir(), ".gobuildkey-"+m)
		key, err := ioutil.ReadFile(keyPath)
		if os.IsNotExist(err) && len(modes) == 1 {
			globalKeyPath := filepath.Join(homedir(), ".gobuildkey")
			key, err = ioutil.ReadFile(globalKeyPath)
			if err != nil {
				log.Fatalf("cannot read either key file %q or %q: %v", keyPath, globalKeyPath, err)
			}
		} else if err != nil {
			log.Fatalf("cannot read key file %q: %v", keyPath, err)
		}
		keys = append(keys, string(key))
	}

	caCert := build.ProdCoordinatorCA
	addr := *coordinator
	if addr == "farmer.golang.org" {
		addr = "farmer.golang.org:443"
	}
	if devMode {
		caCert = build.DevCoordinatorCA
	}

	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM([]byte(caCert)) {
		log.Fatal("failed to append coordinator CA certificate")
	}

	log.Printf("Dialing coordinator %s...", addr)
	tcpConn, err := net.Dial("tcp", addr)
	if err != nil {
		return err // try again
	}
	config := &tls.Config{
		ServerName:         "go",
		RootCAs:            caPool,
		InsecureSkipVerify: devMode,
	}
	conn := tls.Client(tcpConn, config)
	if err := conn.Handshake(); err != nil {
		return fmt.Errorf("failed to handshake with coordinator: %v", err)
	}
	bufr := bufio.NewReader(conn)

	req, err := http.NewRequest("GET", "/reverse", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header["X-Go-Builder-Type"] = modes
	req.Header["X-Go-Builder-Key"] = keys
	req.Header.Set("X-Go-Builder-Hostname", *hostname)
	req.Header.Set("X-Go-Builder-Version", strconv.Itoa(buildletVersion))
	if err := req.Write(conn); err != nil {
		return fmt.Errorf("coordinator /reverse request failed: %v", err)
	}
	resp, err := http.ReadResponse(bufr, req)
	if err != nil {
		return fmt.Errorf("coordinator /reverse response failed: %v", err)
	}
	if resp.StatusCode != 101 {
		msg, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("coordinator registration failed:\n\t%s", msg)
	}

	log.Printf("Connected to coordinator; reverse dialing active")
	srv := &http.Server{}
	err = srv.Serve(revdial.NewListener(bufio.NewReadWriter(
		bufio.NewReader(conn),
		bufio.NewWriter(conn),
	)))
	log.Printf("Reverse buildlet Serve complete; err=%v", err)
	return err
}

func devBuilderKey(builder string) string {
	h := hmac.New(md5.New, []byte("gophers rule"))
	io.WriteString(h, builder)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func homedir() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
	}
	home := os.Getenv("HOME")
	if home != "" {
		return home
	}
	if os.Getuid() == 0 {
		return "/root"
	}
	return "/"
}

// TestDialCoordinator dials the coordinator. Exported for testing.
func TestDialCoordinator() {
	// TODO(crawshaw): move some of this logic out of main to simplify testing hook.
	http.Handle("/status", http.HandlerFunc(handleStatus))
	dialCoordinator()
}
