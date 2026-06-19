import os
from datetime import datetime, timedelta, timezone

from dotenv import load_dotenv
from jose import jwt

load_dotenv()

JWT_SECRET: str = os.getenv("JWT_SECRET", "cambiar_en_produccion")
ALGORITHM: str = "HS256"
EXPIRE_HOURS: int = 24


def create_token(user_id: int, username: str, pin: str) -> str:
    payload = {
        "user_id": str(user_id),
        "username": username,
        "pin": pin,
        "exp": datetime.now(timezone.utc) + timedelta(hours=EXPIRE_HOURS),
    }
    return jwt.encode(payload, JWT_SECRET, algorithm=ALGORITHM)


def decode_token(token: str) -> dict:
    return jwt.decode(token, JWT_SECRET, algorithms=[ALGORITHM])
