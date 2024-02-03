import socket

def client():
    host = "localhost"
    port = 8080
    client_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    client_socket.connect((host, port))

    message = input("Enter the command: ").strip()  # Remove trailing newline character
    client_socket.send(message.encode())

    response = client_socket.recv(1024).decode()
    print("Received:", response)

    client_socket.close()

if __name__ == "__main__":
    client()