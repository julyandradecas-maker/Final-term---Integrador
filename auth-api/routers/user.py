from fastapi import APIRouter, Depends, HTTPException, status
from sqlmodel import Session, select

from models import User, Contact
from services.database import get_session
from services.deps import get_current_user
from models.schemas import (
    UserResponse, ContactCreate, ContactResponse, MessageResponse,
)

router = APIRouter(prefix="/users", tags=["Usuarios"])

@router.get("/me", response_model=UserResponse)
def my_profile(current: User = Depends(get_current_user)):
    return current

@router.post(
    "/contacts",
    response_model=ContactResponse,
    status_code=status.HTTP_201_CREATED,
)
def add_contact(
    body: ContactCreate,
    current: User = Depends(get_current_user),
    session: Session = Depends(get_session),
):
    if body.pin == current.pin:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="No puedes agregarte a ti mismo",
        )

    target = session.exec(select(User).where(User.pin == body.pin)).first()
    if not target or not target.is_active:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="No existe ninguna persona con ese PIN",
        )

    already = session.exec(
        select(Contact).where(
            Contact.owner_id == current.id,
            Contact.contact_id == target.id,
        )

    ).first()
    if not already:
        session.add(Contact(owner_id=current.id, contact_id=target.id))
        session.commit()

    return ContactResponse(
        user_id=target.id,
        username=target.username,
        pin=target.pin,
    )

@router.get("/contacts", response_model=list[ContactResponse])
def list_contacts(
    current: User = Depends(get_current_user),
    session: Session = Depends(get_session),
):
    rows = session.exec(
        select(User)
        .join(Contact, Contact.contact_id == User.id)
        .where(Contact.owner_id == current.id)
    ).all()

    return [
        ContactResponse(user_id=u.id, username=u.username, pin=u.pin)
        for u in rows
    ]

@router.delete("/contacts/{pin}", response_model=MessageResponse)
def remove_contact(
    pin: str,
    current: User = Depends(get_current_user),
    session: Session = Depends(get_session),
):
    
    target = session.exec(select(User).where(User.pin == pin)).first()
    if not target:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="No existe ninguna persona con ese PIN",
        )

    contact = session.exec(
        select(Contact).where(
            Contact.owner_id == current.id,
            Contact.contact_id == target.id,
        )
    ).first()
    if not contact:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Ese contacto no está en tu agenda",
        )

    session.delete(contact)
    session.commit()
    return {"message": f"Contacto '{target.username}' eliminado."}
