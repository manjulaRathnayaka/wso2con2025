FROM ubuntu:22.04

# Install system dependencies
RUN apt update && apt install -y curl python3 python3-pip

# Install Ollama
RUN curl -fsSL https://ollama.com/install.sh | sh

# Install required Python packages
RUN pip3 install requests fastapi uvicorn

# Create a shared model directory that Ollama can access
RUN mkdir -p /opt/ollama/models && chmod -R 777 /opt/ollama

# Start Ollama and pre-download the Mistral model
RUN OLLAMA_MODELS=/opt/ollama ollama serve & sleep 5 && ollama pull mistral

# Create a non-root user with UID 10005
RUN groupadd -g 10005 appgroup && \
    useradd -m -u 10005 -g appgroup -s /bin/bash appuser

# Ensure the new user has access to the shared model directory
RUN chown -R appuser:appgroup /opt/ollama

# Switch to non-root user
USER appuser

# Set up working directory
WORKDIR /home/appuser/app

# Copy the script
COPY app.py .

# Set environment variables
ENV OLLAMA_MODELS=/opt/ollama
ENV OLLAMA_DEBUG=true
ENV OLLAMA_LOG_LEVEL=debug

# Expose API service and Ollama API
EXPOSE 8080 11434

# Start Ollama & FastAPI server
CMD ollama serve & sleep 5 && uvicorn app:app --host 0.0.0.0 --port 8080
