import cv2
from deepface import DeepFace

def get_face_embedding_from_frame(frame, model_name='Facenet'):
    # Convert BGR to RGB
    rgb_frame = cv2.cvtColor(frame, cv2.COLOR_BGR2RGB)
    try:
        embeddings = DeepFace.represent(rgb_frame, model_name=model_name, enforce_detection=True)
    except ValueError:
        print("No face detected with enforce_detection=True. Trying enforce_detection=False.")
        embeddings = DeepFace.represent(rgb_frame, model_name=model_name, enforce_detection=False)

    if embeddings:
        return embeddings[0]['embedding']
    else:
        print("No embeddings returned.")
        return None
