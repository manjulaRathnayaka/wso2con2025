from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
import re
import pickle
import os

app = FastAPI(
    title="Bill Parser API",
    description="API that parses bill text to extract structured data",
    version="1.0.0"
)

# Enable CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Load pre-trained category model and vectorizer
MODEL_PATH = os.getenv("MODEL_PATH", "models/category_model.pkl")
VECTORIZER_PATH = os.getenv("VECTORIZER_PATH", "models/vectorizer.pkl")

if os.path.exists(MODEL_PATH) and os.path.exists(VECTORIZER_PATH):
    with open(MODEL_PATH, "rb") as f:
        category_model = pickle.load(f)
    with open(VECTORIZER_PATH, "rb") as f:
        vectorizer = pickle.load(f)
else:
    raise Exception("Model files not found. Run `train_model.py` first.")

# Pydantic model for request validation
class BillRequest(BaseModel):
    text: str

def extract_amount(text):
    match = re.search(r"(?i)(total|amount due|balance|payable)\s*[:\-\s]?\s*([\$€£]?\d+[.,]?\d*)", text)
    return match.group(2) if match else None

def extract_date(text):
    match = re.search(r"\b(\d{1,2}[\/\-]\d{1,2}[\/\-]\d{2,4})\b", text)
    return match.group(1) if match else None

def extract_merchant(text):
    lines = text.split("\n")
    return lines[0].strip() if lines else "Unknown Merchant"

def predict_category(merchant, text):
    input_text = merchant + " " + text
    vectorized_input = vectorizer.transform([input_text])
    return category_model.predict(vectorized_input)[0]

@app.post("/process_bill")
async def process_bill(request: BillRequest):
    text = request.text  # Extract text from JSON body
    merchant = extract_merchant(text)
    amount = extract_amount(text)
    date = extract_date(text)
    category = predict_category(merchant, text)

    return {
        "merchant_name": merchant,
        "total_amount": amount,
        "date": date,
        "category": category
    }
