import pandas as pd
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.linear_model import LogisticRegression
import pickle
import os

# Load sample data
DATA_FILE = "data/category_data.csv"
MODEL_DIR = "models"

# Ensure models directory exists
if not os.path.exists(MODEL_DIR):
    os.makedirs(MODEL_DIR)

df = pd.read_csv(DATA_FILE)

# Prepare training data
X = df["text"]
y = df["category"]

vectorizer = TfidfVectorizer(max_features=1000)
X_vec = vectorizer.fit_transform(X)

# Train classifier
model = LogisticRegression()
model.fit(X_vec, y)

# Save model and vectorizer
with open(os.path.join(MODEL_DIR, "category_model.pkl"), "wb") as f:
    pickle.dump(model, f)

with open(os.path.join(MODEL_DIR, "vectorizer.pkl"), "wb") as f:
    pickle.dump(vectorizer, f)

print("✅ Model training completed. Files saved in `models/` directory.")
