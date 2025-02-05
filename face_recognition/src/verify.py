import requests
from capture import capture_image_from_camera
from embedding import get_face_embedding_from_frame
from face_key_generator import FaceKeyGenerator
import json

def verify_user(server_url="http://localhost:8080"):
    username = input("Enter username to verify: ")

    # 1. Retrieve proof and helper from server
    resp = requests.get(f"{server_url}/retrieve", params={"username": username})
    resp.raise_for_status()
    data = resp.json()
    
    helper_str = data["helper"]  # JSON-encoded string
    helper_data = json.loads(helper_str)  # parse it into a list of dicts

    fk = FaceKeyGenerator()
    helper = fk.deserialize_helper(helper_data)

    # 2. Capture new image and get embedding
    frame = capture_image_from_camera()
    if frame is None:
        print("Could not capture image from camera.")
        return
    
    embedding = get_face_embedding_from_frame(frame)
    if embedding is None:
        print("No face embedding returned. Please try again.")
        return

    # 3. Reproduce the key
    reproduced_key = fk.reproduce_key(embedding, helper)
    if reproduced_key is None:
        print("Verification failed. Could not reproduce the key. Are you the correct user?")
        return
    
    reproduced_key_hex = reproduced_key.hex()

    # 4. Send username + reproduced_key to verify
    payload = {
        "username": username,
        "key": reproduced_key_hex
    }
    verify_resp = requests.post(f"{server_url}/verify", json=payload)
    verify_resp.raise_for_status()
    print("Verification response:", verify_resp.json())

# # verify.py
# import json
# import requests
# from capture import capture_image_from_camera
# from embedding import get_face_embedding_from_frame
# from face_key_generator import FaceKeyGenerator

# def verify_user(server_url="http://localhost:8080"):
#     username = input("Enter username to verify: ")

#     # 1. Retrieve proof and helper from the server
#     resp = requests.get(f"{server_url}/retrieve", params={"username": username})
#     resp.raise_for_status()
#     data = resp.json()

#     proof_hex = data["proof"]
#     helper_str = data["helper"]  # <-- This is a JSON string, not yet a list/dict

#     # 2. Convert the JSON string into a Python list/dict
#     helper_data = json.loads(helper_str)

#     # 3. Reconstruct helper tuple and reproduce the key
#     fk = FaceKeyGenerator()
#     helper = fk.deserialize_helper(helper_data)

#     frame = capture_image_from_camera()
#     if frame is None:
#         print("Could not capture image.")
#         return

#     embedding = get_face_embedding_from_frame(frame)
#     if embedding is None:
#         print("Could not get embedding.")
#         return

#     reproduced_key = fk.reproduce_key(embedding, helper)
#     reproduced_key_hex = reproduced_key.hex()

#     # 4. Send username and reproduced_key to the verify endpoint
#     payload = {
#         "username": username,
#         "key": reproduced_key_hex
#     }
#     verify_resp = requests.post(f"{server_url}/verify", json=payload)
#     verify_resp.raise_for_status()
#     print("Verification response:", verify_resp.json())


# import requests
# from capture import capture_image_from_camera
# from embedding import get_face_embedding_from_frame
# from face_key_generator import FaceKeyGenerator

# def verify_user(server_url="http://localhost:8080"):
#     username = input("Enter username to verify: ")

#     # 1. Retrieve proof and helper from the server
#     resp = requests.get(f"{server_url}/retrieve", params={"username": username})
#     resp.raise_for_status()
#     data = resp.json()

#     proof_hex = data["proof"]        # We might not need `proof_hex` on the Python side,
#     helper_data = data["helper"]     # but we do need `helper`.

#     # 2. Reconstruct helper and reproduce key
#     fk = FaceKeyGenerator()
#     helper = fk.deserialize_helper(helper_data)

#     frame = capture_image_from_camera()
#     if frame is None:
#         print("Could not capture image.")
#         return
#     embedding = get_face_embedding_from_frame(frame)
#     if embedding is None:
#         print("Could not get embedding.")
#         return

#     reproduced_key = fk.reproduce_key(embedding, helper)
#     reproduced_key_hex = reproduced_key.hex()

#     # 3. Send username and reproduced_key to the verify endpoint
#     payload = {
#         "username": username,
#         "key": reproduced_key_hex
#     }
#     verify_resp = requests.post(f"{server_url}/verify", json=payload)
#     verify_resp.raise_for_status()
#     print("Verification response:", verify_resp.json())
