package main

import (
        "fmt"
        "net/http"
        "log"
        "attestation-server/retrieve_launch_endorsement"
        "attestation-server/attestation"
)

/* This function handles attestation generation */
func handle_attestation(w http.ResponseWriter, r *http.Request){
        attestation.HandleAttestation(w, r)
}

/* This function handles launch endorsement retrieval */
func handle_retrieve_launch_endorsement(w http.ResponseWriter, r *http.Request){
        retrieve_launch_endorsement.HandleRetrieveLaunchEndorsement(w, r)
}

/* Handling incoming requests */
func main(){
        http.HandleFunc("/attest", handle_attestation)
        http.HandleFunc("/retrieve_launch_endorsement", handle_retrieve_launch_endorsement)
        port := 8080
        fmt.Printf("Server listening on port %d\n", port)
        if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
                log.Fatalf("Failed to start server: %v", err)
        }
}