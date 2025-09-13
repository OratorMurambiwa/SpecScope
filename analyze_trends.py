import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns
from datetime import datetime
import argparse
import subprocess
import json

# --- Load data ---
def load_from_csv(path):
    df = pd.read_csv(path)
    return enrich_df(df)

def load_from_simulated(temp_file="simulated_output.json", samples=200):
    subprocess.run(["go", "run", "specscope-backend/simulate/simulate.go", temp_file, str(samples)], check=True)
    with open(temp_file, "r") as f:
        data = json.load(f)
    df = pd.DataFrame(data)
    return enrich_df(df)

# --- Enrichment ---
def enrich_df(df):
    df["timestamp"] = pd.to_datetime(df["timestamp"])
    df["hour"] = df["timestamp"].dt.hour
    df["day_of_week"] = df["timestamp"].dt.day_name()
    df["month"] = df["timestamp"].dt.month_name()

    def get_daypart(hour):
        if 5 <= hour < 12:
            return "Morning"
        elif 12 <= hour < 17:
            return "Afternoon"
        elif 17 <= hour < 21:
            return "Evening"
        else:
            return "Night"

    df["daypart"] = df["hour"].apply(get_daypart)

    # If missing, default interference column to 0
    if "interference" not in df.columns:
        df["interference"] = 0

    return df

# --- Plotting ---
def plot_interference_by_hour(df):
    hourly = df.groupby("hour")["interference"].agg(["sum", "count"])
    hourly["percentage"] = (hourly["sum"] / hourly["count"]) * 100
    plt.figure(figsize=(10, 5))
    sns.barplot(x=hourly.index, y=hourly["percentage"], color='steelblue')
    plt.title("Interference % by Hour of Day")
    plt.ylabel("Percentage (%)")
    plt.xlabel("Hour")
    plt.xticks(range(0, 24))
    plt.tight_layout()
    plt.show()

def plot_interference_by_day(df):
    order = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"]
    daily = df.groupby("day_of_week")["interference"].agg(["sum", "count"])
    daily["percentage"] = (daily["sum"] / daily["count"]) * 100
    plt.figure(figsize=(10, 5))
    sns.barplot(x=daily.index, y=daily["percentage"], order=order, color='orange')
    plt.title("Interference % by Day of Week")
    plt.ylabel("Percentage (%)")
    plt.xlabel("Day")
    plt.tight_layout()
    plt.show()

def plot_interference_by_month(df):
    month_order = [
        "January", "February", "March", "April", "May", "June",
        "July", "August", "September", "October", "November", "December"
    ]
    monthly = df.groupby("month")["interference"].agg(["sum", "count"])
    monthly["percentage"] = (monthly["sum"] / monthly["count"]) * 100
    plt.figure(figsize=(12, 5))
    sns.barplot(x=monthly.index, y=monthly["percentage"], order=month_order, color='green')
    plt.title("Interference % by Month")
    plt.ylabel("Percentage (%)")
    plt.xlabel("Month")
    plt.tight_layout()
    plt.show()

def plot_interference_by_daypart(df):
    part_order = ["Morning", "Afternoon", "Evening", "Night"]
    part = df.groupby("daypart")["interference"].agg(["sum", "count"])
    part["percentage"] = (part["sum"] / part["count"]) * 100
    plt.figure(figsize=(8, 5))
    sns.barplot(x=part.index, y=part["percentage"], order=part_order, color='purple')
    plt.title("Interference % by Daypart")
    plt.ylabel("Percentage (%)")
    plt.xlabel("Time of Day")
    plt.tight_layout()
    plt.show()

# --- Main Execution ---
if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--source", choices=["file", "simulated"], required=True)
    parser.add_argument("--path", type=str, help="CSV file path if using --source file")
    parser.add_argument("--samples", type=int, default=200, help="Number of simulated samples")
    args = parser.parse_args()

    if args.source == "file":
        if not args.path:
            raise ValueError("Please specify --path when using --source file")
        df = load_from_csv(args.path)
    else:
        df = load_from_simulated(samples=args.samples)

    plot_interference_by_hour(df)
    plot_interference_by_day(df)
    plot_interference_by_month(df)
    plot_interference_by_daypart(df)
