## YouTube Search Service

This is a GO-lang service that allows users to search for YouTube videos. 

## Instructions

To run the service, follow these steps:

1. First, navigate to the root directory of the project in your terminal.

2. Build the Docker image using the following command:
    ```bash
    docker build -t yt_search_service .
    ```

3. Run the Docker container, binding port 8090 of the container to port 8090 of your host machine, and giving the container a name:
    ```bash
    docker run -p 8090:8090 --name yt_search_service yt_search_service
    ```

These commands will build the Docker image for your YouTube search service and then run a container from it. You can access the service at `http://localhost:8090` once it's up and running.

## Usage

Once the service is running, you can use it to search for YouTube videos. You can interact with the service via HTTP requests. Here's an example:

- **Endpoint**: `http://localhost:8090/health`
- **Method**: GET
  - This endpoint returns a 200 OK response if the service is healthy.

## Dependencies

This service is built using Go programming language. It uses Docker for containerization. Make sure you have Docker installed on your system before running the service.

