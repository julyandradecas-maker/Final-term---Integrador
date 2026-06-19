import os
from sqlmodel import SQLModel, Session, create_engine

DB_DIR = "/app/data" if os.path.exists("/app/data") else "."
DATABASE_URL = f"sqlite:///{DB_DIR}/database.db"

engine = create_engine(
    DATABASE_URL,
    echo=False,
    connect_args={"check_same_thread": False},
)

def create_db_and_tables() -> None:
    import models 
    SQLModel.metadata.create_all(engine)

def get_session():
    with Session(engine) as session:
        yield session
