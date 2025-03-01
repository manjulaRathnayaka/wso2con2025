from fastapi import FastAPI
from pydantic import BaseModel
import requests
import json

app = FastAPI()

# Define request model
class RequestData(BaseModel):
    text: str

# Define Ollama server details
OLLAMA_URL = "http://localhost:11434/api/generate"

@app.post("/extract")
async def extract_info(data: RequestData):
    prompt = f"""
    Extract the following details from the provided text:
    - Total Billed Amount
    - Merchant Name
    - Category
    - Date
    In case, any of the above details are not present in the provided text, please generate a valid response, do not omit any field.
    Also, if date is missing, use current date. If merchant name is missing, use "Unknown Merchant". If category is missing, try to infer it from the text.
    There could be multiple dollar amounts in the text, consider the one which is most likely to be the total billed amount.
    Text:
    {data.text}

    Provide the result in JSON format.
    """

    # Call Ollama's API with streaming enabled
    response = requests.post(
        OLLAMA_URL,
        json={"model": "mistral", "prompt": prompt},
        stream=True  # Enables streaming response handling
    )

    # Process streamed response
    full_response = ""
    for line in response.iter_lines():
        if line:
            try:
                json_data = json.loads(line.decode("utf-8"))
                if "response" in json_data:
                    full_response += json_data["response"]
            except json.JSONDecodeError:
                continue

    return {"extracted_data": full_response}

# Run FastAPI server when executed directly
if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8080)
