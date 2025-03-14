package edgenode

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	log "github.com/sirupsen/logrus"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	// "linklab/device-control-v2/base-library/cache"
	// "linklab/device-control-v2/base-library/cache/value"
)

// GetPodNamesByNamespace 根据命名空间下查找对应的Pod名称列表
// func (ed *Driver) GetPodNamesByNamespace(namespace string) ([]string, error) {
// 	log.Debugf("开始获取命名空间下的Pod名称列表 - 命名空间: {%v}", namespace)
// 	// 获取redis连接
// 	rdb, err := cache.Cdriver.GetRdb()
// 	if err != nil {
// 		err = fmt.Errorf("获取redis连接失败: {%v}", err)
// 		log.Error(err)
// 		return nil, err
// 	}
// 	log.Debugf("获取redis连接成功")
// 	// 获取所有客户端ID
// 	clientIDs, err := rdb.Keys(context.TODO(), "edgenodes:resource:*").Result()
// 	if err != nil {
// 		err = fmt.Errorf("获取边缘节点资源列表失败: %v", err)
// 		log.Error(err)
// 		return nil, err
// 	}
// 	log.Debugf("找到边缘节点资源列表，总数：{%d}", len(clientIDs))
// 	podNames := []string{}
// 	for _, clientKey := range clientIDs {
// 		// 从键名中国呢提取客户端ID
// 		clientID := strings.TrimPrefix(clientKey, "edgenodes:resource:")
// 		log.Debugf("处理边缘节点资源 - 客户端ID: {%v}", clientID)

// 		//获取该客户端下的所有pod信息
// 		podsMap, err := rdb.HGetAll(context.TODO(), fmt.Sprintf("edgenodes:resource:%s", clientID)).Result()
// 		if err != nil {
// 			log.Warnf("获取客户端Pod资源信息失败 - 客户端ID：{%v},  错误：{%v}", clientID, err)
// 			continue
// 		}

// 		// 遍历所有Pod
// 		for key, _ := range podsMap {
// 			parts := strings.Split(key, ":")
// 			if len(parts) == 2 && parts[0] == namespace{
// 				podName := parts[1]
// 				log.Debugf("找到匹配的Pod - 客户端ID：{%v}, 命名空间：{%v}, Pod名称：{%v}", clientID, namespace, podName)
// 				podNames = append(podNames, podName)
// 			}
// 		}
// 	}
// 	if len(podNames) == 0 {
// 		log.Debugf("未找到命名空间下的Pod - 命名空间: {%v}", namespace)
// 	}else{
// 		log.Debugf("命名空间查找完成 - 命名空间：{%v},找到的Pod数量：{%d}", namespace, len(podNames))
// 	}

// 	return podNames, nil
// }


