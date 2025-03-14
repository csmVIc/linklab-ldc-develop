import torch
from utils.utils import *

class CloudModel(torch.nn.Module):
    def __init__(self, model, partition):
        super(CloudModel, self).__init__()

        layers = get_layers(model)
        print(f"Total num of layers: {len(layers)}")

        self.back = torch.nn.Sequential(*layers[partition:])
        print("Cloud model:")
        print(self.back)

    def forward(self, x):
        x = self.back(x)
        return x
