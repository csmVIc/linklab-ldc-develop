package edgenode

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	apiv1 "k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/yaml"

	"context"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "math/rand"
    "time"
)

// LinkPodIngress 关联PodIngress
func (ed *Driver) LinkPodIngress(podYaml []byte, namespace string) (*map[string]string, error) {
    log.Debugf("文件<linkingress> 开始处理Pod Ingress配置 - 命名空间: %s", namespace)

	podJson, err := yaml.ToJSON(podYaml)
	if err != nil {
		err = fmt.Errorf("yaml.ToJSON error {%v}", err)
		log.Error(err)
		return nil, err
	}
    log.Debugf("YAML转JSON成功")

	podConfig := &apiv1.Pod{}
	if err := json.Unmarshal(podJson, podConfig); err != nil {
		err = fmt.Errorf("json.Unmarshal error {%v}", err)
		log.Error(err)
		return nil, err
	}
    log.Debugf("解析Pod配置成功 - Pod名称: %s", podConfig.Name)

	// Service
	svcConfig := &corev1.Service{}
	svcConfig.Name = podConfig.Name
	svcConfig.Namespace = namespace
	svcConfig.Spec.Type = corev1.ServiceTypeNodePort
	svcConfig.Spec.Ports = []corev1.ServicePort{}
	// 服务没有暴露任何端口号
	// if len(svcConfig.Spec.Ports) < 1 {
	// 	log.Debugf("未检测到任何端口配置，跳过Service和Ingress创建")
	// 	return nil, nil
	// }
	// 打印端口号内容 - len(svcConfig.Spec.Ports) = 0 有待了解这一块
	log.Debugf("len(svcConfig.Spec.Ports): %d", len(svcConfig.Spec.Ports))
	// 选择器使用app
	svcConfig.Spec.Selector = map[string]string{"app": podConfig.Name}

	log.Debugf("创建Service配置:")
    log.Debugf("- 服务名称: %s", svcConfig.Name)
    log.Debugf("- 服务类型: %s", svcConfig.Spec.Type)
    log.Debugf("- 选择器: %+v", svcConfig.Spec.Selector)

	// i用来保存随机端口
	getPort := int32(0)
	// 为每个容器端口配置服务端口
	for _, containerInfo := range podConfig.Spec.Containers {
		// 获取随机可用端口
		randomPort, err := getRandomAvailablePort(ed, namespace, 30000, 32767)
		if err != nil {
			log.Error("获取随机端口失败: %v", err)
			return nil, err
		}
		getPort = randomPort
		// 添加服务端口配置
		svcConfig.Spec.Ports = append(svcConfig.Spec.Ports,corev1.ServicePort{
			Name: fmt.Sprintf("port-%d", randomPort),
			Protocol: corev1.ProtocolTCP,
			Port: int32(containerInfo.Ports[0].ContainerPort), 	// 集群内部访问 Service 的端口。如果没有设置，默认与 TargetPort 相同
			// TargetPort: intstr.FromInt(int(randomPort)),
			TargetPort: intstr.FromInt(int(containerInfo.Ports[0].ContainerPort)), // Service 转发流量到 Pod 的端口。如果没有设置，默认与 Port 相同
			NodePort: randomPort,	// 集群节点上暴露的端口，用于外部访问。如果没有设置，默认不暴露
		})
		log.Debugf("port: %d, targetPort: %d, NodePort: %d", svcConfig.Spec.Ports[0].Port, svcConfig.Spec.Ports[0].TargetPort, svcConfig.Spec.Ports[0].NodePort)
		log.Debugf("Service服务配置详情:")
		log.Debugf("- 服务类型：%s", svcConfig.Spec.Type)
		log.Debugf("- 选择器：%+v", svcConfig.Spec.Selector)
		log.Debugf("- 端口配置：")
		for _, port := range svcConfig.Spec.Ports {
			log.Debugf(" - %s: %d -> %v",port.Name,port.Port, port.TargetPort)
		}
	}
	log.Debugf("getProt: %d", getPort)
	// Ingress
	ingressConfig := &v1beta1.Ingress{}
	ingressConfig.Name = podConfig.Name
	ingressConfig.Namespace = namespace
	ingressConfig.Annotations = map[string]string{"nginx.ingress.kubernetes.io/rewrite-target": "/$2"}
	ingressConfig.Spec.Rules = []v1beta1.IngressRule{}
	ingressConfig.Spec.Rules = append(ingressConfig.Spec.Rules, v1beta1.IngressRule{
		Host: ed.info.Ingress.Domain,
		IngressRuleValue: v1beta1.IngressRuleValue{
			HTTP: &v1beta1.HTTPIngressRuleValue{
				Paths: []v1beta1.HTTPIngressPath{},
			},
		},
	})
	log.Debugf("创建Ingress配置:")
    log.Debugf("- Ingress名称: %s", ingressConfig.Name)
    log.Debugf("- Ingress域名: %s", ed.info.Ingress.Domain)
    log.Debugf("- Ingress注解: %+v", ingressConfig.Annotations)

	// 回复结果
	res := map[string]string{}

	// 端口映射关系
	log.Debugf("开始处理容器端口映射:")
	for i, containerInfo := range podConfig.Spec.Containers {
		log.Debugf("处理容器[%d] - 名称: %s", i, containerInfo.Name)
		for _, portInfo := range containerInfo.Ports {
            log.Debugf("- 处理端口: %d", portInfo.ContainerPort)

            // 构建Ingress路径
            ingressPath := fmt.Sprintf("/%s/%s/port%v(/|$)(.*)", namespace, svcConfig.Name, portInfo.ContainerPort)
			ingressConfig.Spec.Rules[0].HTTP.Paths = append(ingressConfig.Spec.Rules[0].HTTP.Paths, v1beta1.HTTPIngressPath{
				Path: fmt.Sprintf("/%s/%s/port%v(/|$)(.*)", namespace, svcConfig.Name, portInfo.ContainerPort),
				Backend: v1beta1.IngressBackend{
					ServiceName: svcConfig.Name,
					ServicePort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: portInfo.ContainerPort,
					},
				},
			})
            log.Debugf("  - 添加Ingress路径: %s", ingressPath)
			// 添加映射关系 - pod暴露的端口 -> ingress路径
			res[fmt.Sprintf("%v", getPort)] = fmt.Sprintf("http://%v%v", ed.info.Ingress.Domain, fmt.Sprintf("/%s/%s/port%v", namespace, svcConfig.Name, portInfo.ContainerPort))

		}
	}

	// // 服务没有暴露任何端口号
	// if len(svcConfig.Spec.Ports) < 1 {
	// 	log.Debugf("未检测到任何端口配置，跳过Service和Ingress创建")
	// 	return nil, nil
	// }

	// 创建服务
	log.Debugf("开始创建Service - 名称: %s, 端口数量: %d", svcConfig.Name, len(svcConfig.Spec.Ports))
	if err := ed.ServiceApplyByConfig(svcConfig, namespace); err != nil {
		log.Error(err)
		return nil, err
	}
    log.Debugf("Service创建成功")

	// 创建Ingress
	log.Debugf("开始创建Ingress - 名称: %s, 路径数量: %d", ingressConfig.Name, len(ingressConfig.Spec.Rules[0].HTTP.Paths))
	if err := ed.IngressApplyByConfig(ingressConfig, namespace); err != nil {
		log.Error(err)
		return nil, err
	}
	log.Debugf("Ingress创建成功")

    log.Debugf("端口映射配置完成，总计 %d 个映射", len(res))
    log.Debugf("linkgress文件阅读完毕")

	return &res, nil
}

// 获取随机可用端口的函数
func getRandomAvailablePort(ed *Driver, namespace string, min, max int32) (int32, error) {
    // 获取当前命名空间下的所有 service
    services, err := ed.clientset.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
    if err != nil {
        return 0, fmt.Errorf("list services error: %v", err)
    }

    // 创建已使用端口的集合
    usedPorts := make(map[int32]bool)
    for _, svc := range services.Items {
        for _, port := range svc.Spec.Ports {
            usedPorts[port.Port] = true
        }
    }

    // 随机生成端口直到找到未使用的
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    for tries := 0; tries < 50; tries++ { // 设置最大尝试次数避免死循环
        port := min + r.Int31n(max-min+1)
        if !usedPorts[port] {
            return port, nil
        }
    }
    
    return 0, fmt.Errorf("no available ports in range %d-%d", min, max)
}