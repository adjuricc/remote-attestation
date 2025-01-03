package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"go-proba/sevsnp"
	"go-proba/tdx"
	"io"
	"log"
	"net/http"
)

/*
* This function generates a random 64 bytes nonce value
* It's used when generating an attestation report for
* replay attacks

@return 64 bytes nonce value
*/
func generate_nonce() ([64]byte, error) {
	var nonce [64]byte
	_, err := rand.Read(nonce[:])
	if err != nil {
		return [64]byte{}, err
	}
	return nonce, nil
}

/*
* This function sends a POST http request to a specified url
* for fetching an attestation report from Attester side using
* generated nonce value

@param url Specified url
@return Attestation report in bytes representation
*/
func request_attestation(url string) ([]byte, error) {
	nonce, err := generate_nonce()
	if err != nil {
		fmt.Println("Error generating nonce:", err)
		return nil, nil
	}

	nonce_encoded := base64.StdEncoding.EncodeToString(nonce[:])
	body := []byte(fmt.Sprintf(`{"nonce": "%s"}`, nonce_encoded))

	r, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return buf.Bytes(), nil
}

/*
* This function sends a GET http request to a specified url
* for obtaining a hashed expected measurement from Attester
* side.

@param url Specified url
@return Hashed measurement representing initial firmware in string
*/
func request_retrieve_launch_endorsement(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to fetch data: %s\n", resp.Status)
		return ""
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return ""
	}

	return string(body)
}

func main() {
	/* External IP */
	var MY_VM_IP string
	fmt.Print("Enter IP Address: ")
	fmt.Scanln(&MY_VM_IP)

	/* Could be SEV-SNP or TDX */
	var TYPE string
	fmt.Print("Choose technology type (SEV-SNP or TDX): ")
	fmt.Scanln(&TYPE)

	/* Request End points */
	attestation_url := fmt.Sprintf("http://%s:8080/attest", MY_VM_IP)
	endorsement_url := fmt.Sprintf("http://%s:8080/retrieve_launch_endorsement", MY_VM_IP)

	/* Retrieve hashed Firmware measurements for attestation validation */
	hashed_endorsement_measurement := request_retrieve_launch_endorsement(endorsement_url)

	/* Send attestation request to VM */
	attestation_bytes, err := request_attestation(attestation_url)
	if err != nil {
		log.Fatalf("Error getting attestation: %v", err)
	}

	/* Handle attestation response */
	if TYPE == "SEV-SNP" {
		sevsnp.HandleAttestation(attestation_bytes, hashed_endorsement_measurement)
	} else if TYPE == "TDX" {
		tdx.HandleAttestation(attestation_bytes, hashed_endorsement_measurement)
	} else {
		log.Fatalf("Invalid technology type.")
	}
}
