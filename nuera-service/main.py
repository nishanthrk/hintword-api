import os
import json
import logging
import tempfile
from typing import Optional

import numpy as np
import torch
import whisper
import ffmpeg
import wave

from fastapi import FastAPI, WebSocket, WebSocketDisconnect, UploadFile, File, Form
from fastapi.responses import FileResponse, JSONResponse
from fastapi.staticfiles import StaticFiles
from scipy.io import wavfile

# Setup logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

app = FastAPI()

# Serve frontend
app.mount("/static", StaticFiles(directory="static"), name="static")

# Load Whisper model
model = whisper.load_model("base")
# model = whisper.load_model("turbo")

@app.get("/")
async def index():
    return FileResponse("static/index.html")

@app.websocket("/ws/transcribe")
async def websocket_endpoint(websocket: WebSocket):
    await websocket.accept()
    try:
        config = await websocket.receive_json()
        log_id = config.get("log_id")
        language = config.get("language")

        if not log_id:
            await websocket.close(code=1008, reason="log_id is required")
            return

        logger.info(f"Client connected with log_id: {log_id}")
        save_dir = os.path.join("transcriptions", log_id)
        os.makedirs(save_dir, exist_ok=True)
        chunk_counter = 0

        while True:
            try:
                data = await websocket.receive_bytes()

                # Save WebM chunk to temporary file
                with tempfile.NamedTemporaryFile(suffix=".webm", delete=False) as tmp_in:
                    tmp_in.write(data)
                    tmp_in.flush()
                    tmp_in_path = tmp_in.name

                # Save a copy of the incoming webm for debugging
                debug_webm_path = os.path.join(save_dir, f"debug_chunk_{chunk_counter+1}.webm")
                with open(debug_webm_path, "wb") as debug_f:
                    debug_f.write(data)

                tmp_out_path = tmp_in_path.replace(".webm", ".wav")

                try:
                    # Convert webm to wav
                    try:
                        ffmpeg.input(tmp_in_path).output(tmp_out_path, ar=16000, ac=1, format='wav').run(quiet=True, overwrite_output=True)
                    except ffmpeg.Error as ffmpeg_err:
                        logger.error(f"ffmpeg stderr: {ffmpeg_err.stderr.decode() if ffmpeg_err.stderr else 'No stderr'}")
                        raise

                    # Read the wav file
                    sr, audio_int16 = wavfile.read(tmp_out_path)
                    audio_float32 = audio_int16.astype(np.float32) / 32767.0

                    # Run Whisper
                    result = model.transcribe(audio_float32, language=language, fp16=torch.cuda.is_available())

                    # Save files
                    chunk_counter += 1
                    audio_save_path = os.path.join(save_dir, f"chunk_{chunk_counter}.wav")
                    text_save_path = os.path.join(save_dir, f"chunk_{chunk_counter}.txt")
                    os.rename(tmp_out_path, audio_save_path)

                    with open(text_save_path, "w") as f:
                        f.write(result["text"])

                    await websocket.send_json({
                        "log_id": log_id,
                        "text": result["text"],
                        "status": "success"
                    })

                finally:
                    os.remove(tmp_in_path)
                    if os.path.exists(tmp_out_path):
                        os.remove(tmp_out_path)

            except WebSocketDisconnect:
                logger.info(f"Client disconnected: {log_id}")
                break
            except Exception as e:
                logger.error(f"Error: {e}")
                await websocket.send_json({
                    "log_id": log_id,
                    "error": str(e),
                    "status": "error"
                })

    except Exception as e:
        logger.error(f"WebSocket init error: {e}")
        await websocket.close(code=1011, reason=str(e))

