from google import genai
from google.genai import types
from PIL import Image
from io import BytesIO
from dotenv import load_dotenv
import os
import base64
from pathlib import Path




def main() -> None:

    load_dotenv()

    client = genai.Client(api_key= os.getenv("GEMINI_API_KEY"))

    prompt_dir = Path(os.path.join(os.getcwd() + "/Media"))
    prompt_file = prompt_dir / "prompt.txt"
    prompt = ""
    with open(prompt_file, "r") as file:
        prompt = file.read()

    try:
        response = client.models.generate_content(
            model="gemini-2.0-flash-exp-image-generation",
            contents=prompt,
            config=types.GenerateContentConfig(
                response_modalities=['Text', 'Image']
            )
        )

    except Exception as e:
        print(f"Error has occured! {e}")

    else:
        for part in response.candidates[0].content.parts:
            if part.text is not None:
                pass
            elif part.inline_data is not None:
                output_dir = Path(os.path.join(os.getcwd() + "/Media"))
                file = output_dir / "image.png"
                image = Image.open(BytesIO((part.inline_data.data)))
                image.save(file)
                os.remove(os.path.join(os.getcwd(), "Media", "prompt.txt"))

if __name__ == "__main__":
    main()
