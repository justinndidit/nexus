from fastapi import FastAPI
import uvicorn


app = FastAPI()

@app.on_event("startup")
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