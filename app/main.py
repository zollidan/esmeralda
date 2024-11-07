from fastapi import FastAPI
from contextlib import asynccontextmanager
from app.endpoints import router
from app.database import create_db
from art import tprint 

@asynccontextmanager
async def lifespan(app: FastAPI):
    print("-"* 15, "LeFort ESMERALDA 2024", "-"* 15)
    tprint('ESMERALDA')
    # create_db()
    
    yield

app = FastAPI(lifespan=lifespan)
app.include_router(router=router)

