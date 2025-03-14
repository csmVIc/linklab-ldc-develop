import torchvision
from torch.utils.data import DataLoader, Subset

from utils.utils import *
import models.resnet

def load_model(conn, device):
    # 自定义模型加载方式
    weights_path = 'weights/resnet50.pth'
    # weights = recv_data(conn)   # 接收原始模型权重
    # with open(weights_path, 'wb') as f:
    #     f.write(weights)
    
    model = models.resnet.resnet50()
    model.load_state_dict(torch.load(weights_path, map_location=device))
    
    return model

def get_partition():
    # 自定义切分策略
    partition = 0
    
    return partition

def load_data():
    # 自定义dataloader
    transform = torchvision.transforms.Compose([
        torchvision.transforms.ToTensor(),
        torchvision.transforms.Normalize(mean=[0.5071, 0.4867, 0.4408], std=[0.2675, 0.2565, 0.2761]),
    ])
    cifar100_test = torchvision.datasets.CIFAR100(root='./data', train=False, download=True, transform=transform)
    dataloader = DataLoader(Subset(cifar100_test, range(100)), batch_size=1, shuffle=False)

    return dataloader