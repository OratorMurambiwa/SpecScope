import argparse
import pandas as pd
import matplotlib.pyplot as plt
import json
import subprocess
import os

def load_from_simulated(samples=300, outfile="simulated_output.json"):
    subprocess.run(["go", "run", "specscope-backend/simulate/simulate.go", outfile, str(samples)], check=True)
    with open(outfile, 'r') as f:
        data = json.load(f)
    return pd.DataFrame(data)

def load_from_csv(file_path):
    return pd.read_csv(file_path)

def filter_band(df, start_freq, end_freq):
    return df[(df['frequency'] >= start_freq) & (df['frequency'] <= end_freq)]

def detect_interference(df, power_threshold=-10):
    return df[df['power'] > power_threshold]

def plot_band(df, interference_df, band_start, band_end):
    plt.figure(figsize=(12, 6))
    plt.plot(pd.to_datetime(df['timestamp']), df['power'], label='Signal Power')
    if not interference_df.empty:
        plt.scatter(pd.to_datetime(interference_df['timestamp']), interference_df['power'], color='red', label='Interference', zorder=5)
    plt.title(f"Power vs Time in {band_start}-{band_end} MHz Band")
    plt.xlabel("Time")
    plt.ylabel("Power (dBm)")
    plt.legend()
    plt.grid(True)
    plt.tight_layout()
    plt.show()

def main():
    parser = argparse.ArgumentParser(description="Monitor specific RF frequency band for interference.")
    parser.add_argument('--source', choices=['simulated', 'csv'], default='simulated', help="Data source")
    parser.add_argument('--samples', type=int, default=300, help="Number of samples (for simulated)")
    parser.add_argument('--csv-path', type=str, help="Path to CSV file (if source is csv)")
    parser.add_argument('--band-start', type=float, required=True, help="Start frequency of the band in MHz")
    parser.add_argument('--band-end', type=float, required=True, help="End frequency of the band in MHz")
    parser.add_argument('--threshold', type=float, default=-10, help="Power threshold (dBm) for interference")

    args = parser.parse_args()

    # Load data
    if args.source == 'simulated':
        df = load_from_simulated(samples=args.samples)
    elif args.source == 'csv':
        if not args.csv_path:
            print("Error: You must provide --csv-path when using CSV source.")
            return
        df = load_from_csv(args.csv_path)

    # Filter by frequency band
    band_df = filter_band(df, args.band_start, args.band_end)

    # Detect interference
    interference_df = detect_interference(band_df, args.threshold)

    # Plot
    plot_band(band_df, interference_df, args.band_start, args.band_end)

if __name__ == "__main__":
    main()
