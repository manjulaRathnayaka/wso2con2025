import pandas as pd
from sklearn.model_selection import train_test_split
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.ensemble import RandomForestClassifier
import pickle
import os

# Create models directory
os.makedirs("models", exist_ok=True)

# Load training data
data_files = {
    'merchant': 'data/merchant_data.csv',
    'amount': 'data/amount_data.csv',
    'date': 'data/date_data.csv',
    'category': 'data/category_data.csv'
}

models = {}
vectorizers = {}

for field, file_path in data_files.items():
    # Load and prepare data
    df = pd.read_csv(file_path)

    # Create and train vectorizer
    vectorizer = TfidfVectorizer()
    features = vectorizer.fit_transform(df['text'])

    # Create and train model
    model = RandomForestClassifier()
    target_column = f"{field}_name" if field == 'merchant' else field
    model.fit(features, df[target_column])

    # Store models and vectorizers
    models[field] = model
    vectorizers[field] = vectorizer

    # Save models
    with open(f"models/{field}_model.pkl", "wb") as f:
        pickle.dump(model, f)
    with open(f"models/{field}_vectorizer.pkl", "wb") as f:
        pickle.dump(vectorizer, f)

print("All models trained and saved successfully!")