package monitor

// func (md *Driver) podapplyprocess() {
// 	for podapplyinfo := range md.podapplychan {
// 		beginapplytime := time.Now()

// 		// Pod 配置文件下载
// 		podYaml, err := api.ADriver.PodYamlDownload(podapplyinfo.YamlHash)
// 		if err != nil {
// 			err = fmt.Errorf("podyaml download error {%v}", err)
// 			log.Error(err)
// 			if err = topichandler.TDriver.PubPodApplyResult(podapplyinfo.GroupID, false, err.Error(), beginapplytime.UnixNano(), time.Now().UnixNano()); err != nil {
// 				log.Errorf("pod apply result upload error {%v}", err)
// 				md.errchan <- err
// 			}
// 			continue
// 		}

// 		// 创建命名空间
// 		if err := edgenode.EDriver.NamespaceCreateIfNotExist(podapplyinfo.Namespace); err != nil {
// 			err = fmt.Errorf("namespace create error {%v}", err)
// 			log.Error(err)
// 			if err = topichandler.TDriver.PubPodApplyResult(podapplyinfo.GroupID, false, err.Error(), beginapplytime.UnixNano(), time.Now().UnixNano()); err != nil {
// 				log.Errorf("pod apply result upload error {%v}", err)
// 				md.errchan <- err
// 			}
// 			continue
// 		}

// 		// Pod部署
// 		if err := edgenode.EDriver.PodApply(podYaml, podapplyinfo.Namespace); err != nil {
// 			err = fmt.Errorf("pod apply error {%v}", err)
// 			log.Error(err)
// 			if err = topichandler.TDriver.PubPodApplyResult(podapplyinfo.GroupID, false, err.Error(), beginapplytime.UnixNano(), time.Now().UnixNano()); err != nil {
// 				log.Errorf("pod apply result upload error {%v}", err)
// 				md.errchan <- err
// 			}
// 			continue
// 		}
// 	}
// }

// func (md *Driver) podapplystartup() {
// 	for index := 0; index < runtime.NumCPU(); index++ {
// 		log.Debugf("pod apply process {%v} start up", index)
// 		go md.podapplyprocess()
// 	}
// }
