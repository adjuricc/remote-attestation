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

        return stdout.String(), nil
}

func validate_launch_endorsement(endorsement_file string) error {
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

        return nil
}

func get_tdx_measurement(golden_measurement *endorsement.VMGoldenMeasurement) (string, error) {
        tdx := golden_measurement.GetTdx()
        if tdx == nil {
                return "", fmt.Errorf("no TDX data available")
        }

        measurements := tdx.GetMeasurements()
        i := 0
        for _, measurement := range measurements {
                if i == 1 {
                        return hex.EncodeToString(measurement.GetMrtd()), nil
                }
                i += 1
        }

        return "", fmt.Errorf("no valid MRTD found in TDX measurements")
}

func get_hashed_measurement(endorsement_proto *endorsement.VMLaunchEndorsement)(string, error){
        serialized_golden := endorsement_proto.GetSerializedUefiGolden()

        var golden_measurement endorsement.VMGoldenMeasurement
        err := proto.Unmarshal(serialized_golden, &golden_measurement)
        if err != nil {
                return "", fmt.Errorf("failed to unmarshal serialized UEFI golden: %v", err)
        }

        hashed_measurement, err := get_tdx_measurement(&golden_measurement)
        if err != nil {
                return "", fmt.Errorf("failed to get hashed TDX measurement: %v", err)
        }

        return hashed_measurement, nil
}

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