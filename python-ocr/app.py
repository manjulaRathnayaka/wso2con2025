import os
import pytesseract
import cv2
import numpy as np
from fastapi import FastAPI, File, UploadFile
from fastapi.middleware.cors import CORSMiddleware
from PIL import Image
from io import BytesIO

# Set Tesseract data path
os.environ["TESSDATA_PREFIX"] = "/usr/share/tesseract-ocr/5/tessdata"

# Initialize FastAPI app
app = FastAPI()

# Enable CORS for any origin, method, and header
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # Allow any origin
    allow_credentials=True,
    allow_methods=["*"],  # Allow all HTTP methods
    allow_headers=["*"],  # Allow all headers
)

# OCR Function
def extract_text_from_image(image: Image.Image) -> str:
    """Convert PIL image to grayscale and extract text using Tesseract."""
    try:
        # Convert image to OpenCV format
        image_cv = np.array(image)

        # Convert image to grayscale
        gray = cv2.cvtColor(image_cv, cv2.COLOR_RGB2GRAY)

        # Apply thresholding (optional)
        _, thresh = cv2.threshold(gray, 150, 255, cv2.THRESH_BINARY)

        # Convert processed image back to PIL format
        processed_image = Image.fromarray(thresh)

        # Perform OCR
        text = pytesseract.image_to_string(processed_image)

        return text.strip()
    except Exception as e:
        return f"Error: {str(e)}"

# REST API Endpoint
@app.post("/ocr/")
async def ocr_image(file: UploadFile = File(...)):
    """Receives an image file, extracts text using OCR, and returns the text."""
    try:
        # Read image into PIL format
        image = Image.open(BytesIO(await file.read()))

        # Perform OCR
        extracted_text = extract_text_from_image(image)

        return {"filename": file.filename, "text": extracted_text}
    except Exception as e:
        return {"error": str(e)}
