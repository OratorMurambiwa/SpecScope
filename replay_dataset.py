import pandas as pd
import joblib
import numpy as np
from datetime import datetime
import os
import sys

# Load trained model
model = joblib.load("rf_model.pkl")

def predict_from_csv(path):
    if not os.path.exists(path):
        print(f"File not found: {path}")
        return

    try:
        df = pd.read_csv(path)
    except Exception as e:
        print(f"Failed to load CSV: {e}")
        return

    # Check required columns
    required = ["timestamp", "frequency", "power", "latitude", "longitude"]
    if not all(col in df.columns for col in required):
        print(f"Missing required columns. Found: {df.columns.tolist()}")
        return

    # Convert timestamps + extract hour
    df["timestamp"] = pd.to_datetime(df["timestamp"], errors="coerce")
    df = df.dropna(subset=["timestamp"])
    df["hour"] = df["timestamp"].dt.hour

    # Predict using the model
    X = df[["frequency", "power", "latitude", "longitude", "hour"]]
    preds = model.predict(X)
    probs = model.predict_proba(X)[:, 1]

    # Add results
    df["interference"] = preds
    df["confidence"] = probs.round(2)

    # Show basic results
    print("\nDone. Predictions generated.")
    print(f"Total samples: {len(df)}")
    print(f"Interference count: {df['interference'].sum()}")
    print(df[["timestamp", "frequency", "power", "interference", "confidence"]].head(10))

    # Save to file
    out_path = path.replace(".csv", "_with_predictions.csv")
    df.to_csv(out_path, index=False)
    print(f"Output saved to: {out_path}")

    return df[["timestamp", "frequency", "power", "latitude", "longitude", "interference", "confidence"]]

# CLI usage
if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python replay_dataset.py path/to/your_data.csv")
    else:
        predict_from_csv(sys.argv[1])
