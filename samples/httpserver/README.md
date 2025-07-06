# HTTP Server Sample

A simple HTTP server that exposes GPU power metrics via a REST API endpoint.

## Features

- **REST API**: Exposes GPU metrics at `/gpu` endpoint
- **Pretty JSON**: Returns formatted, human-readable JSON responses
- **CORS Support**: Can be accessed from web browsers
- **Real-time Data**: Collects fresh GPU metrics on each request
- **Comprehensive Data**: Includes frequency, idle ratio, DVFM states, and system info

## Usage

### Start the Server

```bash
cd samples
go run ./httpserver
```

The server will start on port 8080.

### Access the API

```bash
# Get GPU metrics
curl http://localhost:8080/gpu
```

### Example Response

```json
{
  "timestamp": "2025-01-06T10:30:45Z",
  "gpu": {
    "frequency_hz": 338000000,
    "frequency_mhz": 338.0,
    "idle_ratio": 0.995324,
    "idle_percent": 99.5324,
    "dvfm_states": [
      {
        "frequency_hz": 338000000,
        "used_ratio": 0.004676,
        "used_percent": 0.4676
      },
      {
        "frequency_hz": 618000000,
        "used_ratio": 0.0,
        "used_percent": 0.0
      }
    ]
  },
  "hardware_model": "Mac16,8",
  "kernel_version": "24F74"
}
```

## API Endpoints

### GET `/gpu`

Returns current GPU power metrics in JSON format.

**Response Fields:**
- `timestamp`: When the metrics were collected
- `gpu.frequency_hz`: Current GPU frequency in Hz
- `gpu.frequency_mhz`: Current GPU frequency in MHz
- `gpu.idle_ratio`: GPU idle ratio (0.0 to 1.0)
- `gpu.idle_percent`: GPU idle percentage (0% to 100%)
- `gpu.dvfm_states`: Array of DVFM (Dynamic Voltage and Frequency Management) states
- `hardware_model`: Mac hardware model
- `kernel_version`: macOS kernel version

## Web Browser Access

The server includes CORS headers, so you can access it directly from a web browser:

```
http://localhost:8080/gpu
```

## Error Handling

- Returns `500 Internal Server Error` if powermetrics fails to execute
- Returns `500 Internal Server Error` if no GPU data is available
- Includes descriptive error messages in the response body

## Dependencies

- `github.com/matiasinsaurralde/powermetrics` - For collecting GPU metrics
- Standard library `net/http` - For HTTP server functionality 