# Use a Go base image
FROM golang:1.23

# Set the working directory inside the container
WORKDIR /cmd/server

# Copy all project files to the container
COPY . .

# Expose the port your server listens on
EXPOSE 8080

# Command to run your server without building
CMD ["go", "run", "."]