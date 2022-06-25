FROM golang:1.18.2

# Select the working directory
WORKDIR GoLang/

# Copy everything from the current directory to the PWD (Present Working Directory) inside the container
COPY . .

# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

# Expose port 8181 to the outside world
EXPOSE 8181

# Run the executable
CMD ["main"]