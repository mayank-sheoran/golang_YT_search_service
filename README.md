# YouTube Search Service

Go-based service allows users to search for YouTube videos efficiently using **elastic search** and **postgres**. It utilizes Docker for containerization, ensuring an isolated and consistent environment for running the service.

## Demo video
https://youtu.be/eNlZf3WAQPA

## Prerequisites

- Docker
- YouTube API key

## Setup

1. **Clone the project:** Clone the repository to your local machine and navigate to the project's root directory in your terminal.

2. **Configure API keys:** Enter your YouTube API key in the `.ENV` file. For multiple API keys, add them in `/internal/service/youtube_data` on line 34.

3. **Build the Docker image:** Use the following command to build the Docker image:
   ```bash
   docker-compose build
4. **Run the Docker containers:**
   Start the containers with Docker Compose
    ```bash
   docker-compose up -d
   ```
## Accessing the Service

Once the service is running, you can access it at `http://localhost:8090`.

### Postman collection
https://drive.google.com/file/d/1be2noJtrrlXaKZHrzzFb6_Uq_tCS3tNI/view

### Endpoints

#### Health Check
- **Endpoint**: `http://localhost:8090/health`
- **Method**: GET
    - Returns a 200 OK response if the service is healthy.

#### Search for Videos
- **Endpoint**: `http://localhost:8090/videos/search?search-query=<QUERY>&page-number=<PAGE_NUMBER>&page-limit=<PAGE_LIMIT>`
- **Method**: GET
    - Searches for videos related to the specified query. If no query is specified, it returns all videos in a paginated format. Default settings are `page-number=1` and `page-limit=100`.

## Dependencies
- **Docker**: Used for creating a consistent running environment for the service.


