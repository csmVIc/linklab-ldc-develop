package edgenode

import (
	"io"

	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/deprecated/scheme"
	"k8s.io/client-go/tools/remotecommand"
)

// PodExec Pod执行
func (ed *Driver) PodExec(namespace string, pod string, container string, commands []string, reader *io.PipeReader, writer *io.PipeWriter) error {

	execReq := ed.clientset.CoreV1().RESTClient().
		Post().
		Resource("pods").
		Namespace(namespace).
		Name(pod).
		SubResource("exec")
	execOption := &v1.PodExecOptions{
		Command: commands,
		Stdin:   true,
		Stdout:  true,
		Stderr:  true,
		TTY:     true,
	}
	if len(container) > 0 {
		execOption.Container = container
	}
	execReq.VersionedParams(
		execOption,
		scheme.ParameterCodec,
	)

	executor, err := remotecommand.NewSPDYExecutor(ed.k8sconfig, "POST", execReq.URL())
	if err != nil {
		// err = fmt.Errorf("remotecommand.NewSPDYExecutor error {%v}", err)
		log.Error(err)
		return err
	}

	if err = executor.Stream(remotecommand.StreamOptions{
		Stdin:             reader,
		Stdout:            writer,
		Stderr:            writer,
		Tty:               true,
		TerminalSizeQueue: nil,
	}); err != nil {
		// err = fmt.Errorf("executor.Stream error {%v}", err)
		log.Error(err)
		return err
	}

	return nil
}
