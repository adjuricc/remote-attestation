package attestation

import (
        "fmt"
        "net/http"
        "log"
        "io/ioutil"
        "encoding/json"
        "encoding/base64"
        "github.com/google/go-sev-guest/client"
        "google.golang.org/protobuf/proto"
)

type Payload struct {
        Nonce string
}

/*
        * This function opens a connection to /dev/sev-guest Linux device
        * and produces an attestation report with a decoded nonce value.

        @param w http Response Writer object
        @param r http Request
        @return Protobuf serialized attestation report
*/

func HandleAttestation(w http.ResponseWriter, r *http.Request){
        device, err := client.OpenDevice()

        if err != nil {
                http.Error(w, fmt.Sprintf("Error opening device: %v", err), http.StatusInternalServerError)
                log.Printf("Error: %v\n", err)
                return
        }

        body, err := ioutil.ReadAll(r.Body)

        if err != nil {
                http.Error(w, "Failed to read request body", http.StatusBadRequest)
                log.Printf("Error reading body: %v\n", err)
                return
        }


        var payload Payload

        err = json.Unmarshal(body, &payload)

        if err != nil {
                http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
                                }

        decoded_nonce, err := base64.StdEncoding.DecodeString(payload.Nonce)
        if err != nil || len(decoded_nonce) != 64 {
                http.Error(w, "Invalid nonce format", http.StatusBadRequest)
                return
        }
        var nonce [64]byte
        copy(nonce[:], decoded_nonce)


        attestation, err := client.GetExtendedReport(device, nonce)
        if err != nil {
                http.Error(w, fmt.Sprintf("Error obtaining attestation report: %v", err), http.StatusInternalServerError)
                log.Printf("Error: %v\n", err)
                return
        }
        defer device.Close()

        data, err := proto.Marshal(attestation)
        if err != nil {
                http.Error(w, fmt.Sprintf("Failed to serialize attestation report: %v", err), http.StatusInternalServerError)
                log.Printf("Error serializing attestation report: %v\n", err)
                return
        }

        w.Header().Set("Content-Type", "application/octet-stream")
        w.Write(data)
}