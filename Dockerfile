FROM alpine:3.14
# Copy our static executable.
COPY ./loggen_linux /main
# Run the binary.
ENTRYPOINT ["/main"]
