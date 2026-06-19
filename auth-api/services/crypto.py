import os
import hmac
import hashlib
import base64
from dotenv import load_dotenv
from cryptography.hazmat.primitives.ciphers import Cipher, algorithms, modes
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives.padding import PKCS7
from cryptography.hazmat.primitives.kdf.pbkdf2 import PBKDF2HMAC

load_dotenv()

PEPPER: str = os.getenv("PEPPER", "")
AES: str = os.getenv("AES", "cambiar_en_produccion")

if not PEPPER:
    print("PEPPER no definido")
if AES == "cambiar_en_produccion":
    print("AES no esta definido")

def aes_llave(salt: bytes) -> bytes:
    kdf = PBKDF2HMAC(
        algorithm=hashes.SHA256(),
        length=32,
        salt=salt,
        iterations=100_000,
    )
    return kdf.derive(AES.encode("utf-8"))

def hash_secreto(password: str) -> tuple[str, str]:
    peppered = (password + PEPPER).encode("utf-8")
    salt = os.urandom(16)
    iv = os.urandom(16)
    aes_key = aes_llave(salt)
    padder = PKCS7(128).padder()
    padded = padder.update(peppered) + padder.finalize()
    cipher = Cipher(algorithms.AES(aes_key), modes.CBC(iv))
    encryptor = cipher.encryptor()
    ciphertext = encryptor.update(padded) + encryptor.finalize()

    mac = hmac.new(
        key=aes_key,
        msg=iv + ciphertext,
        digestmod=hashlib.sha256
    ).digest()

    encrypted = base64.b64encode(iv + mac + ciphertext).decode("utf-8")
    salt_b64 = base64.b64encode(salt).decode("utf-8")

    return encrypted, salt_b64

def verificar(password: str, encrypted_password: str, salt_b64: str) -> bool:
    try:
        raw = base64.b64decode(encrypted_password.encode("utf-8"))
        salt = base64.b64decode(salt_b64.encode("utf-8"))

        iv = raw[:16]
        mac = raw[16:48]
        ciphertext = raw[48:]

        aes_key = aes_llave(salt)

        expected_mac = hmac.new(
            key=aes_key,
            msg=iv + ciphertext,
            digestmod=hashlib.sha256
        ).digest()

        if not hmac.compare_digest(mac, expected_mac):
            return False

        cipher = Cipher(algorithms.AES(aes_key), modes.CBC(iv))
        decryptor = cipher.decryptor()
        padded = decryptor.update(ciphertext) + decryptor.finalize()

        unpadder = PKCS7(128).unpadder()
        decrypted = unpadder.update(padded) + unpadder.finalize()

        peppered_input = (password + PEPPER).encode("utf-8")
        return decrypted == peppered_input

    except Exception:
        return False