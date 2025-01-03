package tdx

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"

	tdxproto "github.com/google/go-tdx-guest/proto/tdx"
	"github.com/google/go-tdx-guest/validate"
	"github.com/google/go-tdx-guest/verify"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
)

/*
* This function verifies TDX attestation report certificate chain and
* validates expected endorsement measurement against a measurement extracted
* from the attestation report.

@param attestation Attestation report as a tdxproto.QuoteV4 object in proto representation
@param hashed_endorsement_measurement Hashed retrieved initial measurement for validation
*/
func verify_attestation(attestation *tdxproto.QuoteV4, hashed_endorsement_measurement string) {
	err := verify.TdxQuote(attestation, verify.DefaultOptions())
	if err != nil {
		fmt.Println("Error verifying attestation:", err)
		return
	}

	decoded_hex, _ := hex.DecodeString(hashed_endorsement_measurement)
	err = validate.TdxQuote(attestation, &validate.Options{
		TdQuoteBodyOptions: validate.TdQuoteBodyOptions{
			MrTd: decoded_hex,
		},
	})
	if err != nil {
		fmt.Println("Error validating attestation:", err)
		return
	}

	fmt.Println("Attestation verified successfully!")
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
	var attestation tdxproto.QuoteV4

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
