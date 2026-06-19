from contextlib import asynccontextmanager
from fastapi import FastAPI
from services.database import create_db_and_tables
from routers import auth, user
 
@asynccontextmanager
async def lifespan(app: FastAPI):
    create_db_and_tables()
    yield
 
app = FastAPI(
    title="API de Autenticación — Sistema de Chat",
    version="1.0.0",
    lifespan=lifespan,
)
 
app.include_router(auth.router)
app.include_router(user.router)
 
 
@app.get("/", tags=["General"])
def root():
    return {"message": "API de autenticación funcionando. Ver /docs"}