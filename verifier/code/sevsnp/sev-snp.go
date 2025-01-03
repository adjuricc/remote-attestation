package sevsnp

import (
	"fmt"
	"log"
	"os"

	sevproto "github.com/google/go-sev-guest/proto/sevsnp"
	"github.com/google/go-sev-guest/verify"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
)

/*
* This function verifies SEV-SNP attestation report certificate chain and
* validates expected endorsement measurement against a measurement extracted
* from the attestation report.

@param attestation Attestation report as a sevproto.Attestation object in proto representation
@param hashed_endorsement_measurement Hashed retrieved initial measurement for validation
*/
func verify_attestation(attestation *sevproto.Attestation, hashed_endorsement_measurement string) {
	err := verify.SnpAttestation(attestation, verify.DefaultOptions())
	if err != nil {
		fmt.Println("Error verifying attestation:", err)
	}

	measurement := attestation.GetReport().GetMeasurement()

	measurement_hex := fmt.Sprintf("%x", measurement)

	if hashed_endorsement_measurement == measurement_hex {
		fmt.Printf("Attestation verified successfully!")
	} else {
		fmt.Printf("Invalid Measurement")
	}

}

/*
* This function handles attestation report decoding bytes to a protocol buffer
* format of attestation report, then writing it to a .textproto file. It calls
* for attestation verification function.
* It is used as an external function.

@param attestation_bytes Attestation report in byte representation
@param hashed_endorsement_measurement Hashed retrieved initial measurement for validation
*/
func HandleAttestation(attestation_bytes []byte, hashed_endorsement_measurement string) {
	var attestation sevproto.Attestation

	err := proto.Unmarshal(attestation_bytes, &attestation)
	if err != nil {
		log.Fatalf("Failed to decode Protobuf attestation: %v", err)
	}

	textData := prototext.Format(&attestation)

	err = os.WriteFile("attestation.textproto", []byte(textData), 0644)
	if err != nil {
		log.Fatalf("Failed to write attestation to file: %v", err)
	}

	verify_attestation(&attestation, hashed_endorsement_measurement)
}
