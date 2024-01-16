# Use an official golang runtime as the base image 
FROM golang:alpine

#Set the working directory inside the container
WORKDIR /app

# Copy the local package files to the container's workspace
COPY ./pScan /app

#Build the go application inside the container
RUN go build -o pscan

#Define the command to run your application
ENTRYPOINT ["./pscan"]