import secrets

from fastapi import APIRouter, Depends, HTTPException, status
from sqlmodel import Session, select

from models import User
from services.database import get_session
from services import crypto
from services.token import create_token
from models.schemas import (
    UserCreate, UserResponse, LoginResponse,
)

router = APIRouter(prefix="/auth", tags=["Autenticación"])

def generar_pin(session: Session) -> str:
    while True:
        pin = "".join(secrets.choice("0123456789") for _ in range(8))
        exists = session.exec(select(User).where(User.pin == pin)).first()
        if not exists:
            return pin

@router.post(
    "/register",
    response_model=UserResponse,
    status_code=status.HTTP_201_CREATED,
)
def register(user_data: UserCreate, session: Session = Depends(get_session)):
    existing = session.exec(
        select(User).where(User.username == user_data.username)
    ).first()

    if existing:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="El nombre de usuario ya está registrado",
        )
    
    encrypted_password, salt = crypto.hash_secreto(user_data.password)

    new_user = User(
        username=user_data.username,
        pin=generar_pin(session),
        hashed_password=encrypted_password,
        salt=salt,
    )
    session.add(new_user)
    session.commit()
    session.refresh(new_user)
    return new_user


@router.post("/login", response_model=LoginResponse)
def login(user_data: UserCreate, session: Session = Depends(get_session)):
    user = session.exec(
        select(User).where(User.username == user_data.username)
    ).first()

    invalid_credentials = HTTPException(
        status_code=status.HTTP_401_UNAUTHORIZED,
        detail="Usuario o contraseña incorrectos",
    )

    if not user:
        raise invalid_credentials
    if not user.is_active:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="La cuenta está desactivada",
        )

    if not crypto.verificar(user_data.password, user.hashed_password, user.salt):
        raise invalid_credentials

    token = create_token(user.id, user.username, user.pin)
    return LoginResponse(
        access_token=token,
        user_id=user.id,
        username=user.username,
        pin=user.pin,
    )
