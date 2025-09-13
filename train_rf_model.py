import pandas as pd
from sklearn.ensemble import RandomForestClassifier
from sklearn.model_selection import train_test_split
from sklearn.metrics import classification_report
import joblib
import json

# Step 1: Load the JSON spectrum data
with open("rf_data.json", "r") as f:
    data = json.load(f)

# Step 2: Convert to pandas DataFrame
df = pd.DataFrame(data)

# Step 3: Convert timestamp to datetime and extract hour
df["timestamp"] = pd.to_datetime(df["timestamp"])
df["hour"] = df["timestamp"].dt.hour

# Step 4: Engineer 'time_period' from hour
def get_time_period(hour):
    if 6 <= hour < 12:
        return "morning"
    elif 12 <= hour < 18:
        return "afternoon"
    elif 18 <= hour < 24:
        return "evening"
    else:
        return "night"

df["time_period"] = df["hour"].apply(get_time_period)

# Step 5: Add 'wifi_proximity' — 1 if near 2400–2485 MHz, else 0
df["wifi_proximity"] = df["frequency"].apply(lambda f: 1 if 2400 <= f <= 2485 else 0)

# Step 6: Create binary target — 1 if power > 30 dBm, else 0
df["interference"] = (df["power"] > 30).astype(int)

# Step 7: Convert categorical 'time_period' to numerical (one-hot encoding)
df = pd.get_dummies(df, columns=["time_period"])

# Step 8: Define features and target
feature_cols = ["frequency", "power", "latitude", "longitude", "hour", "wifi_proximity"] + \
               [col for col in df.columns if col.startswith("time_period_")]
X = df[feature_cols]
y = df["interference"]

# Step 9: Split data into train/test
X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.2, random_state=42)

# Step 10: Train Random Forest model
clf = RandomForestClassifier(n_estimators=100, random_state=42)
clf.fit(X_train, y_train)

# Step 11: Evaluate model
y_pred = clf.predict(X_test)
print("\n📊 Classification Report:\n")
print(classification_report(y_test, y_pred))

# Step 12: Save trained model to disk
joblib.dump(clf, "rf_model.pkl")
print("\n Model saved to rf_model.pkl")
