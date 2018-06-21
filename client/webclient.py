import signal
import sys

import websocket
try:
    import thread
except ImportError:
    import _thread as thread
import time

ws = None

def signal_handler(signal, frame):
    print('You pressed Ctrl+C!')
    #ws.keep_running = False
    #ws.close()
    #sys.exit(1)


def open_ws():
    try:
        ws = websocket.WebSocketApp("ws://192.168.3.140:8080/ws",
                                on_message = on_message,
                                on_error = on_error,
                                on_close = on_close)

        ws.on_open = on_open
        ws.run_forever(ping_interval=60, ping_timeout=10)
    except Exception as expt_msg:
        print(expt_msg)

def on_message(ws, message):
    print(message)

def on_error(ws, error):
    print(error)

def on_close(ws):
    print("### closed ###")
    print("### trying... ###")

    print("re-connecting to server...")
    open_ws()



def on_open(ws):
    def run(*args):
        for i in range(300000):
            time.sleep(1)
            ws.send("Hello %d" % i)
        time.sleep(1)
        ws.close()
        print("thread terminating...")
    thread.start_new_thread(run, ())



if __name__ == "__main__":
    signal.signal(signal.SIGINT, signal_handler)
    websocket.enableTrace(True)

    # ws = websocket.WebSocketApp("ws://192.168.3.140:8080/ws",
    #                           on_message = on_message,
    #                           on_error = on_error,
    #                           on_close = on_close)
    # ws.on_open = on_open
    print("connecting to server...")
    open_ws()
    #ws.run_forever(ping_interval=60, ping_timeout=10)
    #ws.run_forever(ping_interval=60, ping_timeout=10)
