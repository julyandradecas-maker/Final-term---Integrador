from datetime import datetime
from sqlmodel import SQLModel

class UserCreate(SQLModel):
    username: str
    password: str

class UserResponse(SQLModel):
    id: int
    username: str
    pin: str
    is_active: bool
    created_at: datetime

class LoginResponse(SQLModel):
    access_token: str
    token_type: str = "bearer"
    user_id: int
    username: str
    pin: str

class ContactCreate(SQLModel):
    pin: str

class ContactResponse(SQLModel):
    user_id: int
    username: str
    pin: str

class MessageResponse(SQLModel):
    message: str
