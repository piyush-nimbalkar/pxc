/*
Copyright © 2019 Portworx

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package kubernetes

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

// PortForwarder provides a way to forward a local port to Portworx endpoint
type PortForwarder interface {
	Endpoint() string
	Start() error
	Stop() error
}

// KubectlPortForwarder object
type KubectlPortForwarder struct {
	kubeconfig string
	endpoint   string
	cmd        *exec.Cmd
}

// NewKubectlPortForwarder forwards a local port to the Portworx gRPC SDK endpoint
// through the Kubernetes API server using kubectl
// If kubeconfig is not provided, then kubectl will use the default kubeconfig
func NewKubectlPortForwarder(kubeconfig string) PortForwarder {
	return newKubectlPortForwarder(kubeconfig)
}

func newKubectlPortForwarder(kubeconfig string) *KubectlPortForwarder {
	return &KubectlPortForwarder{
		kubeconfig: kubeconfig,
	}
}

// Start creates the portforward using kubectl
func (p *KubectlPortForwarder) Start() error {
	// start kubectl port forward
	// TODO: What if they do not have access to the service?
	var args string
	if len(p.kubeconfig) != 0 {
		args = "--kubeconfig=" + p.kubeconfig + " "
	}
	args = args + "-n kube-system port-forward svc/portworx-api :9020"
	cmd := exec.Command("kubectl", strings.Split(args, " ")...)

	// Setup to read port
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logrus.Errorf("Error while executing [%s]: %v", cmd.String(), err)
		return fmt.Errorf("Unable to setup kubectl: %v", err)
	}

	// Start the port forward process
	err = cmd.Start()
	if err != nil {
		logrus.Errorf("Error while executing [%s]: %v", cmd.String(), err)
		return fmt.Errorf("Unable to execute kubectl. Please make sure kubectl is in your path")
	}
	p.cmd = cmd

	// Read the port
	buf := make([]byte, 1024, 1024)
	n, err := stdout.Read(buf[:])
	if err != nil || n < 0 {
		logrus.Warningf("Error: read[%d] from buffer: %v", n, err)
		return fmt.Errorf("Failed to setup connection to Portworx cluster")
	}
	sbuf := string(buf[:n])
	index := strings.Index(sbuf, "127.0.0.1:")
	if index < 0 {
		p.Stop()
		logrus.Warningf("Unable to find 127.0.0.1: in [%s]", sbuf)
		return fmt.Errorf("Failed to determine endpoint information")
	}

	// Set endpoint
	p.endpoint = strings.Split(sbuf[index:], " ")[0]
	logrus.Infof("Connected to %s", p.endpoint)
	logrus.Debugf("Read %d bytes", n)
	logrus.Debugf("Output: %s", sbuf)

	return nil
}

// Stop ends the session
func (p *KubectlPortForwarder) Stop() error {
	logrus.Debug("Port forwarding stopped")
	if p.cmd != nil {
		return p.cmd.Process.Kill()
	}
	return nil
}

// Endpoint returns the gRPC endpoint to the SDK
func (p *KubectlPortForwarder) Endpoint() string {
	return p.endpoint
}