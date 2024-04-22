from pytube import YouTube
import sys

def download_video(url, path):
    try:
        print(f"Downloading video from URL: {url} to path: {path}")
        yt = YouTube(url)
        stream = yt.streams.get_highest_resolution()
        stream.download(output_path=path)
        print("Download successful")
    except Exception as e:
        print(f"An error occurred: {e}")
        raise  # Optionally re-raise the exception if you want to handle it further

if __name__ == "__main__":
    download_video(sys.argv[1], sys.argv[2])
