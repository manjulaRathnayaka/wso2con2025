from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
import pickle
import os

app = FastAPI(
    title="ML-Enhanced Bill Parser API",
    description="API that uses ML to parse bill text and extract structured data",
    version="1.1.0"
)

# Enable CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Load all models and vectorizers
fields = ['merchant', 'amount', 'date', 'category']
models = {}
vectorizers = {}

for field in fields:
    model_path = f"models/{field}_model.pkl"
    vectorizer_path = f"models/{field}_vectorizer.pkl"

    if os.path.exists(model_path) and os.path.exists(vectorizer_path):
        with open(model_path, "rb") as f:
            models[field] = pickle.load(f)
        with open(vectorizer_path, "rb") as f:
            vectorizers[field] = pickle.load(f)
    else:
        raise Exception(f"Model files for {field} not found. Run train_models.py first.")

class BillRequest(BaseModel):
    text: str

def predict_field(text: str, field: str) -> str:
    """Generic function to predict any field using corresponding model"""
    vectorized_input = vectorizers[field].transform([text])
    return models[field].predict(vectorized_input)[0]

@app.post("/process_bill")
async def process_bill(request: BillRequest):
    text = request.text

    # Use ML models to predict all fields
    predictions = {
        "merchant_name": predict_field(text, "merchant"),
        "total_amount": predict_field(text, "amount"),
        "date": predict_field(text, "date"),
        "category": predict_field(text, "category")
    }

    # Add confidence scores (optional)
    confidence_scores = {
        f"{field}_confidence": float(models[field].predict_proba(
            vectorizers[field].transform([text])
        ).max())
        for field in fields
    }

    return {
        **predictions,
        **confidence_scores,
        "raw_text": text
    }
