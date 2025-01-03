package retrieve_launch_endorsement

import (
        "bytes"
        "encoding/hex"
        "fmt"
        "io/ioutil"
        "os/exec"
        "github.com/google/gce-tcb-verifier/proto/endorsement"
        "google.golang.org/protobuf/proto"
        "net/http"
)

/*
        * This function runs the `gce-tcb-verifier extract` command to generate
        * a launch endorsement file at the specified path.
        * It then returns the contents of that file as a string

        @param file_path Specified path
        @return Content of file at the specified path in string
*/

func retrieve_launch_endorsement(file_path string) (string, error) {
        fmt.Println("Retrieving launch endorsement...")
        cmd := exec.Command("./gce-tcb-verifier", "extract", "--out", file_path)
        var stdout, stderr bytes.Buffer
        cmd.Stdout = &stdout
        cmd.Stderr = &stderr
        cmd.Dir = "/home/anjci011/go_projects/attestation-server/retrieve_launch_endorsement"
        err := cmd.Run()
        if err != nil {
                fmt.Printf("Error: %s\n", stderr.String())
                fmt.Println(err)
                return "", fmt.Errorf("failed to retrieve launch endorsement: output file is missing or empty")
        }

        fmt.Println("Output:")
        fmt.Println(stdout.String())

        return stdout.String(), nil
}

/*
        * This function runs the `gce-tcb-verifier verify` command to verify
        * the launch endorsement file authenticity.

        @param file_path Specified path of the launch endorsement file
        @return Error
*/

func validate_launch_endorsement(endorsement_file string) error {
        fmt.Println("Validating launch endorsement...")

        cmd := exec.Command("./gce-tcb-verifier", "verify", endorsement_file)
        var stdout, stderr bytes.Buffer
        cmd.Stdout = &stdout
        cmd.Stderr = &stderr
        cmd.Dir = "/home/anjci011/go_projects/attestation-server/retrieve_launch_endorsement"
        err := cmd.Run()
        if err != nil {
                fmt.Printf("Error: %s\n", stderr.String())
                fmt.Println(err)
                return fmt.Errorf("failed to verify launch endorsement")
        }

        fmt.Println("Output:")
        fmt.Println(stdout.String())

        return nil
}

/*
        * This function takes a VMGoldenMeasurement object, extracts its SEV-SNP
        * and produces a hashed firmware and hardware measurement

        @param endorsement_proto Golden Measurement in protobuf
        @return Hashed SEV-SNP measurement
*/

func get_sevsnp_measurement(golden_measurement *endorsement.VMGoldenMeasurement) (string, error){
        sev_snp := golden_measurement.GetSevSnp()
        if sev_snp == nil {
                return "", fmt.Errorf("no SEVSNP data available")
        }

        for key, value := range sev_snp.Measurements {
                if key == 8 {
                        return hex.EncodeToString(value), nil
                }
        }

        return "", nil

}

/*
        * This function takes a VMLaunchEndorsement object, extracts its serialized UEFI Golden Measurement
        * and produces a hashed SEV-SNP measurement containing firmware and hardware information.

        @param endorsement_proto Launch Endorsement in protobuf
        @return Hashed SEV-SNP measurement
*/

func get_hashed_measurement(endorsement_proto *endorsement.VMLaunchEndorsement)(string, error){
        serialized_golden := endorsement_proto.GetSerializedUefiGolden()

        var golden_measurement endorsement.VMGoldenMeasurement
        err := proto.Unmarshal(serialized_golden, &golden_measurement)
        if err != nil {
                return "", fmt.Errorf("failed to unmarshal serialized UEFI golden: %v", err)
        }

        hashed_measurement, err := get_sevsnp_measurement(&golden_measurement)
        if err != nil {
                return "", fmt.Errorf("failed to get hashed SEVSNP measurement: %v", err)
        }

        return hashed_measurement, nil
}

/*
        * This function handles launch endorsement retrieval then verifies it's authenticity and
        * produces a hashed SEV-SNP measurement containing firmware and hardware information.
        * It is used as an external function.

        @param w http ResponseWriter object
        @param r http Request
        @return Hashed SEV-SNP measurement
*/

func HandleRetrieveLaunchEndorsement(w http.ResponseWriter, r *http.Request){
        endorsement_file := "endorsement.json"
        _, err := retrieve_launch_endorsement(endorsement_file)
        if err != nil {
                fmt.Println("error Retrieving launch endorsement...")
                http.Error(w, fmt.Sprintf("Error retrieving endorsement: %v", err), http.StatusInternalServerError)
                return
        }

        err = validate_launch_endorsement(endorsement_file)
        if err != nil {
                http.Error(w, fmt.Sprintf("Error validating endorsement: %v", err), http.StatusInternalServerError)
                return
        }

        fmt.Println("Launch endorsement validated successfully!")

        data, err := ioutil.ReadFile("/home/anjci011/go_projects/attestation-server/retrieve_launch_endorsement/endorsement.json")
        if err != nil {
                fmt.Println("Failed to read file")
                fmt.Println(err)
                http.Error(w, fmt.Sprintf("Failed to read file: %v", err), http.StatusInternalServerError)
                return
        }

        fmt.Println("endorsement.json read  successfully!")

        var endorsementProto endorsement.VMLaunchEndorsement
        err = proto.Unmarshal(data, &endorsementProto)
        if err != nil {
                http.Error(w, fmt.Sprintf("Failed to unmarshal proto: %v", err), http.StatusInternalServerError)
                return
        }

        fmt.Println("Proto unmarshaled successfully!")

        hashed_measurement, err := get_hashed_measurement(&endorsementProto)
        if err != nil {
                http.Error(w, fmt.Sprintf("Failed to get hashed measurement: %v", err), http.StatusInternalServerError)
                return
        }

        fmt.Println("Hashed measurement successfully!")

        w.Header().Set("Content-Type", "text/plain")
        w.Write([]byte(hashed_measurement))
}