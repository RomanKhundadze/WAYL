# What Are You Listening? (WAYL)

**WAYL** is a locally hosted browser source that displays what you're listening to on Spotify.

## Installation and Configuration

1. **Download** the project from GitHub and place it in a local directory.

2. **First Launch**:  
   Run the application once. The program will generate a `config.toml` file and then exit with an error since the generated file needs to be configured.

3. **Configure** the `config.toml` file with the following details:
   ```toml
   ClientID     = ""
   ClientSecret = ""
   Port         = ""
   ```
   - **ClientID** and **ClientSecret**: Obtain these from the [Spotify Developer Dashboard](https://developer.spotify.com/).
   - **Port**: Choose a port number (e.g., `8080`).

4. **Second Launch**:
   Launch the application again. Open your browser and go to `http://localhost:PORT/` (replace `PORT` with the port number you configured).
   You will be redirected to a login screen where you need to log in with your Spotify account.
   After a successful login, you will be redirected to `http://localhost:PORT/playback`, where the currently playing music will be displayed.

## Usage

When running `WAYL`, open your browser and navigate to the specified port (e.g., `http://localhost:8080/`).

### Command-Line Arguments

WAYL accepts several command-line arguments:

- **Without Arguments**  
  Runs the server normally. Creates any missing files if they are not found.

- **`-k`**  
  Kills all running instances of the program.  
  *Note: No other arguments can be used with `-k`.*  
  *Example:*  
  ```bash
  ./WAYL -k
  ```

- **`-rI`**  
  Resets all configurations and tokens. Deletes all configuration and token data and reinitializes them.

- **`-rT`**  
  Resets the token data. Deletes the token data and reinitializes the token directory.

- **`-rC`**  
  Resets the configuration files. Deletes the configuration files and reinitializes them.

### Example Usage

- Launch the server normally:  
  ```bash
  ./WAYL
  ```
- Kill all running instances (cannot be combined with other arguments):  
  ```bash
  ./WAYL -k
  ```
- Reset all configurations and tokens:  
  ```bash
  ./WAYL -rI
  ```
- Reset only the token data:  
  ```bash
  ./WAYL -rT
  ```
- Reset only the configuration files:  
  ```bash
  ./WAYL -rC
  ```

## Configuration Files

The `config` directory contains:

- **`playback.html`**: HTML file to display the currently playing track.
- **`styles.css`**: CSS file to style the playback display.
- **`script.js`**: JavaScript file to handle dynamic updates and interactions.

These files can be customized to your preferences. If any file is broken or improperly configured, use the `-rC` argument to reset them to their default state.
