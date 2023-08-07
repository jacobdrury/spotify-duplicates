# Spotify Playlist Duplicate Remover

[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

## Overview

Spotify Playlist Duplicate Remover is a Golang application that helps you remove all duplicate songs from a specified Spotify playlist. This project utilizes the Spotify Web API to access and manipulate playlists. It is designed to run as a standalone application or deployed on a server for automated playlist maintenance.

## Features

- Connects to the Spotify Web API using OAuth2 authentication.
- Retrieves the list of songs in a specified playlist.
- Identifies and removes duplicate songs, keeping only unique tracks in the playlist.
- Supports handling large playlists efficiently through parallelization.

## Installation

### Prerequisites

- Go version 1.16 or later installed.
- Spotify Developer Account: [Sign up here](https://developer.spotify.com/dashboard/applications)

### Getting Started

1. **Clone the repository:**
   ```shell
   git clone https://github.com/yourusername/spotify-playlist-duplicate-remover.git
   ```

2. **Install dependencies:**
    ```shell
    go mod download
    ```
   
3. **Create Spotify App**

   1. Once on the dev dashboard click [Create App](https://developer.spotify.com/dashboard/create)

   2. Fill in the required fields (name, description, redirect uri)
      - RedirectUri is the uri spotify will redirect the user to after successful authentication. \
      Ex: http://localhost:8080/callback
4. **Set up environment variables:**   
   1. Create a `.env` file in the project root. Use `.env.example` as a template
        ```shell
         cp .env.example .env 
        ```
       
   2. Update `.env` with your credentials

   #### Accepted .env variables

   | Key              | Description                                   |
      |------------------|-----------------------------------------------|
   | `SPOTIFY_ID`     | your spotify client id                        |
   | `SPOTIFY_SECRET` | your spotify client secret                    |
   | `BASE_URI`       | baseUri for callback endpoint to be hosted at |
   | `PORT`           | port used for callback endpoint               |

   > *Note*: `BASE_URI` and `PORT` must match the `RedirectUri` set on the spotify app in Step 3

   #### Example `.env`
     ```dotenv
     SPOTIFY_ID=spotify_id
     SPOTIFY_SECRET=spotify_secret
     BASE_URI=localhost
     PORT=8080
     ```

5. **Running the Application:**

   You can run the application with the following flags:

   - To remove duplicates from a specific playlist, use the -playlistIds flag with a comma-separated list of playlist IDs:

   ``` shell
   go run main.go -playlistIds "3cEYpjA9oz9GiPac4AsH4n,..."
   ```

   - To remove duplicates from all playlists, use the -all flag:

   ```shell
   go run main.go -all
   ```

   The application will authenticate with the Spotify Web API using the provided credentials and then proceed to remove duplicate songs from the specified playlists.

   Note: If you only pass one playlist ID to the -playlistIds flag, the application will treat it as a single ID.