@app.post("/transcribe-chunk")
async def transcribe_chunk(audio: UploadFile = File(...), log_id: str = Form(...), language: str = Form("")):
    try:
        if not language:
            language = None
        # Save uploaded audio chunk to a temp file
        suffix = os.path.splitext(audio.filename)[-1]
        with tempfile.NamedTemporaryFile(suffix=suffix, delete=False) as tmp_in:
            tmp_in.write(await audio.read())
            tmp_in.flush()
            tmp_in_path = tmp_in.name

        tmp_out_path = tmp_in_path.replace(suffix, ".wav")
        wav_path = None
        audio_save_path = None
        try:
            # Robust: If input is a .wav file and valid, use it directly; else, convert
            def is_valid_wav(path):
                try:
                    with wave.open(path, 'rb') as wf:
                        return (
                            wf.getnchannels() == 1 and
                            wf.getsampwidth() == 2 and
                            wf.getframerate() == 16000
                        )
                except Exception as e:
                    logger.error(f"WAV validation error: {e}")
                    return False

            if suffix.lower() == ".wav" and is_valid_wav(tmp_in_path):
                logger.info(f"Using uploaded WAV directly: {tmp_in_path}")
                wav_path = tmp_in_path
            else:
                logger.info(f"Converting uploaded file to WAV: {tmp_in_path} -> {tmp_out_path}")
                try:
                    ffmpeg.input(tmp_in_path).output(tmp_out_path, ar=16000, ac=1, format='wav').run(quiet=True, overwrite_output=True)
                    wav_path = tmp_out_path
                except ffmpeg.Error as ffmpeg_err:
                    logger.error(f"ffmpeg stderr: {ffmpeg_err.stderr.decode() if ffmpeg_err.stderr else 'No stderr'}")
                    raise

            # Read the wav file
            sr, audio_int16 = wavfile.read(wav_path)
            audio_float32 = audio_int16.astype(np.float32) / 32767.0

            # Run Whisper
            result = model.transcribe(audio_float32, language=language, fp16=torch.cuda.is_available())

            # Save files (optional, for debugging)
            save_dir = os.path.join("transcriptions", log_id)
            os.makedirs(save_dir, exist_ok=True)
            chunk_counter = len([f for f in os.listdir(save_dir) if f.endswith('.wav')]) + 1
            audio_save_path = os.path.join(save_dir, f"chunk_{chunk_counter}.wav")
            text_save_path = os.path.join(save_dir, f"chunk_{chunk_counter}.txt")
            os.rename(wav_path, audio_save_path)
            with open(text_save_path, "w") as f:
                f.write(result["text"])

            return JSONResponse({
                "log_id": log_id,
                "text": result["text"],
                "status": "success"
            })
        finally:
            # Only delete temp files if they still exist and are not the moved audio_save_path
            for path in [tmp_in_path, tmp_out_path]:
                if path and os.path.exists(path):
                    if audio_save_path is None or os.path.abspath(path) != os.path.abspath(audio_save_path):
                        try:
                            os.remove(path)
                        except Exception as e:
                            logger.error(f"Error deleting temp file {path}: {e}")
    except Exception as e:
        logger.error(f"Transcribe chunk error: {e}")
        return JSONResponse({
            "log_id": log_id,
            "error": str(e),
            "status": "error"
        })

@app.post("/transcribe-file")
async def transcribe_file(audio: UploadFile = File(...), language: str = Form("")):
    try:
        if not language:
            language = None
        suffix = os.path.splitext(audio.filename)[-1]
        with tempfile.NamedTemporaryFile(suffix=suffix, delete=False) as tmp_in:
            tmp_in.write(await audio.read())
            tmp_in.flush()
            tmp_in_path = tmp_in.name
        tmp_out_path = tmp_in_path.replace(suffix, ".wav")
        try:
            try:
                ffmpeg.input(tmp_in_path).output(tmp_out_path, ar=16000, ac=1, format='wav').run(quiet=True, overwrite_output=True)
            except ffmpeg.Error as ffmpeg_err:
                logger.error(f"ffmpeg stderr: {ffmpeg_err.stderr.decode() if ffmpeg_err.stderr else 'No stderr'}")
                raise
            sr, audio_int16 = wavfile.read(tmp_out_path)
            audio_float32 = audio_int16.astype(np.float32) / 32767.0
            result = model.transcribe(audio_float32, language=language, fp16=torch.cuda.is_available())
            return JSONResponse({
                "text": result["text"],
                "status": "success"
            })
        finally:
            os.remove(tmp_in_path)
            if os.path.exists(tmp_out_path):
                os.remove(tmp_out_path)
    except Exception as e:
        logger.error(f"Transcribe file error: {e}")
        return JSONResponse({
            "error": str(e),
            "status": "error"
        })

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8001) 