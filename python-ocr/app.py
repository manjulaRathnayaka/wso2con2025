import os
import pytesseract
import cv2
import numpy as np
from fastapi import FastAPI, File, UploadFile, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from starlette.responses import JSONResponse
from fastapi.responses import Response
from PIL import Image
from io import BytesIO
import logging

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# Set Tesseract data path
os.environ["TESSDATA_PREFIX"] = "/usr/share/tesseract-ocr/5/tessdata"

# Initialize FastAPI app
app = FastAPI()

# Enable CORS with more explicit configuration
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # Consider restricting this in production
    allow_credentials=True,
    allow_methods=["GET", "POST", "OPTIONS"],  # Explicitly list methods
    allow_headers=["*"],  # Allow all headers to be safe
    expose_headers=["Content-Type", "Authorization"],
    max_age=600,  # Cache preflight responses for 10 minutes
)

def extract_text_from_image(image):
    """Extract text from an image using Tesseract OCR"""
    try:
        # Convert PIL image to OpenCV format
        img = np.array(image)

        # Convert to grayscale if the image is in color
        if len(img.shape) > 2:
            img = cv2.cvtColor(img, cv2.COLOR_RGB2GRAY)

        # Apply adaptive thresholding
        img = cv2.adaptiveThreshold(
            img, 255, cv2.ADAPTIVE_THRESH_GAUSSIAN_C, cv2.THRESH_BINARY, 11, 2
        )

        # Use pytesseract to extract text
        text = pytesseract.image_to_string(img)
        return text
    except Exception as e:
        logger.error(f"Error in OCR processing: {str(e)}")
        raise HTTPException(status_code=500, detail=f"OCR processing error: {str(e)}")

@app.post("/ocr")
async def ocr_image(file: UploadFile = File(...)):
    """Receives an image file, extracts text using OCR, and returns the text."""
    try:
        logger.info(f"Processing file: {file.filename}")

        # Read image into PIL format
        contents = await file.read()
        image = Image.open(BytesIO(contents))

        # Perform OCR
        extracted_text = extract_text_from_image(image)

        # Log success and return result
        logger.info(f"Successfully extracted text from {file.filename}")
        return JSONResponse(content={"filename": file.filename, "text": extracted_text})
    except Exception as e:
        logger.error(f"Error processing image: {str(e)}")
        return JSONResponse(
            status_code=500,
            content={"error": f"Failed to process image: {str(e)}"}
        )

@app.options("/ocr")
async def options_ocr():
    """Handle preflight OPTIONS request for CORS"""
    return Response(status_code=200)

@app.get("/health")
async def health():
    """Health check endpoint"""
    return {"status": "ok", "version": "1.0.0"}
