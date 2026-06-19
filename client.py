#!/usr/bin/env python3
import json
import sys
import threading
import urllib.error
import urllib.request

try:
    import websocket 
except ImportError:
    print("Instalar dependencia con pip install websocket-client")
    sys.exit(1)


def pedir(texto, defecto=""):
    valor = input(f"{texto} " + (f"[{defecto}] " if defecto else "")).strip()
    return valor or defecto

def http_json(url, body, token=None):
    data = json.dumps(body).encode() if body is not None else None
    req = urllib.request.Request(url, data=data, method="POST" if data else "GET")
    if data:
        req.add_header("Content-Type", "application/json")
    if token:
        req.add_header("Authorization", "Bearer " + token)
    try:
        with urllib.request.urlopen(req) as r:
            return r.status, json.loads(r.read() or "null")
    except urllib.error.HTTPError as e:
        try:
            return e.code, json.loads(e.read() or "null")
        except Exception:
            return e.code, None

def login_o_registro(api_base):
    while True:
        accion = pedir("¿(l)ogin o (r)egistro?", "l").lower()
        username = pedir("Usuario:")
        password = pedir("Contraseña:")

        if accion.startswith("r"):
            st, resp = http_json(api_base + "/auth/register",
                                 {"username": username, "password": password})
            if st == 201:
                print(f"Registrado. Tu PIN es {resp['pin']} ( tienes que compartirlo para que te agreguen como un numero telefónico).")
            else:
                print("Error en registro:", resp)
                continue

        st, resp = http_json(api_base + "/auth/login",
                             {"username": username, "password": password})
        if st == 200:
            return resp["access_token"], resp["pin"], resp["username"]
        print("Login fallido:", resp)

def main():
    print("Chat 1 a 1 (terminal)")
    host = ""
    while not host:
        host = pedir("IP del servidor:")
        if not host:
            print("Debes escribir la IP del servidor por ejemplo: 192.168.x.x")
    api_port = pedir("Puerto FastAPI:", "8000")
    ws_port = pedir("Puerto backend (WebSocket):", "8080")

    api_base = f"http://{host}:{api_port}"
    token, mi_pin, username = login_o_registro(api_base)

    print(f"\nHola {username}. Tu PIN es: {mi_pin}")
    print("Para agregar contactos por su PIN usa el FastAPI (POST /users/contacts).")
    destino = pedir("\n¿A qué PIN quieres escribir?")

    ws_url = f"ws://{host}:{ws_port}/ws?token={token}"
    ws = websocket.create_connection(ws_url)
    print(f"Conectado. Escribiendo a {destino}.")
    print("Comandos:  /to <pin>  cambia de destinatario   |   /salir  termina\n")

    def recibir():
        while True:
            try:
                raw = ws.recv()
            except Exception:
                break
            if not raw:
                break
            try:
                m = json.loads(raw)
                print(f"\n[{m['from_user']} · {m['from_pin']}] {m['content']}")
                print("> ", end="", flush=True)
            except Exception:
                pass

    threading.Thread(target=recibir, daemon=True).start()

    try:
        while True:
            linea = input("> ").strip()
            if not linea:
                continue
            if linea == "/salir":
                break
            if linea.startswith("/to "):
                destino = linea[4:].strip()
                print(f"Ahora escribes a {destino}.")
                continue
            ws.send(json.dumps({"to": destino, "content": linea}))
    except (KeyboardInterrupt, EOFError):
        pass
    finally:
        ws.close()
        print("\nDesconectado.")

if __name__ == "__main__":
    main()
