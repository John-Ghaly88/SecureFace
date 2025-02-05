import requests
from capture import capture_image_from_camera
from embedding import get_face_embedding_from_frame
from face_key_generator import FaceKeyGenerator

def enroll_user(server_url="http://localhost:8080"):
    username = input("Enter username to enroll: ")

    frame = capture_image_from_camera()
    if frame is None:
        print("Could not capture image.")
        return

    embedding = get_face_embedding_from_frame(frame)
    if embedding is None:
        print("Could not get embedding.")
        return

    fk = FaceKeyGenerator()
    key, helper = fk.generate_key_and_helper(embedding)
    helper_serialized = fk.serialize_helper(helper)

    payload = {
        "username": username,
        "key": key.hex(),
        "helper": helper_serialized
    }

    response = requests.post(f"{server_url}/enroll", json=payload)
    response.raise_for_status()
    print("Enroll response:", response.json())
