package edgenode

import (
    "bufio"
    "context"
    "fmt"
    "time"

    log "github.com/sirupsen/logrus"
    v1 "k8s.io/api/core/v1"
)

// GetPodLog Pod日志
func (ed *Driver) GetPodLog(namespace string, pod string, container string, outputchan chan<- string, exitchan <-chan bool) error {
    log.Debugf("开始获取Pod日志，命名空间: %s, Pod: %s, 容器: %s", namespace, pod, container)

    logoption := &v1.PodLogOptions{
        Follow: true,
    }
    if len(container) > 0 {
        logoption.Container = container
    }

    stream, err := ed.clientset.CoreV1().Pods(namespace).GetLogs(pod, logoption).Stream(context.TODO())
    if err != nil {
        log.Errorf("获取日志流失败: %v", err)
        return err
    }
    defer func() {
        log.Debugf("关闭日志流")
        stream.Close()
    }()
	// linereader.Scan()逐行读取日志内容。
    linereader := bufio.NewScanner(stream)
    // 增加scanner的buffer大小
	// 创建一个64kb的字节切片，作为scanner的初始缓冲区
    buf := make([]byte, 0, 64*1024)
	// 调用scanner方法，设置scanner的缓冲区，第一个参数是初始缓冲区，第二个参数是最大缓冲区(1MB)
    linereader.Buffer(buf, 1024*1024)
    
    donechan := make(chan error, 1)
    go func() {
        lineCount := 0
        defer func() {
            log.Debugf("日志读取协程退出，共读取 %d 行", lineCount)
        }()

        for linereader.Scan() {
            line := linereader.Text()
            lineCount++
            select {
            case outputchan <- line:
                if lineCount%1000 == 0 {
                    log.Debugf("已读取 %d 行日志", lineCount)
                }
                continue
            case <-time.After(time.Second * time.Duration(ed.info.PodLog.ChanTimeOut)):
                err := fmt.Errorf("日志写入通道超时 命名空间: %s, Pod: %s, 容器: %s, 超时时间: %d秒", 
                    namespace, pod, container, ed.info.PodLog.ChanTimeOut)
                log.Error(err)
                donechan <- err
                return
            }
        }

        // 检查Scanner是否有错误
        if err := linereader.Err(); err != nil {
            log.Errorf("日志扫描错误: %v", err)
            donechan <- err
            return
        }

        log.Debugf("日志读取完成，总共读取 %d 行", lineCount)
        donechan <- nil
    }()

    select {
    case err := <-donechan:
        return err
    case <-exitchan:
        log.Debugf("收到退出信号，停止日志读取")
        return nil
    }
}