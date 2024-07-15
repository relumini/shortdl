# ShortDL

## Overview

ShortDL is a tool for downloading YouTube videos, extracting metadata, and storing checksum values to avoid duplicate downloads.

## Features

- Download YouTube videos.
- Extract metadata, including description, captions, title, and transcript.
- Store and check checksum values in a PostgreSQL database to avoid duplicate downloads.
- Handle sensitive content and various errors gracefully.

## Setup

### Prerequisites

- Go 1.16+
- Docker
- PostgreSQL
- Gin framework
- GORM
- `youtube/v2` library

### Database Setup

1. **Run PostgreSQL in Docker:**

   ```sh
   docker run --rm --name dev-postgres -p 5432:5432 -e POSTGRES_PASSWORD=12345678 -d postgres
   ```

2. **Initialize the Database:**

   Ensure the database is running and update the connection string in `database/connect.go` if necessary.

### Environment Configuration

- Update the `DSN` in `database/connect.go` with your PostgreSQL credentials.

### Project Structure

- `main.go`: Application entry point.
- `routes/routes.go`: Defines API routes.
- `services/syoutube.go`: Contains the logic for downloading videos and extracting metadata.
- `database/database.go`: Database connection and initialization.
- `models/models.go`: Defines the database models.
- `utils/utils.go`: Utility functions.
- `handler/handler.go`: Error handling.

### Installation

1. **Clone the repository:**

   ```sh
   git clone https://github.com/yourusername/shortdl.git
   cd shortdl
   ```

2. **Install dependencies:**

   ```sh
   go mod tidy
   ```

3. **Run the application:**

   ```sh
   go run main.go
   ```

## Usage

### API Endpoint

- **GET /yshort**

  Extracts the video ID from the URL, downloads the video, and stores metadata.

  **Request:**

  ```http
  GET /yshort?url=https://www.youtube.com/watch?v=YOUR_VIDEO_ID
  ```

  **Response:**

  ```json
  {
    "message": "Successfully downloaded youtube",
    "data": {
      "Description": "Video description",
      "Caption": {
        "BaseUrl": "https://...",
        "Name": {
          "Language": "English",
          "Value": "English"
        }
      },
      "Title": "Video title",
      "Transcript": "Transcript text"
    }
  }
  ```

## Notes

- Ensure the `download` directory exists in the root of the project to store downloaded videos.
- The application checks for existing checksums to avoid downloading duplicate videos.
- Update `utils/utils.go` and `handler/handler.go` as needed to match your project structure.

## License

MIT License. See [LICENSE](LICENSE) for more information.

---
