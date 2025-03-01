# Use a lightweight Python image
FROM python:3.9-slim

# Install Tesseract OCR and dependencies
RUN apt-get update && apt-get install -y \
    tesseract-ocr \
    libtesseract-dev \
    tesseract-ocr-eng \
    && apt-get clean

# Set correct Tesseract data path
ENV TESSDATA_PREFIX=/usr/share/tesseract-ocr/5/tessdata

# Create a non-root user with UID > 10000
RUN groupadd -g 10100 ocrgroup && \
    useradd -m -u 10100 -g ocrgroup ocruser

# Set working directory
WORKDIR /app

# Copy requirements and install dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy application files
COPY . .

# Change ownership of the application directory
RUN chown -R ocruser:ocrgroup /app

# Switch to the non-root user
USER ocruser

# Expose port 8000
EXPOSE 8000

# Run FastAPI app with Uvicorn
CMD ["uvicorn", "app:app", "--host", "0.0.0.0", "--port", "8000"]
