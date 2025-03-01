import { useState } from "react";
import './OCRUploader.css';

export default function OCRUploader() {
  const [image, setImage] = useState(null);
  const [imagePreview, setImagePreview] = useState(null);
  const [loading, setLoading] = useState(false);
  const [expense, setExpense] = useState({
    text: "",
    amount: "",
    category: "",
    date: "",
    name: "",
    notes: ""
  });

  const handleImageUpload = (event) => {
    const file = event.target.files[0];
    setImage(file);
    // Create preview URL
    const previewUrl = URL.createObjectURL(file);
    setImagePreview(previewUrl);
  };

  const handleSubmit = async () => {
    if (!image) {
      alert("Please upload an image first");
      return;
    }

    setLoading(true);
    const formData = new FormData();
    formData.append("file", image);

    try {
      // Get OCR text
      const ocrResponse = await fetch("http://localhost:8000/ocr/", {
        method: "POST",
        body: formData,
      });
      const ocrData = await ocrResponse.json();

      // Get category classification and other fields
      const classifyResponse = await fetch("http://localhost:8001/classify", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ text: ocrData.text }),
      });
      const classifyData = await classifyResponse.json();

      // Update expense state with all extracted fields
      setExpense({
        ...expense,
        text: ocrData.text,
        category: classifyData.category,
        amount: classifyData.amount?.toString() || "",
        date: classifyData.date || "",
        name: classifyData.merchant || ""
      });
    } catch (error) {
      console.error("Error processing image:", error);
      alert("Failed to process image");
    }
    setLoading(false);
  };

  const handleExpenseSubmit = (e) => {
    e.preventDefault();
    // Here you would submit the complete expense data to your backend
    console.log("Submitting expense:", expense);
  };

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setExpense(prev => ({
      ...prev,
      [name]: value
    }));
  };

  return (
    <div className="expense-container">
      <div className="expense-card">
        <h2 className="expense-title">Add New Expense</h2>

        <div className="upload-section">
          <input
            type="file"
            accept="image/*"
            onChange={handleImageUpload}
            className="hidden"
            id="file-upload"
          />
          <label htmlFor="file-upload" className="upload-button">
            Choose Receipt Image
          </label>

          {imagePreview && (
            <div>
              <img
                src={imagePreview}
                alt="Receipt preview"
                className="preview-image"
              />
              <button
                className="process-button"
                onClick={handleSubmit}
                disabled={loading}
              >
                {loading ? (
                  <>
                    <span className="loading"></span>
                    Processing...
                  </>
                ) : (
                  "Extract Details"
                )}
              </button>
            </div>
          )}
        </div>

        <form onSubmit={handleExpenseSubmit}>
          <div className="form-group">
            <label className="form-label">Extracted Text</label>
            <textarea
              name="text"
              value={expense.text}
              onChange={handleInputChange}
              className="form-textarea"
              rows="4"
              readOnly
            />
          </div>

          <div className="grid-2">
            <div className="form-group">
              <label className="form-label">Amount</label>
              <div className="amount-input">
                <span className="amount-symbol">$</span>
                <input
                  type="number"
                  name="amount"
                  value={expense.amount}
                  onChange={handleInputChange}
                  className="form-input"
                  style={{ paddingLeft: '25px' }}
                  step="0.01"
                  required
                />
              </div>
            </div>

            <div className="form-group">
              <label className="form-label">Date</label>
              <input
                type="date"
                name="date"
                value={expense.date}
                onChange={handleInputChange}
                className="form-input"
                required
              />
            </div>
          </div>

          <div className="form-group">
            <label className="form-label">Category</label>
            <select
              name="category"
              value={expense.category}
              onChange={handleInputChange}
              className="form-select"
              required
            >
              <option value="">Select category</option>
              <option value="Groceries">🛒 Groceries</option>
              <option value="Restaurant">🍽️ Restaurant</option>
              <option value="Transport">🚗 Transport</option>
              <option value="Utilities">🏠 Utilities</option>
              <option value="Entertainment">🎭 Entertainment</option>
            </select>
          </div>

          <div className="form-group">
            <label className="form-label">Merchant Name</label>
            <input
              type="text"
              name="name"
              value={expense.name}
              onChange={handleInputChange}
              className="form-input"
              required
            />
          </div>

          <div className="form-group">
            <label className="form-label">Notes</label>
            <textarea
              name="notes"
              value={expense.notes}
              onChange={handleInputChange}
              className="form-textarea"
              rows="2"
            />
          </div>

          <button type="submit" className="submit-button">
            Save Expense
          </button>
        </form>
      </div>
    </div>
  );
}
