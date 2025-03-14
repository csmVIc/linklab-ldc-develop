import argparse
import socket
import torch

from models.EdgeModel import EdgeModel
from utils.utils import *
from utils.edge_utils import *

if __name__ == '__main__':

    # by podname+namespace
    # 在同一namepsace下，podname可以唯一确定一个pod
    parser = argparse.ArgumentParser()
    parser.add_argument('--podname','-pn',default="edge-pod",type=str,help="k8s pod name")
    parser.add_argument('--namespace','-ns',default='default',type=str,help="k8s namespace")
    parser.add_argument('--port','-p',default=99,type=int,help="k8s Port number")
    parser.add_argument('--partition','-pt',default=None,type=int)
    args = parser.parse_args()
    # 通过pod名称和namespace构建DNS名称
    POD_NAME = args.podname
    NAMESPACE = args.namespace
    PORT = args.port
    DNS_NAME = f"{POD_NAME}.{NAMESPACE}.svc.cluster.local"
    # 创建连接
    edge_client = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    edge_client.connect((DNS_NAME, PORT))
    print(f"Connect to {DNS_NAME}:{PORT}")

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
    