import torch
from utils.utils import *

class EdgeModel(torch.nn.Module):
    def __init__(self, model, partition):
        super(EdgeModel, self).__init__()

        layers = get_layers(model)
        print(f"Total num of layers: {len(layers)}")

        self.front = torch.nn.Sequential(*layers[:partition])
        print("Edge model:")
        print(self.front)

    def forward(self, x):
        x = self.front(x)
        return x
