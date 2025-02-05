import cv2
import time

def capture_image_from_camera():
    cap = cv2.VideoCapture(0)
    if not cap.isOpened():
        print("Error: Could not access the webcam.")
        return None
    try:
        time.sleep(2)  # Let the camera adjust
        ret, frame = cap.read()
        if not ret:
            print("Error: Failed to capture an image from webcam.")
            return None
        return frame
    finally:
        cap.release()
        cv2.destroyAllWindows()