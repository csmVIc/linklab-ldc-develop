import socket
import pickle
import torch
import time

def send_data(conn, data):
    data = pickle.dumps(data)
    data_size = len(data)
    conn.sendall(pickle.dumps(data_size))    # 发送数据长度
    conn.recv(1024).decode()

    conn.sendall(data)

def recv_data(conn):
    data_size = pickle.loads(conn.recv(1024))    # 接收数据长度
    print(f"Receive data size: {data_size} bytes")
    conn.sendall('ready'.encode())
    
    data = b''
    trans_latency = 0   # 传输时延
    while len(data) < data_size:
        start_time = time.perf_counter()
        chunk = conn.recv(4096)
        end_time = time.perf_counter()
        if not chunk:
            break
        trans_latency += (end_time - start_time) * 1000
        data += chunk
    print(f"Transmission latency: {trans_latency:.3f} ms")
    return pickle.loads(data), trans_latency

# 递归遍历所有层
def get_layers(model):
    layers = []
    for name, module in model.named_children():
        if isinstance(module, torch.nn.Sequential): # 只能展开Sequential层
            layers.extend(get_layers(module))
        else:
            layers.append(module)
    return layers

# GPU预热
def warmup(model, inputs, device, epoch=10):
    dummy_inputs = torch.rand(inputs.shape).to(device)
    with torch.no_grad():
        for i in range(epoch):
            _ = model(dummy_inputs)

def inference(model, inputs, device):
    inputs = inputs.to(device)
    # starter = torch.cuda.Event(enable_timing=True)
    # ender = torch.cuda.Event(enable_timing=True)
    with torch.no_grad():
        # starter.record()
        outputs = model(inputs)
        # ender.record()
    # 同步GPU时间
    # torch.cuda.synchronize()
    # inference_latency = starter.elapsed_time(ender) # 推理时延
    inference_latency = 0
    return outputs, inference_latency