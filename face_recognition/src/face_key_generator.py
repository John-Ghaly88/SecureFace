import numpy as np
import base64
from fuzzy_extractor import FuzzyExtractor

class FaceKeyGenerator:
    def __init__(self, n=16, k=8):
        self.extractor = FuzzyExtractor(n, k)

    def embedding_to_binary(self, embedding: np.ndarray) -> bytes:
        # Convert to float32
        embedding = np.array(embedding, dtype=np.float32)
        
        # Normalize
        norm = np.linalg.norm(embedding)
        if norm > 0:
            embedding = embedding / norm

        # Binarize (sign-based)
        bits = (embedding >= 0).astype(np.uint8)
        
        length = len(bits)
        num_bytes = (length + 7) // 8
        padded_bits = np.concatenate([bits, np.zeros(num_bytes*8 - length, dtype=np.uint8)])
        byte_vals = np.packbits(padded_bits, bitorder='big')
        return byte_vals.tobytes()

    def generate_key_and_helper(self, embedding: np.ndarray):
        binary_data = self.embedding_to_binary(embedding)
        key, helper = self.extractor.generate(binary_data)
        return key, helper

    def reproduce_key(self, embedding: np.ndarray, helper: tuple):
        binary_data = self.embedding_to_binary(embedding)
        reproduced_key = self.extractor.reproduce(binary_data, helper)
        return reproduced_key

    def serialize_helper(self, helper: tuple):
        # Turn each numpy array in the tuple into a base64-encoded JSON-serializable form
        helper_serialized = []
        for arr in helper:
            encoded = base64.b64encode(arr.tobytes()).decode('utf-8')
            helper_serialized.append({
                "data": encoded,
                "shape": arr.shape,
                "dtype": str(arr.dtype)
            })
        return helper_serialized

    def deserialize_helper(self, helper_data):
        import numpy as np
        helper_arrays = []
        for h in helper_data:
            shape = tuple(h["shape"])
            dtype = np.dtype(h["dtype"])
            arr_bytes = base64.b64decode(h["data"])
            arr = np.frombuffer(arr_bytes, dtype=dtype).reshape(shape)
            helper_arrays.append(arr)
        return tuple(helper_arrays)
