import time
import datetime
import uuid
from websocket import create_connection
#ws = create_connection("ws://192.168.3.140:8080/ws")

def my_random_string(string_length=10):
    """Returns a random string of length string_length."""
    random = str(uuid.uuid4()) # Convert UUID format to a Python string.
    random = random.upper() # Make all characters uppercase.
    random = random.replace("-","") # Remove the UUID '-'.
    return random[0:string_length] # Return the random string.


def main():
    ws = create_connection("ws://192.168.3.140:8080/ws")


    msg_count = 1000

    for i in range(msg_count):
        m_data = datetime.datetime.now().strftime("%I:%M%p on %B %d, %Y")
        ts = time.time()
        msg = "Hello, World..." + str(ts) + " ==> client 1" + " (" + str(i) + ")"

        try:
            print("Sending ..." + msg)
            ws.send(msg)
        except Exception as msg_expt:
            pass

        print("Sent")

        try:
            print("Receiving...")
            result = ws.recv()
            print("Received '%s'" % result)
        except Exception as msg_expt:
            pass

        time.sleep(0.1)
        print("========================")

    ws.close()


if __name__ == "__main__":
    main()
