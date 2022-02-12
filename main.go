package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"gopkg.in/yaml.v2"
)

type DNSCred map[string]string

func main() {
	http.HandleFunc("/udpate", updateIP)
	http.HandleFunc("/_ping", ping)
	log.Println("Server started")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func updateIP(w http.ResponseWriter, r *http.Request) {
	// fmt.Println(r)
	ip := r.FormValue("ip")

	log.Printf("Updating with IP: %s\n", ip)

	errCount := 0
	for _, cred := range dnsSecrets() {
		params := url.Values{}
		params.Add("hostname", cred["hostname"])
		params.Add("myip", ip)
		log.Printf("Updating %s\n", cred["hostname"])

		googUrl := fmt.Sprintf("https://%s:%s@domains.google.com/nic/update?%s", cred["username"], cred["password"], params.Encode())
		_, err := http.Get(googUrl)
		if err != nil {
			log.Fatalln(err)
			errCount += 1
			log.Printf("Error with %s\n", cred["hostname"])
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if errCount > 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(fmt.Sprintln(`{"status":"error updating `, errCount, ` hostnames"}`)))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}

func dnsSecrets() []DNSCred {
	f, err := os.Open("/data/dns_secrets.yaml")

	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	dec := yaml.NewDecoder(f)

	var dnsCreds []DNSCred

	err = dec.Decode(&dnsCreds)

	if err != nil {
		log.Fatal(err)
	}
	return dnsCreds
}
