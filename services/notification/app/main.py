from fastapi import FastAPI
import uvicorn
from services.notification.app.services.push_service import get_push_service
from services.notification.app.services.email_service import get_email_service
from services.notification.app.services.sms_service import get_sms_service
from services.notification.app.core.logger import get_custom_logger

app = FastAPI()

@app.on_event("startuo")
async def start_up_event():
  pass

@app.get("/")
async def root():
  return {
    "message" : "Hello World!!"
  }

@app.get("/favicon.ico")
async def favicon():
  return {
    "favicon": "point"
  }
if __name__ == "__main__":
  uvicorn.run("main:app", port=8000, log_level="info", reload=True)