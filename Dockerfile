FROM alpine

ARG MAINNET_HTTP_ENDPOINT_ARG
ENV MAINNET_HTTP_ENDPOINT=$MAINNET_HTTP_ENDPOINT_ARG

ARG PROJECT_ID_ARG 
ENV PROJECT_ID=$PROJECT_ID_ARG

ARG MAINNET_WEBSOCKET_ENDPOINT_ARG
ENV MAINNET_WEBSOCKET_ENDPOINT=$MAINNET_WEBSOCKET_ENDPOINT_ARG

ARG PROJECT_SECRET_ARG
ENV PROJECT_SECRET=$PROJECT_SECRET_ARG


# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Copy the code into the container
COPY ./build/infura-server-bin-linux /bin

# Export necessary port
EXPOSE 8000

# Command to run when starting the container
CMD ["/bin/infura-server-bin-linux"]