// 部署和更新 Pod 的功能
// PodApply Pod部署
func (ed *Driver) PodApply(podYaml []byte, namespace string, env []apiv1.EnvVar, edgeregistry bool, nodeselector map[string]string, nodeaddselector string) error {
    log.Debugf("开始处理Pod YAML配置 - 命名空间: {%v}", namespace)
	podJson, err := yaml.ToJSON(podYaml)
	if err != nil {
		err = fmt.Errorf("yaml.ToJSON error {%v}", err)
		log.Error(err)
		return err
	}
	//解析pod配置
	podConfig := &apiv1.Pod{}
	if err := json.Unmarshal(podJson, podConfig); err != nil {
		err = fmt.Errorf("json.Unmarshal error {%v}", err)
		log.Error(err)
		return err
	}
    log.Debugf("Pod基本信息 - 名称: {%v}, 容器数量: {%d}", 
        podConfig.Name, 
        len(podConfig.Spec.Containers))

	// 检查是否需要启用edge_socket
	enableEdgeSocket := false
	// 环境变量
	if env != nil && len(env) > 0 {
		log.Debugf("添加环境变量，数量: {%d}", len(env))
		// podConfig.Spec.Containers[0].Env = append(podConfig.Spec.Containers[0].Env, env...)
		for index := range podConfig.Spec.Containers {
			log.Debugf("容器[%d] {%v} 添加环境变量 - env:%v", 
				index, 
				podConfig.Spec.Containers[index].Name,
				env)
			podConfig.Spec.Containers[index].Env = append(podConfig.Spec.Containers[index].Env, env...)
			// 检查是否有 ENABLE_EDGE_SOCKET 环境变量
			for _, envVar := range podConfig.Spec.Containers[index].Env {
				if envVar.Name == "ENABLE_EDGE_SOCKET" {
					enableEdgeSocket = true
					break
				}
			}
		}
	}
	//如果启用的情况下，动态添加args
	if enableEdgeSocket {
		// podNames, err := ed.GetPodNamesByNamespace(namespace)
		// if err != nil {
		// 	log.Warnf("获取命名空间下pod列表失败：{%v}", err)
		// }
		// log.Debugf("命名空间 {%s} 下找到的Pod列表: %v", namespace, podNames)
		for index := range podConfig.Spec.Containers {
			podConfig.Spec.Containers[index].Args = []string{
				"/bin/bash",
				"-c",
				fmt.Sprintf("python3 edge_socket.py --podname cloud-test --namespace %s --port 99", namespace),
			}
			log.Debugf("容器[%d] 启动参数已经设置: %v",
				index,
				podConfig.Spec.Containers[index].Args)
		}
	}
	// 镜像地址
	if edgeregistry {
		for index := range podConfig.Spec.Containers {
			originalImage := podConfig.Spec.Containers[index].Image
			notReplace := false
			for _, envVal := range podConfig.Spec.Containers[index].Env {
				if envVal.Name == "USE_PUB_REGISTRY" && envVal.Value == "true" {
					log.Debugf("容器[%d]配置使用公共仓库，跳过镜像地址替换", index)
					notReplace = true
					break
				}
			}
			if notReplace {
				break
			}
			// 替换镜像地址
			podConfig.Spec.Containers[index].Image = fmt.Sprintf("%v/%v/%v", ed.info.ImageBuild.RegistryAddress, namespace, podConfig.Spec.Containers[index].Image)
			log.Debugf("容器[%d]镜像地址更新 - 原地址: {%v}, 新地址: {%v}", 
				index, 
				originalImage, 
				podConfig.Spec.Containers[index].Image)
		}
	}

	if podConfig.Labels == nil {
		podConfig.Labels = map[string]string{}
	}
	podConfig.Labels["name"] = podConfig.Name
	// 节点选择器
	if nodeselector != nil {
		log.Debugf("使用自定义节点选择器: %v", nodeselector)
		if podConfig.Spec.NodeSelector == nil {
			podConfig.Spec.NodeSelector = map[string]string{}
		}
		for key, val := range nodeselector {
			podConfig.Spec.NodeSelector[key] = val
		}
	}else{
		log.Debugf("使用默认节点选择器配置")
		if podConfig.Spec.NodeSelector == nil {
			log.Debugf("设置默认树莓派节点选择器")
			podConfig.Spec.NodeSelector = map[string]string{}
			podConfig.Spec.NodeSelector["linklab.edgetype"] = "raspberrypi4bextend"
		}else{
			if strings.ToLower(podConfig.Spec.NodeSelector["linklab.edgetype"]) == "jetsonnano"{
				podConfig.Spec.NodeSelector["linklab.edgetype"] = nodeaddselector
			}
		}
	}
    log.Debugf("最终节点选择器配置-----: %v", podConfig.Spec.NodeSelector)

	podClient := ed.clientset.CoreV1().Pods(namespace)
	if _, err := podClient.Get(context.TODO(), podConfig.Name, metav1.GetOptions{}); err == nil {
		// 更新Pod
		log.Debugf("检测到已存在的Pod，准备更新 - 命名空间: {%v}, Pod名称: {%v}", 
			namespace, 
			podConfig.Name)
		podApi, err := podClient.Update(context.TODO(), podConfig, metav1.UpdateOptions{})
		if err != nil {
			// err = fmt.Errorf("podClient.Update error {%v}", err)
			log.Error(err)
			return err
		}
		log.Debugf("namespace {%v} pod {%v} update success", podApi.Namespace, podApi.Name)
	} else {
		// 创建Pod
		log.Debugf("未检测到已存在的Pod，准备创建 - 命名空间: {%v}, Pod名称: {%v}", 
			namespace, 
			podConfig.Name)
		podApi, err := podClient.Create(context.TODO(), podConfig, metav1.CreateOptions{})
		if err != nil {
			// err = fmt.Errorf("podClient.Create error {%v}", err)
			log.Error(err)
			return err
		}
		log.Debugf("Pod创建成功 - 命名空间: {%v}, Pod名称: {%v}", 
			podApi.Namespace, 
			podApi.Name)
	}
	return nil
}
