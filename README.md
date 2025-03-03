# Smart Expense Management System

A microservices-based expense management system that uses AI/ML to automatically process and categorize receipts and bills.

##### At a very high level
1. Smart Expense Frontend -> Python OCR Service (Image to Text conversion)
2. Smart Expense Frontend -> Bill Parser Service (Derive data from text input)

##### ToDo
1. Store data in a db and show historical bill information.
2. Enable login to app
## Services Overview

### 1. Bill Parser Service (Python/FastAPI)
Located in [bill-parser-service/](bill-parser-service/), this service:
- Uses ML models to extract structured data from bill text
- Predicts merchant name, amount, date, and expense category
- Built with FastAPI and scikit-learn
- Runs on port 8000

Key files:
- [main.py](bill-parser-service/main.py) - Main FastAPI application
- [train_models.py](bill-parser-service/train_models.py) - ML model training script
- [.choreo/](bill-parser-service/.choreo/) - Choreo deployment configuration

### 2. Python OCR Service
Located in [python-ocr/](python-ocr/), this service:
- Extracts text from receipt/bill images using Tesseract OCR
- Provides REST API for image-to-text conversion
- Built with FastAPI and pytesseract
- Runs on port 8000

Key files:
- [app.py](python-ocr/app.py) - FastAPI application with OCR logic
- [Dockerfile](python-ocr/Dockerfile) - Container configuration with Tesseract


### 3. Smart Expense Frontend
Located in [smartexpense-app/](smartexpense-app/), this is:
- React-based web application
- Provides UI for receipt uploads and expense management
- Built with React and Tailwind CSS
- Runs on port 3000

Key files:
- [src/OCRUploader.js](smartexpense-app/src/OCRUploader.js) - Main receipt upload component
- [src/App.js](smartexpense-app/src/App.js) - Root React component

## Getting Started

Each service can be run independently using Docker:

```bash
# Build and run bill parser service
cd bill-parser-service
docker build -t bill-parser .
docker run -p 8000:8000 bill-parser

# Build and run OCR service
cd python-ocr
docker build -t ocr-service .
docker run -p 8000:8000 ocr-service

# Run frontend application
cd smartexpense-app
npm install
npm start