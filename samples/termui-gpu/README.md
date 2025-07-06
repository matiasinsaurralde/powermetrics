# Terminal GPU Usage Dashboard

A real-time terminal dashboard that graphs GPU idle ratio using [termui](https://github.com/gizak/termui).

## Features

- **Real-time Updates**: Continuously collects GPU power samples every second
- **Sparkline Chart**: Displays GPU idle ratio as a live-updating sparkline
- **Full Terminal**: Uses the entire terminal window for maximum visibility
- **Responsive**: Automatically resizes when terminal window changes
- **Clean Exit**: Gracefully exits with Ctrl+C or 'q' key

## Usage

### Prerequisites

- macOS (required for powermetrics)
- Go 1.24 or later

### Running the Dashboard

```bash
cd samples
go run ./term-gpu-usage
```

### Controls

- **Any key**: Exit the dashboard
- **Ctrl+C**: Exit the dashboard
- **Terminal resize**: Chart automatically adapts to new window size

## What You'll See

The dashboard displays:
- **Title**: Shows current number of samples and latest idle ratio percentage
- **Sparkline**: Real-time graph of GPU idle ratio over time
- **Green line**: GPU idle ratio trend (higher = more idle)

## Technical Details

- **Sample Rate**: 1 sample per second
- **Data Window**: Last 30 samples (30 seconds of data)
- **Update Frequency**: Every second
- **Chart Type**: Sparkline with dot markers
- **Color**: Green line on dark background

## Example Output

```
GPU Idle Ratio (15 samples) - Latest: 98.25%
┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                             │
│                                                                             │
│                                                                             │
│                                                                             │
│                                                                             │
│                                                                             │
│                                                                             │
│                                                                             │
│                                                                             │
│                                                                             │
│                                                                             │
│                                                                             │
│                                                                             │
│                                                                             │
│                                                                             │
│                                                                             │
│                                                                             │
│                                                                             │
│                                                                             │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Dependencies

- `github.com/matiasinsaurralde/powermetrics` - For collecting GPU metrics
- `github.com/gizak/termui/v3` - For terminal UI and sparkline chart 