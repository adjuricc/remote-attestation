module attestation-server

go 1.22.0

require (
        github.com/google/gce-tcb-verifier v0.2.2
        github.com/google/go-sev-guest v0.11.2-0.20241017023127-f94d851ddd48
        google.golang.org/protobuf v1.33.0
)

require (
        github.com/google/go-configfs-tsm v0.2.2 // indirect
        github.com/google/logger v1.1.1 // indirect
        github.com/google/uuid v1.6.0 // indirect
        go.uber.org/multierr v1.11.0 // indirect
        golang.org/x/crypto v0.21.0 // indirect
        golang.org/x/sys v0.18.0 // indirect
)

replace github.com/google/gce-tcb-verifier => ../gce-tcb-verifier