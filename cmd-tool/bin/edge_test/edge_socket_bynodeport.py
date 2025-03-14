import argparse
import socket
import torch

from models.EdgeModel import EdgeModel
from utils.utils import *
from utils.edge_utils import *

if __name__ == '__main__':

    parser = argparse.ArgumentParser()
    parser.add_argument('--ip','-i',default="127.0.0.1",type=str,help="Ingress URL or IP")
    parser.add_argument('--port','-p',default=9999,type=int)
    parser.add_argument('--partition','-pt',default=None,type=int)
    args = parser.parse_args()
    IP = args.ip
    PORT = args.port
    # 创建连接
    edge_client = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    edge_client.connect((IP, PORT))
    print(f"Connect to {IP}:{PORT}")

    if torch.cuda.is_available() == True:
        device = torch.device("cuda:0")
    else:
        device = torch.device("cpu")

    model = load_model(edge_client, device) # 获取原始模型

    if args.partition is None:
        partition = get_partition()
    else:
        partition = args.partition
    print(f"Partition point: {partition}")
    edge_client.sendall(pickle.dumps(partition))    # 发送划分点
    
    edge_model = EdgeModel(model=model, partition=partition)    # 划分边缘模型
    edge_model.to(device)
    edge_model.eval()

    dataloader = load_data()
    total_samples = 0
    correct_samples = 0
    total_edge_latency = 0
    # 边缘推理过程
    for batch_idx, (inputs, labels) in enumerate(dataloader):
        # warmup(edge_model, inputs, device)
        features, edge_latency = inference(edge_model, inputs, device)
        print(f"Edge inference latency: {edge_latency:.3f} ms")
        total_edge_latency += edge_latency
        send_data(edge_client, features)    # 发送中间结果

        predictions = pickle.loads(edge_client.recv(1024)).to(device)    # 等待云端传回推理结果
        labels = labels.to(device)
        for i in range(len(labels)):
            print(f"Prediction: {predictions[i]}, Ground Truth: {labels[i]}")
        correct_samples += torch.sum(predictions == labels).item()
        total_samples += len(labels)
    
    send_data(edge_client, 'END')

    acc = correct_samples / total_samples * 100
    print(f"Accuracy: {acc:.2f}%")

    avg_edge_latency = total_edge_latency / total_samples
    print(f"Avg edge inference latency: {avg_edge_latency:.3f} ms")

    total_trans_latency = pickle.loads(edge_client.recv(1024))
    avg_trans_latency = total_trans_latency / total_samples
    print(f"Avg transmission latency: {avg_trans_latency:.3f} ms")

    edge_client.sendall('avoid sticky'.encode())

    total_cloud_latency = pickle.loads(edge_client.recv(1024))
    avg_cloud_latency = total_cloud_latency / total_samples
    print(f"Avg cloud inference latency: {avg_cloud_latency:.3f} ms")

    edge_client.close()
    