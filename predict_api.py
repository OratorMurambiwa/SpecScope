from fastapi import FastAPI
from pydantic import BaseModel
import joblib
import numpy as np

# Load the trained Random Forest model from file
# Make sure rf_model.pkl is in the same folder as this script
model = joblib.load("rf_model.pkl")

# Define the expected input structure using Pydantic
class Reading(BaseModel):
    frequency: float     
    power: float         
    latitude: float      # GPS coordinate
    longitude: float     # GPS coordinate
    hour: int            

# Create the FastAPI instance
app = FastAPI(title="RF Interference Detection API")

# 🔍 Define a POST route for predictions
@app.post("/predict")
def predict(reading: Reading):
    """
    Accepts a signal reading and returns:
    - Whether interference is likely (True/False)
    - Confidence level (0 to 1)
    """
    # Convert input into a 2D array 
    X = np.array([[reading.frequency, reading.power, reading.latitude, reading.longitude, reading.hour]])

    # Run the prediction using the trained model
    pred = model.predict(X)[0]                    # 0 = no interference, 1 = interference
    proba = model.predict_proba(X)[0][1]          # Confidence score for class 1 (interference)

    # Return result as a JSON
    return {
        "interference": bool(pred),               # Convert to True/False
        "confidence": round(float(proba), 2)      # Keep 2 decimal places
    }

# Home route to check if the API is running
@app.get("/")
def root():
    return {"message": "This is the RF Prediction API. Use POST /predict."}
