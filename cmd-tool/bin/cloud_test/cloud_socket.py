import argparse
import socketserver
import torch

from models.CloudModel import CloudModel
from utils.utils import *
from utils.cloud_utils import *

class CloudServerHandler(socketserver.BaseRequestHandler):
    def handle(self):   # 与客户端连接时调用
        print(f"Connection from {self.client_address}")

        if torch.cuda.is_available() == True:
            device = torch.device("cuda:0")
        else:
            device = torch.device("cpu")

        model = load_model(self.request, device) # 加载原始模型

        partition = pickle.loads(self.request.recv(1024))   # 接收划分点
        print(f"Partition point: {partition}")

        cloud_model = CloudModel(model=model, partition=partition)  # 划分云端模型
        cloud_model.to(device)
        cloud_model.eval()

        total_trans_latency = 0
        total_cloud_latency = 0
        while True: # 循环等待中间结果
            features, trans_latency = recv_data(self.request)
            if features == 'END':
                break
            else:
                total_trans_latency += trans_latency
                # 云端继续推理
                # warmup(cloud_model, features, device)
                outputs, cloud_lantency = inference(cloud_model, features, device)
                print(f"Cloud inference latency: {cloud_lantency:.3f} ms")
                total_cloud_latency += cloud_lantency

                predictions = torch.argmax(outputs, dim=1)  # 默认选择概率最大的预测结果
                self.request.sendall(pickle.dumps(predictions))  # 传回推理结果

        self.request.sendall(pickle.dumps(total_trans_latency))
        self.request.recv(1024) # 防止粘包
        self.request.sendall(pickle.dumps(total_cloud_latency))


if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('--ip', '-i', default="127.0.0.1", type=str)
    parser.add_argument('--port', '-p', default=9999, type=int)
    args = parser.parse_args()
    IP = args.ip
    PORT = args.port

    cloud_server = socketserver.ThreadingTCPServer((IP, PORT), CloudServerHandler)
    print(f"Server started on {IP}:{PORT}")
    cloud_server.serve_forever()