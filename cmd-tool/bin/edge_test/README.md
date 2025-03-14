# Cloud-Edge Collaborative Inference

本项目基于Socket通信实现云端主机和边缘设备的协同推理。将深度神经网络（DNN）线性地分割成两部分，分别部署在云端和边缘设备上。在边缘上进行前半部分的推理，将中间结果传到云端进行后半部分推理，得到最终结果。

目前提供了[ResNet-50](https://github.com/KaimingHe/deep-residual-networks)模型在[CIFAR-100](https://www.cs.toronto.edu/~kriz/cifar.html)数据集上的分类任务进行实验。

## 0 运行环境

本项目在Ubuntu 20.04服务器和Jetson Nano 4GB（Ubuntu 18.04）上进行了测试。服务器主机可以使用WSL2 + CUDA + Pytorch，具体可以参考[这篇文章](https://zhuanlan.zhihu.com/p/683058297)。Jetson的环境配置参考以下步骤：

(1) Jetson上需要安装GPU版的Pytorch。一般Jetson已经原装了CUDA，可以验证CUDA版本：

```bash
nvcc -V
```

如果没有看到正常输出，需要添加环境变量。在`~/.bashrc`末尾添加：

```bash
export CUDA_HOME=$CUDA_HOME:/usr/local/cuda
export LD_LIBRARY_PATH=/usr/local/cuda/lib64:$LD_LIBRARY_PATH
export PATH=/usr/local/cuda/bin:$PATH
```

保存后执行`source ~/.bashrc`使其生效。再输入`nvcc -V`，看到如下输出说明正常：

> nvcc: NVIDIA (R) Cuda compiler driver  
> Copyright (c) 2005-2021 NVIDIA Corporation  
> Built on Sun_Feb_28_22:34:44_PST_2021  
> Cuda compilation tools, release 10.2, V10.2.300  
> Build cuda_10.2_r440.TC440_70.29663091_0  

(2) 由于Jetson的架构是aarch64，不支持Anaconda，需要安装Archiconda：

```
wget https://github.com/Archiconda/build-tools/releases/download/0.2.3/Archiconda3-0.2.3-Linux-aarch64.sh
bash Archiconda3-0.2.3-Linux-aarch64.sh
```

修改环境变量，在`~/.bashrc`末尾添加

```bash
export PATH=~/archiconda3/bin:$PATH
```

之后就可以用`conda`命令管理虚拟环境。

(3) 在Jestson上安装GPU版的Pytorch需要去NVIDIA官网[PyTorch for Jetson](https://forums.developer.nvidia.com/t/pytorch-for-jetson/72048)查找对应版本。

![pytorch_for_jetson](figs/pytorch_for_jetson.png)

可以用`jtop`工具查看JetPack版本：

```bash
sudo pip3 install jetson-stats
jtop
```

`jtop`不仅能输出Jetpack版本，还能动态显示设备的使用信息。根据JetPack版本查找对应的Pytorch和Python版本，比如JetPack 4.6.5对应的PyTorch v1.10.0和Python 3.6。下载`torch-1.10.0-cp36-cp36m-linux_aarch64.whl`。

![pytorch_v1.10.0](figs/pytorch_v1.10.0.png)

创建一个conda环境：

```bash
conda create -n env_name python=3.6
conda activate env_name
```

运行以下命令安装Pytorch：

```bash
sudo apt-get install libopenblas-base libopenmpi-dev libomp-dev
pip3 install 'Cython<3'
pip3 install numpy torch-1.10.0-cp36-cp36m-linux_aarch64.whl
```

官网拉到下方在Installation下找到对应的torchvision版本：

![torchvision](figs/torchvision.png)

运行以下命令安装torchvision：

```bash
sudo apt-get install libjpeg-dev zlib1g-dev libpython3-dev libopenblas-dev libavcodec-dev libavformat-dev libswscale-dev
git clone --branch v0.11.1 https://github.com/pytorch/vision torchvision
cd torchvision
export BUILD_VERSION=0.11.1
python3 setup.py install --user
```

验证安装是否成功：

```bash
python3
>>> import torch, torchvision
>>> print(torch.__version__)
>>> print(torch.cuda.is_available())
>>> print(torchvision.__version__)
```

## 1 代码准备

配置好环境后，将仓库分别clone到云端主机和边缘设备：

```bash
git clone https://gitee.com/zhang-yang24/cloud-edge-collaborative-inference.git
```

### 2.1 云端

#### 模型加载

在`utils/cloud_utils.py`的`load_model`函数下自定义云端模型加载的方法。已经提供了在CIFAR-100数据集上训练好的ResNet-50模型权重文件和结构（[来源](https://github.com/weiaicunzai/pytorch-cifar100)），可以直接加载。

```python
def load_model(conn):
    # 自定义模型加载方式
    weights_path = 'weights/resnet50-best.pth'
    model = models.resnet.resnet50()
    model.load_state_dict(torch.load(weights_path))

    # 发送模型权重到边缘
    with open(weights_path, 'rb') as f:
        weights = f.read()
    send_data(conn, weights)

    return model
```

`send_data(conn, weights)`将权重文件发送到边缘，以模拟部署过程。如果边缘上已有权重文件，可以省略这一步。

### 2.2 边缘

#### 模型加载

在`utils/edge_utils.py`的`load_model`函数下自定义边缘模型加载的方法，与云端类似。可以从云端接收权重文件来加载模型，但如果传输的.pt/.pth文件只包含权重而不包含结构，还需要边缘上有相同的模型结构脚本来加载模型。

#### 划分策略

在`utils/edge_utils.py`的`get_partition`函数下自定义确定划分点的策略，返回要划分的层数`partition`。模型将会在指定层数之后被分割成前后两个部分，分别部署在两端。

```python
def get_partition():
    # 自定义切分策略
    # 如根据网络带宽、模型结构、延迟要求等选择最优的划分点
    # 也可以简单返回指定划分点
    partition = 0
    return partition
```

目前只支持链式展开`torch.nn.Sequential`类模型，并且`forward`函数只能按顺序逐层进行模型中定义好的操作。

#### 数据集加载

在`utils/edge_utils.py`的`load_data`函数下自定义数据集的加载和预处理过程。已经提供了torchvision的CIFAR-100数据集的加载。

```python
def load_data():
    # 自定义dataloader
    transform = torchvision.transforms.Compose([
        torchvision.transforms.ToTensor(),
        torchvision.transforms.Normalize(mean=[0.5071, 0.4867, 0.4408], std=[0.2675, 0.2565, 0.2761]),
    ])
    cifar100_test = torchvision.datasets.CIFAR100(root='./data', train=False, download=True, transform=transform)
    dataloader = DataLoader(Subset(cifar100_test, range(100)), batch_size=1, shuffle=False)

    return dataloader
```

这里返回`DataLoader`对象以方便数据处理，`Subset(cifar100_test, range(100))`控制了测试数据量，`batch_size`表示每次推理的样本数，`shuffle`表示是否打乱。

## 2 项目运行

在云端运行`cloud_socket.py`开启服务器，参数`-i`为主机的ip地址，可以运行`ifconfig`命令查看ip地址。`-p`为开放端口，默认9999。示例：

```bash
python3 cloud_socket.py -i 127.0.0.1 -p 9999
```

边缘设备运行`edge_socket.py`连接服务器，`-i`和`-p`与主机保持一致，两者要处于同一局域网下（校园网即可）。`-pt`表示模型的划分点，即在某一层之后划分模型。如果不指定划分点，将根据边缘上的划分策略确定。示例：

```bash
python3 edge_socket.py -i 127.0.0.1 -p 9999 -pt 20
```

边缘客户端连接服务器后会自动开始推理过程。边缘会将划分点传至云端，两端根据划分点各自部署模型的前后部分。然后边缘会自动读取数据开始推理，将中间结果传至云端继续推理。云端传回最终预测结果，并与真实标签一起显示在边缘上。

## 3 运行结果

云端和边缘会分别统计推理时延。全部数据推理完成后，云端会传回传输时延和推理时延，并在边缘上展示准确率、平均传输时延和平均推理时延。

> Accuracy: 80.00%  
> Avg edge inference latency: 49.663 ms  
> Avg transmission latency: 7.930 ms  
> Avg cloud inference latency: 0.244 ms  


- 指定不同的划分点，测试在不同层数后切分的平均边缘推理时延、传输时延、云端推理时延和准确率。

| 划分点 | 平均边缘<br>推理时延 (ms) | 平均传输<br>时延 (ms) | 平均云端<br>推理时延 (ms) | 准确率 |
|:---:|:---:|:---:|:---:|:---:|
| 0 | | | | |
| 3 | | | | |
| 10 | | | | |
| 16 | | | | |
| 20 | | | | |
| 22 | | | | |

- 实现自动划分策略，比如根据[[ASPLOS '17] Neurosurgeon](https://dl.acm.org/doi/10.1145/3037697.3037698)实现静态网络下寻找最优的线性划分点。

- 可以尝试不同的模型和数据集。
