<<<<<<< HEAD
# Sistema de gestión de chat 1 a 1

Servidor de mensajería persona‑a‑persona pensado para operar dentro de una red:
una máquina se conecta desde cualquier punto, se loguea contra el FastAPI y
conversa con otra persona desde su terminal (cmd, PowerShell, bash...).

No hay frontend web. Cada persona se identifica con un **PIN único de 8 dígitos**
(como un número telefónico); para escribirle a alguien necesitas su PIN.

- **auth-api** (Python / FastAPI): registro, login, perfil y agenda de
  contactos. Guarda usuarios en SQLite y firma un JWT que incluye el PIN.
- **backend** (Go / Gin + Gorilla WebSocket): chat en tiempo real 1 a 1.
  Enruta cada mensaje al PIN destino; si la persona está desconectada, lo
  guarda en Redis y se lo entrega cuando vuelve.
- **redis**: sesiones activas y cola de mensajes pendientes (offline).

## Cómo levantarlo

```bash
docker compose up --build
```

- FastAPI: http://localhost:8000  (documentación interactiva en `/docs`)
- Backend WebSocket: ws://localhost:8080/ws

Las variables (`PEPPER`, `JWT_SECRET`, `AES`) se leen de `.env`.

## Flujo de uso

1. **Registro / login** (FastAPI). El registro devuelve tu PIN único:

   ```bash
   curl -X POST http://localhost:8000/auth/register \
        -H "Content-Type: application/json" \
        -d '{"username":"ana","password":"secret"}'
   ```

2. **Agregar a alguien por su PIN**:

   ```bash
   curl -X POST http://localhost:8000/users/contacts \
        -H "Authorization: Bearer <TOKEN>" \
        -H "Content-Type: application/json" \
        -d '{"pin":"25076009"}'
   ```

3. **Chatear desde la terminal**:

   ```bash
   pip install websocket-client
   python client.py
   ```

   El cliente te loguea, te muestra tu PIN y te pide el PIN del destinatario.
   Comandos dentro del chat: `/to <pin>` cambia de destinatario, `/salir` termina.

## Endpoints del FastAPI

| Método | Ruta                    | Descripción                              |
|--------|-------------------------|------------------------------------------|
| POST   | `/auth/register`        | Crea usuario y genera su PIN único       |
| POST   | `/auth/login`           | Devuelve JWT + PIN                        |
| GET    | `/users/me`             | Perfil del usuario logueado              |
| POST   | `/users/contacts`       | Agrega un contacto por su PIN            |
| GET    | `/users/contacts`       | Lista tus contactos                      |
| DELETE | `/users/contacts/{pin}` | Quita un contacto                        |

## Protocolo WebSocket

Conexión: `GET ws://<host>:8080/ws?token=<jwt>`

Mensaje que envía el cliente:

```json
{ "to": "25076009", "content": "hola" }
```

Mensaje que recibe el destinatario:

```json
{
  "id": "uuid",
  "from_pin": "67170182",
  "from_user": "ana",
  "to_pin": "25076009",
  "content": "hola",
  "created_at": "2026-06-15T18:00:00Z"
}
```
=======
# Final-term---Integrador
>>>>>>> 2f10add72bf9675732e4605e0fa429cffe8dc7b7
