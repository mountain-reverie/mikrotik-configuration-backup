FROM scratch

# GoReleaser v2 places binaries in platform-specific directories
ARG TARGETPLATFORM
COPY $TARGETPLATFORM/mikrotik-backup /mikrotik-backup

# Run as non-root user (nobody:nobody = 65534:65534)
USER 65534:65534

ENTRYPOINT ["/mikrotik-backup"]
