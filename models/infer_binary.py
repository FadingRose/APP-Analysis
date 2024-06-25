import tensorflow as tf
import numpy as np
import json
from sklearn.feature_extraction.text import TfidfVectorizer
from tensorflow.keras.preprocessing.text import Tokenizer
from tensorflow.keras.preprocessing.sequence import pad_sequences
import sys

# print current working directory
model_file = sys.argv[1]
target_file = sys.argv[2]
output_dir = sys.argv[3]

# Load the trained model
model = tf.keras.models.load_model(model_file)

print(tf.__version__)

# Check the model's input shapes
print("Model Input Shapes:")
for layer in model.input:
    print(f"{layer.name}: {layer.shape}")
    if layer.name == 'permissions_input':
        expected_permissions_shape = layer.shape[1]
    else:
        expected_shape = layer.shape[1:]

# Define the same max length used during training
MAX_LEN = 250

# Load and parse the JSON data for inference
with open(target_file, 'r') as f:
    input_data = json.load(f)

# Extract data from JSON
permissions_list = input_data['permissions']
package_list = input_data['package']
endpoints_list = input_data['endpoints']
unique_urls_list = input_data['unique_urls']
activities_list = input_data['Activities']
services_list = input_data['Services']
receivers_list = input_data['Receivers']
providers_list = input_data['Providers']

# Combine the lists into single strings for each attribute
permissions_combined = ' '.join(permissions_list)
package_combined = ' '.join(package_list)
endpoints_combined = ' '.join(endpoints_list)
unique_urls_combined = ' '.join(unique_urls_list)
activities_combined = ' '.join(activities_list)
services_combined = ' '.join(services_list)
receivers_combined = ' '.join(receivers_list)
providers_combined = ' '.join(providers_list)

# Fit and transform the TF-IDF Vectorizer for permissions
permissions_vectorizer = TfidfVectorizer(max_features=2500)
permissions_tfidf = permissions_vectorizer.fit_transform([permissions_combined]).toarray()



# Function to fit and transform a tokenizer and pad sequences
def text_to_padded_sequence(text_combined, max_len=MAX_LEN):
    tokenizer = Tokenizer()
    tokenizer.fit_on_texts([text_combined])
    seq = tokenizer.texts_to_sequences([text_combined])
    seq_padded = pad_sequences(seq, maxlen=max_len, padding='post')
    return seq_padded, len(tokenizer.word_index) + 1

# Fit and transform tokenizers and pad sequences for each input
endpoints_seq_padded, _ = text_to_padded_sequence(endpoints_combined)
unique_urls_seq_padded, _ = text_to_padded_sequence(unique_urls_combined)
activities_seq_padded, _ = text_to_padded_sequence(activities_combined)
services_seq_padded, _ = text_to_padded_sequence(services_combined)
receivers_seq_padded, _ = text_to_padded_sequence(receivers_combined)
providers_seq_padded, _ = text_to_padded_sequence(providers_combined)

# Reshape inputs to ensure consistent cardinality
num_samples = 1  # As we are inferring on a single input set
# Ensure TF-IDF vector matches the expected shape
permissions_tfidf = np.resize(permissions_tfidf, (num_samples,expected_permissions_shape))
endpoints_seq_padded = np.resize(endpoints_seq_padded, (num_samples, MAX_LEN))
unique_urls_seq_padded = np.resize(unique_urls_seq_padded, (num_samples, MAX_LEN))
activities_seq_padded = np.resize(activities_seq_padded, (num_samples, MAX_LEN))
services_seq_padded = np.resize(services_seq_padded, (num_samples, MAX_LEN))
receivers_seq_padded = np.resize(receivers_seq_padded, (num_samples, MAX_LEN))
providers_seq_padded = np.resize(providers_seq_padded, (num_samples, MAX_LEN))

# Prepare the processed data dictionary
processed_data = {
    'permissions_input': permissions_tfidf,
    'endpoints_input': endpoints_seq_padded,
    'unique_urls_input': unique_urls_seq_padded,
    'activities_input': activities_seq_padded,
    'services_input': services_seq_padded,
    'receivers_input': receivers_seq_padded,
    'providers_input': providers_seq_padded
}

# Perform inference
predictions = model.predict([
    processed_data['permissions_input'],
    processed_data['endpoints_input'],
    processed_data['unique_urls_input'],
    processed_data['activities_input'],
    processed_data['services_input'],
    processed_data['receivers_input'],
    processed_data['providers_input']
])

# Process and print the predictions
print("Predictions:", predictions)

# save the predictions to a file
output_file = output_dir + '/predictions_binary.txt'
np.savetxt(output_file, predictions)