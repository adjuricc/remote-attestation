package attestation

import (
        "fmt"
        "net/http"
        "log"
        "github.com/google/go-tdx-guest/client"
        "io/ioutil"
        "encoding/json"
        "encoding/base64"
        "google.golang.org/protobuf/proto"
        "google.golang.org/protobuf/reflect/protoreflect"
)

type Payload struct {
        Nonce string
}

func HandleAttestation(w http.ResponseWriter, r *http.Request){
        quote_provider, err := client.GetQuoteProvider()

        if err != nil {
                http.Error(w, fmt.Sprintf("Failed to get quote provider: %v", err), http.StatusInternalServerError)
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

        attestation, err := client.GetQuote(quote_provider, nonce)

        if err != nil {
                http.Error(w, fmt.Sprintf("Error obtaining attestation report: %v", err), http.StatusInternalServerError)
                log.Printf("Error: %v\n", err)
                return
        }

        attestation_proto, ok := attestation.(protoreflect.ProtoMessage)
        if !ok {
                http.Error(w, "Invalid attestation type", http.StatusInternalServerError)
                log.Printf("Error: attestation does not implement ProtoMessage\n")
                return
        }

        w.Header().Set("Content-Type", "application/octet-stream")


        data, err := proto.Marshal(attestation_proto)
        if err != nil {
                http.Error(w, fmt.Sprintf("Failed to serialize attestation report: %v", err), http.StatusInternalServerError)
                log.Printf("Error serializing attestation report: %v\n", err)
                return
        }

        w.Write(data)
}