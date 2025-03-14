import torch

from utils.utils import *
import models.resnet

def load_model(conn, device):
    # 自定义模型加载方式
    weights_path = 'weights/resnet50.pth'
    model = models.resnet.resnet50()
    model.load_state_dict(torch.load(weights_path, map_location=device))

    # 发送模型权重到边缘
    # with open(weights_path, 'rb') as f:
    #     weights = f.read()
    # send_data(conn, weights)

    return model