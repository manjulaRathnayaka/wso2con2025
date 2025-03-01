FROM python:3.9-slim

# Create a non-root user
RUN groupadd -g 10001 appgroup && \
    useradd -u 10001 -g appgroup -ms /bin/bash appuser

WORKDIR /app

# Copy requirements and install dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy application files
COPY . .

# Train models and ensure correct permissions
RUN python train_models.py && \
    chown -R appuser:appgroup /app

# Switch to non-root user
USER 10001

# Expose port
EXPOSE 8000

# Run FastAPI app
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]
